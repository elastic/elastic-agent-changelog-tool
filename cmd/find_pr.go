// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License
// 2.0; you may not use this file except in compliance with the Elastic License
// 2.0.

package cmd

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/elastic/elastic-agent-changelog/internal/github"
)

const commitHashLen = 40
const defaultOwner = "elastic"
const defaultRepo = "beats"

const repoFlagName = "repo"
const repoFlagDescription = "target repository"
const findPRLongDescription = `Use this command to find the original PR that included the commit in the repository.

argument with commit hash is required
--repo flag is optional and will default to elastic/beats if left unspecified.`

func setupFindPRCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find-pr",
		Short: "Find the original Pull Request",
		Long:  findPRLongDescription,
		RunE:  findPRCommandAction,
	}
	cmd.Flags().String(repoFlagName, "", repoFlagDescription)

	return cmd
}

func findPRCommandAction(cmd *cobra.Command, args []string) error {
	var commit string

	repo, err := cmd.Flags().GetString(repoFlagName)
	if err != nil {
		return errors.Wrap(err, "repo flag malformed")
	}

	if len(args) > 0 {
		if len(args[0]) != commitHashLen {
			return errors.Wrap(err, "commit hash malformed")
		}
		commit = args[0]
	} else {
		return errors.Wrap(err, "commit hash argument not found")
	}

	if repo == "" {
		repo = defaultRepo
	}

	// Setup GitHub
	err = github.EnsureAuthConfigured()
	if err != nil {
		return errors.Wrap(err, "GitHub auth configuration failed")
	}

	githubClient, err := github.Client()
	if err != nil {
		return errors.Wrap(err, "creating GitHub client failed")
	}

	//TODO heuristics in case of backport label
	pullRequest, _, err := githubClient.PullRequests.ListPullRequestsWithCommit(context.Background(), defaultOwner, repo, commit, nil)
	if err != nil {
		return errors.Wrap(err, "fetching GitHub API data failed")
	}

	cmd.Printf("%s:%d\n", commit, *pullRequest[0].Number)

	return nil
}
