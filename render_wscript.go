package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/hwaf/hwaf/hlib"
)

type wscript_t struct {
	Package   wpackage_t
	Options   woptions_t
	Configure wconfigure_t
	Build     wbuild_t
}

type wpackage_t struct {
	Name     string
	Authors  []string
	Managers []string
	Version  string
	Deps     wdeps_t
}

type wdeps_t struct {
	Public  []string
	Private []string
	Runtime []string
}

type woptions_t struct {
	Tools []string
	//Stmts []Stmt
	//HwafCall []string
}

type wconfigure_t struct {
	Tools []string
	Env   wenv_t
	Tag   []string
	Stmts []Stmt
}

type wenv_t map[string]interface{} //FIXME: map[string][]interface{} instead ??

// type hbuild_t struct {
// 	Targets  htargets_t `yaml:"tgts,flow,omitempty"`
// 	HwafCall []string   `yaml:"hwaf-call,flow,omitempty"`
// }

type wbuild_t struct {
	Tools   []string
	Env     wenv_t
	Targets wtargets_t
	Stmts   []Stmt
}

type wtargets_t []wtarget_t

type wtarget_t struct {
	Name     string
	Features string
	Source   []string
	Use      []string
	Defines  []string
}

func (r *Renderer) render_wscript() error {
	var err error

	_, err = fmt.Fprintf(
		r.w,
		`## automatically generated by cmt2yml
## do NOT edit

## waf imports
import waflib.Logs as msg

`,
	)
	handle_err(err)

	basedir := filepath.Dir(filepath.Dir(r.req.Filename))

	wscript := wscript_t{
		Package:   wpackage_t{Name: basedir},
		Configure: wconfigure_t{Env: make(wenv_t)},
		Build:     wbuild_t{Env: make(wenv_t)},
	}

	//complibs := map[string]struct{}{}
	linklibs := map[string]*Library{}
	apps := map[string]*Application{}
	//dictlibs := map[string]struct{}{}

	// first pass to detect targets
	for _, stmt := range r.req.Stmts {
		switch stmt.(type) {
		case *Application:
			x := stmt.(*Application)
			apps[x.Name] = x

		case *Library:
			x := stmt.(*Library)
			linklibs[x.Name] = x
		}
	}

	// second pass to collect
	for _, stmt := range r.req.Stmts {
		wpkg := &wscript.Package
		wbld := &wscript.Build
		wcfg := &wscript.Configure

		if author, ok := stmt.(*Author); ok {
			wpkg.Authors = append(wpkg.Authors, author.Name)
		}

		if mgr, ok := stmt.(*Manager); ok {
			wpkg.Managers = append(wpkg.Managers, mgr.Name)
		}

		if version, ok := stmt.(*Version); ok {
			wpkg.Version = version.Value
		}

		if use, ok := stmt.(*UsePkg); ok {
			deps := &wpkg.Deps
			if use.IsPrivate {
				deps.Private = append(deps.Private, path.Join(use.Path, use.Package))
			} else {
				deps.Public = append(deps.Public, path.Join(use.Path, use.Package))
			}
		}

		if lib, ok := stmt.(*Library); ok {
			linklibs[lib.Name] = lib
			tgt := wtarget_t{Name: lib.Name}
			sanitize_srcs(lib.Source)
			for _, src := range lib.Source {
				tgt.Source = append(tgt.Source, src)
			}
			if features, ok := g_profile.features["library"]; ok {
				tgt.Features = features
			}
			wbld.Targets = append(wbld.Targets, tgt)
		}

		if app, ok := stmt.(*Application); ok {
			apps[app.Name] = app
			tgt := wtarget_t{Name: app.Name}
			sanitize_srcs(app.Source)
			for _, src := range app.Source {
				tgt.Source = append(tgt.Source, src)
			}
			if features, ok := g_profile.features["application"]; ok {
				tgt.Features = features
			}
			wbld.Targets = append(wbld.Targets, tgt)
			wbld.Targets[len(wbld.Targets)-1].Name = tgt.Name
		}

		switch stmt.(type) {
		case *Macro,
			*MacroAppend, *MacroRemove, *MacroPrepend,
			*Tag,
			*PathAppend, *PathPrepend, *PathRemove:
			wcfg.Stmts = append(wcfg.Stmts, stmt)
		}
	}

	// generate package header
	const pkg_hdr_tmpl = `
PACKAGE = {
    "name":    "{{.Name}}",
    "authors": {{.Authors | as_pylist}},
{{if .Managers}}    "managers": {{.Managers | as_pylist}},{{end}}
{{if .Version}}    "version":  "{{.Version}}",{{end}}
}

### ---------------------------------------------------------------------------
def pkg_deps(ctx):
    {{with .Deps}}## public dependencies
    {{if .Public}}{{range .Public}}ctx.use_pkg("{{.}}", public=True)
    {{end}}{{else}}## => none{{end}}

    ## private dependencies
    {{if .Private}}{{range .Private}}ctx.use_pkg("{{.}}", private=True){{end}}{{else}}## => none{{end}}

    ## runtime dependencies
    {{if .Runtime}}{{range .Runtime}}ctx.use_pkg("{{.}}", runtime=True){{end}}{{else}}## => none{{end}}{{end}}

    return # pkg_deps
`
	err = w_tmpl(r.w, pkg_hdr_tmpl, wscript.Package)
	handle_err(err)

	// generate options - section
	err = w_tmpl(
		r.w,
		`

### ---------------------------------------------------------------------------
def options(ctx):
    {{range .Tools}}ctx.load("{{.}}")
    {{end}}
    return # options
`,
		wscript.Options,
	)
	handle_err(err)

	// generate configure - section
	err = w_tmpl(
		r.w,
		`

### ---------------------------------------------------------------------------
def configure(ctx):
    {{range .Tools}}ctx.load("{{.}}")
    {{end}}
    {{range .Stmts}}##{{. | gen_wscript_stmts}}
    {{end}}
    return # configure
`,
		wscript.Configure,
	)
	handle_err(err)

	// generate build - section
	err = w_tmpl(
		r.w,
		`

### ---------------------------------------------------------------------------
def build(ctx):
    {{range .Tools}}ctx.load("{{.}}")
    {{end}}
    {{range .Stmts}}##{{. | gen_wscript_stmts}}
    {{end}}
    return # configure
`,
		wscript.Build,
	)
	handle_err(err)

	_, err = fmt.Fprintf(
		r.w,
		"\n## EOF ##\n",
	)
	handle_err(err)

	return err
}

func w_tmpl(w *os.File, text string, data interface{}) error {
	t := template.New("wscript")
	t.Funcs(template.FuncMap{
		"trim": strings.TrimSpace,
		"as_pylist": func(list []string) string {
			str := make([]string, 0, len(list))
			for _, s := range list {
				str = append(str, fmt.Sprintf("%q", s))
			}
			return "[" + strings.Join(str, ", ") + "]"
		},
		"gen_wscript_stmts": gen_wscript_stmts,
	})
	template.Must(t.Parse(text))
	return t.Execute(w, data)
}

func w_gen_taglist(tags string) []string {
	str := make([]string, 0, strings.Count(tags, "&"))
	for _, v := range strings.Split(tags, "&") {
		v = strings.Trim(v, " ")
		if len(v) > 0 {
			str = append(str, v)
		}
	}
	return str
}

func w_py_strlist(str []string) string {
	o := make([]string, 0, len(str))
	for _, v := range str {
		o = append(o, fmt.Sprintf("%q", v))
	}
	return strings.Join(o, ", ")
}

func w_gen_valdict_switch_str(indent string, values [][2]string) string {
	o := make([]string, 0, len(values)+2)
	o = append(o, "(")
	for _, v := range values {
		tags := w_gen_taglist(v[0])
		key_fmt := "(%s)"
		if strings.Count(v[0], "&") <= 0 {
			key_fmt = "%s"
		}
		o = append(o,
			fmt.Sprintf(
				"%s  {%s: %q},",
				indent,
				fmt.Sprintf(key_fmt, w_py_strlist(tags)),
				v[1],
			),
		)
	}
	o = append(o, indent+")")
	return strings.Join(o, "\n")
}

func w_py_hlib_value(indent string, fctname string, x hlib.Value) []string {
	str := make([]string, 0)

	values := make([][2]string, 0, len(x.Set))
	for _, v := range x.Set {
		k := v.Tag
		values = append(values, [2]string{k, w_py_strlist(v.Value)})
	}
	str = append(
		str,
		fmt.Sprintf(
			"ctx.%s(%q, %s)",
			fctname,
			x.Name,
			w_gen_valdict_switch_str(indent, values),
		),
	)

	return str
}

func gen_wscript_stmts(stmt Stmt) string {
	const indent = "    "
	var str []string
	switch xx := stmt.(type) {
	case *Macro:
		str = []string{fmt.Sprintf("## macro %v", stmt)}
		x := hlib.Value(*xx)
		str = append(
			str,
			w_py_hlib_value(indent, "hwaf_declare_macro", x)...,
		)

	case *MacroAppend:
		str = []string{fmt.Sprintf("## macro_append %v", stmt)}
		x := hlib.Value(*xx)
		str = append(
			str,
			w_py_hlib_value(indent, "hwaf_macro_append", x)...,
		)

	case *MacroPrepend:
		str = []string{fmt.Sprintf("## macro_prepend %v", stmt)}
		x := hlib.Value(*xx)
		str = append(
			str,
			w_py_hlib_value(indent, "hwaf_macro_prepend", x)...,
		)

	case *MacroRemove:
		str = []string{fmt.Sprintf("## macro_remove %v", stmt)}
		x := hlib.Value(*xx)
		str = append(
			str,
			w_py_hlib_value(indent, "hwaf_macro_remove", x)...,
		)

	case *Tag:
		str = []string{fmt.Sprintf("## tag %v", stmt)}
		x := xx
		values := w_py_strlist(x.Content)
		str = append(str,
			"ctx.hwaf_declare_tag(",
			fmt.Sprintf("%s%q,", indent, x.Name),
			fmt.Sprintf("%scontent=[%s]", indent, values),
			")",
		)

	case *Path:
		str = []string{fmt.Sprintf("## path %v", stmt)}
		x := hlib.Value(*xx)
		str = append(
			str,
			w_py_hlib_value(indent, "hwaf_declare_path", x)...,
		)

	case *PathAppend:
		str = []string{fmt.Sprintf("## path_append %v", stmt)}
		x := hlib.Value(*xx)
		str = append(
			str,
			w_py_hlib_value(indent, "hwaf_path_append", x)...,
		)

	case *PathPrepend:
		str = []string{fmt.Sprintf("## path_prepend %v", stmt)}
		x := hlib.Value(*xx)
		str = append(
			str,
			w_py_hlib_value(indent, "hwaf_path_prepend", x)...,
		)

	case *PathRemove:
		str = []string{fmt.Sprintf("## path_remove %v", stmt)}
		x := hlib.Value(*xx)
		str = append(
			str,
			w_py_hlib_value(indent, "hwaf_path_remove", x)...,
		)

	default:
		str = []string{fmt.Sprintf("### **** statement %T (%v)", stmt, stmt)}
	}

	// reindent:
	for i, s := range str[1:] {
		str[i+1] = indent + s
	}

	return strings.Join(str, "\n")
}

// EOF
