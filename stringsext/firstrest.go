package stringsext

import (
	"strings"
)

func First(s string, sep string) string {
	// return first word
	i := strings.Index(s, sep)
	if i != -1 {
		return s[:i]
	}
	return s
}

func Rest(s string, sep string) string {
	// return all words after first word
	i := strings.Index(s, sep)
	if i != -1 {
		return s[i+1:]
	}
	return ""
}
