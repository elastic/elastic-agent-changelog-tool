package github

import (
	"errors"
	"fmt"
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

// EnsureAuthConfigured method ensures that GitHub auth token is available.
func EnsureAuthConfigured(tk AuthToken) (bool, error) {
	tkloc, err := TokenLocation()
	if err != nil {
		return false, fmt.Errorf("cannot determine token location: %w", err)
	}

	val, err := tk.AuthToken()
	if err != nil {
		return false, fmt.Errorf("GitHub authorization token is missing. Please use either environment variable %s or ~/%s: %w",
			envAuth, tkloc, err)
	}

	if val == "" {
		return false, errors.New("GitHub authorization token is empty. Make sure a value is provided")
	}

	return true, nil
}
