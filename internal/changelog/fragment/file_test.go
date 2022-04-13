// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package fragment_test

import (
	"path"
	"testing"

	"github.com/elastic/elastic-agent-changelog-tool/internal/changelog/fragment"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestFile(t *testing.T) {
	fs := afero.NewCopyOnWriteFs(afero.NewOsFs(), afero.NewMemMapFs())

	viper.Set("fragment_location", "testdata")

	f1, err := fragment.Load(fs, path.Join(viper.GetString("fragment_location"), "1648040924-breaking-change.yaml"))
	require.Nil(t, err)

	require.Equal(t, f1.Title, "breaking-change")
	require.Equal(t, f1.Timestamp, int64(1648040924))
	require.Equal(t, f1.Fragment.Kind, "breaking-change")
	require.Equal(t, f1.Fragment.Summary, "a change with breaking changes")

	f2, err := fragment.Load(fs, path.Join(viper.GetString("fragment_location"), "1648040928-enhancement.yaml"))
	require.Nil(t, err)

	require.Equal(t, f2.Title, "enhancement")
	require.Equal(t, f2.Timestamp, int64(1648040928))
	require.Equal(t, f2.Fragment.Kind, "enhancement")
	require.Equal(t, f2.Fragment.Summary, "a new feature")

}

func TestFileChecksum(t *testing.T) {
	fs := afero.NewCopyOnWriteFs(afero.NewOsFs(), afero.NewMemMapFs())

	viper.Set("fragment_location", "testdata")

	f, err := fragment.Load(fs, path.Join(viper.GetString("fragment_location"), "1648040924-breaking-change.yaml"))
	require.Nil(t, err)

	require.Equal(t, f.Checksum(), "f085dd487a7c88d165235546cdb528cc4e94a046")
}
