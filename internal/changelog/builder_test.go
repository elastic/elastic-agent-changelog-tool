// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package changelog_test

import (
	"testing"

	"github.com/elastic/elastic-agent-changelog-tool/internal/changelog"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestBuilder_Build(t *testing.T) {
	fs := afero.NewCopyOnWriteFs(afero.NewOsFs(), afero.NewMemMapFs())
	viper.Set("fragment_location", "testdata")

	b := changelog.NewBuilder(fs, "filename", "0.0.0", "testdata", "testdata")

	err := b.Build()
	require.NoError(t, err)

	// FIXME: built changelog is not inspectable as b.changelog is not updated &
	// there is no way to access it anyway
	// fmt.Println(b.Changelog())
	// require.Len(t, b.Changelog().Entries, 10)
}
