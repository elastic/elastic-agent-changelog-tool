// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package cmd_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/elastic/elastic-agent-changelog-tool/cmd"
	"github.com/elastic/elastic-agent-changelog-tool/internal/changelog"
	"github.com/elastic/elastic-agent-changelog-tool/internal/changelog/fragment"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestBuildCmd_default(t *testing.T) {
	testFs := afero.NewMemMapFs()
	c := fragment.NewCreator(testFs, viper.GetString("fragment_location"))
	err := c.Create("foo")
	require.Nil(t, err)
	// NOTE: sleeping to produce different fragment's timestamps
	time.Sleep(1 * time.Second)
	err = c.Create("bar")
	require.Nil(t, err)

	cmd := cmd.BuildCmd(testFs)

	b := new(bytes.Buffer)
	cmd.SetOut(b)

	err = cmd.Execute()
	require.Nil(t, err)

	content, err := afero.ReadFile(testFs, viper.GetString("changelog_destination"))
	require.Nil(t, err)

	ch := changelog.Changelog{}
	err = yaml.Unmarshal(content, &ch)
	require.Nil(t, err)

	fmt.Println(ch)

	require.Equal(t, "8.2.1", ch.Version)
	require.Len(t, ch.Entries, 2)
}
