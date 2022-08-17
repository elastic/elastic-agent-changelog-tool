// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CleanupCmd(fs afero.Fs) *cobra.Command {
	cleanupCmd := &cobra.Command{
		Use:  "cleanup",
		Long: "Delete all fragments",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fragmentLocation := viper.GetString("fragment_location")

			fragments, err := os.ReadDir(fragmentLocation)
			if err != nil {
				return fmt.Errorf("could not get fragments folder: %w", err)
			}

			for _, f := range fragments {
				ext := filepath.Ext(f.Name())

				if ext == ".yaml" || ext == ".yml" {
					err = fs.Remove(filepath.Join(fragmentLocation, f.Name()))
					if err != nil {
						return fmt.Errorf("could not remove fragment: %w", err)
					}
				}
			}

			return nil
		},
	}

	return cleanupCmd
}
