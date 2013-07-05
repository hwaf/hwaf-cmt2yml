package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
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

	wscript := hlib.Wscript_t{
		Package:   hlib.Package_t{Name: basedir},
		Configure: hlib.Configure_t{Env: make(hlib.Env_t)},
		Build: hlib.Build_t{Env: make(hlib.Env_t)},
	}

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
	fmt.Printf("foo=%s\n", wscript.Package.Name)

	// second pass to collect
	for _, stmt := range r.req.Stmts {
		wpkg := &wscript.Package
		wbld := &wscript.Build
		//wcfg := &wscript.Configure
		switch x := stmt.(type) {
		case *Author:
			wpkg.Authors = append(wpkg.Authors, x.Name)
		case *Manager:
			wpkg.Managers = append(wpkg.Managers, x.Name)
		case *Version:
			wpkg.Version = x.Value
		case *UsePkg:
			deptype := hlib.PrivateDep
			if !x.IsPrivate {
				deptype = hlib.PublicDep
			}
			if str_is_in_slice(x.Switches, "-no_auto_imports") {
				deptype |= hlib.RuntimeDep
			}
			wpkg.Deps = append(
				wpkg.Deps, 
				hlib.Dep_t{
					Name: path.Join(x.Path, x.Package),
					Version: x.Version,
					Type: deptype,
				},
			)
			
		case *Library:
			tgt := hlib.Target_t{Name: x.Name}
			srcs, rest := sanitize_srcs(x.Source)
			// FIXME: handle -s=some/dir
			if len(rest) > 0 {
			}
			val := hlib.Value{Name: x.Name, Default:srcs}
			tgt.Source = append(tgt.Source, val)
			if features, ok := g_profile.features["library"]; ok {
				tgt.Features = features
			}
			wbld.Targets = append(wbld.Targets, tgt)
		}
	}

	for _, stmt := range r.req.Stmts {
		switch stmt.(type) {
		case *PathRemove, *MakeFragment, *Pattern:
			r.wscript = true
			return nil
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

func render_yaml(req *ReqFile) error {
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

func init_env_map_from(env henv_t, key string) map[string]interface{} {
	vv := map[string]interface{}{}
	old, haskey := env[key]
	if haskey {
		switch old.(type) {
		case string:
			vv["default"] = old
			panic("boo")
		case map[string]interface{}:
			old := env[key].(map[string]interface{})
			for k, _ := range old {
				vk := sanitize_env_string(k)
				vk = strings.Trim(vk, " ")
				vv[vk] = old[k]
			}
		default:
			panic(fmt.Sprintf("unknown type: %T", old))
		}
	}
	return vv
}

// EOF
