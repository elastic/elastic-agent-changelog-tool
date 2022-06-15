// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package changelog

import (
	"log"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_collectFragment(t *testing.T) {
	testFs := afero.NewCopyOnWriteFs(afero.NewOsFs(), afero.NewMemMapFs())

	var files []string
	err := afero.Walk(testFs, "testdata", func(path string, info os.FileInfo, err error) error {
		return collectFragment(testFs, path, info, err, &files)
	})

	require.NoError(t, err)

	log.Println(files)
	require.NotContains(t, files, "testdata/.gitkeep")
	require.Contains(t, files, "testdata/1648040928-breaking-change.yaml")
	require.Contains(t, files, "testdata/1648040928-bug-fix.yaml")
	require.Contains(t, files, "testdata/1648040928-deprecation.yaml")
	require.Contains(t, files, "testdata/1648040928-enhancement.yaml")
	require.Contains(t, files, "testdata/1648040928-feature.yaml")
	require.Contains(t, files, "testdata/1648040928-known-issue.yaml")
	require.Contains(t, files, "testdata/1648040928-other.yaml")
	require.Contains(t, files, "testdata/1648040928-security.yaml")
	require.Contains(t, files, "testdata/1648040928-unknown.yaml")
	require.Contains(t, files, "testdata/1648040928-upgrade.yaml")
}
