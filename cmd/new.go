// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package cmd

import (
	"fmt"
	"github.com/elastic/elastic-agent-changelog-tool/internal/changelog/fragment"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os/exec"
	"strings"
)

func NewCmd() *cobra.Command {

	newCmd := &cobra.Command{
		Use:   "new title",
		Short: "Create a new changelog fragment",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Determine changelog fragment title
			var title string
			if len(args) > 0 {
				// Title is the string that follows "elastic-agent-changelog-tool new"
				title = strings.Join(args, " ")
			} else {
				// Title is the name of the git branch
				var err error
				title, err = gitBranchName()
				if err != nil {
					return fmt.Errorf("error setting default title: %w", err)
				}
			}

			title = strings.ToLower(strings.TrimSpace(title))

			location := viper.GetString("fragment_location")
			fc := fragment.NewCreator(afero.NewOsFs(), location)

			if err := fc.Create(title); err != nil {
				return err
			}

			return nil
		},
	}

	newCmd.Flags().String("title", "", "The title for the changelog")

	return newCmd
}

func gitBranchName() (string, error) {
	// git symbolic-ref --short HEAD
	cmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	stdout, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("unable to determine git branch: %w", err)
	}

	return string(stdout), nil
}
