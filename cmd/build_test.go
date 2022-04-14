// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package cmd_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/elastic/elastic-agent-changelog-tool/cmd"
	"github.com/elastic/elastic-agent-changelog-tool/internal/changelog"
	"github.com/elastic/elastic-agent-changelog-tool/internal/changelog/fragment"
	"github.com/elastic/elastic-agent-changelog-tool/internal/settings"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestBuildCmd(t *testing.T) {
	settings.Init()

	testFs := afero.NewMemMapFs()
	c := fragment.NewCreator(testFs, viper.GetString("fragment_location"))
	err := c.Create("foo")
	require.Nil(t, err)
	// NOTE: sleeping to produce different fragment's timestamps
	time.Sleep(1 * time.Second)
	err = c.Create("bar")
	require.Nil(t, err)

	cmd := cmd.BuildCmd(testFs)

	expectedVersion := "0.0.0"
	cmd.SetArgs([]string{
		fmt.Sprintf("--version=%s", expectedVersion),
	})

	b := new(bytes.Buffer)
	cmd.SetOut(b)

	err = cmd.Execute()
	require.Nil(t, err)

	content, err := afero.ReadFile(testFs, viper.GetString("changelog_destination"))
	require.Nil(t, err)

	ch := changelog.Changelog{}
	err = yaml.Unmarshal(content, &ch)
	require.Nil(t, err)

	require.Equal(t, expectedVersion, ch.Version)
	require.Len(t, ch.Entries, 2)
}

func TestBuildCmd_missingFlag(t *testing.T) {
	testFs := afero.NewMemMapFs()
	cmd := cmd.BuildCmd(testFs)

	cmd.SetOut(ioutil.Discard)
	cmd.SetErr(ioutil.Discard)

	err := cmd.Execute()
	expectedErr := errors.New("required flag(s) \"version\" not set")
	require.Error(t, expectedErr, err)
}
