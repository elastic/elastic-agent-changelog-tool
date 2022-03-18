// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package fragment

import (
	"path"
	"reflect"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilename(t *testing.T) {
	fc := TestNewCreator()

	expected := "1647345675-foobar.yaml"
	got := fc.filename("foobar")
	assert.Equal(t, expected, got)
}

func TestSanitizeFilename(t *testing.T) {
	type test struct {
		input string
		want  string
	}

	tests := []test{
		{input: "foo bar", want: "foo-bar"},
		{input: "foo bar foobar", want: "foo-bar-foobar"},
		{input: "foo/bar", want: "foo-bar"},
		{input: "foo\\bar", want: "foo-bar"},
		{input: "foo bar/foobar\\", want: "foo-bar-foobar-"},
	}

	for _, tc := range tests {
		got := sanitizeFilename(tc.input)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}

func TestCreate(t *testing.T) {
	fc := TestNewCreator()
	location := afero.GetTempDir(fc.fs, "testdata")

	fc.Create(location, "foobar")

	content, err := afero.ReadFile(fc.fs, path.Join(location, fc.filename("foobar")))
	require.Nil(t, err)

	expected := `breaking_changes: []
enhancements: []
bugfixes: []
security_fixes: []
known_issues: []
deprecations: []
`
	got := string(content)
	assert.Equal(t, expected, got)
}
