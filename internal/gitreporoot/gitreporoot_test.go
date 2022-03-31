// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package gitreporoot_test

import (
	"os"
	"path"
	"testing"

	"github.com/elastic/elastic-agent-changelog-tool/internal/gitreporoot"
	"github.com/stretchr/testify/require"
)

func TestFind(t *testing.T) {
	p, err := gitreporoot.Find()
	require.Nil(t, err)

	// NOTE: git rev-parse --show-toplevel doesn't seem to play out nice with repo in repo structures
	// to avoid making this more complex than needed, the current repo is leveraged, but this test
	// is a bit brittle.
	require.Equal(t, path.Join(os.Getenv("PWD"), "..", ".."), p)
}
