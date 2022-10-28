// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package changelog

import (
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func NewLinter(fs afero.Fs) Linter {
	return newLinter(fs)
}

type Linter struct {
	fs         afero.Fs
	validators map[string]func(entry Entry) error
	errors     []error
}

func newLinter(fs afero.Fs) Linter {
	return Linter{
		fs: fs,
		validators: map[string]func(entry Entry) error{
			"multiplePRIDs": func(entry Entry) error {
				if len(entry.LinkedPR) > 1 {
					return fmt.Errorf("changelog entry: %s has multiple PR ids", entry.File.Name)
				}

				return nil
			},
			"noPRIDs": func(entry Entry) error {
				if len(entry.LinkedPR) == 0 {
					return fmt.Errorf("changelog entry: %s has no PR id", entry.File.Name)
				}

				return nil
			},
			"noIssueID": func(entry Entry) error {
				if len(entry.LinkedIssue) == 0 {
					return fmt.Errorf("changelog entry: %s has no issue id", entry.File.Name)
				}

				return nil
			},
			// multiple issues check?
			"validComponent": func(entry Entry) error {
				configComponents := viper.GetStringSlice("components")

				switch len(configComponents) {
				case 0:
					return nil
				case 1:
					c := configComponents[0]

					if c != entry.Component && len(entry.Component) > 0 {
						return fmt.Errorf("Component [%s] not found in config", entry.Component)
					}
				default:
					var match string

					if entry.Component == "" {
						return fmt.Errorf("Component cannot be assumed, choose it from config values: %s", entry.File.Name)
					}

					match = ""
					for _, c := range configComponents {
						if entry.Component != c {
							continue
						}
						match = entry.Component
					}

					if match == "" {
						return fmt.Errorf("Component [%s] not found in config", entry.Component)
					}
				}

				return nil
			},
		},
	}
}

type LinterErrors []error

func (l Linter) Lint(dest, version string) []error {
	c, err := FromFile(l.fs, fmt.Sprintf("./%s/%s.yaml", dest, version))
	if err != nil {
		return []error{fmt.Errorf("error loading changelog from file: %w", err)}
	}

	for _, entry := range c.Entries {
		for _, validator := range l.validators {
			err := validator(entry)
			if err != nil {
				l.errors = append(l.errors, err)
			}
		}
	}

	return l.errors
}
