// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package fragment

import (
	"bytes"
	"embed"
	"fmt"
	txttempl "text/template"
)

//go:embed template.yaml
var template embed.FS

func Template(slug string) ([]byte, error) {
	data, err := template.ReadFile("template.yaml")
	if err != nil {
		return nil, fmt.Errorf("cannot read embedded template: %w", err)
	}

	tmpl, err := txttempl.New("template").Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("cannot parse template: %w", err)
	}

	vars := make(map[string]interface{})
	vars["Summary"] = slug

	buf := bytes.NewBuffer(nil)

	err = tmpl.Execute(buf, vars)
	if err != nil {
		return nil, fmt.Errorf("cannot execute template: %w", err)
	}

	return buf.Bytes(), nil
}
