package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hwaf/hwaf/hlib"
)

type Renderer struct {
	req     *ReqFile
	wscript bool
	w       *os.File
	pkg     hlib.Wscript_t
}

func NewRenderer(req *ReqFile) (*Renderer, error) {
	var err error
	var r *Renderer

	r = &Renderer{req: req, wscript: false}
	return r, err
}

func (r *Renderer) Render() error {
	var err error
	err = r.analyze()
	if err != nil {
		return err
	}
	err = r.render()
	if err != nil {
		return err
	}
	return err
}

func (r *Renderer) analyze() error {
	var err error

	basedir := filepath.Dir(filepath.Dir(r.req.Filename))

	r.pkg = hlib.Wscript_t{
		Package:   hlib.Package_t{Name: basedir},
		Configure: hlib.Configure_t{Env: make(hlib.Env_t)},
		Build:     hlib.Build_t{Env: make(hlib.Env_t)},
	}
	wscript := &r.pkg

	// targets
	apps := make(map[string]*Application)
	libs := make(map[string]*Library)

	// first pass: discover targets
	for _, stmt := range r.req.Stmts {
		switch stmt.(type) {
		case *Application:
			x := stmt.(*Application)
			apps[x.Name] = x

		case *Library:
			x := stmt.(*Library)
			libs[x.Name] = x
		}
	}

	// list of macros related to targets.
	// this will be used to:
	//  - fold them together
	//  - pre-process macro_append, macro_remove, ...
	//  - dispatch to wscript equivalents. e.g.:
	//     - <name>linkopts -> ctx(use=[...], cxxshlibflags=[...])
	//     - <name>_dependencies -> ctx(depends_on=[...])
	//     - includes -> ctx(includes=[..])
	macros := make(map[string][]Stmt)

	tgt_names := make([]string, 0, len(apps)+len(libs))
	for k, _ := range apps {
		tgt_names = append(tgt_names, k)
	}
	for k, _ := range libs {
		tgt_names = append(tgt_names, k)
	}
	sort.Strings(tgt_names)

	//fmt.Printf("+++ tgt_names: %v\n", tgt_names)

	// second pass: collect macros
	for _, stmt := range r.req.Stmts {
		switch x := stmt.(type) {
		default:
			continue
		case *Macro:
			//fmt.Printf("== [%s] ==\n", x.Name)
			//pat := x.Name+"(_dependencies|linkopts)"
			pat := ".*?"
			if !re_is_in_slice_suffix(tgt_names, x.Name, pat) {
				continue
			}
			macros[x.Name] = append(macros[x.Name], x)

		case *MacroAppend:
			pat := ".*?"
			if !re_is_in_slice_suffix(tgt_names, x.Name, pat) {
				continue
			}
			macros[x.Name] = append(macros[x.Name], x)

		case *MacroRemove:
			pat := ".*?"
			if !re_is_in_slice_suffix(tgt_names, x.Name, pat) {
				continue
			}
			macros[x.Name] = append(macros[x.Name], x)

		}
	}

	// models private/public, end_private/end_public
	ctx_visible := []bool{true}

	// 3rd pass to collect
	for _, stmt := range r.req.Stmts {
		wpkg := &wscript.Package
		wbld := &wscript.Build
		wcfg := &wscript.Configure
		switch x := stmt.(type) {

		case *BeginPublic:
			ctx_visible = append(ctx_visible, true)
		case *EndPublic:
			ctx_visible = ctx_visible[:len(ctx_visible)-1]

		case *BeginPrivate:
			ctx_visible = append(ctx_visible, false)
		case *EndPrivate:
			ctx_visible = ctx_visible[:len(ctx_visible)-1]

		case *Author:
			wpkg.Authors = append(wpkg.Authors, hlib.Author(x.Name))

		case *Manager:
			wpkg.Managers = append(wpkg.Managers, hlib.Manager(x.Name))

		case *Version:
			wpkg.Version = hlib.Version(x.Value)

		case *UsePkg:
			deptype := hlib.PrivateDep
			if ctx_visible[len(ctx_visible)-1] {
				deptype = hlib.PublicDep
			}
			if str_is_in_slice(x.Switches, "-no_auto_imports") {
				deptype = hlib.RuntimeDep | deptype
			}
			wpkg.Deps = append(
				wpkg.Deps,
				hlib.Dep_t{
					Name:    path.Join(x.Path, x.Package),
					Version: hlib.Version(x.Version),
					Type:    deptype,
				},
			)

		case *Library:
			tgt := hlib.Target_t{Name: x.Name}
			srcs, rest := sanitize_srcs(x.Source)
			// FIXME: handle -s=some/dir
			if len(rest) > 0 {
			}
			val := hlib.Value{
				Name: x.Name,
				Set: []hlib.KeyValue{
					{Tag: "default", Value: srcs},
				},
			}
			tgt.Source = append(tgt.Source, val)
			if features, ok := g_profile.features["library"]; ok {
				tgt.Features = features
			}
			w_distill_tgt(&tgt, macros)
			wbld.Targets = append(wbld.Targets, tgt)

		case *Application:
			tgt := hlib.Target_t{Name: x.Name}
			srcs, rest := sanitize_srcs(x.Source)
			// FIXME: handle -s=some/dir
			if len(rest) > 0 {
			}
			val := hlib.Value{
				Name: x.Name,
				Set: []hlib.KeyValue{
					{Tag: "default", Value: srcs},
				},
			}

			tgt.Source = append(tgt.Source, val)
			if features, ok := g_profile.features["application"]; ok {
				tgt.Features = features
			}
			w_distill_tgt(&tgt, macros)
			wbld.Targets = append(wbld.Targets, tgt)

		case *Alias:
			val := hlib.Value(*x)
			wcfg.Stmts = append(wcfg.Stmts, &hlib.AliasStmt{Value: val})

		case *Macro:
			if _, ok := macros[x.Name]; ok {
				// this will be used by a library or application
				continue
			}
			val := hlib.Value(*x)
			wcfg.Stmts = append(wcfg.Stmts, &hlib.MacroStmt{Value: val})

		case *MacroAppend:
			if _, ok := macros[x.Name]; ok {
				// this will be used by a library or application
				continue
			}
			val := hlib.Value(*x)
			wcfg.Stmts = append(wcfg.Stmts, &hlib.MacroAppendStmt{Value: val})

		case *MacroPrepend:
			if _, ok := macros[x.Name]; ok {
				// this will be used by a library or application
				continue
			}
			val := hlib.Value(*x)
			wcfg.Stmts = append(wcfg.Stmts, &hlib.MacroPrependStmt{Value: val})

		case *MacroRemove:
			if _, ok := macros[x.Name]; ok {
				// this will be used by a library or application
				continue
			}
			val := hlib.Value(*x)
			wcfg.Stmts = append(wcfg.Stmts, &hlib.MacroRemoveStmt{Value: val})

		case *Path:
			val := hlib.Value(*x)
			wcfg.Stmts = append(wcfg.Stmts, &hlib.PathStmt{Value: val})

		case *PathAppend:
			val := hlib.Value(*x)
			wcfg.Stmts = append(wcfg.Stmts, &hlib.PathAppendStmt{Value: val})

		case *PathPrepend:
			val := hlib.Value(*x)
			wcfg.Stmts = append(wcfg.Stmts, &hlib.PathPrependStmt{Value: val})

		case *PathRemove:
			val := hlib.Value(*x)
			wcfg.Stmts = append(wcfg.Stmts, &hlib.PathRemoveStmt{Value: val})

		case *Pattern:
			wcfg.Stmts = append(wcfg.Stmts, (*hlib.PatternStmt)(x))

		case *ApplyPattern:
			if cnv, ok := g_profile.cnvs[x.Name]; ok {
				err = cnv(wscript, x)
				if err != nil {
					return err
				}
			} else {
				wbld.Stmts = append(wbld.Stmts, (*hlib.ApplyPatternStmt)(x))
			}

		case *Tag:
			wcfg.Stmts = append(wcfg.Stmts, (*hlib.TagStmt)(x))

		case *ApplyTag:
			val := hlib.Value(*x)
			wcfg.Stmts = append(wcfg.Stmts, &hlib.ApplyTagStmt{Value: val})

		case *TagExclude:
			wcfg.Stmts = append(wcfg.Stmts, (*hlib.TagExcludeStmt)(x))

		case *MakeFragment:
			wcfg.Stmts = append(wcfg.Stmts, (*hlib.MakeFragmentStmt)(x))

		case *SetEnv:
			val := hlib.Value(*x)
			wcfg.Stmts = append(wcfg.Stmts, &hlib.SetStmt{Value: val})

		case *SetAppend:
			val := hlib.Value(*x)
			wcfg.Stmts = append(wcfg.Stmts, &hlib.SetAppendStmt{Value: val})

		case *SetRemove:
			val := hlib.Value(*x)
			wcfg.Stmts = append(wcfg.Stmts, &hlib.SetRemoveStmt{Value: val})

		case *Package:
			// already dealt with

		case *Action:
			// FIXME

		case *IncludePaths:
			wcfg.Stmts = append(wcfg.Stmts, (*hlib.IncludePathStmt)(x))

		case *IncludeDirs:
			wcfg.Stmts = append(wcfg.Stmts, (*hlib.IncludeDirsStmt)(x))

		case *CmtPathPattern:
			// FIXME

		case *CmtPathPatternReverse:
			// FIXME

		case *IgnorePattern:
			// FIXME

		case *Document:
			wbld.Stmts = append(wbld.Stmts, (*hlib.DocumentStmt)(x))

		default:
			return fmt.Errorf("unhandled statement [%v] (type=%T)\ndir=%v", x, x, r.req.Filename)
		}
	}

	for _, stmt := range r.req.Stmts {
		switch stmt := stmt.(type) {
		case *PathRemove, *MakeFragment, *Pattern, *MacroRemove:
			r.wscript = true
			break
		case *Macro:
			if len(stmt.Set) > 1 {
				r.wscript = true
				break
			}
		}
	}

	// fixups for boost
	for _, tgt := range wscript.Build.Targets {
		for _, use := range tgt.Use {
			for _, kv := range use.Set {
				for i, vv := range kv.Value {
					vv = strings.Replace(vv, "-${boost_libsuffix}", "", -1)
					vv = strings.Replace(vv, "boost_", "boost-", -1)
					kv.Value[i] = vv
				}
			}
		}
	}
	return err
}

func (r *Renderer) render() error {
	var err error
	pkgdir := filepath.Dir(filepath.Dir(r.req.Filename))
	fname := ""
	render := r.render_hscript
	if r.wscript {
		fname = filepath.Join(pkgdir, "wscript")
		render = r.render_wscript
	} else {
		fname = filepath.Join(pkgdir, "hscript.yml")
		render = r.render_hscript
	}

	if path_exists(fname) {
		f, err := os.Open(fname)
		if err == nil {
			buf := make([]byte, 64)
			f.Read(buf)
			if !bytes.HasPrefix(buf, []byte(`## automatically generated by cmt2yml`)) {
				// user generated file.
				// keep it.
				fmt.Printf("**warning** file [%s] already present\n", fname)
				return nil
			}
		}
	}

	r.w, err = os.Create(fname)
	if err != nil {
		return err
	}
	defer func() {
		r.w.Sync()
		r.w.Close()
	}()

	err = render()
	return err
}

func render_script(req *ReqFile) error {
	var err error

	renderer, err := NewRenderer(req)
	if err != nil {
		return err
	}

	err = renderer.Render()
	if err != nil {
		return err
	}

	// if false {
	// 	hscript, err = os.Open(fname)
	// 	handle_err(err)
	// 	hprint, err := os.Create(fname + ".ok")
	// 	handle_err(err)

	// 	pprint := exec.Command("python", "-c", "import yaml, sys; o = yaml.load(sys.stdin); yaml.dump(o, stream=sys.stdout)")
	// 	pprint.Stdin = hscript
	// 	pprint.Stdout = hprint
	// 	err = pprint.Run()
	// }

	return err
}

func sanitize_srcs(sources []string) (srcs []string, rest []string) {
	srcs = make([]string, 0, len(sources))
	rest = make([]string, 0)
	for _, src := range sources {
		if strings.HasPrefix(src, "../") {
			src = src[len("../"):]
		}
		if strings.HasPrefix(src, "-") {
			// discard -globals -no_prototypes -s=$(some)/src
			rest = append(rest, src)
			continue
		}
		srcs = append(srcs, src)
	}
	return srcs, rest
}

func sanitize_env_string(v string) string {
	v = strings.Replace(v, "$(", "${", -1)
	v = strings.Replace(v, ")", "}", -1)
	if strings.HasPrefix(v, `"`) {
		v = v[1:]
	}
	if strings.HasSuffix(v, `"`) {
		v = v[0 : len(v)-1]
	}
	return v
}

func sanitize_env_strings(v []string) []string {
	o := make([]string, 0, len(v))
	for _, vv := range v {
		vv = sanitize_env_string(vv)
		o = append(o, vv)
	}
	return o
}

// w_distill_tgt inspects a list of CMT macro statements and
// converts these macros into their corresponding waf syntax,
// directly adding these to the hlib.Target_t target.
//
// Note: we only do that for macros whose values are simple
//       ie: no cmt-tag is involved.
func w_distill_tgt(tgt *hlib.Target_t, macros map[string][]Stmt) {
	type mungefct_t func(s string) string

	env_munge := func(s string) string {
		out := s
		out = strings.Replace(out, "$(", "${", -1)
		out = strings.Replace(out, ")", "}", -1)
		return out
	}

	linkopts_munge := func(s string) string {
		if strings.HasPrefix(s, "-l") {
			s = env_munge(s[len("-l"):])
		}
		return s
	}

	noop_munge := func(s string) string {
		return s
	}

	type munger_ctx struct {
		suffix string
		fct    mungefct_t
		out    *[]hlib.Value
	}

	mungers := []munger_ctx{
		{
			suffix: "_shlibflags",
			fct:    linkopts_munge,
			out:    &tgt.Use,
		},
		{
			suffix: "linkopts",
			fct:    linkopts_munge,
			out:    &tgt.Use,
		},
		{
			suffix: "_pp_cppflags",
			fct:    noop_munge,
			out:    &tgt.CxxFlags,
		},
		{
			suffix: "_cxxflags",
			fct:    noop_munge,
			out:    &tgt.CxxFlags,
		},
		{
			suffix: "_cflags",
			fct:    noop_munge,
			out:    &tgt.CFlags,
		},
	}

	// defines_munge := func(s string) string {
	// 	if strings.HasPrefix(s, "-D") {
	// 		s = s[len("-D"):]
	// 	}
	// 	return s
	// }

	for n, stmts := range macros {
		if !strings.HasPrefix(n, tgt.Name) {
			continue
		}
		// fmt.Printf(">>> [%s]:(%s) %v: [", n, tgt.Name, len(stmts))
		// for _, stmt := range stmts {
		// 	fmt.Printf("%v (%T), ", stmt, stmt)
		// }
		// fmt.Printf("]\n")

		// n_stmts := len(stmts)
		tgt_decl_stmts := make([]Stmt, 0, len(stmts))
		tgt_app_stmts := make([]Stmt, 0, len(stmts))
		tgt_rem_stmts := make([]Stmt, 0, len(stmts))

		for _, stmt := range stmts {
			switch x := stmt.(type) {
			case *Macro:
				if len(x.Set) == 1 {
					tgt_decl_stmts = append(tgt_decl_stmts, x)
				}
			case *MacroAppend:
				if len(x.Set) == 1 {
					tgt_app_stmts = append(tgt_app_stmts, x)
				}
			case *MacroRemove:
				if len(x.Set) == 1 {
					tgt_rem_stmts = append(tgt_rem_stmts, x)
				}
			}
		}

		stmts = make([]Stmt, 0, len(stmts))
		stmts = append(stmts, tgt_decl_stmts...)
		stmts = append(stmts, tgt_app_stmts...)
		stmts = append(stmts, tgt_rem_stmts...)

		// fmt.Printf("+++ [%s]: %d\n", n, len(stmts))
		// if n_stmts != len(stmts) {
		// 	panic(fmt.Errorf("boo: %s: %d -> %d", n, n_stmts, len(stmts)))
		// }

		// do_select := func(name string) bool {
		// 	for _, str := range []string{
		// 		"linkopts",
		// 		"_dependencies",
		// 		"_cflags",
		// 		"_cxxflags",
		// 		"_shlibflags",
		// 	} {
		// 		if strings.HasSuffix(name, str) {
		// 			return true
		// 		}
		// 	}
		// 	return false
		// }

		for _, stmt := range stmts {
			switch x := stmt.(type) {
			case *Macro:
				for _, munger := range mungers {
					if x.Name == tgt.Name+munger.suffix {
						for i, str := range x.Set[0].Value {
							x.Set[0].Value[i] = munger.fct(str)
						}
						*munger.out = append(*munger.out, *(*hlib.Value)(x))
					}
				}
			case *MacroAppend:
				for _, munger := range mungers {
					if x.Name == tgt.Name+munger.suffix {
						for i, str := range x.Set[0].Value {
							x.Set[0].Value[i] = munger.fct(str)
						}
						*munger.out = append(*munger.out, *(*hlib.Value)(x))
					}
				}
			case *MacroRemove:
				for _, munger := range mungers {
					if x.Name == tgt.Name+munger.suffix {
						for i, str := range x.Set[0].Value {
							x.Set[0].Value[i] = munger.fct(str)
						}
						*munger.out = append(*munger.out, *(*hlib.Value)(x))
					}
				}
			}
		}
	}

}

// EOF
