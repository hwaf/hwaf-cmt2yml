package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hwaf/hwaf/hlib"
)

// map of pkgname -> libname
//  if empty => ignore dep.
var g_pkg_map = map[string][]string{
	"AtlasAIDA":          []string{"AIDA"},
	"AtlasBoost":         []string{"AtlasBoost"},
	"AtlasCLHEP":         []string{"CLHEP"},
	"AtlasCOOL":          []string{"COOL"},
	"AtlasCORAL":         []string{"CORAL"},
	"AtlasCppUnit":       []string{"CppUnit"},
	"AtlasCxxPolicy":     nil,
	"AtlasFortranPolicy": nil,
	"AtlasGdb":           []string{"bfd"},
	"AtlasPOOL":          []string{"POOL"},
	"AtlasPolicy":        nil,
	"AtlasPython":        []string{"AtlasPython"},
	"AtlasPyROOT":        []string{"PyROOT"},
	"AtlasROOT":          []string{"ROOT"},
	"AtlasReflex":        []string{"Reflex"},
	"AtlasTBB":           []string{"tbb"},
	"AtlasValgrind":      []string{"valgrind"},
	"DetCommonPolicy":    nil,
	"ExternalPolicy":     nil,
	"GaudiInterface":     []string{"GaudiKernel"},
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
			use_pkg = []string{pkg}
		}
		if len(use_pkg) > 0 {
			uses = append(uses, use_pkg...)
		}
	}
	return uses
}

func cmt_arg_map(args []string) map[string]string {
	o := make(map[string]string, len(args))
	for _, v := range args {
		idx := strings.Index(v, "=")
		if idx < 0 {
			panic(fmt.Errorf("cmt2yml: could not find '=' in string [%s]", v))
		}
		if idx < 1 {
			panic(fmt.Errorf("cmt2yml: malformed string [%s]", v))
		}
		kk := v[:idx]
		vv := v[idx+1:]
		if vv == "" {
			panic(fmt.Errorf("cmt2yml: malformed string [%s]", v))
		}
		if vv[0] == '"' {
			vv = vv[1:]
		}
		if strings.HasPrefix(vv, "../") {
			vv = vv[len("../"):]
		}
		o[kk] = vv
	}
	return o
}

func cnv_atlas_library(wscript *hlib.Wscript_t, stmt Stmt) error {
	x := stmt.(*ApplyPattern)
	libname := ""
	switch len(x.Args) {
	case 0:
		// installed_library pattern
		libname = filepath.Base(wscript.Package.Name)
	default:
		// named_installed_library pattern
		margs := cmt_arg_map(x.Args)
		libname = margs["library"]
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

	//fmt.Printf(">>> [%v] \n", *tgt)
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
		margs := cmt_arg_map(x.Args)
		libname = margs["library"]
	}
	if libname == "" {
		return fmt.Errorf(
			"cmt2yml: empty atlas_component name (package=%s, args=%v)",
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
	tgt.Features = []string{"atlas_component"}
	uses := use_list(wscript)
	if len(uses) > 0 {
		tgt.Use = []hlib.Value{hlib.DefaultValue("uses", uses)}
	}

	//fmt.Printf(">>> component [%v]...\n", *tgt)
	return nil
}

func cnv_atlas_dual_use_library(wscript *hlib.Wscript_t, stmt Stmt) error {
	x := stmt.(*ApplyPattern)
	libname := ""
	switch len(x.Args) {
	case 0:
		// dual_use_library pattern
		libname = filepath.Base(wscript.Package.Name)
	default:
		// named_dual_use_library pattern
		margs := cmt_arg_map(x.Args)
		if _, ok := margs["library"]; ok {
			libname = margs["library"]
		} else {
			libname = filepath.Base(wscript.Package.Name)
		}

	}
	if libname == "" {
		return fmt.Errorf(
			"cmt2yml: empty atlas_dual_use_library name (package=%s, args=%v)",
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
	tgt.Features = []string{"atlas_dual_use_library"}
	uses := use_list(wscript)
	if len(uses) > 0 {
		tgt.Use = []hlib.Value{hlib.DefaultValue("uses", uses)}
	}

	//fmt.Printf(">>> [%v] \n", *tgt)
	return nil
}

func cnv_atlas_tpcnv_library(wscript *hlib.Wscript_t, stmt Stmt) error {
	x := stmt.(*ApplyPattern)
	libname := ""
	switch len(x.Args) {
	case 0:
		// tpcnv_library pattern
		libname = filepath.Base(wscript.Package.Name)
	default:
		// named_tpcnv_library pattern
		margs := cmt_arg_map(x.Args)
		libname = margs["name"]
	}
	if libname == "" {
		return fmt.Errorf(
			"cmt2yml: empty atlas_tpcnv name (package=%s, args=%v)",
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
	tgt.Features = []string{"atlas_tpcnv"}
	uses := use_list(wscript)
	if len(uses) > 0 {
		tgt.Use = []hlib.Value{hlib.DefaultValue("uses", uses)}
	}

	//fmt.Printf(">>> [%v] \n", *tgt)
	return nil
}

func cnv_atlas_install_joboptions(wscript *hlib.Wscript_t, stmt Stmt) error {
	//x := stmt.(*ApplyPattern)
	//fmt.Printf(">>> [%s] \n", x.Name)
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
	//x := stmt.(*ApplyPattern)
	//fmt.Printf(">>> [%s] \n", x.Name)
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
	//x := stmt.(*ApplyPattern)
	//fmt.Printf(">>> [%s] \n", x.Name)
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
	//x := stmt.(*ApplyPattern)
	//fmt.Printf(">>> [%s] \n", x.Name)
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
	//x := stmt.(*ApplyPattern)
	//fmt.Printf(">>> [%s] \n", x.Name)
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
	margs := cmt_arg_map(x.Args)
	pkgname := filepath.Base(wscript.Package.Name)
	libname := margs["dict"] + "Dict"
	selfile := pkgname + "/" + margs["selectionfile"]
	hdrfiles := strings.Split(margs["headerfiles"], " ")
	hdrfiles, _ = sanitize_srcs(hdrfiles, "")

	itgt, tgt := find_tgt(wscript, libname)
	if itgt < 0 {
		wscript.Build.Targets = append(
			wscript.Build.Targets,
			hlib.Target_t{Name: libname},
		)
		itgt, tgt = find_tgt(wscript, libname)
	}
	tgt.Features = []string{"atlas_dictionary"}
	tgt.Source = []hlib.Value{hlib.DefaultValue("source", hdrfiles)}
	if tgt.KwArgs == nil {
		tgt.KwArgs = make(map[string][]hlib.Value)
	}
	tgt.KwArgs["selection_file"] = []hlib.Value{hlib.DefaultValue("selfile", []string{selfile})}
	//tgt.Use = []hlib.Value{hlib.DefaultValue("uses", use_list(wscript))}
	uses := use_list(wscript)
	if len(uses) > 0 {
		tgt.Use = []hlib.Value{hlib.DefaultValue("uses", uses)}
	}
	//fmt.Printf(">>> %v\n", *tgt)
	return nil
}

func cnv_atlas_unittest(wscript *hlib.Wscript_t, stmt Stmt) error {
	x := stmt.(*ApplyPattern)
	margs := cmt_arg_map(x.Args)
	pkgname := filepath.Base(wscript.Package.Name)
	name := margs["unit_test"]
	tgtname := fmt.Sprintf("%s-test-%s", pkgname, name)
	extra := margs["extrapatterns"]
	source := fmt.Sprintf("test/%s_test.cxx", name)

	itgt, tgt := find_tgt(wscript, tgtname)
	if itgt < 0 {
		wscript.Build.Targets = append(
			wscript.Build.Targets,
			hlib.Target_t{Name: tgtname},
		)
		itgt, tgt = find_tgt(wscript, tgtname)
	}
	tgt.Features = []string{"atlas_unittest"}
	tgt.Source = []hlib.Value{hlib.DefaultValue("source", []string{source})}
	if tgt.KwArgs == nil {
		tgt.KwArgs = make(map[string][]hlib.Value)
	}
	if extra != "" {
		tgt.KwArgs["extrapatterns"] = []hlib.Value{
			hlib.DefaultValue("extrapatterns", []string{extra}),
		}
	}
	//tgt.Use = []hlib.Value{hlib.DefaultValue("uses", use_list(wscript))}
	uses := use_list(wscript)
	if len(uses) > 0 {
		tgt.Use = []hlib.Value{hlib.DefaultValue("uses", uses)}
	}
	//fmt.Printf(">>> %v\n", *tgt)
	return nil
}

func cnv_atlas_athenarun_test(wscript *hlib.Wscript_t, stmt Stmt) error {
	x := stmt.(*ApplyPattern)
	margs := cmt_arg_map(x.Args)
	pkgname := filepath.Base(wscript.Package.Name)
	name := margs["name"]
	tgtname := fmt.Sprintf("%s-runtest-%s", pkgname, name)
	options := margs["options"]
	post := margs["post_script"]

	itgt, tgt := find_tgt(wscript, tgtname)
	if itgt < 0 {
		wscript.Build.Targets = append(
			wscript.Build.Targets,
			hlib.Target_t{Name: tgtname},
		)
		itgt, tgt = find_tgt(wscript, tgtname)
	}
	tgt.Features = []string{"atlas_athenarun_test"}
	if tgt.KwArgs == nil {
		tgt.KwArgs = make(map[string][]hlib.Value)
	}
	if options != "" {
		tgt.KwArgs["joboptions"] = []hlib.Value{
			hlib.DefaultValue("options", []string{options}),
		}
	}
	if post != "" {
		tgt.KwArgs["post_script"] = []hlib.Value{
			hlib.DefaultValue("post", []string{post}),
		}
	}
	tgt.Use = []hlib.Value{hlib.DefaultValue("uses", []string{pkgname})}

	//fmt.Printf(">>> %v\n", *tgt)
	return nil
}

func cnv_atlas_generic_install(wscript *hlib.Wscript_t, stmt Stmt) error {
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
	tgt.Features = []string{"atlas_generic_install"}
	tgt.Source = []hlib.Value{hlib.DefaultValue("source", []string{source})}
	if prefix != "" {
		if tgt.KwArgs == nil {
			tgt.KwArgs = make(map[string][]hlib.Value)
		}
		tgt.KwArgs["install_prefix"] = []hlib.Value{hlib.DefaultValue("prefix", []string{prefix})}
	}
	return nil
}

func cnv_atlas_install_trfs(wscript *hlib.Wscript_t, stmt Stmt) error {
	x := stmt.(*ApplyPattern)
	margs := cmt_arg_map(x.Args)
	jo := margs["jo"]
	tfs := margs["tfs"]
	pkgname := filepath.Base(wscript.Package.Name)
	tgtname := fmt.Sprintf("%s-install-trfs", pkgname)

	itgt, tgt := find_tgt(wscript, tgtname)
	if itgt < 0 {
		wscript.Build.Targets = append(
			wscript.Build.Targets,
			hlib.Target_t{Name: tgtname},
		)
		itgt, tgt = find_tgt(wscript, tgtname)
	}
	tgt.Features = []string{"atlas_install_trfs"}
	tgt.Source = nil
	if tgt.KwArgs == nil {
		tgt.KwArgs = make(map[string][]hlib.Value)
	}
	if jo != "" {
		tgt.KwArgs["trf_jo"] = []hlib.Value{hlib.DefaultValue("trf_jo", []string{jo})}
	}
	if tfs != "" {
		tgt.KwArgs["trf_tfs"] = []hlib.Value{hlib.DefaultValue("trf_tfs", []string{tfs})}
	}
	return nil
}

// EOF
