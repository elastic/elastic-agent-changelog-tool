// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package github

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/v32/github"
)

type PRForCommit struct {
	CommitHash    string `json:"commit"`
	PullRequestID int    `json:"pull-request"`
}

type FoundPRs struct {
	Items []PRForCommit `json:"items"`
}

func findPRID(pr *github.PullRequest) int {
	rBackportTitle, _ := regexp.Compile(`backport #(\d+)`)
	rDigit, _ := regexp.Compile(`(\d+)`)

	for _, label := range pr.Labels {
		if strings.Contains(label.GetName(), "backport") {
			backport := rBackportTitle.FindString(pr.GetTitle())
			PRNumber, err := strconv.Atoi(rDigit.FindString(backport))

			if err == nil {
				return PRNumber
			}
		}
	}
	return pr.GetNumber()
}

func FindPR(ctx context.Context, c *Client, owner, repo, commit string) (FoundPRs, error) {
	prs, _, err := c.PullRequests.ListPullRequestsWithCommit(
		context.Background(), owner, repo, commit, nil)
	if err != nil {
		return FoundPRs{}, fmt.Errorf("failed listing prs with commit: %w", err)
	}

	respData := FoundPRs{
		Items: make([]PRForCommit, len(prs)),
	}

	for i, pr := range prs {
		respData.Items[i] = PRForCommit{
			CommitHash:    commit,
			PullRequestID: findPRID(pr),
		}
	}

	return respData, nil
}
