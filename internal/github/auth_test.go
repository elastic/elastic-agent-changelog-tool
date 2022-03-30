// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package github

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestAuthToken(t *testing.T) {
	expectedToken := "ghp_tuQprmeVXWdaMhatQiw8pJdEXPxHWm9tkTJb"

	AuthToken := NewAuthToken(&afero.MemMapFs{})

	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	githubTokenPath := filepath.Join(homeDir, authTokenFile)

	err = AuthToken.fs.WriteFile(githubTokenPath, []byte(expectedToken), fs.ModeAppend)
	require.NoError(t, err)

	token, err := AuthToken.AuthToken()
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Equal(t, expectedToken, token)
}

func TestEnsureAuthConfigured(t *testing.T) {
	expectedToken := "ghp_tuQprmeVXWdaMhatQiw8pJdEXPxHWm9tkTJb"

	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	githubTokenPath := filepath.Join(homeDir, authTokenFile)

	AuthToken := NewAuthToken(&afero.MemMapFs{})

	err = AuthToken.fs.WriteFile(githubTokenPath, []byte(expectedToken), fs.ModeAppend)
	require.NoError(t, err)

	err = AuthToken.EnsureAuthConfigured()
	require.NoError(t, err)
}
