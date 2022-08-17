// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package changelog

import (
	"fmt"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

func FromFile(fs afero.Fs, inFile string) (Changelog, error) {
	b, err := afero.ReadFile(fs, inFile)
	if err != nil {
		return Changelog{}, fmt.Errorf("cannot read file: %s", inFile)
	}

	c := Changelog{}
	if err := yaml.Unmarshal(b, &c); err != nil {
		return Changelog{}, fmt.Errorf("cannot parse changelog YAML from file: %s", inFile)
	}

	return c, nil
}
