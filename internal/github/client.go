// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License
// 2.0; you may not use this file except in compliance with the Elastic License
// 2.0.

package github

import (
	"context"

	"github.com/pkg/errors"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

// Client function creates new instance of the GitHub API client.
func Client() (*github.Client, error) {
	authToken, err := AuthToken()
	if err != nil {
		return nil, errors.Wrap(err, "reading auth token failed")
	}
	return github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: authToken},
	))), nil
}
