package github_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/elastic/elastic-agent-changelog-tool/internal/github"
	"github.com/stretchr/testify/require"
)

func TestTokenLocation(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	tkloc, err := github.TokenLocation()
	require.NoError(t, err, "cannot get token location")
	require.Equal(t, filepath.Join(home, ".elastic", "github.token"), tkloc)
}
