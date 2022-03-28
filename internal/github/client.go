// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License
// 2.0; you may not use this file except in compliance with the Elastic License
// 2.0.

package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v32/github"
)

type GithubClient struct {
	client Client
}

func NewClient(gclient Client) (*GithubClient, error) {
	return &GithubClient{
		client: gclient,
	}, nil
}

func NewUnauthorizedClient() *GithubClient {
	return &GithubClient{
		client: NewWrapper(github.NewClient(new(http.Client))),
	}
}

func (c *GithubClient) User() (string, error) {
	user, _, err := c.client.UsersGet(context.Background(), "")
	if err != nil {
		return "", fmt.Errorf("fetching authenticated user failed: %w", err)
	}

	return *user.Login, nil
}

func (c *GithubClient) ListPullRequestsWithCommit(
	ctx context.Context,
	owner string,
	repo string,
	sha string,
	opts *github.PullRequestListOptions,
) ([]*github.PullRequest, *github.Response, error) {
	pr, resp, err := c.client.ListPullRequestsWithCommit(ctx, owner, repo, sha, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list pull requests with commit: %w", err)
	}

	return pr, resp, nil
}
