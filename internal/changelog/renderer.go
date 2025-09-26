// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package changelog

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/elastic/elastic-agent-changelog-tool/internal/assets"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type Renderer struct {
	changelog Changelog
	fs        afero.Fs
	// dest is the destination location where the changelog is written to
	dest  string
	templ string
	repo  string
}

func NewRenderer(fs afero.Fs, c Changelog, dest string, templ string, repo string) *Renderer {
	return &Renderer{
		changelog: c,
		fs:        fs,
		dest:      dest,
		templ:     templ,
		repo:      repo,
	}
}

func (r Renderer) Render() error {
	log.Printf("render %s for version: %s\n", r.templ, r.changelog.Version)

	tpl, err := r.Template()
	if err != nil {
		log.Fatal(err)
	}

	type TemplateData struct {
		Component string
		Version   string
		Repo      string
		Changelog Changelog
		Kinds     map[Kind]bool

		// In Markdown, this goes to release notes
		Enhancement map[string][]Entry
		Feature     map[string][]Entry
		Security    map[string][]Entry
		BugFix      map[string][]Entry
		// In Markdown, this goes to breaking changes
		BreakingChange map[string][]Entry
		// In Markdown, this goes to deprecations
		Deprecation map[string][]Entry
		// In Markdown, this goes to known issues
		KnownIssue map[string][]Entry
		// In Markdown... TBD
		Upgrade map[string][]Entry
		Other   map[string][]Entry
	}

	td := TemplateData{
		buildTitleByComponents(r.changelog.Entries),
		r.changelog.Version,
		r.repo,
		r.changelog,
		collectKinds(r.changelog.Entries),
		// In Markdown, this goes to release notes
		collectByKindMap(r.changelog.Entries, Enhancement),
		collectByKindMap(r.changelog.Entries, Feature),
		collectByKindMap(r.changelog.Entries, Security),
		collectByKindMap(r.changelog.Entries, BugFix),
		// In Markdown, this goes to breaking changes
		collectByKindMap(r.changelog.Entries, BreakingChange),
		// In Markdown, this goes to deprecations
		collectByKindMap(r.changelog.Entries, Deprecation),
		// In Markdown, this goes to known issues
		collectByKindMap(r.changelog.Entries, KnownIssue),
		// In Markdown... TBD
		collectByKindMap(r.changelog.Entries, Upgrade),
		collectByKindMap(r.changelog.Entries, Other),
	}

	tmpl, err := template.New("release-notes").
		Funcs(template.FuncMap{
			"crossreferenceList": func(ids []string) string {
				return strings.Join(ids, "-")
			},
			// nolint:staticcheck // ignoring for now, supports for multiple component is not implemented
			"linkPRSource": func(repo string, ids []string) string {
				res := make([]string, len(ids))
				for i, id := range ids {
					res[i] = getLink(id, r.repo, "pull", r.templ)
				}
				return strings.Join(res, " ")
			},
			// nolint:staticcheck // ignoring for now, supports for multiple component is not implemented
			"linkIssueSource": func(repo string, ids []string) string {
				res := make([]string, len(ids))
				for i, id := range ids {
					res[i] = getLink(id, r.repo, "issues", r.templ)
				}
				return strings.Join(res, " ")
			},
			// Capitalize sentence and ensure ends with .
			"beautify": func(s string) string {
				if s == "" {
					return ""
				}
				s = strings.ToUpper(string(s[0])) + s[1:]
				if !strings.HasSuffix(s, ".") {
					s += "."
				}
				return s
			},
			// Indent lines
			"indent": func(s string) string {
				re := regexp.MustCompile(`\n|\r|^`)
				return re.ReplaceAllString(s, "\n  ")
			},
			"other_links": func() string {
				var links []string
				if len(td.KnownIssue) > 0 {
					links = append(
						links,
						"[Known issues](/release-notes/known-issues.md)",
					)
				}
				if len(td.BreakingChange) > 0 {
					links = append(
						links,
						fmt.Sprintf(
							"[Breaking changes](/release-notes/breaking-changes.md#%s-%s-breaking-changes)",
							r.repo,
							r.changelog.Version,
						),
					)
				}
				if len(td.Deprecation) > 0 {
					links = append(
						links,
						fmt.Sprintf(
							"[Deprecations](/release-notes/deprecations.md#%s-%s-deprecations)",
							r.repo,
							r.changelog.Version,
						),
					)
				}
				if len(links) > 0 {
					return fmt.Sprintf(
						"_This release also includes: %s._",
						strings.Join(links, " and"),
					)
				} else {
					return ""
				}
			},
			// Ensure components have section styling
			"header2": func(s1 string) string {
				return fmt.Sprintf("**%s**", s1)
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

	outFile := func(template string) string {
		if template == "markdown-index" {
			return path.Join(r.dest, r.changelog.Version, "index.md")
		} else if template == "markdown-breaking" {
			return path.Join(r.dest, r.changelog.Version, "breaking-changes.md")
		} else if template == "markdown-deprecations" {
			return path.Join(r.dest, r.changelog.Version, "deprecations.md")
		} else {
			return path.Join(r.dest, fmt.Sprintf("%s.asciidoc", r.changelog.Version))
		}
	}
	if r.templ != "asciidoc-embedded" {
		dir := path.Join(r.dest, r.changelog.Version)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory %s: %v", dir, err)
		}
	}
	return afero.WriteFile(r.fs, outFile(r.templ), data.Bytes(), changelogFilePerm)
}

func (r Renderer) Template() ([]byte, error) {
	var data []byte
	var err error

	if embeddedFileName, ok := assets.GetEmbeddedTemplates()[r.templ]; ok {
		var readFunc func(string) ([]byte, error)
		switch r.templ {
		case "markdown-index":
			readFunc = assets.MarkdownIndexTemplate.ReadFile
		case "markdown-breaking":
			readFunc = assets.MarkdownBreakingTemplate.ReadFile
		case "markdown-deprecations":
			readFunc = assets.MarkdownDeprecationsTemplate.ReadFile
		case "asciidoc-embedded":
			readFunc = assets.AsciidocTemplate.ReadFile
		}
		if readFunc != nil {
			data, err = readFunc(embeddedFileName)
			if err != nil {
				return nil, fmt.Errorf("cannot read embedded template: %s %w", embeddedFileName, err)
			}
			// If using the snippet/include model, update the includes
			if strings.Contains(r.dest, "release-notes/_snippets") {
				switch r.templ {
				case "markdown-index":
					addInclude(r.fs, r.changelog.Version, r.dest, "index")
				case "markdown-breaking":
					addInclude(r.fs, r.changelog.Version, r.dest, "breaking-changes")
				case "markdown-deprecations":
					addInclude(r.fs, r.changelog.Version, r.dest, "deprecations")
				}
			}
			return data, nil
		}
	}

	data, err = afero.ReadFile(r.fs, r.templ)
	if err != nil {
		return []byte{}, fmt.Errorf("cannot read custom template: %w", err)
	}

	return data, nil
}

func getLink(id string, repo string, ghType string, templ string) string {
	re := regexp.MustCompile(`\d+$`)
	number := re.FindString(id)
	if id == number {
		id = fmt.Sprintf("https://github.com/elastic/%s/%s/%s", repo, ghType, id)
	}
	if templ == "asciidoc-embedded" {
		// Format as AsciiDoc links
		return fmt.Sprintf("%s[#%s]", id, number)
	} else {
		// Format as Markdown links
		return fmt.Sprintf("[#%s](%s)", number, id)
	}
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

func addInclude(fs afero.Fs, version, dest, templ string) {
	// Extract minor version (e.g., "8.12" from "8.12.1")
	minorVersion := regexp.MustCompile(`^\d+\.\d+`).FindString(version)
	if minorVersion == "" {
		fmt.Printf("Could not get minor version from: %v\n", version)
		return
	}

	// Extract include directory (e.g., "/release-notes/...")
	includeDir := regexp.MustCompile(`/release-notes/.+$`).FindString(dest)
	if includeDir == "" {
		fmt.Printf("Could not derive include directory from: %v\n", dest)
		return
	}

	minorFilePath := fmt.Sprintf("%s/%s/%s.md", dest, templ, minorVersion)
	templateTypeFilePath := fmt.Sprintf("%s/%s.md", dest, templ)

	// Read or create the minor file
	minorFileContent, err := afero.ReadFile(fs, minorFilePath)
	if err != nil {
		// Create the file
		if err := afero.WriteFile(fs, minorFilePath, nil, changelogFilePerm); err == nil {
			fmt.Printf("Created new include file: %s\n", minorFilePath)
		}
		// Prepend new minor version include to the template type file (e.g. "breaking-changes")
		if templateTypeFileContent, err := afero.ReadFile(fs, templateTypeFilePath); err == nil {
			newMinorInclude := fmt.Sprintf(":::{include} %s/%s/%s.md\n:::", includeDir, templ, minorVersion)
			newContent := fmt.Sprintf("%s\n\n%s", newMinorInclude, templateTypeFileContent)
			afero.WriteFile(fs, templateTypeFilePath, []byte(newContent), changelogFilePerm)
		}
		minorFileContent = nil // ensure it's empty for next step
	}

	// Prepend new patch version include to the minor file
	newPatchInclude := fmt.Sprintf(":::{include} %s/%s/%s.md\n:::", includeDir, version, templ)
	newContent := fmt.Sprintf("%s\n\n%s", newPatchInclude, minorFileContent)
	afero.WriteFile(fs, minorFilePath, []byte(newContent), changelogFilePerm)
}
