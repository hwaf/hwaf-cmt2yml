package main

import (
	"io"
	"strings"

	"github.com/hwaf/hwaf/hlib"
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

func (req *ReqFile) ToYaml(w io.Writer) error {
	var err error
	return err
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
	"macro_remove":    parseMacroRemove,
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

type Macro hlib.Value

func (s *Macro) ToYaml(w io.Writer) error {
	return nil
}

func parseMacro(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := Macro(hlib_value_from_slice(tokens[1], tokens[2:]))
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type MacroAppend hlib.Value

func (s *MacroAppend) ToYaml(w io.Writer) error {
	return nil
}

func parseMacroAppend(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := MacroAppend(hlib_value_from_slice(tokens[1], tokens[2:]))
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type MacroPrepend hlib.Value

func (s *MacroPrepend) ToYaml(w io.Writer) error {
	return nil
}

func parseMacroPrepend(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := MacroPrepend(hlib_value_from_slice(tokens[1], tokens[2:]))
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type MacroRemove hlib.Value

func (s *MacroRemove) ToYaml(w io.Writer) error {
	return nil
}

func parseMacroRemove(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := MacroRemove(hlib_value_from_slice(tokens[1], tokens[2:]))
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

type SetEnv hlib.Value

func (s *SetEnv) ToYaml(w io.Writer) error {
	return nil
}

func parseSet(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := SetEnv(hlib_value_from_slice(tokens[1], tokens[2:]))
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
	if tokens[1][0] == '-' {
		tokens[1], tokens[2] = tokens[2], tokens[1]
	}
	vv := Pattern{
		Name: tokens[1],
		Def:  strings.Join(tokens[2:], " "),
	}
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type ApplyPattern hlib.Value

func (s *ApplyPattern) ToYaml(w io.Writer) error {
	return nil
}

func parseApplyPattern(p *Parser) error {
	var err error
	tokens := p.tokens
	if tokens[1][0] == '-' {
		tokens[1], tokens[2] = tokens[2], tokens[1]
	}
	vv := ApplyPattern(hlib_value_from_slice(tokens[1], tokens[2:]))
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type IgnorePattern hlib.Value

func (s *IgnorePattern) ToYaml(w io.Writer) error {
	return nil
}

func parseIgnorePattern(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := IgnorePattern(hlib_value_from_slice(tokens[1], nil))
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type Path hlib.Value

func (s *Path) ToYaml(w io.Writer) error {
	return nil
}

func parsePath(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := Path(hlib_value_from_slice(tokens[1], tokens[2:]))
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type PathAppend hlib.Value

func (s *PathAppend) ToYaml(w io.Writer) error {
	return nil
}

func parsePathAppend(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := PathAppend(hlib_value_from_slice(tokens[1], tokens[2:]))
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type PathRemove hlib.Value

func (s *PathRemove) ToYaml(w io.Writer) error {
	return nil
}

func parsePathRemove(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := PathRemove(hlib_value_from_slice(tokens[1], tokens[2:]))
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type PathPrepend hlib.Value

func (s *PathPrepend) ToYaml(w io.Writer) error {
	return nil
}

func parsePathPrepend(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := PathPrepend(hlib_value_from_slice(tokens[1], tokens[2:]))
	p.req.Stmts = append(p.req.Stmts, &vv)
	return err
}

type Tag hlib.TagStmt

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

type ApplyTag hlib.Value

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

type Action hlib.Value

func (s *Action) ToYaml(w io.Writer) error {
	return nil
}

func parseAction(p *Parser) error {
	var err error
	tokens := p.tokens
	vv := Action(hlib_value_from_slice(tokens[1], tokens[2:]))
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

type MakeFragment hlib.MakeFragmentStmt

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
