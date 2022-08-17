// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package cmd

import (
	"errors"

	"github.com/elastic/elastic-agent-changelog-tool/internal/changelog/fragment"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var errNewCmdMissingArg = errors.New("new requires title argument")

func NewCmd() *cobra.Command {

	newCmd := &cobra.Command{
		Use:   "new title",
		Short: "Create a new changelog fragment",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errNewCmdMissingArg
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			title := args[0]
			location := viper.GetString("fragment_location")
			fc := fragment.NewCreator(afero.NewOsFs(), location)

			if err := fc.Create(title); err != nil {
				return err
			}

			return nil
		},
	}

	return newCmd
}
