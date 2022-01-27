// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License
// 2.0; you may not use this file except in compliance with the Elastic License
// 2.0.

package cmd

import (
	"sort"

	"github.com/elastic/elastic-agent-changelog/internal/cobraext"

	"github.com/spf13/cobra"
)

var commands = []*cobraext.Command{}

// RootCmd creates and returns root cmd for elastic-agent-changelog.
func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "elastic-agent-changelog",
		Short:        "elastic-agent-changelog - Command line tool used for managing the change log for Elastic Agent and related components, including Beats.",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return cobraext.ComposeCommandActions(cmd, args)
		},
	}

	for _, cmd := range commands {
		rootCmd.AddCommand(cmd.Command)
	}
	return rootCmd
}

// Commands returns the list of commands that have been setup for elastic-agent-changelog.
func Commands() []*cobraext.Command {
	sort.SliceStable(commands, func(i, j int) bool {
		return commands[i].Name() < commands[j].Name()
	})

	return commands
}
