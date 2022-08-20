package tables

import (
	"fmt"
    "github.com/nboughton/go-roll"
)

type Table struct {
	Name      string
	Path      string
	Size      int
	Err       int
	Variables map[string]string // Keyword/value pairs
	Groups    map[string]Group
}

func NewTable(name string) *Table{
    return &Table{
        Name: name,
        Variables: make(map[string]string),
        Groups: make(map[string]Group),
}
}


func (t *Table) AddVariable(name string, value string) error {
    if t.Variables == nil{
        t.Variables = make(map[string]string)
    }
	_, exists := t.Variables[name]
	if exists {
		return fmt.Errorf("Variable already exists %s", name)
	}
	t.Variables[name] = value
	return nil
}

func (t *Table) AddGroup(g *Group) error {
    if t.Groups == nil{
        t.Groups = make(map[string]Group)
    }
	name := g.Name
	_, exists := t.Groups[name]
	if exists {
		return fmt.Errorf("Variable already exists %s", name)
	}
	t.Groups[name] = *g
	return nil
}

type Group struct {
	//Name   string
	Prefix string
	Suffix string
    roll.Table
}

func (g *Group)Close(){
    if g == nil{
        return
    }
}

func (g *Group)AddItem(min,max int,l string) {
    g.Items = append(g.Items, 
        roll.TableItem{ Match: []int{min,max}, Text: l })
}

