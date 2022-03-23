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

// Init initalize settings and default values
func Init() {
	viper.AutomaticEnv()
	// NOTE: err value is ignored as it only checks for missing argument
	_ = viper.BindEnv("ELASTIC_AGENT_CHANGELOG")

	setDefaults()
	setConstants()
}

func setDefaults() {
	viper.SetDefault("cache_dir", xdg.CacheHome())
	viper.SetDefault("config_dir", xdg.ConfigHome())
	viper.SetDefault("data_dir", xdg.DataHome())

	root, err := gitreporoot.Find()
	if err != nil {
		log.Printf("git repo root not found, $GIT_REPO_ROOT will be empty: %v\n", err)
	} else {
		os.Setenv("GIT_REPO_ROOT", root)
	}

	// fragment_root supports env var expansion
	viper.SetDefault("fragment_root", "$GIT_REPO_ROOT")
	viper.SetDefault("fragment_path", "changelog/fragments")
	viper.SetDefault("fragment_location", path.Join(
		os.ExpandEnv(viper.GetString("fragment_root")),
		viper.GetString("fragment_path")))

	viper.SetDefault("changelog_destination", ".")
	viper.SetDefault("changelog_filename", "changelog.yaml")
}

func setConstants() {
	// viper.Set()
}
