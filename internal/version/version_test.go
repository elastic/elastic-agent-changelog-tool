// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package version_test

import (
	"testing"

	"github.com/elastic/elastic-agent-changelog-tool/internal/version"
	"github.com/stretchr/testify/require"
)

func TestCommitHashDefault(t *testing.T) {
	require.Equal(t, "undefined", version.CommitHash)
}

func TestBuildTimeFormattedDefault(t *testing.T) {
	// NOTE: this test is order sensitive, as it tests the default value
	v := version.BuildTimeFormatted()
	require.Equal(t, "unknown", v)
}

func TestBuildTimeFormatted_invalid(t *testing.T) {
	version.BuildTime = "foobar"
	v := version.BuildTimeFormatted()
	require.Equal(t, "invalid", v)
	// NOTE: reset value to default to avoid test order issues
	version.BuildTime = "unknown"
}

func TestBuildTimeFormatted_valid(t *testing.T) {
	version.BuildTime = "1648570012"
	v := version.BuildTimeFormatted()
	require.Equal(t, "2022-03-29T16:06:52Z", v)
	// NOTE: reset value to default to avoid test order issues
	version.BuildTime = "unknown"
}

func TestTagDefault(t *testing.T) {
	require.Empty(t, version.Tag)
}
