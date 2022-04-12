package github_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/elastic/elastic-agent-changelog-tool/internal/github"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestTokenLocation(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	tkloc, err := github.TokenLocation()
	require.NoError(t, err, "cannot get token location")
	require.Equal(t, filepath.Join(home, ".elastic", "github.token"), tkloc)
}

func TestEnsureAuthConfigured(t *testing.T) {
	expectedToken := "ghp_tuQprmeVXWdaMhatQiw8pJdEXPxHWm9tkTJb"
	testFs := afero.MemMapFs{}

	tokenLocation, err := github.TokenLocation()
	require.NoError(t, err, "cannot get token location")
	afero.WriteFile(&testFs, tokenLocation, []byte(expectedToken), fs.ModeAppend)

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
	afero.WriteFile(&testFs, tokenLocation, []byte(expectedToken), fs.ModeAppend)

	tk := github.NewAuthToken(&testFs, tokenLocation)

	ok, err := github.EnsureAuthConfigured(tk)
	require.Error(t, err)
	require.False(t, ok)
}
