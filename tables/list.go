package tables

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	//"errors"
	"path/filepath"
)

func QuickParse(filePath string) ([]string, error) {
	// read the whole content of file and pass it to file variable, in case of error pass it to err variable

	var groups []string = make([]string, 0, 20)

	file, err := os.Open(filePath)
	if err != nil {
		return groups, fmt.Errorf("Could not open the file due to this %s error \n", err)
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {

		line := fileScanner.Text()
		if len(line) == 0 {
			continue
		}
		if line[0] == ';' || line[0] == ':' {
			name := line[1:]
			idx := strings.Index(name, "#")
			if idx != -1 {
				name = name[:idx]
			}
			groups = append(groups, name)
		}
	}
	return groups, nil
}

func FindTables(root string) (paths []string, err error) {

	var result []string

	err = filepath.Walk(root, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if !info.IsDir() {
			if strings.HasSuffix(filePath, ".tab") {
				result = append(result, filePath)
			}
		}
		return nil
	})
	return result, err
}

type TablesByCategory map[string][]string

func NewTablesByCategory(paths []string) TablesByCategory {

	tables := make(TablesByCategory)

	// make a map of all tables in all catagories
	// tables are paths that end in .tab
	// categories are the containing directory
	// e.g. Names/Greek.tab
	for _, filepath := range paths {
		if strings.HasSuffix(filepath, ".tab") {
			name := path.Base(filepath)
			name = strings.TrimSuffix(name, ".tab")
			dir := path.Base(path.Dir(filepath))
			tables[dir] = append(tables[dir], name)
		}
	}
	return tables
}

type LoadedTable struct {
	path  string
	table *Table
}

type TablePathsByName map[string]*LoadedTable // map table name to paths

func NewTableList(paths []string) TablePathsByName {

	tables := make(TablePathsByName)

	// make a map of all tables in all catagories
	// tables are paths that end in .tab
	// categories are the containing directory
	// e.g. Names/Greek.tab
	for _, filepath := range paths {
		if strings.HasSuffix(filepath, ".tab") {
			name := path.Base(filepath)
			name = strings.TrimSuffix(name, ".tab")
			name = strings.ToLower(name) // hold all names as lower case
			tables[name] = &LoadedTable{filepath, nil}
		}
	}
	return tables
}

func PrintPaths(tables TablesByCategory, showAll bool) {

	// sort the keys(categories) so they
	// show 'in order'
	var keys = make([]string, 0, len(tables))
	for k := range tables {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// now print the structured categories and tables
	for _, cat := range keys {
		fmt.Println(cat)
		sort.Strings(tables[cat])
		for _, table := range tables[cat] {
			fmt.Println("  ", table)
		}
	}
}

// Most import Variable -- holds all table references where
// Parse/Lookup can find tables ....
var TableRegistry TablePathsByName

func LoadAllTables(rootpath string) error {
	paths, err := FindTables(rootpath)
	if err != nil {
		return err
	}
	if len(paths) == 0 {
		return fmt.Errorf("No tables found")
	}
	TableRegistry = NewTableList(paths)
	return nil
}
