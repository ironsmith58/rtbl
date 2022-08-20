package tables

import (
	"strings"
)

/*
 * Make a Table Name from the filename of
 * the file path
 */
func makeName(path, root string) string {
	name := strings.TrimSuffix(path, ".tab")
	name = strings.TrimPrefix(name, root)
	name = strings.TrimPrefix(name, "/")
	name = strings.Replace(name, "/", ".", -1)
	return name
}
