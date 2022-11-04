// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package assets

import (
	"embed"
	"fmt"
	"strings"
)

// Binds strings to the actual template file name
// These strings can be used in the config template field or renderer template flag
func GetEmbeddedTemplates() embeddedTemplates {
	return map[string]string{
		"asciidoc-embedded": "asciidoc-template.asciidoc",
	}
}

//go:embed asciidoc-template.asciidoc
var AsciidocTemplate embed.FS

type embeddedTemplates map[string]string

func (t embeddedTemplates) String() string {
	var sb strings.Builder

	for k := range t {
		sb.WriteString(fmt.Sprintf("- %s\n", k))
	}

	return sb.String()
}
