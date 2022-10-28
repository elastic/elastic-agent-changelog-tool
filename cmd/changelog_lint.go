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

func ChangelogLintCmd(fs afero.Fs) *cobra.Command {

	lintCmd := &cobra.Command{
		Use:   "changelog_lint",
		Short: "Lint the consolidated changelog",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dest := viper.GetString("changelog_destination")

			version, err := cmd.Flags().GetString("version")
			if err != nil {
				return fmt.Errorf("error parsing flag 'version': %w", err)
			}

			relaxed, err := cmd.Flags().GetBool("relaxed")
			if err != nil {
				return fmt.Errorf("error parsing flag 'relaxed': %w", err)
			}

			linter := changelog.NewLinter(fs)
			errs := linter.Lint(dest, version)

			for _, err := range errs {
				log.Println(err)
			}

			if !relaxed && len(errs) > 0 {
				log.Fatal("Linting failed.")
			}

			log.Println("Linting done.")

			return nil
		},
	}

	lintCmd.Flags().VisitAll(viperOverrides(lintCmd))

	lintCmd.Flags().String("version", "", "The version of the consolidated changelog being created")
	lintCmd.Flags().Bool("relaxed", false, "Relaxed mode will only log erros, without terminating execution")
	err := lintCmd.MarkFlagRequired("version")
	if err != nil {
		// NOTE: the only case this error appear is when the flag is not defined
		log.Fatal(err)
	}

	return lintCmd
}
