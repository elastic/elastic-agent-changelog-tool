// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package main

import (
	"os"

	"github.com/elastic/elastic-agent-changelog-tool/cmd"
	"github.com/elastic/elastic-agent-changelog-tool/internal/settings"
	"github.com/spf13/afero"
)

func main() {
	settings.Init()

	appFs := afero.NewOsFs()

	rootCmd := cmd.RootCmd()
	rootCmd.AddCommand(cmd.BuildCmd(appFs))
	rootCmd.AddCommand(cmd.CleanupCmd(appFs))
	rootCmd.AddCommand(cmd.FindPRCommand(appFs))
	rootCmd.AddCommand(cmd.NewCmd())
	rootCmd.AddCommand(cmd.PrHasFragmentCommand(appFs))
	rootCmd.AddCommand(cmd.RenderCmd(appFs))
	rootCmd.AddCommand(cmd.VersionCmd())

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
