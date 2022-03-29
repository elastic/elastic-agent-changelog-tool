// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/elastic/elastic-agent-changelog-tool/internal/version"
)

const versionLongDescription = `Use this command to print the version of elastic-package that you have installed. This is especially useful when reporting bugs.`

func VersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show application version",
		Long:  versionLongDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			var sb strings.Builder
			sb.WriteString("elastic-agent-changelog-tool ")
			if version.Tag != "" {
				sb.WriteString(version.Tag)
				sb.WriteString(" ")
			} else {
				sb.WriteString("devel ")
			}
			sb.WriteString(fmt.Sprintf("version-hash %s ", version.CommitHash))
			sb.WriteString(fmt.Sprintf("(build time: %s)", version.BuildTimeFormatted()))

			// NOTE: allow replacing stdout for testing
			fmt.Fprintln(cmd.OutOrStdout(), sb.String())

			return nil
		},
	}

	return cmd
}
