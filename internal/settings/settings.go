// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package settings

import (
	"log"
	"os"
	"path"

	"github.com/OpenPeeDeeP/xdg"
	"github.com/elastic/elastic-agent-changelog-tool/internal/gitreporoot"
	"github.com/spf13/viper"
)

const envPrefix = "ELASTIC_AGENT_CHANGELOG"
const configFileFolder = "elastic-agent-changelog-tool"

// Init initalize settings and default values
func Init() {
	viper.AutomaticEnv()
	// NOTE: err value is ignored as it only checks for missing argument
	_ = viper.BindEnv(envPrefix)

	setDefaults()
	setConstants()

	viper.AddConfigPath(viper.GetString("config_file"))

	// TODO: better error handling (skip missing file error)
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		// NOTE: we do not fail in this case as it's ok not to have the config file
		// but we want to alert the user that the file has not been found
		log.Println(err)
	}
}

func setDefaults() {
	viper.SetDefault("config_file", path.Join(xdg.ConfigHome(), configFileFolder))

	// try to compute GIT_REPO_ROOT value if empty
	if os.Getenv("GIT_REPO_ROOT") == "" {
		root, err := gitreporoot.Find()
		if err != nil {
			log.Printf("git repo root not found, $GIT_REPO_ROOT will be empty: %v\n", err)
		} else {
			os.Setenv("GIT_REPO_ROOT", root)
		}
	}

	// fragment_root supports env var expansion
	viper.SetDefault("fragment_root", "$GIT_REPO_ROOT")
	viper.SetDefault("fragment_path", "changelog/fragments")
	viper.SetDefault("fragment_location", path.Join(
		os.ExpandEnv(viper.GetString("fragment_root")),
		viper.GetString("fragment_path")))

	viper.SetDefault("changelog_destination", ".")
	viper.SetDefault("rendered_changelog_destination", ".")
}

func setConstants() {
	// viper.Set()
}
