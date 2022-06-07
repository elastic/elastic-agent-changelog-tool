package github

import (
	"context"
	"fmt"
)

// PRHasLabels returns if a PR (given by it's number) has specified labels or not
func PRHasLabels(ctx context.Context, c *Client, owner, repo string, pr int, labels []string) (bool, error) {
	prData, resp, err := c.PullRequests.Get(ctx, owner, repo, pr)
	if err != nil {
		return false, fmt.Errorf("cannot retrieve PR information: %w", err)
	}

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("response HTTP status code is not 200 while retrieving PR information: actual status code %d", resp.StatusCode)
	}

	for _, l := range labels {
		// TODO: replace this with slice.Contains in go1.18
		for _, prl := range prData.Labels {
			if l == *prl.Name {
				return true, nil
			}
		}
	}

	return false, nil
}
