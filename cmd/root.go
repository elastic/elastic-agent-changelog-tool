// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package cmd

import (
	"log"

	"github.com/elastic/elastic-agent-changelog-tool/internal/settings"
	"github.com/spf13/cobra"
)

const defaultOwner = "elastic"
const defaultRepo = "elastic-agent"

var config *settings.Config

func GetOwner(flagOwner string) string {
	switch {
	case config.Owner != "":
		return config.Owner
	default:
		return flagOwner
	}
}

func GetRepo(flagRepo string) string {
	switch {
	case config.Repo != "":
		return config.Repo
	default:
		return flagRepo
	}
}

// RootCmd creates and returns root cmd for elastic-agent-changelog-tool.
func RootCmd() *cobra.Command {
	var err error

	config, err = settings.LoadConfig()
	if err != nil {
		log.Printf("could not load config: %s", err)
	}

	rootCmd := &cobra.Command{
		Use:          "elastic-agent-changelog-tool",
		Short:        "elastic-agent-changelog-tool - Command line tool used for managing the change log for Elastic Agent and related components, including Beats.",
		SilenceUsage: true,
	}

	return rootCmd
}
