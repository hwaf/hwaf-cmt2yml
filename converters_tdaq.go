package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hwaf/hwaf/hlib"
)

func cnv_tdaq_library(wscript *hlib.Wscript_t, stmt Stmt) error {
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
			"cmt2yml: empty tdaq_library name (package=%s, args=%v)",
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
	tgt.Features = []string{"tdaq_library"}
	uses := use_list(wscript)
	if len(uses) > 0 {
		tgt.Use = []hlib.Value{hlib.DefaultValue("uses", uses)}
	}

	tgt.Source = []hlib.Value{hlib.DefaultValue("source", source)}
	return nil
}

func cnv_tdaq_application(wscript *hlib.Wscript_t, stmt Stmt) error {
	x := stmt.(*ApplyPattern)
	margs := cmt_arg_map(x.Args)
	appname := margs["name"]
	if appname == "" {
		return fmt.Errorf(
			"cmt2yml: empty tdaq_application name (package=%s, args=%v)",
			wscript.Package.Name,
			x.Args,
		)
	}
	tgtname := appname

	itgt, tgt := find_tgt(wscript, tgtname)
	if itgt < 0 {
		wscript.Build.Targets = append(
			wscript.Build.Targets,
			hlib.Target_t{Name: tgtname},
		)
		itgt, tgt = find_tgt(wscript, tgtname)
	}
	tgt.Features = []string{"tdaq_application"}

	source := []string{fmt.Sprintf("src/%s.cxx", appname)}
	tgt.Source = []hlib.Value{hlib.DefaultValue("source", source)}

	uses := []string{filepath.Base(wscript.Package.Name), "boost-thread"}
	uses = append(uses, use_list(wscript)...)
	tgt.Use = []hlib.Value{hlib.DefaultValue("uses", uses)}

	//fmt.Printf(">>> [%v] \n", *tgt)
	return nil
}

func cnv_tdaq_declare_lcg_mapping(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_declare_lcg_mapping] not implemented\n")
	return nil
}

func cnv_tdaq_external_rpm_package(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_external_rpm_package] not implemented\n")
	return nil
}

func cnv_tdaq_external_rpm_post(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_external_rpm_post] not implemented\n")
	return nil
}

func cnv_tdaq_external_rpm_preun(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_external_rpm_preun] not implemented\n")
	return nil
}

func cnv_tdaq_make_external_slinks(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_make_external_slinks] not implemented\n")
	return nil
}

func cnv_tdaq_check_target(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_check_target] not implemented\n")
	return nil
}

func cnv_tdaq_copy_file(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_copy_file] not implemented\n")
	return nil
}

func cnv_tdaq_global_install_dirs(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_global_install_dirs] not implemented\n")
	return nil
}

func cnv_tdaq_global_rpms(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_global_rpms] not implemented\n")
	return nil
}

func cnv_tdaq_global_rpms_macros(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_global_rpms_macros] not implemented\n")
	return nil
}

func cnv_tdaq_include_path_1(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_include_path_1] not implemented\n")
	return nil
}

func cnv_tdaq_inst_docs_auto(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_inst_docs_auto] not implemented\n")
	return nil
}

func cnv_tdaq_inst_headers_auto(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_inst_headers_auto] not implemented\n")
	return nil
}

func cnv_tdaq_inst_headers_bin_auto(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_inst_headers_bin_auto] not implemented\n")
	return nil
}

func cnv_tdaq_inst_idl_auto(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_inst_idl_auto] not implemented\n")
	return nil
}

func cnv_tdaq_inst_scripts_auto(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_inst_scripts_auto] not implemented\n")
	return nil
}

func cnv_tdaq_install_apps(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_install_apps] not implemented\n")
	return nil
}

func cnv_tdaq_install_data(wscript *hlib.Wscript_t, stmt Stmt) error {
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

func cnv_tdaq_install_dir(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_install_dir] not implemented\n")
	return nil
}

func cnv_tdaq_install_docs(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_install_docs] not implemented\n")
	return nil
}

func cnv_tdaq_install_examples(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_install_examples] not implemented\n")
	return nil
}

func cnv_tdaq_install_headers(wscript *hlib.Wscript_t, stmt Stmt) error {
	pkgname := filepath.Base(wscript.Package.Name)
	tgt := hlib.Target_t{Name: pkgname + "-install-headers"}
	tgt.Features = []string{"tdaq_install_headers"}
	wscript.Build.Targets = append(wscript.Build.Targets, tgt)
	return nil
}

func cnv_tdaq_install_libs(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_install_libs] not implemented\n")
	return nil
}

func cnv_tdaq_install_scripts(wscript *hlib.Wscript_t, stmt Stmt) error {
	pkgname := filepath.Base(wscript.Package.Name)
	tgt := hlib.Target_t{Name: pkgname + "-install-scripts"}
	tgt.Features = []string{"tdaq_install_scripts"}
	tgt.Source = []hlib.Value{hlib.DefaultValue(
		"script-files",
		[]string{"scripts/*"},
	)}
	wscript.Build.Targets = append(wscript.Build.Targets, tgt)
	return nil
}

func cnv_tdaq_release_inst_path(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_release_inst_path] not implemented\n")
	return nil
}

func cnv_tdaq_set_cmtpath(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_set_cmtpath] not implemented\n")
	return nil
}

func cnv_tdaq_set_release_package(wscript *hlib.Wscript_t, stmt Stmt) error {
	fmt.Fprintf(os.Stderr, "** error ** [cnv_tdaq_set_release_package] not implemented\n")
	return nil
}

// EOF
