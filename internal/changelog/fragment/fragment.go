// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package fragment

type Fragment struct {
	BreakingChanges []string `yaml:"breaking_changes"`
	Deprecations    []string `yaml:"deprecations"`
	Bugfixes        []string `yaml:"bugfixes"`
	Enhancements    []string `yaml:"enhancements"`
	Features        []string `yaml:"features"`
	KnownIssues     []string `yaml:"known_issues"`
	Security        []string `yaml:"security"`
	Upgrades        []string `yaml:"upgrades"`
	Other           []string `yaml:"other"`
}
