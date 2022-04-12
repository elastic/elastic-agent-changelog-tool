package github

import (
	"context"
	"fmt"
)

type PRForCommit struct {
	CommitHash    string `json:"commit"`
	PullRequestID int    `json:"pull-request"`
}

type FoundPRs struct {
	Items []PRForCommit `json:"items"`
}

func FindPR(ctx context.Context, c *Client, owner, repo, commit string) (FoundPRs, error) {
	prs, _, err := c.PullRequests.ListPullRequestsWithCommit(
		context.Background(), owner, repo, commit, nil)
	if err != nil {
		return FoundPRs{}, fmt.Errorf("failed listing prs with commit: %w", err)
	}

	respData := FoundPRs{
		Items: make([]PRForCommit, len(prs)),
	}

	for i, pr := range prs {
		respData.Items[i] = PRForCommit{
			CommitHash:    commit,
			PullRequestID: pr.GetNumber(),
		}
	}

	return respData, nil
}
