// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package fragment

import (
	"embed"
	"fmt"
)

//go:embed template.yaml
var template embed.FS

func Template() ([]byte, error) {
	data, err := template.ReadFile("template.yaml")
	if err != nil {
		return []byte{}, fmt.Errorf("cannot read embedded template: %w", err)
	}

	return data, nil
}
