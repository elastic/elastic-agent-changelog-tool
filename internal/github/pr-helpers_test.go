package github_test

import (
	"testing"

	"github.com/elastic/elastic-agent-changelog-tool/internal/github"
	"github.com/stretchr/testify/require"
)

func TestFillEmptyPRField(t *testing.T) {
	r, hc := getHttpClient(t)
	defer r.Stop() //nolint:errcheck

	c := github.NewClient(hc)

	prIDs, err := github.FillEmptyPRField("fe25c73907336fc462d5e6e059d3cd86512484fe", c)
	require.NoError(t, err)
	require.Len(t, prIDs, 3)
	require.NotEmpty(t, prIDs)
	require.ElementsMatch(t, []int{30979, 31531, 31279}, prIDs)
}

func TestFillEmptyPRFieldBadHash(t *testing.T) {
	r, hc := getHttpClient(t)
	defer r.Stop() //nolint:errcheck

	c := github.NewClient(hc)

	prIDs, err := github.FillEmptyPRField("123", c)
	require.Error(t, err)
	require.Empty(t, prIDs)
}
