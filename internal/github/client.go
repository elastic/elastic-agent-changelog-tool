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

type githubClient struct {
	client Client
}

func NewClient(gclient Client) (*githubClient, error) {
	return &githubClient{
		client: gclient,
	}, nil
}

func NewUnauthorizedClient() *githubClient {
	return &githubClient{
		client: NewWrapper(github.NewClient(new(http.Client))),
	}
}

func (c *githubClient) User() (string, error) {
	user, _, err := c.client.UsersGet(context.Background(), "")
	if err != nil {
		return "", fmt.Errorf("fetching authenticated user failed: %w", err)
	}

	return *user.Login, nil
}
