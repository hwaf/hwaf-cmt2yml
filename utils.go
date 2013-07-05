package main

import (
	"strings"
)

// str_split slices s into all (non-empty) substrings separated by sep
func str_split(s, sep string) []string {
	strs := strings.Split(s, sep)
	out := make([]string, 0, len(strs))
	for _, str := range strs {
		str = strings.Trim(str, " \t")
		if len(str) == 0 {
			continue
		}
		out = append(out, str)
	}
	return out
}

// EOF
