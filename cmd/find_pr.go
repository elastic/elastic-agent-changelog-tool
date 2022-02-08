// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License
// 2.0; you may not use this file except in compliance with the Elastic License
// 2.0.

package cmd

import "github.com/spf13/cobra"

const findPRLongDescription = `Use this command to find the original PR that included the commit in the repository.`

func setupFindPRCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find-pr",
		Short: "Find the original Pull Request",
		Long:  findPRLongDescription,
	}

	return cmd
}
