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

	// This is a PR merged on trunk:
	// https://github.com/elastic/beats/pull/30859
	// https://github.com/elastic/beats/commit/56df883ca93b11816206dad401c49a2c96fa268d
	res, err := github.FindPR(ctx, c, "elastic", "beats", "56df883ca93b11816206dad401c49a2c96fa268d")
	require.NoError(t, err)
	require.Len(t, res.Items, 1)
	require.Equal(t, 30859, res.Items[0].PullRequestID)
}

func TestFindPR_backport(t *testing.T) {
	r, hc := getHttpClient(t)
	defer r.Stop() //nolint:errcheck

	c := github.NewClient(hc)
	ctx := context.Background()

	// This is a merge commit of a backport from main to 8.1:
	// from https://github.com/elastic/beats/pull/31279 to https://github.com/elastic/beats/pull/31343
	// https://github.com/elastic/beats/commit/fe25c73907336fc462d5e6e059d3cd86512484fe
	res, err := github.FindPR(ctx, c, "elastic", "beats", "fe25c73907336fc462d5e6e059d3cd86512484fe")
	require.NoError(t, err)
	require.Len(t, res.Items, 4)
	// not a backport: https://github.com/elastic/beats/pull/31396
	require.Equal(t, 31396, res.Items[0].PullRequestID)
	// backport: https://github.com/elastic/beats/pull/31417 => source: https://github.com/elastic/beats/issues/31013
	require.Equal(t, 31013, res.Items[1].PullRequestID)
	// backport: https://github.com/elastic/beats/pull/31396 => source: https://github.com/elastic/beats/issues/31369
	require.Equal(t, 31369, res.Items[2].PullRequestID)
	// backport: https://github.com/elastic/beats/pull/31343 => source: https://github.com/elastic/beats/pull/31279
	require.Equal(t, 31279, res.Items[3].PullRequestID)
}

func TestFindPR_forwardport(t *testing.T) {
	r, hc := getHttpClient(t)
	defer r.Stop() //nolint:errcheck

	c := github.NewClient(hc)
	ctx := context.Background()

	// This is a merge commit of a forwardport from 8.0 to main:
	// from https://github.com/elastic/beats/issues/29209 to https://github.com/elastic/beats/pull/30626
	// https://github.com/elastic/beats/commit/8800e5f6ad5beb024dee141a2639630b79a99a37
	res, err := github.FindPR(ctx, c, "elastic", "beats", "8800e5f6ad5beb024dee141a2639630b79a99a37")
	require.NoError(t, err)
	require.Len(t, res.Items, 1)
	require.Equal(t, 29209, res.Items[0].PullRequestID)
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
