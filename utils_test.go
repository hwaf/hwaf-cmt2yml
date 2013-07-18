package main

import (
	"testing"
)

func toto_is_tata() {
	println("bobo")
}

func TestReInSlice(t *testing.T) {
	macros := []string{
		"ErsBaseStreams",
		"ers",
		"ers_psc",
		"ers_receiver",
		"ers_test",
		"erspy",
	}

	for _, table := range []struct {
		pattern  string
		macro    string
		expected bool
	}{
		{
			pattern:  ".*?",
			macro:    "erspy",
			expected: true,
		},
		{
			pattern:  "(_shlibflags|.*?)",
			macro:    "erspy",
			expected: true,
		},
		{
			pattern:  "_shlibflags",
			macro:    "erspy_shlibflags",
			expected: true,
		},
		{
			pattern:  "_dependencies",
			macro:    "erspy_dependencies",
			expected: true,
		},
		{
			pattern:  ".*?",
			macro:    "ErsBaseStreams",
			expected: true,
		},
		{
			pattern:  "_dependencies",
			macro:    "ErsBaseStreams_dependencies",
			expected: true,
		},
		{
			pattern:  ".*?",
			macro:    "ErsBaseStreamS",
			expected: false,
		},
	} {
		ret := re_is_in_slice_suffix(macros, table.macro, table.pattern)
		if ret != table.expected {
			t.Fatalf(
				"expected %v for value %v. got %v",
				table.expected,
				table.pattern,
				ret,
			)
		}
	}
}

// EOF
