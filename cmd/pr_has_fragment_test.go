// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package cmd_test

import (
	"bytes"
	"io/ioutil"
	"log"
	"testing"

	"github.com/elastic/elastic-agent-changelog-tool/cmd"
	"github.com/elastic/elastic-agent-changelog-tool/internal/settings"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestPrHasFragmentCmd_noArgs(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	fs := afero.NewMemMapFs()
	cmd := cmd.PrHasFragmentCommand(fs)

	cmd.SetOut(ioutil.Discard)
	cmd.SetErr(ioutil.Discard)

	err := cmd.Execute()
	require.Error(t, err)
}

func TestPrHasFragmentCmd_oneArg(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	settings.Init()

	fs := afero.NewMemMapFs()
	cmd := cmd.PrHasFragmentCommand(fs)

	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(ioutil.Discard)

	cmd.SetArgs([]string{"--repo", "elastic-agent-changelog-tool", "29"})

	err := cmd.Execute()
	require.Nil(t, err)
}

func TestPrHasFragmentCmd_oneArgFailCase(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	settings.Init()

	fs := afero.NewMemMapFs()
	cmd := cmd.PrHasFragmentCommand(fs)

	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(ioutil.Discard)

	cmd.SetArgs([]string{"--repo", "elastic-agent-changelog-tool", "33"})

	err := cmd.Execute()
	require.Error(t, err)
}
