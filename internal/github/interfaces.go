// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package github

import (
	"context"

	gh "github.com/google/go-github/v32/github"
)

type githubPullRequestsService interface {
	// https://pkg.go.dev/github.com/google/go-github/v32/github#PullRequestsService.Get
	Get(ctx context.Context, owner string, repo string, number int) (*gh.PullRequest, *gh.Response, error)
	// https://pkg.go.dev/github.com/google/go-github/v32/github#PullRequestsService.ListPullRequestsWithCommit
	ListPullRequestsWithCommit(
		ctx context.Context,
		owner string,
		repo string,
		sha string,
		opts *gh.PullRequestListOptions,
	) ([]*gh.PullRequest, *gh.Response, error)
	// https://pkg.go.dev/github.com/google/go-github/v32/github#PullRequestsService.ListFiles
	ListFiles(ctx context.Context, owner string, repo string, number int, opts *gh.ListOptions) ([]*gh.CommitFile, *gh.Response, error)
	Get(
		ctx context.Context,
		owner string,
		repo string,
		number int,
	) (*gh.PullRequest, *gh.Response, error)
}

type githubUsersService interface {
	Get(ctx context.Context, user string) (*gh.User, *gh.Response, error)
}
