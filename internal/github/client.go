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
	"golang.org/x/oauth2"
)

type GithubClient interface {
	User() (string, error)
}

type githubClient struct {
	client *github.Client
}

func NewClient() (*githubClient, error) {
	authToken, err := AuthToken()
	if err != nil {
		return nil, fmt.Errorf("reading auth token failed: %w", err)
	}
	return &githubClient{
		client: github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: authToken},
		))),
	}, nil
}

func NewUnauthorizedClient() *githubClient {
	return &githubClient{
		client: github.NewClient(new(http.Client)),
	}
}

func (c *githubClient) User() (string, error) {
	user, _, err := c.client.Users.Get(context.Background(), "")
	if err != nil {
		return "", fmt.Errorf("fetching authenticated user failed: %w", err)
	}
	return *user.Login, nil
}
