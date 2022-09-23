package tables

import (
	"fmt"
)

type Table struct {
	Name      string
	Path      string
	Size      int
	Err       int
	Header    string            // set by /OutputHeader directive
	Footer    string            // set by /OutputFooter directive
	Variables map[string]string // Keyword/value pairs
	Groups    map[string]*Group
}

func NewTable(name string) *Table {
	return &Table{
		Name:      name,
		Variables: make(map[string]string),
		Groups:    make(map[string]*Group),
	}
}

func (t *Table) AddVariable(name string, value string) error {
	if t.Variables == nil {
		t.Variables = make(map[string]string)
	}
	/*
		_, exists := t.Variables[name]
		if exists {
			return fmt.Errorf("Variable already exists %s", name)
		}
	*/
	t.Variables[name] = value
	return nil
}

func (t *Table) GetVariable(name string) (value string, exists bool) {
	val, exists := t.Variables[name]
	return val, exists
}

func (t *Table) AddGroup(g *Group) error {
	t.Groups[g.Name] = g
	g.Close()
	return nil
}

func (t *Table) GetGroup(name string) (*Group, error) {
	g, ok := t.Groups[name]
	if !ok {
		return nil, fmt.Errorf("table %s has no group named %s", t.Name, name)
	}
	return g, nil
}
