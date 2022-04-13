// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package github

import (
	"os"
	"path/filepath"
)

const (
	authTokenLocation = ".elastic"
	authTokenFilename = "github.token"
)

// TokenLocation returns the expected location of the GitHub Token file.
func TokenLocation() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, authTokenLocation, authTokenFilename), nil
}
