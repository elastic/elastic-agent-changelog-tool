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
	dest string
}

func NewRenderer(fs afero.Fs, c Changelog, dest string) *Renderer {
	return &Renderer{
		changelog: c,
		fs:        fs,
		dest:      dest,
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

	tmpl, err := template.New("asciidoc-release-notes").
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

//go:embed asciidoc-template.asciidoc
var asciidocTemplate embed.FS

func (r Renderer) Template() ([]byte, error) {
	data, err := asciidocTemplate.ReadFile("asciidoc-template.asciidoc")
	if err != nil {
		return []byte{}, fmt.Errorf("cannot read embedded template: %w", err)
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

func collectByKindMap(items []Entry, k Kind) map[string][]Entry {
	componentEntries := map[string][]Entry{}

	for _, e := range items {
		if e.Kind == k {
			for _, c := range viper.GetStringSlice("components") {
				if c != e.Component {
					continue
				}

				componentEntries[e.Component] = append(componentEntries[e.Component], e)
			}
		}
	}

	if len(componentEntries) == 0 {
		for _, e := range items {
			if e.Kind == k {
				componentEntries["unidentified component"] = append(componentEntries[e.Component], e)
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
	matchingComponents := map[string]struct{}{}
	entriesComponents := map[string]struct{}{}
	var componentNotFound bool

	for _, e := range entries {
		if e.Component == "" {
			componentNotFound = true
		}
		for _, c := range viper.GetStringSlice("components") {
			if c != e.Component {
				componentNotFound = true
				continue
			}
			matchingComponents[c] = struct{}{}
		}
	}

	for k := range entriesComponents {
		if _, ok := matchingComponents[k]; !ok {
			log.Printf("Component [%s] not found in config", k)
		}
	}

	components := []string{}
	for k := range matchingComponents {
		components = append(components, k)
	}

	switch {
	case len(components) == 0:
		return "unspecified component"
	case componentNotFound == true:
		return fmt.Sprintf("%s %s", strings.Join(components, " and "), "and unidentified component")
	default:
		return strings.Join(components, " and ")
	}
}
