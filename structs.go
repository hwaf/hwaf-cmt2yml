package main

const (
	tok_PACKAGE        = "package"
	tok_AUTHOR         = "author"
	tok_MANAGER        = "manager"
	tok_USE            = "use"
	tok_MACRO          = "macro"
	tok_MACRO_APPEND   = "macro_append"
	tok_MACRO_PREPEND  = "macro_prepend"
	tok_INCLUDE_DIRS   = "include_dirs"
	tok_INCLUDE_PATH   = "include_path"
	tok_VERSION        = "version"
	tok_SET            = "set"
	tok_PATTERN        = "pattern"
	tok_APPLY_PATTERN  = "apply_pattern"
	tok_IGNORE_PATTERN = "ignore_pattern"
	tok_PATH           = "path"
	tok_PATH_APPEND    = "path_append"
	tok_PATH_PREPEND   = "path_prepend"
	tok_PATH_REMOVE    = "path_remove"
	tok_TAG            = "tag"
	tok_APPLY_TAG      = "apply_tag"
	tok_LIBRARY        = "library"
	tok_ACTION         = "action"
	tok_APPLICATION    = "application"
	tok_DOCUMENT       = "document"

	tok_CMTPATH_PATTERN = "cmtpath_pattern"
	tok_MAKE_FRAGMENT   = "make_fragment"

	tok_PRIVATE     = "private"
	tok_END_PRIVATE = "end_private"
	tok_PUBLIC      = "public"
	tok_END_PUBLIC  = "end_public"
)

type ReqFile struct {
	Package string
	//	Stmts []Statement
	Authors         []string
	Managers        []string
	Uses            []UsePkg
	Macros          []Macro
	MacroAppends    []MacroAppend
	MacroPrepends   []MacroPrepend
	IncludeDirs     []IncludeDirs
	IncludePaths    []IncludePaths
	Version         *Version
	Sets            []SetEnv
	Patterns        []Pattern
	ApplyPatterns   []ApplyPattern
	IgnorePatterns  []IgnorePattern
	Paths           []Path
	PathAppends     []PathAppend
	PathPrepends    []PathPrepend
	PathRemoves     []PathRemove
	Tags            []Tag
	ApplyTags       []ApplyTag
	Libraries       []Library
	Actions         []Action
	Applications    []Application
	Documents       []Document
	CmtPathPatterns []CmtPathPattern
	MakeFragments   []MakeFragment
}

func NewReqFile(name string) ReqFile {
	return ReqFile{
		Package: name,
	}
}

type UsePkg struct {
	Package   string
	Version   string
	Path      string
	Switches  []string
	IsPrivate bool
}

type Macro struct {
	Name  string
	Value map[string]string
}

type MacroAppend struct {
	Name  string
	Value map[string]string
}

type MacroPrepend struct {
	Name  string
	Value map[string]string
}

type IncludeDirs struct {
	Value string
}

type IncludePaths struct {
	Value string
}

type Version struct {
	Value string
}

type SetEnv struct {
	Name  string
	Value map[string]string
}

type Pattern struct {
	Name string
	Def  string
}

type ApplyPattern struct {
	Name string
	Args []string
}

type IgnorePattern struct {
	Name string
}

type Path struct {
	Name  string
	Value string
}

type PathAppend struct {
	Name  string
	Value map[string]string
}

type PathRemove struct {
	Name  string
	Value map[string]string
}

type PathPrepend struct {
	Name  string
	Value map[string]string
}

type Tag struct {
	Name    string
	Content []string
}

type ApplyTag struct {
	Name string
	Args []string
}

type Library struct {
	Name   string
	Source []string
}

type Action struct {
	Name  string
	Value map[string]string
}

type Application struct {
	Name   string
	Source []string
}

type Document struct {
	Name   string
	Group  string
	S      string
	Source []string
}

type CmtPathPattern struct {
	Cmd []string
}

type MakeFragment struct {
	Name string
}

// EOF
