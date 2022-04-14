// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package github_test

import (
	"io/fs"
	"testing"

	"github.com/elastic/elastic-agent-changelog-tool/internal/github"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestEnsureAuthConfigured(t *testing.T) {
	expectedToken := "ghp_tuQprmeVXWdaMhatQiw8pJdEXPxHWm9tkTJb"
	testFs := afero.MemMapFs{}

	tokenLocation, err := github.TokenLocation()
	require.NoError(t, err, "cannot get token location")
	err = afero.WriteFile(&testFs, tokenLocation, []byte(expectedToken), fs.ModeAppend)
	require.NoError(t, err)

	tk := github.NewAuthToken(&testFs, tokenLocation)

	ok, err := github.EnsureAuthConfigured(tk)
	require.NoError(t, err)
	require.True(t, ok)
}

func TestEnsureAuthConfigured_failure(t *testing.T) {
	testFs := afero.MemMapFs{}

	tokenLocation, err := github.TokenLocation()
	require.NoError(t, err, "cannot get token location")

	tk := github.NewAuthToken(&testFs, tokenLocation)

	ok, err := github.EnsureAuthConfigured(tk)
	require.Error(t, err)
	require.False(t, ok)
}

func TestEnsureAuthConfigured_empty(t *testing.T) {
	expectedToken := ""
	testFs := afero.MemMapFs{}

	tokenLocation, err := github.TokenLocation()
	require.NoError(t, err, "cannot get token location")
	err = afero.WriteFile(&testFs, tokenLocation, []byte(expectedToken), fs.ModeAppend)
	require.NoError(t, err)

	tk := github.NewAuthToken(&testFs, tokenLocation)

	ok, err := github.EnsureAuthConfigured(tk)
	require.Error(t, err)
	require.False(t, ok)
}
