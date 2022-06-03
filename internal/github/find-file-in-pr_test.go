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

func TestFindFileInPR(t *testing.T) {
	r, hc := getHttpClient(t)
	defer r.Stop() //nolint:errcheck

	c := github.NewClient(hc)
	ctx := context.Background()

	// This is a PR merged with changelog fragments:
	// https://github.com/elastic/elastic-agent-changelog-tool/pull/29
	res, err := github.FindFileInPR(ctx, c, "elastic", "elastic-agent-changelog-tool", 29, "changelog/fragments/*")
	require.NoError(t, err)
	require.True(t, res)
}

func TestFindFileInPR_failureCase(t *testing.T) {
	r, hc := getHttpClient(t)
	defer r.Stop() //nolint:errcheck

	c := github.NewClient(hc)
	ctx := context.Background()

	// This is a PR without changelog fragments:
	// https://github.com/elastic/elastic-agent-changelog-tool/pull/33
	res, err := github.FindFileInPR(ctx, c, "elastic", "elastic-agent-changelog-tool", 33, "changelog/fragments/*")
	require.NoError(t, err)
	require.False(t, res)
}

func TestFindFileInPR_PrIsNotFound(t *testing.T) {
	r, hc := getHttpClient(t)
	defer r.Stop() //nolint:errcheck

	c := github.NewClient(hc)
	ctx := context.Background()

	// This is a purposedly a PR that do not exists
	res, err := github.FindFileInPR(ctx, c, "elastic", "elastic-agent-changelog-tool", -1, "changelog/fragments/*")
	require.Error(t, err)
	require.False(t, res)
}
