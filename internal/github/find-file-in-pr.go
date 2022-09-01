// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package github

import (
	"context"
	"fmt"
	"path/filepath"
)

// FindFilesInPR searches for changes files in a PR that match a given pattern.
// NOTE: it does not check for multiple files matching, a single match is enough.
func FindFileInPR(ctx context.Context, c *Client, owner, repo string, pr int, pattern string) (bool, error) {
	files, resp, err := c.PullRequests.ListFiles(ctx, owner, repo, pr, nil)
	if err != nil {
		return false, fmt.Errorf("cannot list files in pr: %w", err)
	}

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("response not OK while listing files in PR(%s): %+v",
			fmt.Sprintf("%s/%s#%d", owner, repo, pr),
			resp)
	}

	for _, f := range files {
		if f.Filename != nil {
			found, err := filepath.Match(pattern, *f.Filename)
			if err != nil {
				continue
			}

			if found && *f.Status == "removed" {
				continue
			}

			if found {
				return true, nil
			}
		}
	}

	return false, nil
}
