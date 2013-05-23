package main

import (
	"bufio"
	"bytes"
	"fmt"
	//"io"
	"os"
	"strings"
)

const dbg_parse_line = false

func fmt_line(data []string) string {
	s := bytes.NewBufferString("[")
	for i, v := range data {
		if i == 0 {
			fmt.Fprintf(s, "%q", v)
		} else {
			fmt.Fprintf(s, ", %q", v)
		}
	}
	fmt.Fprintf(s, "]")
	return string(s.Bytes())
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func scan_line(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanLines(data, atEOF)
	return
	// sz := len(token)
	// if sz > 0 && token[sz-1] == '\\' {
	// 	return
	// }
}

func parse_file(fname string) (*ReqFile, error) {
	fmt.Printf("req=%q\n", fname)
	var err error
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(bufio.NewReader(f))
	if scanner == nil {
		return nil, fmt.Errorf("cmt2yml: nil bufio.Scanner")
	}

	req := ReqFile{Filename: fname}

	bline := []byte{}
	ctx := []string{"public"}
	for scanner.Scan() {
		data := scanner.Bytes()
		data = bytes.TrimSpace(data)

		if len(data) == 0 {
			continue
		}
		if data[0] == '#' {
			continue
		}

		idx := len(data) - 1
		if data[idx] == '\\' {
			bline = append(bline, data[:idx-1]...)
			continue
		} else {
			bline = append(bline, data...)
		}
		//fmt.Printf("%q\n", string(bline))

		var tokens []string
		tokens, err = parse_line(bline)
		if err != nil {
			return nil, err
		}

		switch tokens[0] {
		case tok_PACKAGE:
			req.Package = tokens[1]

		case tok_AUTHOR:
			req.Authors = append(
				req.Authors,
				tokens[1:]...,
			)

		case tok_MANAGER:
			req.Managers = append(
				req.Managers,
				tokens[1:]...,
			)

		case tok_USE:
			use := UsePkg{Package: tokens[1]}
			if len(tokens) > 2 {
				use.Version = tokens[2]
			}
			if len(tokens) > 3 {
				use.Path = tokens[3]
			}
			if len(tokens) > 4 {
				use.Switches = append(use.Switches, tokens[4:]...)
			}
			req.Uses = append(req.Uses, use)

		case tok_MACRO:
			//fmt.Printf("macro: %v\n", fmt_line(tokens))
			vv := Macro{Name: tokens[1]}
			vv.Value = make(map[string]string)
			vv.Value["default"] = tokens[2]
			if len(tokens) > 3 {
				toks := tokens[3:]
				for i := 0; i+1 < len(toks); i += 2 {
					vv.Value[toks[i]] = toks[i+1]
				}
			}
			req.Macros = append(req.Macros, vv)
			//fmt.Printf("macro: %v\n", macro)

		case tok_MACRO_APPEND:
			vv := MacroAppend{Name: tokens[1]}
			vv.Value = make(map[string]string)
			vv.Value["default"] = tokens[2]
			if len(tokens) > 3 {
				toks := tokens[3:]
				for i := 0; i+1 < len(toks); i += 2 {
					vv.Value[toks[i]] = toks[i+1]
				}
			}
			req.MacroAppends = append(req.MacroAppends, vv)

		case tok_MACRO_PREPEND:
			vv := MacroPrepend{Name: tokens[1]}
			vv.Value = make(map[string]string)
			vv.Value["default"] = tokens[2]
			if len(tokens) > 3 {
				toks := tokens[3:]
				for i := 0; i+1 < len(toks); i += 2 {
					vv.Value[toks[i]] = toks[i+1]
				}
			}
			req.MacroPrepends = append(req.MacroPrepends, vv)

		case tok_PATH:
			vv := Path{Name: tokens[1], Value: tokens[2]}
			req.Paths = append(req.Paths, vv)

		case tok_PATH_PREPEND:
			vv := PathPrepend{Name: tokens[1]}
			vv.Value = make(map[string]string)
			vv.Value["default"] = tokens[2]
			if len(tokens) > 3 {
				toks := tokens[3:]
				for i := 0; i+1 < len(toks); i += 2 {
					vv.Value[toks[i]] = toks[i+1]
				}
			}
			req.PathPrepends = append(req.PathPrepends, vv)

		case tok_PATH_APPEND:
			vv := PathAppend{Name: tokens[1]}
			vv.Value = make(map[string]string)
			vv.Value["default"] = tokens[2]
			if len(tokens) > 3 {
				toks := tokens[3:]
				for i := 0; i+1 < len(toks); i += 2 {
					vv.Value[toks[i]] = toks[i+1]
				}
			}
			req.PathAppends = append(req.PathAppends, vv)

		case tok_PATH_REMOVE:
			vv := PathRemove{Name: tokens[1]}
			vv.Value = make(map[string]string)
			vv.Value["default"] = tokens[2]
			if len(tokens) > 3 {
				toks := tokens[3:]
				for i := 0; i+1 < len(toks); i += 2 {
					vv.Value[toks[i]] = toks[i+1]
				}
			}
			req.PathRemoves = append(req.PathRemoves, vv)

		case tok_PATTERN:
			//fmt.Printf("pattern: %v (%d)\n", fmt_line(tokens), len(tokens))
			vv := Pattern{
				Name: tokens[1],
				Def:  strings.Join(tokens[2:], " "),
			}
			req.Patterns = append(req.Patterns, vv)

		case tok_APPLY_PATTERN:
			vv := ApplyPattern{Name: tokens[1]}
			if len(tokens) > 2 {
				vv.Args = append(vv.Args, tokens[2:]...)
			}
			req.ApplyPatterns = append(req.ApplyPatterns, vv)

		case tok_IGNORE_PATTERN:
			vv := IgnorePattern{Name: tokens[1]}
			req.IgnorePatterns = append(req.IgnorePatterns, vv)

		case tok_INCLUDE_DIRS:
			vv := IncludeDirs{Value: tokens[1]}
			req.IncludeDirs = append(req.IncludeDirs, vv)

		case tok_INCLUDE_PATH:
			vv := IncludePaths{Value: tokens[1]}
			req.IncludePaths = append(req.IncludePaths, vv)

		case tok_PRIVATE:
			ctx = append(ctx, tok_PRIVATE)

		case tok_END_PRIVATE:
			ctx = ctx[:len(ctx)-1]

		case tok_PUBLIC:
			ctx = append(ctx, tok_PUBLIC)

		case tok_END_PUBLIC:
			ctx = ctx[:len(ctx)-1]

		case tok_APPLICATION:
			vv := Application{Name: tokens[1]}
			if len(tokens) > 2 {
				vv.Source = append(vv.Source, tokens[2:]...)
			}
			req.Applications = append(req.Applications, vv)

		case tok_LIBRARY:
			vv := Library{Name: tokens[1]}
			if len(tokens) > 2 {
				vv.Source = append(vv.Source, tokens[2:]...)
			}
			req.Libraries = append(req.Libraries, vv)

		case tok_DOCUMENT:
			vv := Document{
				Name:   tokens[1],
				Source: make([]string, 0, len(tokens[2:])), //FIXME
			}
			vv.Source = append(vv.Source, tokens[2:]...)
			req.Documents = append(req.Documents, vv)

		case tok_SET:
			vv := SetEnv{Name: tokens[1]}
			vv.Value = make(map[string]string)
			vv.Value["default"] = tokens[2]
			if len(tokens) > 3 {
				toks := tokens[3:]
				for i := 0; i+1 < len(toks); i += 2 {
					vv.Value[toks[i]] = toks[i+1]
				}
			}
			req.Sets = append(req.Sets, vv)

		case tok_TAG:
			vv := Tag{Name: tokens[1]}
			vv.Content = append(vv.Content, tokens[2:]...)
			req.Tags = append(req.Tags, vv)

		case tok_VERSION:
			vv := Version{Value: tokens[1]}
			req.Version = &vv

		case tok_CMTPATH_PATTERN:
			vv := CmtPathPattern{}
			vv.Cmd = append(vv.Cmd, tokens[2:]...)
			req.CmtPathPatterns = append(req.CmtPathPatterns, vv)

		case tok_MAKE_FRAGMENT:
			vv := MakeFragment{Name: tokens[1]}
			req.MakeFragments = append(req.MakeFragments, vv)

		case tok_ACTION:
			vv := Action{Name: tokens[1]}
			vv.Value = make(map[string]string)
			vv.Value["default"] = tokens[2]
			if len(tokens) > 3 {
				toks := tokens[3:]
				for i := 0; i+1 < len(toks); i += 2 {
					vv.Value[toks[i]] = toks[i+1]
				}
			}
			req.Actions = append(req.Actions, vv)

		default:
			return nil, fmt.Errorf("cmt2yml: unknown token [%v]", tokens[0])
		}
		bline = nil
	}

	fmt.Printf("req=%q [done]\n", fname)
	return &req, err
}

func parse_line(data []byte) ([]string, error) {
	var err error
	line := []string{}

	worder := bufio.NewScanner(bytes.NewBuffer(data))
	worder.Split(bufio.ScanWords)
	tokens := []string{}
	for worder.Scan() {
		tok := worder.Text()
		if tok != "" {
			tokens = append(tokens, worder.Text())
		}
	}

	my_printf := func(format string, args ...interface{}) (int, error) {
		return 0, nil
	}
	if dbg_parse_line {
		my_printf = func(format string, args ...interface{}) (int, error) {
			return fmt.Printf(format, args...)
		}
	}

	my_printf("===============\n")
	my_printf("tokens: %v\n", fmt_line(tokens))

	in_dquote := false
	in_squote := false
	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]
		my_printf("tok[%d]=%q\n", i, tok)
		if in_squote || in_dquote {
			if len(line) > 0 {
				line[len(line)-1] += " " + tok
			} else {
				panic("logic error")
			}
		} else {
			line = append(line, tok)
		}
		if strings.HasPrefix(tok, `"`) && !strings.HasSuffix(tok, `"`) {
			in_dquote = !in_dquote
			my_printf("--> dquote: %v -> %v\n", !in_dquote, in_dquote)
		}
		if strings.HasPrefix(tok, `'`) && !strings.HasSuffix(tok, `'`) {
			in_squote = !in_squote
			my_printf("--> squote: %v -> %v\n", !in_squote, in_squote)
		}
		if in_dquote && strings.HasSuffix(tok, `"`) {
			in_dquote = !in_dquote
			my_printf("<-- dquote: %v -> %v\n", !in_dquote, in_dquote)
		}
		if in_squote && strings.HasSuffix(tok, `'`) {
			in_squote = !in_squote
			my_printf("<-- squote: %v -> %v\n", !in_squote, in_squote)
		}
	}

	return line, err
}

// EOF
