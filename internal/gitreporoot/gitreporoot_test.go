package gitreporoot_test

import (
	"os"
	"path"
	"testing"

	"github.com/elastic/elastic-agent-changelog-tool/internal/gitreporoot"
	"github.com/stretchr/testify/require"
)

func TestFind(t *testing.T) {
	p, err := gitreporoot.Find()
	require.Nil(t, err)

	// NOTE: git rev-parse --show-toplevel doesn't seem to play out nice with repo in repo structures
	// to avoid making this more complex than needed, the current repo is leveraged, but this test
	// is a bit brittle.
	require.Equal(t, path.Join(os.Getenv("PWD"), "..", ".."), p)
}
