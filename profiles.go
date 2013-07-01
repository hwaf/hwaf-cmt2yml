package main

type Profile struct {
	patterns map[string]string
	features map[string]string
}

var (
	g_profile *Profile = nil
)

func init() {
	g_tdaq_profile := &Profile{
		patterns: map[string]string{
			// TDAQCExternal
			"declare_lcg_mapping":  "tdaq_declare_lcg_mapping",
			"external.RPM.package": "tdaq_external_rpm_package",
			"external.RPM.post":    "tdaq_external_rpm_post",
			"external.RPM.preun":   "tdaq_external_rpm_preun",
			"make_external_slinks": "tdaq_make_external_slinks",

			// TDAQCPolicy/TDAQCPolicyInt
			"check_target":          "tdaq_check_target",
			"copy_file":             "tdaq_copy_file",
			"global_install_dirs":   "tdaq_global_install_dirs",
			"global_rpms":           "tdaq_global_rpms",
			"global_rpms_macros":    "tdaq_global_rpms_macros",
			"include_path_1":        "tdaq_include_path_1",
			"inst_docs_auto":        "tdaq_inst_docs_auto",
			"inst_headers_auto":     "tdaq_inst_headers_auto",
			"inst_headers_bin_auto": "tdaq_inst_headers_bin_auto",
			"inst_idl_auto":         "tdaq_inst_idl_auto",
			"inst_scripts_auto":     "tdaq_inst_scripts_auto",
			"install_apps":          "tdaq_install_apps",
			"install_data":          "tdaq_install_data",
			"install_dir":           "tdaq_install_dir",
			"install_docs":          "tdaq_install_docs",
			"install_examples":      "tdaq_install_examples",
			"install_headers":       "tdaq_install_headers",
			"install_libs":          "tdaq_install_libs",
			"install_scripts":       "tdaq_install_scripts",
			"release_inst_path":     "tdaq_release_inst_path",
			"set_cmtpath":           "tdaq_set_cmtpath",
			"set_release_package":   "tdaq_set_release_package",
		},
		features: map[string]string{
			"application": "tdaq_application",
			"library":     "tdaq_library",
		},
	}

	g_profile = g_tdaq_profile
}

// EOF
