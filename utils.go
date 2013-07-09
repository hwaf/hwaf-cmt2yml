package main

import (
	"os"
	"regexp"
	"strings"
)

func path_exists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

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

// str_is_in_slice returns true if str is in the given slice of strings
func str_is_in_slice(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// re_is_in_slice_suffix returns true if an element in the given slice of strings is a prefix of value.
func re_is_in_slice_suffix(slice []string, macro, pattern string) bool {
	for _, s := range slice {
		pat := regexp.MustCompile(s + pattern)
		if pat.MatchString(macro) {
			return true
		}
	}
	return false
}

// EOF
