package tables

import (
	"fmt"
	"io/ioutil"
	"os"

	//"errors"
	"path/filepath"
)

func List(root string) ([]Table, error) {

	var result []Table

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() {
			name := makeName(path, root)
			lines, err := ioutil.ReadFile(path)
			if err == nil {
				table := NewTable(name)
				table.Path = path
				table.Size = len(lines)
				result = append(result, *table)
			} else {
				table := NewTable(name)
				table.Path = path
				table.Err = 1
				result = append(result, *table)
			}
		}
		return nil
	})
	return result, err
}
