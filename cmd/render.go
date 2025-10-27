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

var RenderLongDescription = `Use this command to render the consolidated changelog.

--version is required. Consolidated changelog version (x.y.z) in 'changelogs' folder
--file_type is optional. Specify the file_type: 'asciidoc' or 'markdown'. Default: markdown
--template is optional. Specify full path to your template file or use predefined templates. Default: asciidoc-embedded`

func RenderCmd(fs afero.Fs) *cobra.Command {
	renderCmd := &cobra.Command{
		Use:   "render",
		Short: "Render a changelog in an AsciiDoc or Markdown file",
		Long:  RenderLongDescription,
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dest := viper.GetString("changelog_destination")
			renderedDest := viper.GetString("rendered_changelog_destination")
			repo := viper.GetString("repo")
			subsections := viper.GetBool("subsections")

			version, err := cmd.Flags().GetString("version")
			if err != nil {
				return fmt.Errorf("error parsing flag 'version': %w", err)
			}

			file_type, err := cmd.Flags().GetString("file_type")
			if err != nil {
				return fmt.Errorf("error parsing flag 'file_type': %w", err)
			}

			template, err := cmd.Flags().GetString("template")
			if err != nil {
				return fmt.Errorf("error parsing flag 'template': %w", err)
			}

			c, err := changelog.FromFile(fs, fmt.Sprintf("./%s/%s.yaml", dest, version))
			if err != nil {
				return fmt.Errorf("error loading changelog from file: %w", err)
			}

			if file_type == "asciidoc" {
				r := changelog.NewRenderer(fs, c, renderedDest, "asciidoc-embedded", repo, subsections)
				if err := r.Render(); err != nil {
					return fmt.Errorf("cannot build asciidoc file: %w", err)
				}
			} else if file_type == "markdown" {
				r_index := changelog.NewRenderer(fs, c, renderedDest, "markdown-index", repo, subsections)
				if err := r_index.Render(); err != nil {
					return fmt.Errorf("cannot build markdown file: %w", err)
				}
				r_breaking := changelog.NewRenderer(fs, c, renderedDest, "markdown-breaking", repo, subsections)
				if err := r_breaking.Render(); err != nil {
					return fmt.Errorf("cannot build markdown file: %w", err)
				}
				r_deprecations := changelog.NewRenderer(fs, c, renderedDest, "markdown-deprecations", repo, subsections)
				if err := r_deprecations.Render(); err != nil {
					return fmt.Errorf("cannot build markdown file: %w", err)
				}
			} else {
				r := changelog.NewRenderer(fs, c, renderedDest, template, repo, subsections)
				if err := r.Render(); err != nil {
					return fmt.Errorf("cannot build file: %w", err)
				}
			}

			return nil
		},
	}

	renderCmd.Flags().String("file_type", viper.GetString("file_type"), "The file type of the rendered release notes: `asciidoc` or `markdown`")
	renderCmd.Flags().String("template", viper.GetString("template"), "The template used to generate the changelog")
	renderCmd.Flags().String("version", "", "The version of the consolidated changelog being created")
	err := renderCmd.MarkFlagRequired("version")
	if err != nil {
		// NOTE: the only case this error appear is when the flag is not defined
		log.Fatal(err)
	}

	return renderCmd
}
