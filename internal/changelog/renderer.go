// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package changelog

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"path"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Renderer struct {
	changelog Changelog
	fs        afero.Fs
	// dest is the destination location where the changelog is written to
	dest  string
	templ string
}

func NewRenderer(fs afero.Fs, c Changelog, dest string, templ string) *Renderer {
	return &Renderer{
		changelog: c,
		fs:        fs,
		dest:      dest,
		templ:     templ,
	}
}

func (r Renderer) Render() error {
	log.Printf("render changelog for version: %s\n", r.changelog.Version)

	tpl, err := r.Template()
	if err != nil {
		log.Fatal(err)
	}

	type TemplateData struct {
		Component string
		Version   string
		Changelog Changelog
		Kinds     map[Kind]bool

		BreakingChange map[string][]Entry
		Deprecation    map[string][]Entry
		BugFix         map[string][]Entry
		Enhancement    map[string][]Entry
		Feature        map[string][]Entry
		KnownIssue     map[string][]Entry
		Security       map[string][]Entry
		Upgrade        map[string][]Entry
		Other          map[string][]Entry
	}

	td := TemplateData{
		buildTitleByComponents(r.changelog.Entries), r.changelog.Version, r.changelog,
		collectKinds(r.changelog.Entries),
		collectByKindMap(r.changelog.Entries, BreakingChange),
		collectByKindMap(r.changelog.Entries, Deprecation),
		collectByKindMap(r.changelog.Entries, BugFix),
		collectByKindMap(r.changelog.Entries, Enhancement),
		collectByKindMap(r.changelog.Entries, Feature),
		collectByKindMap(r.changelog.Entries, KnownIssue),
		collectByKindMap(r.changelog.Entries, Security),
		collectByKindMap(r.changelog.Entries, Upgrade),
		collectByKindMap(r.changelog.Entries, Other),
	}

	tmpl, err := template.New("release-notes").
		Funcs(template.FuncMap{
			// nolint:staticcheck // ignoring for now, supports for multiple component is not implemented
			"linkPRSource": func(component string, ids []string) string {
				res := make([]string, len(ids))

				for i, id := range ids {
					res[i] = fmt.Sprintf("{%s-pull}%v[#%v]", component, id, id)
				}

				return strings.Join(res, " ")
			},
			// nolint:staticcheck // ignoring for now, supports for multiple component is not implemented
			"linkIssueSource": func(component string, ids []string) string {
				res := make([]string, len(ids))

				for i, id := range ids {
					res[i] = fmt.Sprintf("{%s-issue}%v[#%v]", component, id, id)
				}

				return strings.Join(res, " ")
			},
			// Capitalize sentence and ensure ends with .
			"beautify": func(s1 string) string {
				s2 := strings.Builder{}
				s2.WriteString(cases.Title(language.English).String(s1))
				if !strings.HasSuffix(s1, ".") {
					s2.WriteString(".")
				}
				return s2.String()
			},
			// Ensure components have section styling
			"header2": func(s1 string) string {
				s2 := strings.Builder{}
				s2.WriteString(s1)
				if !strings.HasSuffix(s1, "::") && s1 != "" {
					s2.WriteString("::")
				}
				return s2.String()
			},
		}).
		Parse(string(tpl))
	if err != nil {
		panic(err)
	}

	var data bytes.Buffer

	err = tmpl.Execute(&data, td)
	if err != nil {
		panic(err)
	}

	outFile := path.Join(r.dest, fmt.Sprintf("%s.asciidoc", r.changelog.Version))
	log.Printf("saving changelog in %s\n", outFile)

	return afero.WriteFile(r.fs, outFile, data.Bytes(), changelogFilePerm)
}

// Binds strings to the actual template file name
// These strings can be used in the config template field or renderer template flag
func getEmbeddedTemplates() map[string]string {
	return map[string]string{
		"asciidoc-embedded": "asciidoc-template.asciidoc",
	}
}

//go:embed asciidoc-template.asciidoc
var asciidocTemplate embed.FS

func (r Renderer) Template() ([]byte, error) {
	var data []byte
	var err error

	embeddedFileName, ok := getEmbeddedTemplates()[r.templ]
	if ok {
		data, err = asciidocTemplate.ReadFile(embeddedFileName)
		if err != nil {
			return []byte{}, fmt.Errorf("cannot read embedded template: %s %w", embeddedFileName, err)
		}

		return data, nil
	}

	data, err = afero.ReadFile(r.fs, fmt.Sprintf("./assets/%s", r.templ))
	if err != nil {
		return []byte{}, fmt.Errorf("cannot read custom template: %w", err)
	}

	return data, nil
}

func collectKinds(items []Entry) map[Kind]bool {
	// NOTE: collect kinds in a set-like map to avoid duplicates
	kinds := map[Kind]bool{}

	for _, e := range items {
		kinds[e.Kind] = true
	}

	return kinds
}

func collectByKindMap(entries []Entry, k Kind) map[string][]Entry {
	componentEntries := map[string][]Entry{}

	for _, e := range entries {
		if e.Kind == k {
			if len(e.Component) > 0 {
				componentEntries[e.Component] = append(componentEntries[e.Component], e)
			} else {
				componentEntries[""] = append(componentEntries[""], e)
			}
		}
	}

	return componentEntries
}

func collectByKind(items []Entry, k Kind) []Entry {
	entries := []Entry{}

	for _, e := range items {
		if e.Kind == k {
			entries = append(entries, e)
		}
	}

	return entries
}

func buildTitleByComponents(entries []Entry) string {
	configComponents := viper.GetStringSlice("components")

	switch len(configComponents) {
	case 0:
		return ""
	case 1:
		c := configComponents[0]
		for _, e := range entries {
			if c != e.Component && len(e.Component) > 0 {
				log.Fatalf("Component [%s] not found in config", e.Component)
			}
		}
		return c
	default:
		var match string
		for _, e := range entries {
			if e.Component == "" {
				log.Fatalf("Component cannot be assumed, choose it from config values: %s", e.File.Name)
			}

			match = ""
			for _, c := range configComponents {
				if e.Component != c {
					continue
				}
				match = e.Component
			}

			if match == "" {
				log.Fatalf("Component [%s] not found in config", e.Component)
			}
		}
		return match
	}
}
