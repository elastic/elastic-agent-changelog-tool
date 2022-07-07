// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func CleanupCmd(fs afero.Fs) *cobra.Command {
	cleanupCmd := &cobra.Command{
		Use:  "cleanup",
		Long: "Delete all fragments",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := cmd.Flags().GetString("path")
			if err != nil {
				return fmt.Errorf("error parsing flag 'path': %w", err)
			}

			err = fs.RemoveAll(path)
			if err != nil {
				return fmt.Errorf("error deleting fragments: %w", err)
			}

			err = fs.Mkdir(path, os.ModePerm)
			if err != nil {
				return fmt.Errorf("error creating fragments folder: %w", err)
			}

			return nil
		},
	}

	cleanupCmd.Flags().String("path", "changelog/fragments", "The folder where fragments are stored.")

	return cleanupCmd
}
