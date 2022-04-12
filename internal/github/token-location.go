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
