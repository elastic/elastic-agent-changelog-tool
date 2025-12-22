// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package cmd_test

import (
	"bytes"
	"io"
	"log"
	"testing"

	"github.com/elastic/elastic-agent-changelog-tool/cmd"
	"github.com/elastic/elastic-agent-changelog-tool/internal/settings"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestPrHasFragmentCmd_noArgs(t *testing.T) {
	log.SetOutput(io.Discard)

	fs := afero.NewMemMapFs()
	c := cmd.PrHasFragmentCommand(fs)

	c.SetOut(io.Discard)
	c.SetErr(io.Discard)

	err := c.Execute()
	require.Error(t, err)
}

func TestPrHasFragmentCmd_oneArg(t *testing.T) {
	log.SetOutput(io.Discard)

	settings.Init()

	fs := afero.NewMemMapFs()
	c := cmd.PrHasFragmentCommand(fs)

	b := new(bytes.Buffer)
	c.SetOut(b)
	c.SetErr(io.Discard)

	c.SetArgs([]string{"--repo", "elastic-agent-changelog-tool", "145"})

	err := c.Execute()
	require.Nil(t, err)
}

// Fail because no fragment present
func TestPrHasFragmentCmd_oneArgFailCaseNoPresent(t *testing.T) {
	log.SetOutput(io.Discard)

	settings.Init()

	fs := afero.NewMemMapFs()
	c := cmd.PrHasFragmentCommand(fs)

	b := new(bytes.Buffer)
	c.SetOut(b)
	c.SetErr(io.Discard)

	c.SetArgs([]string{"--repo", "elastic-agent-changelog-tool", "33"})

	err := c.Execute()
	require.Error(t, err)
}

// Fail because fragment missing required field
func TestPrHasFragmentCmd_oneArgFailCaseField(t *testing.T) {
	log.SetOutput(io.Discard)

	settings.Init()

	fs := afero.NewMemMapFs()
	c := cmd.PrHasFragmentCommand(fs)

	b := new(bytes.Buffer)
	c.SetOut(b)
	c.SetErr(io.Discard)

	c.SetArgs([]string{"--repo", "elastic-agent-changelog-tool", "29"})

	err := c.Execute()
	require.Error(t, err)
}

func TestPrHasFragmentCmd_hasSkipChangelogLabel(t *testing.T) {
	settings.Init()

	fs := afero.NewMemMapFs()
	c := cmd.PrHasFragmentCommand(fs)

	b := new(bytes.Buffer)
	c.SetOut(b)
	c.SetErr(io.Discard)

	c.SetArgs([]string{"--repo", "elastic-agent-changelog-tool", "47"})

	err := c.Execute()
	require.Nil(t, err)
	require.Contains(t, b.String(), "PR requires no changelog")
}
