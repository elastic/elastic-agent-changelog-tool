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

func TestAuthToken(t *testing.T) {
	expectedToken := "ghp_tuQprmeVXWdaMhatQiw8pJdEXPxHWm9tkTJb"
	testFs := afero.MemMapFs{}

	tokenLocation, err := github.TokenLocation()
	require.NoError(t, err, "cannot get token location")
	err = afero.WriteFile(&testFs, tokenLocation, []byte(expectedToken), fs.ModeAppend)
	require.NoError(t, err)

	tk := github.NewAuthToken(&testFs, tokenLocation)

	token, err := tk.AuthToken()
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Equal(t, expectedToken, token)
}

func TestAuthToken_fromEnv(t *testing.T) {
	expectedToken := "ghp_tuQprmeVXWdaMhatQiw8pJdEXPxHWm9tkTJb"
	testFs := afero.MemMapFs{}

	tk := github.NewTestAuthToken(&testFs, "using env", func(key string) string {
		return expectedToken
	})

	token, err := tk.AuthToken()
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Equal(t, expectedToken, token)
}

func TestAuthToken_readFailure(t *testing.T) {
	testFs := afero.MemMapFs{}

	tk := github.NewAuthToken(&testFs, "foobar")

	token, err := tk.AuthToken()
	require.Error(t, err)
	require.Empty(t, token)
}
