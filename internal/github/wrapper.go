// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package github

import (
	"context"

	"github.com/google/go-github/v32/github"
)

type Client interface {
	UsersGet(ctx context.Context, user string) (*github.User, *github.Response, error)
}

type Wrapper struct {
	client *github.Client
}

func NewWrapper(client *github.Client) *Wrapper {
	return &Wrapper{
		client: client,
	}
}

func (gw *Wrapper) UsersGet(ctx context.Context, user string) (*github.User, *github.Response, error) {
	return gw.client.Users.Get(ctx, user)
}
