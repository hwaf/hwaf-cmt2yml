package main

import (
	"fmt"

	"github.com/hwaf/hwaf/hlib"
)

func (r *Renderer) render_wscript() error {
	var err error

	enc := hlib.NewHscriptPyEncoder(r.w)
	if enc == nil {
		return fmt.Errorf("got invalid hscript-py encoder")
	}

	err = enc.Encode(&r.pkg)
	if err != nil {
		return err
	}

	return err
}

// EOF
