// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/elastic/elastic-agent-changelog-tool/internal/github"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var errListPRCmdMissingCommitHash = errors.New("find-pr requires commit hash argument")

const defaultOwner = "elastic"
const defaultRepo = "beats"

const findPRLongDescription = `Use this command to find the original PR that included the commit in the repository.

argument with commit hash is required
--repo flag is optional and will default to elastic/beats if left unspecified.`

func FindPRCommand(appFs afero.Fs) *cobra.Command {
	findPRCommand := &cobra.Command{
		Use:  "find-pr",
		Long: findPRLongDescription,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errListPRCmdMissingCommitHash
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			hc, err := github.GetHTTPClient(appFs)
			if err != nil {
				return fmt.Errorf("cannot initialize http client: %w", err)
			}

			c := github.NewClient(hc)

			repo, err := cmd.Flags().GetString("repo")
			if err != nil {
				return fmt.Errorf("repo flag malformed: %w", err)
			}

			owner, err := cmd.Flags().GetString("owner")
			if err != nil {
				return fmt.Errorf("owner flag malformed: %w", err)
			}

			commit := args[0]
			ctx := context.Background()

			res, err := github.FindPR(ctx, c, owner, repo, commit)
			if err != nil {
				return fmt.Errorf("failed listing prs with commit: %w", err)
			}

			respJSON, err := json.Marshal(res)
			if err != nil {
				return fmt.Errorf("failed marshalling JSON output: %w", err)
			}

			cmd.Println(string(respJSON))

			return nil
		},
	}

	findPRCommand.Flags().String("repo", defaultRepo, "target repository")
	findPRCommand.Flags().String("owner", defaultOwner, "target repository owner")

	return findPRCommand
}
