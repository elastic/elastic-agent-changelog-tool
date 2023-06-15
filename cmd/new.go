// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package cmd

import (
	"bytes"
	"fmt"
	"github.com/elastic/elastic-agent-changelog-tool/internal/changelog/fragment"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os/exec"
)

func NewCmd() *cobra.Command {

	newCmd := &cobra.Command{
		Use:   "new title",
		Short: "Create a new changelog fragment",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			title, err := cmd.Flags().GetString("title")
			if err != nil {
				return fmt.Errorf("error parsing flag 'title': %w", err)
			}

			if title == "" {
				title, err = defaultTitle()
				if err != nil {
					return fmt.Errorf("error setting default title: %w", err)
				}
			}

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

func defaultTitle() (string, error) {
	// git symbolic-ref --short HEAD
	cmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	stdout, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("unable to determine git branch: %w", err)
	}

	branchName := bytes.ToLower(bytes.TrimSpace(stdout))
	return string(branchName), nil
}
