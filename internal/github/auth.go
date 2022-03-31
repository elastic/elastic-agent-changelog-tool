// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package github

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

const (
	envAuth       = "GITHUB_TOKEN"
	authTokenFile = ".elastic/github.token"
)

// EnsureAuthConfigured method ensures that GitHub auth token is available.
func (f *AuthToken) EnsureAuthConfigured() error {
	if _, err := f.AuthToken(); err != nil {
		return fmt.Errorf("GitHub authorization token is missing. Please use either environment variable %s or ~/%s: %w",
			envAuth, authTokenFile, err)
	}

	return nil
}

type AuthToken struct {
	fs *afero.Afero
}

func NewAuthToken(fs afero.Fs) *AuthToken {
	return &AuthToken{
		fs: &afero.Afero{
			Fs: fs,
		},
	}
}

// AuthToken method finds and returns the GitHub authorization token.
func (f *AuthToken) AuthToken() (string, error) {
	githubTokenVar := os.Getenv(envAuth)

	if githubTokenVar != "" {
		fmt.Println("Using GitHub token from environment variable.")
		return githubTokenVar, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("reading user home directory failed: %w", err)
	}

	githubTokenPath := filepath.Join(homeDir, authTokenFile)

	token, err := f.fs.ReadFile(githubTokenPath)
	if err != nil {
		return "", fmt.Errorf("reading Github token file failed (path: %s): %w", githubTokenPath, err)
	}

	return strings.TrimSpace(string(token)), nil
}
