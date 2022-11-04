// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package changelog_test

import (
	"path"
	"strings"
	"testing"

	"github.com/elastic/elastic-agent-changelog-tool/internal/changelog"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestRenderer(t *testing.T) {
	fs := afero.NewCopyOnWriteFs(afero.NewOsFs(), afero.NewMemMapFs())

	viper.Set("changelog_destination", "testdata")

	filename := "0.0.0.yaml"
	src := "testdata"
	dest := viper.GetString("changelog_destination")

	t.Log("building changelog from test fragments")
	builder := changelog.NewBuilder(fs, filename, "0.0.0", src, dest)
	err := builder.Build("elastic", "elastic-agent-changelog-tool")
	require.NoError(t, err)

	t.Log("loading generated changelog")
	inFile := path.Join(src, filename)
	c, err := changelog.FromFile(fs, inFile)
	require.NoError(t, err)

	r := changelog.NewRenderer(fs, c, dest, "asciidoc-embedded")

	err = r.Render()
	require.Nil(t, err)

	out, err := afero.ReadFile(fs, "testdata/0.0.0.asciidoc")
	require.NoError(t, err)
	t.Log(string(out))
	require.NotEmpty(t, out)

	for _, e := range c.Entries {
		switch e.Kind {
		// NOTE: this is the list of kinds of entries we expect to see
		// in the rendered changelog (not all kinds are expected)
		case changelog.BreakingChange, changelog.Deprecation,
			changelog.BugFix, changelog.Enhancement,
			changelog.Feature, changelog.KnownIssue,
			changelog.Security:
			require.Contains(t, strings.ToLower(string(out)), e.Summary)
		default:
			require.NotContains(t, strings.ToLower(string(out)), e.Summary)
		}
	}
}
