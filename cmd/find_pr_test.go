// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package cmd_test

import (
	"bytes"
	"testing"

	"github.com/elastic/elastic-agent-changelog-tool/cmd"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestFindPrCmd_noArgs(t *testing.T) {
	fs := afero.NewMemMapFs()
	cmd := cmd.FindPRCommand(fs)

	err := cmd.Execute()
	require.Error(t, err)
}

func TestFindPrCmd_oneArg(t *testing.T) {
	fs := afero.NewMemMapFs()
	cmd := cmd.FindPRCommand(fs)

	b := new(bytes.Buffer)
	cmd.SetOut(b)

	cmd.SetArgs([]string{"--repo", "elastic-agent-changelog-tool", "fe56b2e"})

	err := cmd.Execute()
	require.Nil(t, err)
	expected := `{"items":[{"commit":"fe56b2e","pull-request":17}]}
`
	require.Equal(t, expected, b.String())

}
