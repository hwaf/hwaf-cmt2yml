package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Renderer struct {
	req     *ReqFile
	wscript bool
	w       *os.File
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

func sanitize_srcs(sources []string) []string {
	for i, src := range sources {
		if strings.HasPrefix(src, "../") {
			sources[i] = src[len("../"):]
		}
	}
	return sources
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
