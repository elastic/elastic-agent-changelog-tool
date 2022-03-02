// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License
// 2.0; you may not use this file except in compliance with the Elastic License
// 2.0.

package github

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

const (
	envAuth       = "GITHUB_TOKEN"
	authTokenFile = ".elastic/github.token"
)

// EnsureAuthConfigured method ensures that GitHub auth token is available.
func EnsureAuthConfigured() error {
	_, err := AuthToken()
	if err != nil {
		return errors.Wrapf(err, "GitHub authorization token is missing. Please use either environment variable %s or ~/%s",
			envAuth, authTokenFile)
	}
	return nil
}

// AuthToken method finds and returns the GitHub authorization token.
func AuthToken() (string, error) {
	githubTokenVar := os.Getenv(envAuth)
	if githubTokenVar != "" {
		fmt.Println("Using GitHub token from environment variable.")
		return githubTokenVar, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "reading user home directory failed")
	}

	githubTokenPath := filepath.Join(homeDir, authTokenFile)
	token, err := os.ReadFile(githubTokenPath)
	if err != nil {
		return "", errors.Wrapf(err, "reading Github token file failed (path: %s)", githubTokenPath)
	}
	return strings.TrimSpace(string(token)), nil
}
