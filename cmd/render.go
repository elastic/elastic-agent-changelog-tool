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

func RenderCmd(fs afero.Fs) *cobra.Command {
	renderCmd := &cobra.Command{
		Use:   "render",
		Short: "Render a changelog in an asciidoc file",
		Long:  "Render a changelog in an asciidoc file",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dest := viper.GetString("changelog_destination")
			renderedDest := viper.GetString("rendered_changelog_destination")

			version, err := cmd.Flags().GetString("version")
			if err != nil {
				return fmt.Errorf("error parsing flag 'version': %w", err)
			}

			c, err := changelog.FromFile(fs, fmt.Sprintf("./%s/%s.yaml", dest, version))
			if err != nil {
				return fmt.Errorf("error loading changelog from file: %w", err)
			}

			r := changelog.NewRenderer(fs, c, renderedDest)

			if err := r.Render(); err != nil {
				return fmt.Errorf("cannot build asciidoc file: %w", err)
			}

			return nil
		},
	}

	renderCmd.Flags().String("version", "", "The version of the consolidated changelog being created")
	err := renderCmd.MarkFlagRequired("version")
	if err != nil {
		// NOTE: the only case this error appear is when the flag is not defined
		log.Fatal(err)
	}

	return renderCmd
}
