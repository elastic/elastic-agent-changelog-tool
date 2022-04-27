// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package github

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/google/go-github/v32/github"
)

type Strategy interface {
	FindPullRequestID(pr *github.PullRequest) (int, error)
}

var ErrStrategyFailed = errors.New("strategy failed")

type BackportPRNumber struct {
	Strategy
}

func (s *BackportPRNumber) FindPullRequestID(pr *github.PullRequest) (int, error) {
	patterns := []string{`backport #(\d+)`, `cherry-pick of #(\d+)`, `cherry-pick of PR #(\d+)`}
	rDigit, _ := regexp.Compile(`(\d+)`)

	for _, label := range pr.Labels {
		if label.GetName() == "backport" {
			for _, p := range patterns {
				regexPattern, _ := regexp.Compile(p)
				backport := regexPattern.FindString(pr.GetTitle())
				if backport == "" {
					backport = regexPattern.FindString(*pr.Body)
				}

				PRNumber, err := strconv.Atoi(rDigit.FindString(backport))
				if err == nil {
					return PRNumber, err
				}
			}
		}
	}
	return -1, ErrStrategyFailed
}

type PRNumber struct {
	Strategy
}

func (s *PRNumber) FindPullRequestID(pr *github.PullRequest) (int, error) {
	if pr.Number != nil {
		return pr.GetNumber(), nil
	}
	return -1, ErrStrategyFailed
}

func TestStrategies(pr *github.PullRequest, strategies ...Strategy) (int, error) {
	var (
		prID int
		err  error
	)

	for _, s := range strategies {
		prID, err = s.FindPullRequestID(pr)
		if err == nil {
			break
		}
	}

	return prID, err
}

type PRForCommit struct {
	CommitHash    string `json:"commit"`
	PullRequestID int    `json:"pull-request"`
}

type FoundPRs struct {
	Items []PRForCommit `json:"items"`
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

	backportStrategy := &BackportPRNumber{}
	prNumberStrategy := &PRNumber{}

	for i, pr := range prs {
		prID, err := TestStrategies(pr, backportStrategy, prNumberStrategy)
		if err != nil {
			return FoundPRs{}, fmt.Errorf("failed testing strategies: %w", err)
		}

		respData.Items[i] = PRForCommit{
			CommitHash:    commit,
			PullRequestID: prID,
		}
	}

	return respData, nil
}
