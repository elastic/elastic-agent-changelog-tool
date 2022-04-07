// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package settings_test

import (
	"testing"

	"github.com/OpenPeeDeeP/xdg"
	"github.com/stretchr/testify/assert"

	"github.com/elastic/elastic-agent-changelog-tool/internal/settings"
)

func TestCacheDir(t *testing.T) {
	settings.Init()

	expected := xdg.CacheHome()
	got := settings.CacheDir()

	assert.Equal(t, expected, got)
}

func TestConfigDir(t *testing.T) {
	settings.Init()

	expected := xdg.ConfigHome()
	got := settings.ConfigDir()

	assert.Equal(t, expected, got)
}

func TestDataDir(t *testing.T) {
	settings.Init()

	expected := xdg.DataHome()
	got := settings.DataDir()

	assert.Equal(t, expected, got)
}
