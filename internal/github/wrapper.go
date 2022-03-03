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
