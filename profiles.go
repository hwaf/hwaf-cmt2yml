package main

import (
	"github.com/hwaf/hwaf/hlib"
)

type cnvfct_t func(wscript *hlib.Wscript_t, stmt Stmt) error

type Profile struct {
	features map[string][]string
	cnvs     map[string]cnvfct_t
}

var (
	g_profile  *Profile = nil
	g_profiles map[string]*Profile
)

func init() {
	g_profiles = make(map[string]*Profile)
	g_profiles["tdaq"] = &Profile{
		features: map[string][]string{
			"application": []string{"tdaq_application"},
			"library":     []string{"tdaq_library"},
		},
		cnvs: map[string]cnvfct_t{
			// TDAQCExternal
			"declare_lcg_mapping":  cnv_tdaq_declare_lcg_mapping,
			"external.RPM.package": cnv_tdaq_external_rpm_package,
			"external.RPM.post":    cnv_tdaq_external_rpm_post,
			"external.RPM.preun":   cnv_tdaq_external_rpm_preun,
			"make_external_slinks": cnv_tdaq_make_external_slinks,

			// TDAQCPolicy/TDAQCPolicyInt
			"check_target":          cnv_tdaq_check_target,
			"copy_file":             cnv_tdaq_copy_file,
			"global_install_dirs":   cnv_tdaq_global_install_dirs,
			"global_rpms":           cnv_tdaq_global_rpms,
			"global_rpms_macros":    cnv_tdaq_global_rpms_macros,
			"include_path_1":        cnv_tdaq_include_path_1,
			"inst_docs_auto":        cnv_tdaq_inst_docs_auto,
			"inst_headers_auto":     cnv_tdaq_inst_headers_auto,
			"inst_headers_bin_auto": cnv_tdaq_inst_headers_bin_auto,
			"inst_idl_auto":         cnv_tdaq_inst_idl_auto,
			"inst_scripts_auto":     cnv_tdaq_inst_scripts_auto,
			"install_apps":          cnv_tdaq_install_apps,
			"install_data":          cnv_tdaq_install_data,
			"install_dir":           cnv_tdaq_install_dir,
			"install_docs":          cnv_tdaq_install_docs,
			"install_examples":      cnv_tdaq_install_examples,
			"install_headers":       cnv_tdaq_install_headers,
			"install_libs":          cnv_tdaq_install_libs,
			"install_scripts":       cnv_tdaq_install_scripts,
			"release_inst_path":     cnv_tdaq_release_inst_path,
			"set_cmtpath":           cnv_tdaq_set_cmtpath,
			"set_release_package":   cnv_tdaq_set_release_package,
		},
	}

	g_profiles["atlasoff"] = &Profile{
		features: map[string][]string{
			"application": []string{"atlas_application"},
			"library":     []string{"atlas_library"},
		},
		cnvs: map[string]cnvfct_t{
			// DetCommonPolicy
			"detcommon_shared_library":         cnv_detcommon_shared_library,
			"detcommon_shared_generic_library": cnv_detcommon_shared_library,
			"detcommon_shared_named_library":   cnv_detcommon_shared_library,
			"detcommon_header_installer":       cnv_detcommon_install_headers,
			"trigconf_application":             cnv_trigconf_application,
			"trigconf_generic_application":     cnv_trigconf_application,
			"detcommon_generic_install":        cnv_detcommon_generic_install,
			"detcommon_link_files":             cnv_detcommon_generic_install,
			"detcommon_copy_files":             cnv_detcommon_generic_install,
			"detcommon_install_docs":           cnv_detcommon_generic_install,

			// AtlasPolicy
			"installed_library":        cnv_atlas_library,
			"named_installed_library":  cnv_atlas_library,
			"component_library":        cnv_atlas_component_library,
			"named_component_library":  cnv_atlas_component_library,
			"dual_use_library":         cnv_atlas_dual_use_library,
			"named_dual_use_library":   cnv_atlas_dual_use_library,
			"tpcnv_library":            cnv_atlas_tpcnv_library,
			"named_tpcnv_library":      cnv_atlas_tpcnv_library,
			"declare_joboptions":       cnv_atlas_install_joboptions,
			"declare_data":             cnv_atlas_install_data,
			"declare_python_modules":   cnv_atlas_install_python_modules,
			"declare_scripts":          cnv_atlas_install_scripts,
			"declare_xmls":             cnv_atlas_install_xmls,
			"declare_java":             cnv_atlas_install_java,
			"generic_declare_for_link": cnv_atlas_generic_install,

			// TestPolicy
			"UnitTest_run":   cnv_atlas_unittest,
			"athenarun_test": cnv_atlas_athenarun_test,

			// AtlasReflex
			"lcgdict": cnv_atlas_dictionary,

			// PyJobTransforms
			"declare_job_transforms": cnv_atlas_install_trfs,
		},
	}

	g_profile = g_profiles["tdaq"]
}

// EOF
