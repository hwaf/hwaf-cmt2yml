package main

import (
	"fmt"
	"path/filepath"

	"github.com/hwaf/hwaf/hlib"
)

// map of pkgname -> libname
//  if empty => ignore dep.
var g_pkg_map = map[string]string{
	"AtlasPolicy":        "",
	"AtlasCxxPolicy":     "",
	"AtlasFortranPolicy": "",
	"ExternalPolicy":     "",
	"GaudiInterface":     "GaudiKernel",
	"AtlasROOT":          "ROOT",
	"AtlasReflex":        "Reflex",
	"AtlasCLHEP":         "CLHEP",
	"AtlasPOOL":          "POOL",
	"AtlasCOOL":          "COOL",
	"AtlasCORAL":         "CORAL",
}

func find_tgt(wscript *hlib.Wscript_t, name string) (int, *hlib.Target_t) {
	wbld := &wscript.Build
	for i := range wbld.Targets {
		if wbld.Targets[i].Name == name {
			return i, &wbld.Targets[i]
		}
	}
	return -1, nil
}

func use_list(wscript *hlib.Wscript_t) []string {
	uses := []string{}
	for _, dep := range wscript.Package.Deps {
		pkg := filepath.Base(dep.Name)
		use_pkg, ok := g_pkg_map[pkg]
		if !ok {
			use_pkg = pkg
		}
		if use_pkg != "" {
			uses = append(uses, pkg)
		}
	}
	return uses
}

func cnv_atlas_library(wscript *hlib.Wscript_t, stmt Stmt) error {
	x := stmt.(*ApplyPattern)
	fmt.Printf(">>> [%s] \n", x.Name)
	return nil
}

func cnv_atlas_component_library(wscript *hlib.Wscript_t, stmt Stmt) error {
	x := stmt.(*ApplyPattern)
	libname := ""
	switch len(x.Args) {
	case 0:
		// component_library pattern
		libname = filepath.Base(wscript.Package.Name)
	default:
		// named_component_library pattern
		libname = x.Args[1]
	}
	itgt, tgt := find_tgt(wscript, libname)
	if itgt < 0 {
		wscript.Build.Targets = append(
			wscript.Build.Targets,
			hlib.Target_t{Name: libname},
		)
		itgt, tgt = find_tgt(wscript, libname)
	}
	tgt.Features = []string{"atlas_component"}
	tgt.Use = []hlib.Value{hlib.DefaultValue("uses", use_list(wscript))}

	fmt.Printf(">>> component [%v]...\n", *tgt)
	return nil
}

func cnv_atlas_dual_use_library(wscript *hlib.Wscript_t, stmt Stmt) error {
	x := stmt.(*ApplyPattern)
	fmt.Printf(">>> [%s] \n", x.Name)
	return nil
}

func cnv_atlas_tpcnv_library(wscript *hlib.Wscript_t, stmt Stmt) error {
	x := stmt.(*ApplyPattern)
	fmt.Printf(">>> [%s] \n", x.Name)
	return nil
}

func cnv_atlas_install_joboptions(wscript *hlib.Wscript_t, stmt Stmt) error {
	x := stmt.(*ApplyPattern)
	fmt.Printf(">>> [%s] \n", x.Name)
	pkgname := filepath.Base(wscript.Package.Name)
	tgt := hlib.Target_t{Name: pkgname + "-install-jobos"}
	tgt.Features = []string{"atlas_install_joboptions"}
	tgt.Source = []hlib.Value{hlib.DefaultValue(
		"jobos",
		[]string{"share/*.py", "share/*.txt"},
	)}
	wscript.Build.Targets = append(wscript.Build.Targets, tgt)
	return nil
}

func cnv_atlas_install_python_modules(wscript *hlib.Wscript_t, stmt Stmt) error {
	x := stmt.(*ApplyPattern)
	fmt.Printf(">>> [%s] \n", x.Name)
	pkgname := filepath.Base(wscript.Package.Name)
	tgt := hlib.Target_t{Name: pkgname + "-install-py"}
	tgt.Features = []string{"atlas_install_python_modules"}
	tgt.Source = []hlib.Value{hlib.DefaultValue(
		"python-files",
		[]string{"python/*.py"},
	)}
	wscript.Build.Targets = append(wscript.Build.Targets, tgt)
	return nil
}

func cnv_atlas_install_scripts(wscript *hlib.Wscript_t, stmt Stmt) error {
	x := stmt.(*ApplyPattern)
	fmt.Printf(">>> [%s] \n", x.Name)
	pkgname := filepath.Base(wscript.Package.Name)
	tgt := hlib.Target_t{Name: pkgname + "-install-scripts"}
	tgt.Features = []string{"atlas_install_scripts"}
	tgt.Source = []hlib.Value{hlib.DefaultValue(
		"script-files",
		[]string{"scripts/*"},
	)}
	wscript.Build.Targets = append(wscript.Build.Targets, tgt)
	return nil
}

func cnv_atlas_install_xmls(wscript *hlib.Wscript_t, stmt Stmt) error {
	x := stmt.(*ApplyPattern)
	fmt.Printf(">>> [%s] \n", x.Name)
	pkgname := filepath.Base(wscript.Package.Name)
	tgt := hlib.Target_t{Name: pkgname + "-install-xmls"}
	tgt.Features = []string{"atlas_install_xmls"}
	tgt.Source = []hlib.Value{hlib.DefaultValue(
		"xml-files",
		[]string{"xml/*"},
	)}
	wscript.Build.Targets = append(wscript.Build.Targets, tgt)
	return nil
}

func cnv_atlas_install_data(wscript *hlib.Wscript_t, stmt Stmt) error {
	x := stmt.(*ApplyPattern)
	fmt.Printf(">>> [%s] \n", x.Name)
	pkgname := filepath.Base(wscript.Package.Name)
	tgt := hlib.Target_t{Name: pkgname + "-install-data"}
	tgt.Features = []string{"atlas_install_data"}
	tgt.Source = []hlib.Value{hlib.DefaultValue(
		"data-files",
		[]string{"data/*"},
	)}
	wscript.Build.Targets = append(wscript.Build.Targets, tgt)
	return nil
}

func cnv_atlas_install_java(wscript *hlib.Wscript_t, stmt Stmt) error {
	x := stmt.(*ApplyPattern)
	fmt.Printf(">>> [%s] \n", x.Name)
	return nil
}

func cnv_atlas_dictionary(wscript *hlib.Wscript_t, stmt Stmt) error {
	x := stmt.(*ApplyPattern)
	fmt.Printf(">>> %v\n", x)
	return nil
}

// EOF
