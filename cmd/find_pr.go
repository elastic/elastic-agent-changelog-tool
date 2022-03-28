// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License
// 2.0; you may not use this file except in compliance with the Elastic License
// 2.0.

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	gb "github.com/elastic/elastic-agent-changelog-tool/internal/github"
	"github.com/google/go-github/v32/github"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var errListPRCmdMissingCommitHash = errors.New("find-pr requires commit hash argument")

const defaultOwner = "elastic"
const defaultRepo = "beats"

const repoFlagName = "repo"
const repoFlagDescription = "target repository"
const findPRLongDescription = `Use this command to find the original PR that included the commit in the repository.

argument with commit hash is required
--repo flag is optional and will default to elastic/beats if left unspecified.`

type PRInfo struct {
	CommitHash    string `json:"commit"`
	PullRequestID string `json:"pull-request"`
}

func prToDomain(commitHash string, pr *github.PullRequest) PRInfo {
	return PRInfo{
		CommitHash:    commitHash,
		PullRequestID: fmt.Sprintf("%d", *pr.Number),
	}
}

func FindPRCommand() *cobra.Command {
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
			authToken := gb.NewAuthToken(&afero.OsFs{})

			githubAccessToken, err := authToken.AuthToken()
			if err != nil {
				log.Fatal(err)
			}

			GithubClient, err := gb.NewClient(gb.NewWrapper(github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: githubAccessToken,
				}),
			))))
			if err != nil {
				log.Fatal(err)
			}

			var commit string

			repo, err := cmd.Flags().GetString(repoFlagName)
			if err != nil {
				return fmt.Errorf("repo flag malformed: %w", err)
			}

			commit = args[0]

			if repo == "" {
				repo = defaultRepo
			}

			prs, _, err := GithubClient.ListPullRequestsWithCommit(context.Background(), defaultOwner, repo, commit, nil)
			if err != nil {
				return fmt.Errorf("failed listing prs with commit: %w", err)
			}

			type resp struct {
				Items []PRInfo `json:"items"`
			}

			respData := resp{
				Items: make([]PRInfo, len(prs)),
			}

			for i, pr := range prs {
				respData.Items[i] = prToDomain(commit, pr)
			}

			respJSON, err := json.Marshal(respData)
			if err != nil {
				return fmt.Errorf("failed listing prs with commit: %w", err)
			}

			cmd.Println(string(respJSON))

			return nil
		},
	}

	findPRCommand.Flags().String(repoFlagName, "", repoFlagDescription)

	return findPRCommand
}
