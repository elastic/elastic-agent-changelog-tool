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
