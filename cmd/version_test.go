// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package cmd_test

import (
	"bytes"
	"testing"

	"github.com/elastic/elastic-agent-changelog-tool/cmd"
	"github.com/elastic/elastic-agent-changelog-tool/internal/version"
	"github.com/stretchr/testify/require"
)

func TestVersionCmd_default(t *testing.T) {
	cmd := cmd.VersionCmd()

	b := new(bytes.Buffer)
	cmd.SetOut(b)

	err := cmd.Execute()
	require.Nil(t, err)

	const expected = "elastic-agent-changelog-tool devel version-hash undefined (build time: unknown)"
	require.Equal(t, expected, b.String())
}

func TestVersionCmd_withValues(t *testing.T) {
	cmd := cmd.VersionCmd()

	b := new(bytes.Buffer)
	cmd.SetOut(b)

	version.Tag = "v0.1.0"
	version.BuildTime = "1648570012"
	version.CommitHash = "5561aef"

	err := cmd.Execute()
	require.Nil(t, err)

	const expected = "elastic-agent-changelog-tool v0.1.0 version-hash 5561aef (build time: 2022-03-29T16:06:52Z)"
	require.Equal(t, expected, b.String())
}
