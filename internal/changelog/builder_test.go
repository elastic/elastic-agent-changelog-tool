// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package changelog_test

import (
	"testing"

	"github.com/elastic/elastic-agent-changelog-tool/internal/changelog"
	"github.com/elastic/elastic-agent-changelog-tool/internal/github"
	"github.com/elastic/elastic-agent-changelog-tool/internal/githubtest"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestBuilder_Build(t *testing.T) {
	fs := afero.NewCopyOnWriteFs(afero.NewOsFs(), afero.NewMemMapFs())
	viper.Set("fragment_location", "testdata")

	b := changelog.NewBuilder(fs, "filename", "0.0.0", "testdata", "testdata")

	err := b.Build("elastic", "beats")
	require.NoError(t, err)

	// FIXME: built changelog is not inspectable as b.changelog is not updated &
	// there is no way to access it anyway
	// fmt.Println(b.Changelog())
	// require.Len(t, b.Changelog().Entries, 10)
}

func TestFillEmptyPRField(t *testing.T) {
	r, hc := githubtest.GetHttpClient(t)
	defer r.Stop() //nolint:errcheck

	c := github.NewClient(hc)

	prIDs, err := changelog.FillEmptyPRField("fe25c73907336fc462d5e6e059d3cd86512484fe", "elastic", "beats", c)
	require.NoError(t, err)
	require.Len(t, prIDs, 2)
	require.NotEmpty(t, prIDs)
	require.ElementsMatch(t, []int{30979, 31279}, prIDs)
}

func TestFillEmptyPRFieldBadHash(t *testing.T) {
	r, hc := githubtest.GetHttpClient(t)
	defer r.Stop() //nolint:errcheck

	c := github.NewClient(hc)

	prIDs, err := changelog.FillEmptyPRField("123", "elastic", "beats", c)
	require.Error(t, err)
	require.Empty(t, prIDs)
}
