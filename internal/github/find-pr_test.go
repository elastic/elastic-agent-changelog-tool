// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package github_test

import (
	"context"
	"testing"

	"github.com/elastic/elastic-agent-changelog-tool/internal/github"
	"github.com/stretchr/testify/require"
)

func TestFindPR(t *testing.T) {
	r, hc := getHttpClient(t)
	defer r.Stop() //nolint:errcheck

	c := github.NewClient(hc)
	ctx := context.Background()

	// https://github.com/elastic/elastic-agent-changelog-tool/commit/9c995a2e397d346e68ea5052c54bcbd0f8b142ad
	res, err := github.FindPR(ctx, c, "elastic", "elastic-agent-changelog-tool", "9c995a2e397d346e68ea5052c54bcbd0f8b142ad")
	require.NoError(t, err)
	require.Len(t, res.Items, 1)
	require.Equal(t, res.Items[0].PullRequestID, 30)
}

func TestFindPR_missingCommit(t *testing.T) {
	r, hc := getHttpClient(t)
	defer r.Stop() //nolint:errcheck

	c := github.NewClient(hc)
	ctx := context.Background()

	res, err := github.FindPR(ctx, c, "elastic", "elastic-agent-changelog-tool", "does-not-exists")
	require.Error(t, err)
	require.Len(t, res.Items, 0)
}

func TestFindPR_missingRepo(t *testing.T) {
	r, hc := getHttpClient(t)
	defer r.Stop() //nolint:errcheck

	c := github.NewClient(hc)
	ctx := context.Background()

	res, err := github.FindPR(ctx, c, "elastic", "does-not-exists", "9c995a2e397d346e68ea5052c54bcbd0f8b142ad")
	require.Error(t, err)
	require.Len(t, res.Items, 0)
}
