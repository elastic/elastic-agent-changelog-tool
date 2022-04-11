// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package cmd

import (
	"fmt"

	"github.com/elastic/elastic-agent-changelog-tool/internal/changelog"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func BuildCmd(fs afero.Fs) *cobra.Command {

	buildCmd := &cobra.Command{
		Use:  "build",
		Long: "Create changelog from fragments",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := viper.GetString("changelog_filename")
			src := viper.GetString("fragment_location")
			dest := viper.GetString("changelog_destination")

			b := changelog.NewBuilder(fs, filename, "8.2.1", src, dest)

			if err := b.Build(); err != nil {
				return fmt.Errorf("cannot build changelog: %w", err)
			}

			return nil
		},
	}

	return buildCmd
}
