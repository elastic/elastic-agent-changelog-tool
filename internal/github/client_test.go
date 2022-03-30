// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package github

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/elastic/elastic-agent-changelog-tool/internal/mocks"
	"github.com/google/go-github/v32/github"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestNewClient(t *testing.T) {
	expectedToken := "ghp_tuQprmeVXWdaMhatQiw8pJdEXPxHWm9tkTJb"

	AuthToken := NewAuthToken(&afero.MemMapFs{})

	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	githubTokenPath := filepath.Join(homeDir, authTokenFile)

	err = AuthToken.fs.WriteFile(githubTokenPath, []byte(expectedToken), fs.ModeAppend)
	require.NoError(t, err)

	bb, err := AuthToken.AuthToken()
	require.Nil(t, err)

	githubClient, err := NewClient(NewWrapper(github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: bb}),
	))))
	require.NotNil(t, githubClient)
	require.NoError(t, err)
}

func TestNewUnauthorizedClient(t *testing.T) {
	unauthorizedGithubClient := NewUnauthorizedClient()
	require.NotNil(t, unauthorizedGithubClient)
}

func TestGithubUsersWithValidToken(t *testing.T) {
	expectedToken := "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"

	AuthToken := NewAuthToken(&afero.MemMapFs{})

	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	githubTokenPath := filepath.Join(homeDir, authTokenFile)

	err = AuthToken.fs.WriteFile(githubTokenPath, []byte(expectedToken), fs.ModeAppend)
	require.NoError(t, err)

	mockedGithubClient := new(mocks.GithubClient)

	githubClient, err := NewClient(mockedGithubClient)
	require.NotNil(t, githubClient)
	require.NoError(t, err)
	require.NotNil(t, githubClient.client)

	githubLogin := "user1"

	githubUserResponse := &github.User{
		Login: &githubLogin,
	}

	githubResponse := &github.Response{}

	mockedGithubClient.On("UsersGet", mock.Anything, "").Return(githubUserResponse, githubResponse, nil).Once()

	login, err := githubClient.User()
	require.NoError(t, err)
	require.NotEmpty(t, login)
	require.Equal(t, login, githubLogin)
}
