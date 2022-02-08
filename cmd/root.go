// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License
// 2.0; you may not use this file except in compliance with the Elastic License
// 2.0.

package cmd

import (
	"github.com/spf13/cobra"
)

// commands holds all commands created for elastic-agent-changelog
var commands = []*cobra.Command{
	setupFindPRCommand(),
}

// RootCmd creates and returns root cmd for elastic-agent-changelog.
func RootCmd() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:          "elastic-agent-changelog",
		Short:        "elastic-agent-changelog - Command line tool used for managing the change log for Elastic Agent and related components, including Beats.",
		SilenceUsage: true,
	}

	for _, cmd := range commands {
		rootCmd.AddCommand(cmd)
	}

	return rootCmd
}
