// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package github

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/afero"
)

const (
	envAuth = "GITHUB_TOKEN"
)

type envProvider func(env string) string

type AuthToken struct {
	fs       afero.Fs
	location string

	envProvider envProvider
}

func NewAuthToken(fs afero.Fs, tkLoc string) AuthToken {
	return AuthToken{
		fs:          fs,
		location:    tkLoc,
		envProvider: os.Getenv,
	}
}

func NewTestAuthToken(fs afero.Fs, tkLoc string, ep envProvider) AuthToken {
	at := NewAuthToken(fs, tkLoc)
	at.envProvider = ep

	return at
}

// AuthToken method finds and returns the GitHub authorization token.
func (f AuthToken) AuthToken() (string, error) {
	if githubTokenVar := f.envProvider(envAuth); githubTokenVar != "" {
		log.Println("Using GitHub token from environment variable")
		return githubTokenVar, nil
	}

	token, err := afero.ReadFile(f.fs, f.location)
	if err != nil {
		return "", fmt.Errorf("cannot read Github token file failed (path: %s): %w", f.location, err)
	}

	return strings.TrimSpace(string(token)), nil
}
