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
