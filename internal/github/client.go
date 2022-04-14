// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package github

import (
	"net/http"

	gh "github.com/google/go-github/v32/github"
)

type Client struct {
	PullRequests githubPullRequestsService
	Users        githubUsersService
}

func NewClient(httpClient *http.Client) *Client {
	client := gh.NewClient(httpClient)

	return &Client{
		PullRequests: client.PullRequests,
		Users:        client.Users,
	}
}
