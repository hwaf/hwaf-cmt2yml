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

type Parser struct {
	req       *ReqFile
	isPrivate bool

	table   map[string]ParseFunc
	f       *os.File
	scanner *bufio.Scanner
	ctx     []string
	tokens  []string
}

func (p *Parser) Close() error {
	if p.f == nil {
		return nil
	}
	err := p.f.Close()
	p.f = nil
	return err
}

func NewParser(fname string) (*Parser, error) {

	var err error
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bufio.NewReader(f))
	if scanner == nil {
		return nil, fmt.Errorf("cmt2yml: nil bufio.Scanner")
	}

	p := &Parser{
		table:   g_dispatch,
		f:       f,
		scanner: scanner,
		req:     &ReqFile{Filename: fname},
		tokens:  nil,
		ctx:     []string{tok_PUBLIC},
	}
	return p, nil
}

func (p *Parser) run() error {
	var err error
	bline := []byte{}
	for p.scanner.Scan() {
		data := p.scanner.Bytes()
		data = bytes.TrimSpace(data)
		data = bytes.Trim(data, " \t\r\n")

		if len(data) == 0 {
			continue
		}

		if data[0] == '#' {
			continue
		}

		idx := len(data) - 1
		if data[idx] == '\\' {
			bline = append(bline, ' ')
			bline = append(bline, data[:idx-1]...)
			continue
		} else {
			bline = append(bline, ' ')
			bline = append(bline, data...)
		}

		var tokens []string
		tokens, err = parse_line(bline)
		if err != nil {
			return err
		}
		p.tokens = tokens

		fct, ok := p.table[p.tokens[0]]
		if !ok {
			return fmt.Errorf("cmt2yml: unknown token [%v]", tokens[0])
		}
		err = fct(p)
		if err != nil {
			return err
		}
		bline = nil
	}

	return err
}

func parse_file(fname string) (*ReqFile, error) {
	fmt.Printf("req=%q\n", fname)
	p, err := NewParser(fname)
	if err != nil {
		return nil, err
	}
	defer p.Close()

	err = p.run()
	fmt.Printf("req=%q [done]\n", fname)
	return p.req, err
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
		tok := strings.Trim(tokens[i], " \t")
		my_printf("tok[%d]=%q    (q=%v)\n", i, tok, in_squote || in_dquote)
		if in_squote || in_dquote {
			if len(line) > 0 {
				ttok := tok
				if strings.HasPrefix(ttok, `"`) || strings.HasPrefix(ttok, "'") {
					ttok = ttok[1:]
				}
				if strings.HasSuffix(ttok, `"`) || strings.HasSuffix(ttok, "'") {
					if !strings.HasSuffix(ttok, `\"`) {
						ttok = ttok[:len(ttok)-1]
					}
				}
				ttok = strings.Trim(ttok, " \t")
				if len(ttok) > 0 {
					line_val := line[len(line)-1]
					line_sep := ""
					if len(line_val) > 0 {
						line_sep = " "
					}
					ttok = strings.Replace(ttok, `\"`, `"`, -1)
					line[len(line)-1] += line_sep + ttok
				}
			} else {
				panic("logic error")
			}
		} else {
			ttok := tok
			if strings.HasPrefix(ttok, `"`) || strings.HasPrefix(ttok, "'") {
				ttok = ttok[1:]
			}
			if strings.HasSuffix(ttok, `"`) || strings.HasSuffix(ttok, "'") {
				if !strings.HasSuffix(ttok, `\"`) {
					ttok = ttok[:len(ttok)-1]
				}
			}
			ttok = strings.Replace(ttok, `\"`, `"`, -1)
			line = append(line, strings.Trim(ttok, " \t"))
		}
		if len(tok) == 1 && strings.HasPrefix(tok, "\"") {
			in_dquote = !in_dquote
			continue
		}
		if len(tok) == 1 && strings.HasPrefix(tok, "'") {
			in_squote = !in_squote
			continue
		}
		if strings.HasPrefix(tok, "\"") && !strings.HasSuffix(tok, "\"") {
			in_dquote = !in_dquote
			my_printf("--> dquote: %v -> %v\n", !in_dquote, in_dquote)
		}
		if strings.HasPrefix(tok, "'") && !strings.HasSuffix(tok, "'") {
			in_squote = !in_squote
			my_printf("--> squote: %v -> %v\n", !in_squote, in_squote)
		}
		if in_dquote && strings.HasSuffix(tok, "\"") && !strings.HasSuffix(tok, `\""`) {
			in_dquote = !in_dquote
			my_printf("<-- dquote: %v -> %v\n", !in_dquote, in_dquote)
		}
		if in_squote && strings.HasSuffix(tok, "'") {
			in_squote = !in_squote
			my_printf("<-- squote: %v -> %v\n", !in_squote, in_squote)
		}
	}

	return line, err
}

// EOF
