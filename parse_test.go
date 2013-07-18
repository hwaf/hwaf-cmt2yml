package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParseLine(t *testing.T) {

	for _, v := range []struct {
		fname    string
		expected []string
	}{
		{
			fname: "testdata/correct_gcc_config_tag.txt",
			expected: []string{
				"macro", "correct_gcc_config_tag", "i686-slc5-gcc43-opt",
				"x86_64", "x86_64-slc5",
			},
		},
		{
			fname: "testdata/buggy_gcc_config_tag.txt",
			expected: []string{
				"macro", "buggy_gcc_config_tag", "i686-slc5-gcc43-opt",
				"x86_64", "x86_64-slc5",
			},
		},
		// FIXME!!
		{
			fname: "testdata/pp_cppflags.txt",
			expected: []string{
				"macro_append", "pp_cppflags", `-DTDAQ_PACKAGE_NAME="$(package)"`,
			},
		},
		{
			fname: "testdata/lib_suffix.txt",
			expected: []string{
				"macro", "lib_suffix", ".so",
				"ppc-rtems-rce405", ".a",
			},
		},
		{
			fname: "testdata/app_inst_group.txt",
			expected: []string{"macro", "app_inst_group", "inst",
				"ppc-rtems-rce405", "____"},
		},
		{
			fname: "testdata/c_opt_flags.txt",
			expected: []string{
				"macro", "c_opt_flags", "",
				"x86_64-slc4&gcc34", "-O2 -mtune=nocona -fomit-frame-pointer -fno-exceptions",
				"x86_64-slc4&gcc43", "-O2 -mtune=nocona -fomit-frame-pointer -fno-exceptions",
				"x86_64-slc5&gcc43", "-O2 -mtune=nocona -fomit-frame-pointer -fno-exceptions",
				"x86_64-slc5&gcc46", "-O2 -mtune=core2 -ftree-vectorize -ftree-vectorizer-verbose=2 -fomit-frame-pointer",
				"x86_64-slc6&gcc4", "-O2 -mtune=nocona -fomit-frame-pointer -fno-exceptions",
				"gcc32", "-O2 -march=i686 -mcpu=i686 -funroll-loops -falign-loops -falign-jumps -falign-functions -fno-exceptions",
				"gcc323", "-O2 -mcpu=pentium4 -funroll-loops -falign-loops -falign-jumps -falign-functions -fno-exceptions",
				"i686-slc4&gcc34", "-O2 -mtune=pentium4 -funroll-loops -falign-loops -falign-jumps -falign-functions -fno-exceptions",
				"i686-slc4&gcc43", "-O2 -mtune=pentium4 -funroll-loops -falign-loops -falign-jumps -falign-functions -fno-exceptions",
				"i686-slc5&gcc4", "-O2 -mtune=pentium4 -funroll-loops -falign-loops -falign-jumps -falign-functions -fno-exceptions",
				"i686-slc6&gcc4", "-O2 -mtune=pentium4 -funroll-loops -falign-loops -falign-jumps -falign-functions -fno-exceptions",
				"icc8", "-O2 -mtune=pentium4",
				"icc11", "-O2 -axSSSE3 -vec-report1 -par-report1 -parallel",
				"powerpc-rtems-gcc43", "-O4 -mlongcall -msoft-float",
			},
		},
		{
			fname: "testdata/cpp_opt_flags.txt",
			expected: []string{
				"macro", "cpp_opt_flags", "",
				"x86_64-slc4&gcc34", "-O2 -mtune=nocona -fomit-frame-pointer",
				"x86_64-slc4&gcc43", "-O2 -mtune=nocona -fomit-frame-pointer",
				"x86_64-slc6&gcc46", "-O2 -mtune=core2 -std=c++0x -ftree-vectorize -ftree-vectorizer-verbose=2 -fomit-frame-pointer",
				"x86_64-slc6&gcc47", "-O2 -mtune=core2 -std=c++11 -ftree-vectorize -ftree-vectorizer-verbose=2 -fomit-frame-pointer",
				"x86_64-slc6&gcc4", "-O2 -mtune=nocona -fomit-frame-pointer",
				"x86_64-slc5&gcc43", "-O2 -mtune=nocona -fomit-frame-pointer",
				"x86_64-slc5&gcc46", "-O2 -mtune=core2 -std=c++0x -ftree-vectorize -ftree-vectorizer-verbose=2 -fomit-frame-pointer",
				"x86_64-slc5&gcc47", "-O2 -mtune=core2 -std=c++11 -ftree-vectorize -ftree-vectorizer-verbose=2 -fomit-frame-pointer",
				"gcc32", "-O2 -march=i686 -mcpu=i686 -funroll-loops -falign-loops -falign-jumps -falign-functions",
				"gcc323", "-O2 -mcpu=pentium4 -funroll-loops -falign-loops -falign-jumps -falign-functions",
				"i686-slc4&gcc34", "-O2 -mtune=pentium4 -funroll-loops -falign-loops -falign-jumps -falign-functions",
				"i686-slc4&gcc43", "-O2 -mtune=pentium4 -funroll-loops -falign-loops -falign-jumps -falign-functions",
				"i686-slc5&gcc4", "-O2 -mtune=pentium4 -funroll-loops -falign-loops -falign-jumps -falign-functions",
				"i686-slc6&gcc4", "-O2 -mtune=pentium4 -funroll-loops -falign-loops -falign-jumps -falign-functions",
				"icc8", "-O2 -mtune=pentium4",
				"icc11", "-O2 -axSSSE3 -vec-report1 -par-report1 -parallel",
				"powerpc-rtems-gcc43", "-O4 -mlongcall -msoft-float",
			},
		},
		{
			fname: "testdata/cppappcomp.txt",
			expected: []string{
				"macro", "cppappcomp", "$(cppcomp)",
				"ppc-rtems-rce405", "echo",
			},
		},
	} {
		p, err := NewParser(v.fname)
		if err != nil {
			t.Fatalf(err.Error())
		}
		err = p.run()
		out := p.tokens
		if !reflect.DeepEqual(out, v.expected) {
			s_expected := make([]string, 0, len(v.expected))
			for _, vv := range v.expected {
				s_expected = append(s_expected, fmt.Sprintf("%q", vv))
			}
			s_out := make([]string, 0, len(out))
			for _, vv := range out {
				s_out = append(s_out, fmt.Sprintf("%q", vv))
			}
			t.Fatalf(
				"\nexpected: %v\ngot:      %v\n",
				strings.Join(s_expected, ", "),
				strings.Join(s_out, ", "),
			)
		}
	}
}
