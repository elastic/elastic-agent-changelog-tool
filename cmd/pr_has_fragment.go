// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package cmd

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/elastic/elastic-agent-changelog-tool/internal/github"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var errPrCheckCmdMissingArg = errors.New("pr-has-fragment command requires pr number argument")

func PrHasFragmentCommand(appFs afero.Fs) *cobra.Command {
	prCheckCmd := &cobra.Command{
		Use:   "pr-has-fragment <pr-number>",
		Short: "Check changelog fragment presence in the given PR.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errPrCheckCmdMissingArg
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

			repo = GetRepo(repo)
			owner = GetOwner(owner)

			pr, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			ctx := context.Background()
			// TODO: move this to configuration or flag
			labels := []string{"skip-changelog", "backport"}

			shouldSkip, err := github.PRHasLabels(ctx, c, owner, repo, pr, labels)
			if err != nil {
				return err
			}
			if shouldSkip {
				fmt.Fprintf(cmd.OutOrStdout(), "PR requires no changelog because it has one of these labels: %q\n", labels)
				return nil
			}

			pattern := fmt.Sprintf("%s/*", viper.GetString("fragment_path"))

			found, err := github.FindFileInPR(ctx, c, owner, repo, pr, pattern)
			if err != nil {
				return err
			}
			if !found {
				return fmt.Errorf("fragment not present in PR %d, to resolve this do one of the following:\n"+
					"1) add a fragment using the 'new' command\n"+
					"2) add a label (one of: %q) to skip validation", pr, labels)
			}

			return nil
		},
	}

	prCheckCmd.Flags().String("repo", defaultRepo, "target repository")
	prCheckCmd.Flags().String("owner", defaultOwner, "target repository owner")

	return prCheckCmd
}
