// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package fragment

import (
	"log"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilename(t *testing.T) {
	fc := TestNewCreator()

	expected := "1136239445-foobar.yaml"
	got := fc.filename("foobar")
	assert.Equal(t, expected, got)
}

func TestTimestamp_default(t *testing.T) {
	// NOTE: using sleep to test timestamp default function ability to return different values.
	// Sleeping 1 second to test UNIX timestamp (with second resolution).
	log.Println("SLOW TEST, it sleeps")
	testFs := afero.NewMemMapFs()
	fc := NewCreator(testFs, "foobar")

	t1 := fc.timestamp()
	time.Sleep(1 * time.Second)
	t2 := fc.timestamp()
	time.Sleep(1 * time.Second)
	t3 := fc.timestamp()

	require.Greater(t, t2.Unix(), t1.Unix())
	require.Greater(t, t3.Unix(), t2.Unix())
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

	err := fc.Create("foobar")
	require.Nil(t, err)

	content, err := afero.ReadFile(fc.fs, path.Join(fc.location, fc.filename("foobar")))
	require.Nil(t, err)

	expected := `# Kind can be one of:
# - breaking-change: a change to previously-documented behavior
# - deprecation: functionality that is being removed in a later release
# - bug-fix: fixes a problem in a previous version
# - enhancement: extends functionality but does not break or fix existing behavior
# - feature: new functionality
# - known-issue: problems that we are aware of in a given version
# - security: impacts on the security of a product or a user’s deployment.
# - upgrade: important information for someone upgrading from a prior version
# - other: does not fit into any of the other categories
kind: feature

# Change summary; a 80ish characters long description of the change.
summary: foobar

# Long description; in case the summary is not enough to describe the change
# this field accommodate a description without length limits.
#description:

# Affected component; a word indicating the component this changeset affects.
component:

# PR number; optional; the PR number that added the changeset.
# If not present is automatically filled by the tooling finding the PR where this changelog fragment has been added.
# NOTE: the tooling supports backports, so it's able to fill the original PR number instead of the backport PR number.
# Please provide it if you are adding a fragment for a different PR.
#pr: 1234

# Issue number; optional; the GitHub issue related to this changeset (either closes or is part of).
# If not present is automatically filled by the tooling with the issue linked to the PR number.
#issue: 1234
`
	got := string(content)
	assert.Equal(t, expected, got)
}
