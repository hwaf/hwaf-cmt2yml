package main

import (
	"io"
	"strings"
)

const (
	tok_PRIVATE = "private"
	tok_PUBLIC  = "public"
)

type ReqFile struct {
	Filename string
	Package  Package
	Stmts    []Stmt
}

func NewReqFile(name string) ReqFile {
	return ReqFile{
		Package: Package{name},
	}
}

type Stmt interface {
	ToYaml(w io.Writer) error
}

type ParseFunc func(p *Parser) error

var g_dispatch = map[string]ParseFunc{
	"package":         parsePackage,
	"author":          parseAuthor,
	"manager":         parseManager,
	"use":             parseUse,
	"macro":           parseMacro,
	"macro_append":    parseMacroAppend,
	"macro_prepend":   parseMacroPrepend,
	"private":         parsePrivate,
	"end_private":     parseEndPrivate,
	"public":          parsePublic,
	"end_public":      parseEndPublic,
	"application":     parseApplication,
	"pattern":         parsePattern,
	"ignore_pattern":  parseIgnorePattern,
	"apply_pattern":   parseApplyPattern,
	"library":         parseLibrary,
	"version":         parseVersion,
	"path":            parsePath,
	"path_append":     parsePathAppend,
	"path_prepend":    parsePathPrepend,
	"path_remove":     parsePathRemove,
	"include_dirs":    parseIncludeDirs,
	"include_path":    parseIncludePaths,
	"set":             parseSet,
	"tag":             parseTag,
	"document":        parseDocument,
	"cmtpath_pattern": parseCmtPathPattern,
	"make_fragment":   parseMakeFragment,
	"action":          parseAction,
	//"apply_tag":      parseApplyTag, //FIXME
}

type Package struct {
	Name string
}

func (s *Package) ToYaml(w io.Writer) error {
	return nil
}

func parsePackage(p *Parser) error {
	var err error
	p.req.Package = Package{Name: p.tokens[1]}
	p.req.Stmts = append(p.req.Stmts, &p.req.Package)
	return err
}

type Author struct {
	Name string
}

func (s *Author) ToYaml(w io.Writer) error {
	return nil
}

func parseAuthor(p *Parser) error {
	var err error
	for _, tok := range p.tokens[1:] {
		p.req.Stmts = append(p.req.Stmts, &Author{Name: tok})
	}
	return err
}

type Manager struct {
	Name string
}

func (s *Manager) ToYaml(w io.Writer) error {
	return nil
}

func parseManager(p *Parser) error {
	var err error
	for _, tok := range p.tokens[1:] {
		p.req.Stmts = append(p.req.Stmts, &Manager{Name: tok})
	}
	return err
}

type UsePkg struct {
	Package   string
	Version   string
	Path      string
	Switches  []string
	IsPrivate bool
}

func (s *UsePkg) ToYaml(w io.Writer) error {
	return nil
}

func parseUse(p *Parser) error {
	var err error
	tokens := p.tokens
	use := &UsePkg{Package: tokens[1]}
	if len(tokens) > 2 {
		use.Version = tokens[2]
	}
	if len(tokens) > 3 {
		use.Path = tokens[3]
	}
	if len(tokens) > 4 {
		use.Switches = append(use.Switches, tokens[4:]...)
	}
	p.req.Stmts = append(p.req.Stmts, use)
	return err
}

type Macro struct {
	Name  string
	Value map[string]string
}

func (s *Macro) ToYaml(w io.Writer) error {
	return nil
}

func parseMacro(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := Macro{Name: tokens[1]}
	vv.Value = make(map[string]string)
	vv.Value["default"] = tokens[2]
	if len(tokens) > 3 {
		toks := tokens[3:]
		for i := 0; i+1 < len(toks); i += 2 {
			vv.Value[toks[i]] = toks[i+1]
		}
	}
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type MacroAppend struct {
	Name  string
	Value map[string]string
}

func (s *MacroAppend) ToYaml(w io.Writer) error {
	return nil
}

func parseMacroAppend(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := MacroAppend{Name: tokens[1]}
	vv.Value = make(map[string]string)
	vv.Value["default"] = tokens[2]
	if len(tokens) > 3 {
		toks := tokens[3:]
		for i := 0; i+1 < len(toks); i += 2 {
			vv.Value[toks[i]] = toks[i+1]
		}
	}
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type MacroPrepend struct {
	Name  string
	Value map[string]string
}

func (s *MacroPrepend) ToYaml(w io.Writer) error {
	return nil
}

func parseMacroPrepend(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := MacroPrepend{Name: tokens[1]}
	vv.Value = make(map[string]string)
	vv.Value["default"] = tokens[2]
	if len(tokens) > 3 {
		toks := tokens[3:]
		for i := 0; i+1 < len(toks); i += 2 {
			vv.Value[toks[i]] = toks[i+1]
		}
	}
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type IncludeDirs struct {
	Value string
}

func (s *IncludeDirs) ToYaml(w io.Writer) error {
	return nil
}

func parseIncludeDirs(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := IncludeDirs{Value: tokens[1]}
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type IncludePaths struct {
	Value string
}

func (s *IncludePaths) ToYaml(w io.Writer) error {
	return nil
}

func parseIncludePaths(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := IncludePaths{Value: tokens[1]}
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type Version struct {
	Value string
}

func (s *Version) ToYaml(w io.Writer) error {
	return nil
}

func parseVersion(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := Version{Value: tokens[1]}
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type SetEnv struct {
	Name  string
	Value map[string]string
}

func (s *SetEnv) ToYaml(w io.Writer) error {
	return nil
}

func parseSet(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := SetEnv{Name: tokens[1]}
	vv.Value = make(map[string]string)
	vv.Value["default"] = tokens[2]
	if len(tokens) > 3 {
		toks := tokens[3:]
		for i := 0; i+1 < len(toks); i += 2 {
			vv.Value[toks[i]] = toks[i+1]
		}
	}
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type Pattern struct {
	Name string
	Def  string
}

func (s *Pattern) ToYaml(w io.Writer) error {
	return nil
}

func parsePattern(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := Pattern{
		Name: tokens[1],
		Def:  strings.Join(tokens[2:], " "),
	}
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type ApplyPattern struct {
	Name string
	Args []string
}

func (s *ApplyPattern) ToYaml(w io.Writer) error {
	return nil
}

func parseApplyPattern(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := ApplyPattern{Name: tokens[1]}
	if len(tokens) > 2 {
		vv.Args = append(vv.Args, tokens[2:]...)
	}
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type IgnorePattern struct {
	Name string
}

func (s *IgnorePattern) ToYaml(w io.Writer) error {
	return nil
}

func parseIgnorePattern(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := IgnorePattern{Name: tokens[1]}
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type Path struct {
	Name  string
	Value string
}

func (s *Path) ToYaml(w io.Writer) error {
	return nil
}

func parsePath(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := Path{Name: tokens[1], Value: tokens[2]}
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type PathAppend struct {
	Name  string
	Value map[string]string
}

func (s *PathAppend) ToYaml(w io.Writer) error {
	return nil
}

func parsePathAppend(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := PathAppend{Name: tokens[1]}
	vv.Value = make(map[string]string)
	vv.Value["default"] = tokens[2]
	if len(tokens) > 3 {
		toks := tokens[3:]
		for i := 0; i+1 < len(toks); i += 2 {
			vv.Value[toks[i]] = toks[i+1]
		}
	}
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type PathRemove struct {
	Name  string
	Value map[string]string
}

func (s *PathRemove) ToYaml(w io.Writer) error {
	return nil
}

func parsePathRemove(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := PathRemove{Name: tokens[1]}
	vv.Value = make(map[string]string)
	vv.Value["default"] = tokens[2]
	if len(tokens) > 3 {
		toks := tokens[3:]
		for i := 0; i+1 < len(toks); i += 2 {
			vv.Value[toks[i]] = toks[i+1]
		}
	}
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type PathPrepend struct {
	Name  string
	Value map[string]string
}

func (s *PathPrepend) ToYaml(w io.Writer) error {
	return nil
}

func parsePathPrepend(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := PathPrepend{Name: tokens[1]}
	vv.Value = make(map[string]string)
	vv.Value["default"] = tokens[2]
	if len(tokens) > 3 {
		toks := tokens[3:]
		for i := 0; i+1 < len(toks); i += 2 {
			vv.Value[toks[i]] = toks[i+1]
		}
	}
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type Tag struct {
	Name    string
	Content []string
}

func (s *Tag) ToYaml(w io.Writer) error {
	return nil
}

func parseTag(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := Tag{Name: tokens[1]}
	vv.Content = append(vv.Content, tokens[2:]...)
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type ApplyTag struct {
	Name string
	Args []string
}

func (s *ApplyTag) ToYaml(w io.Writer) error {
	return nil
}

type Library struct {
	Name   string
	Source []string
}

func (s *Library) ToYaml(w io.Writer) error {
	return nil
}

func parseLibrary(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := Library{Name: tokens[1]}
	if len(tokens) > 2 {
		vv.Source = append(vv.Source, tokens[2:]...)
	}
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type Action struct {
	Name  string
	Value map[string]string
}

func (s *Action) ToYaml(w io.Writer) error {
	return nil
}

func parseAction(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := Action{Name: tokens[1]}
	vv.Value = make(map[string]string)
	vv.Value["default"] = tokens[2]
	if len(tokens) > 3 {
		toks := tokens[3:]
		for i := 0; i+1 < len(toks); i += 2 {
			vv.Value[toks[i]] = toks[i+1]
		}
	}
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type Application struct {
	Name   string
	Source []string
}

func (s *Application) ToYaml(w io.Writer) error {
	return nil
}

func parseApplication(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := Application{Name: tokens[1]}
	if len(tokens) > 2 {
		vv.Source = append(vv.Source, tokens[2:]...)
	}
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type Document struct {
	Name   string
	Group  string
	S      string
	Source []string
}

func (s *Document) ToYaml(w io.Writer) error {
	return nil
}

func parseDocument(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := Document{
		Name:   tokens[1],
		Source: make([]string, 0, len(tokens[2:])), //FIXME
	}
	vv.Source = append(vv.Source, tokens[2:]...)
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type CmtPathPattern struct {
	Cmd []string
}

func (s *CmtPathPattern) ToYaml(w io.Writer) error {
	return nil
}

func parseCmtPathPattern(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := CmtPathPattern{}
	vv.Cmd = append(vv.Cmd, tokens[2:]...)
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type MakeFragment struct {
	Name string
}

func (s *MakeFragment) ToYaml(w io.Writer) error {
	return nil
}

func parseMakeFragment(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := MakeFragment{Name: tokens[1]}
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

func parsePrivate(p *Parser) error {
	var err error
	p.ctx = append(p.ctx, tok_PRIVATE)
	return err
}

func parseEndPrivate(p *Parser) error {
	var err error
	p.ctx = p.ctx[:len(p.ctx)-1]
	return err
}

func parsePublic(p *Parser) error {
	var err error
	p.ctx = append(p.ctx, tok_PUBLIC)
	return err
}

func parseEndPublic(p *Parser) error {
	var err error
	p.ctx = p.ctx[:len(p.ctx)-1]
	return err
}

// EOF
