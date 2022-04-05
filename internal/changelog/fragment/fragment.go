// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package fragment

type Fragment struct {
	Kind        string `yaml:"kind"`
	Summary     string `yaml:"summary"`
	Description string `yaml:"description"`
	Component   string `yaml:"component"`
	Pr          int    `yaml:"pr"`
	Issue       int    `yaml:"issue"`
	Repository  string `yaml:"repository"`
}
