// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package github

import (
	"context"
	"fmt"
	"net/http"

	gh "github.com/google/go-github/v32/github"
	"github.com/shurcooL/githubv4"
)

type Client struct {
	PullRequests githubPullRequestsService
	Users        githubUsersService
}

type ClientGraphQL struct {
	PR githubGraphQLPRService
}

func NewClient(httpClient *http.Client) *Client {
	client := gh.NewClient(httpClient)

	return &Client{
		PullRequests: client.PullRequests,
		Users:        client.Users,
	}
}

func NewGraphQLClient(c *http.Client) *ClientGraphQL {
	client := githubv4.NewClient(c)

	s := &graphqlService{client: client}

	return &ClientGraphQL{
		PR: s,
	}
}

type graphqlService struct {
	client *githubv4.Client
}

func (s *graphqlService) FindIssues(ctx context.Context, owner, repo string, prID, issuesLen int) ([]int, error) {
	variables := map[string]interface{}{
		"issuesLen": githubv4.Int(issuesLen),
		"prID":      githubv4.Int(prID),
		"owner":     githubv4.String(owner),
		"repo":      githubv4.String(repo),
	}

	var q struct {
		Repository struct {
			PullRequest struct {
				ClosingIssuesReferences struct {
					Edges []struct {
						Node struct {
							Number int64
						}
					}
				} `graphql:"closingIssuesReferences (first: $issuesLen)"`
			} `graphql:"pullRequest(number: $prID)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	err := s.client.Query(ctx, &q, variables)
	if err != nil {
		return nil, fmt.Errorf("graphql mutate: %w", err)
	}

	issues := make([]int, len(q.Repository.PullRequest.ClosingIssuesReferences.Edges))

	for i, e := range q.Repository.PullRequest.ClosingIssuesReferences.Edges {
		issues[i] = int(e.Node.Number)
	}

	return issues, nil
}
