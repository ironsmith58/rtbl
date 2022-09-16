package tables

import (
	"bufio"
	"fmt"
	"os"
	"path"
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

func List(root string, showAll bool) ([]Table, error) {

	var result []Table

	err := filepath.Walk(root, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if filePath == root { // lets not print our starting point
			return nil
		}
		if !info.IsDir() {
			if strings.HasSuffix(filePath, ".tab") {
				name := path.Base(filePath)
				name = strings.TrimSuffix(name, ".tab")
				if name[0] != '~' || showAll {
					fmt.Println("    ", name)
				}
			}
		} else {
			name := path.Base(filePath)
			fmt.Println("Category: ", name)
		}

		return nil
	})
	return result, err
}
