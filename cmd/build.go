// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package cmd

import (
	"fmt"
	"log"

	"github.com/elastic/elastic-agent-changelog-tool/internal/changelog"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func BuildCmd(fs afero.Fs) *cobra.Command {

	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Create changelog from fragments",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := cmd.Flags().GetString("repo")
			if err != nil {
				return fmt.Errorf("repo flag malformed: %w", err)
			}

			owner, err := cmd.Flags().GetString("owner")
			if err != nil {
				return fmt.Errorf("owner flag malformed: %w", err)
			}

			filename := viper.GetString("changelog_filename")
			src := viper.GetString("fragment_location")
			dest := viper.GetString("changelog_destination")

			version, err := cmd.Flags().GetString("version")
			if err != nil {
				return fmt.Errorf("error parsing flag 'version': %w", err)
			}

			b := changelog.NewBuilder(fs, filename, version, src, dest)

			if err := b.Build(owner, repo); err != nil {
				return fmt.Errorf("cannot build changelog: %w", err)
			}

			return nil
		},
	}

	buildCmd.Flags().String("repo", defaultRepo, "target repository")
	buildCmd.Flags().String("owner", defaultOwner, "target repository owner")

	buildCmd.Flags().String("version", "", "The version of the consolidated changelog being created")
	err := buildCmd.MarkFlagRequired("version")
	if err != nil {
		// NOTE: the only case this error appear is when the flag is not defined
		log.Fatal(err)
	}

	return buildCmd
}
