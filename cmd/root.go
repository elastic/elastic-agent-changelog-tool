// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package cmd

import (
	"github.com/spf13/cobra"
)

const defaultOwner = "elastic"
const defaultRepo = "beats"

// RootCmd creates and returns root cmd for elastic-agent-changelog-tool.
func RootCmd() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:          "elastic-agent-changelog-tool",
		Short:        "elastic-agent-changelog-tool - Command line tool used for managing the change log for Elastic Agent and related components, including Beats.",
		SilenceUsage: true,
	}

	return rootCmd
}
