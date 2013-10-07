package main

import (
	"fmt"
	"path/filepath"
	//"strings"

	"github.com/hwaf/hwaf/hlib"
)

func cnv_detcommon_shared_library(wscript *hlib.Wscript_t, stmt Stmt) error {
	x := stmt.(*ApplyPattern)
	margs := cmt_arg_map(x.Args)
	libname := ""
	if _, ok := margs["library"]; ok {
		libname = margs["library"]
	} else {
		libname = filepath.Base(wscript.Package.Name)
	}
	source := []string{"src/*.cxx"}
	if _, ok := margs["files"]; ok {
		source = []string{margs["files"]}
	}

	if libname == "" {
		return fmt.Errorf(
			"cmt2yml: empty atlas_library name (package=%s, args=%v)",
			wscript.Package.Name,
			x.Args,
		)
	}
	itgt, tgt := find_tgt(wscript, libname)
	if itgt < 0 {
		wscript.Build.Targets = append(
			wscript.Build.Targets,
			hlib.Target_t{Name: libname},
		)
		itgt, tgt = find_tgt(wscript, libname)
	}
	tgt.Features = []string{"atlas_library"}
	uses := use_list(wscript)
	if len(uses) > 0 {
		tgt.Use = []hlib.Value{hlib.DefaultValue("uses", uses)}
	}

	tgt.Source = []hlib.Value{hlib.DefaultValue("source", source)}
	return nil
}

func cnv_detcommon_install_headers(wscript *hlib.Wscript_t, stmt Stmt) error {
	pkgname := filepath.Base(wscript.Package.Name)
	tgt := hlib.Target_t{Name: pkgname + "-install-headers"}
	tgt.Features = []string{"detcommon_install_headers"}
	wscript.Build.Targets = append(wscript.Build.Targets, tgt)
	return nil
}

func cnv_trigconf_application(wscript *hlib.Wscript_t, stmt Stmt) error {
	x := stmt.(*ApplyPattern)
	margs := cmt_arg_map(x.Args)
	appname := margs["name"]
	if appname == "" {
		return fmt.Errorf(
			"cmt2yml: empty trigconf_application name (package=%s, args=%v)",
			wscript.Package.Name,
			x.Args,
		)
	}
	tgtname := "TrigConf" + appname

	itgt, tgt := find_tgt(wscript, tgtname)
	if itgt < 0 {
		wscript.Build.Targets = append(
			wscript.Build.Targets,
			hlib.Target_t{Name: tgtname},
		)
		itgt, tgt = find_tgt(wscript, tgtname)
	}
	tgt.Features = []string{"trigconf_application"}

	source := []string{fmt.Sprintf("src/test/%s.cxx", appname)}
	tgt.Source = []hlib.Value{hlib.DefaultValue("source", source)}

	uses := []string{filepath.Base(wscript.Package.Name), "boost-thread"}
	uses = append(uses, use_list(wscript)...)
	tgt.Use = []hlib.Value{hlib.DefaultValue("uses", uses)}

	//fmt.Printf(">>> [%v] \n", *tgt)
	return nil
}

func cnv_detcommon_generic_install(wscript *hlib.Wscript_t, stmt Stmt) error {
	x := stmt.(*ApplyPattern)
	margs := cmt_arg_map(x.Args)
	name := margs["name"]
	source := margs["files"]
	kind := margs["kind"]
	prefix := margs["prefix"]
	pkgname := filepath.Base(wscript.Package.Name)
	tgtname := fmt.Sprintf("%s-generic-install-%s-%s", pkgname, name, kind)

	itgt, tgt := find_tgt(wscript, tgtname)
	if itgt < 0 {
		wscript.Build.Targets = append(
			wscript.Build.Targets,
			hlib.Target_t{Name: tgtname},
		)
		itgt, tgt = find_tgt(wscript, tgtname)
	}
	tgt.Features = []string{"detcommon_generic_install"}
	tgt.Source = []hlib.Value{hlib.DefaultValue("source", []string{source})}
	if prefix != "" {
		if tgt.KwArgs == nil {
			tgt.KwArgs = make(map[string][]hlib.Value)
		}
		tgt.KwArgs["install_prefix"] = []hlib.Value{hlib.DefaultValue("prefix", []string{prefix})}
	}
	return nil
}

// EOF
