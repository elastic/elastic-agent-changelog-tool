// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package fragment

type Fragment struct {
	Major         []string `yaml:"breaking_changes"`
	Minor         []string `yaml:"enhancements"`
	Patch         []string `yaml:"bugfixes"`
	SecurityFixes []string `yaml:"security_fixes"`
	KnownIssues   []string `yaml:"known_issues"`
	Deprecations  []string `yaml:"deprecations"`
}
