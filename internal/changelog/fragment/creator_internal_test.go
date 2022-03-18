// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package fragment

import (
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

	fc.Create("foobar")

	content, err := afero.ReadFile(fc.fs, fc.filename("foobar"))
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
