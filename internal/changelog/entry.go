// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package changelog

import (
	"github.com/elastic/elastic-agent-changelog-tool/internal/changelog/fragment"
)

type FragmentFileInfo struct {
	Name     string `yaml:"name"`
	Checksum string `yaml:"checksum"`
}

type Entry struct {
	Kind        Kind     `yaml:"kind"`
	Summary     string   `yaml:"summary"`
	Description string   `yaml:"description"`
	Component   string   `yaml:"component"`
	LinkedPR    []string `yaml:"pr"`
	LinkedIssue []string `yaml:"issue"`
	Repository  string   `yaml:"repository"`

	Timestamp int64            `yaml:"timestamp"`
	File      FragmentFileInfo `yaml:"file"`
}

// EntriesFromFragment returns one or more entries based on the fragment File.
// A single Fragment can contain multiple Changelog entries.
func EntryFromFragment(f fragment.File) Entry {
	e := Entry{
		Kind:        kind2kind(f),
		Summary:     f.Fragment.Summary,
		Description: f.Fragment.Description,
		Component:   f.Fragment.Component,
		LinkedPR:    []string{},
		LinkedIssue: []string{},
		Repository:  f.Fragment.Repository,
		Timestamp:   f.Timestamp,
		File: FragmentFileInfo{
			Name:     f.Name,
			Checksum: f.Checksum(),
		},
	}

	if len(f.Fragment.Pr) > 0 {
		e.LinkedPR = []string{f.Fragment.Pr}
	}

	if len(f.Fragment.Issue) > 0 {
		e.LinkedIssue = []string{f.Fragment.Issue}
	}

	return e
}

func kind2kind(f fragment.File) Kind {
	switch f.Fragment.Kind {
	case string(BreakingChange):
		return BreakingChange
	case string(BugFix):
		return BugFix
	case string(Deprecation):
		return Deprecation
	case string(Enhancement):
		return Enhancement
	case string(Feature):
		return Feature
	case string(KnownIssue):
		return KnownIssue
	case string(Security):
		return Security
	case string(Upgrade):
		return Upgrade
	case string(Other):
		return Other
	default:
		return Unknown
	}
}
