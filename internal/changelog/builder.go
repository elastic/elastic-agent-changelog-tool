// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package changelog

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/elastic/elastic-agent-changelog-tool/internal/changelog/fragment"
	"github.com/elastic/elastic-agent-changelog-tool/internal/github"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type Builder struct {
	changelog Changelog
	filename  string
	fs        afero.Fs
	// src is the source location to gather changelog fragments
	src string
	// dest is the destination location where the changelog is written to
	dest string
}

func NewBuilder(fs afero.Fs, filename, version, src, dest string) *Builder {
	return &Builder{
		changelog: Changelog{Version: version},
		filename:  filename,
		fs:        fs,
		src:       src,
		dest:      dest,
	}
}

var changelogFilePerm = os.FileMode(0660)
var errNoFragments = errors.New("no fragments found in the source folder")

func (b Builder) Build() error {
	log.Printf("building changelog for version: %s\n", b.changelog.Version)
	log.Printf("collecting fragments from %s\n", b.src)

	var files []string
	err := afero.Walk(b.fs, b.src, func(path string, info os.FileInfo, err error) error {
		if info, err := b.fs.Stat(path); err == nil && !info.IsDir() {
			files = append(files, path)
		} else {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("cannot walk path %s: %w", b.src, err)
	}

	if len(files) == 0 {
		return errNoFragments
	}

	for _, file := range files {
		log.Printf("parsing %s", file)

		f, err := fragment.Load(b.fs, file)
		if err != nil {
			return fmt.Errorf("cannot load fragment from file %s: %w", file, err)
		}

		b.changelog.Entries = append(b.changelog.Entries, EntryFromFragment(f))
	}

	hc, err := github.GetHTTPClient(b.fs)
	if err != nil {
		return fmt.Errorf("cannot initialize http client: %w", err)
	}

	c := github.NewClient(hc)

	for i, entry := range b.changelog.Entries {
		// Filling empty PR fields
		if entry.LinkedPR[0] == 0 {
			commitHash, err := github.GetLatestCommitHash(entry.File.Name)
			if err != nil {
				log.Printf("cannot find commit hash, fill the PR field in changelog: %s", entry.File.Name)
				continue
			}

			prIDs, err := github.FillEmptyPRField(commitHash, c)
			if err != nil {
				log.Printf("fill the PR field in changelog: %s", entry.File.Name)
				continue
			}

			if len(prIDs) > 1 {
				log.Printf("delete the additional PRs in changelog: %s", entry.File.Name)
			}

			b.changelog.Entries[i].LinkedPR = prIDs
		} else {
			// Applying heuristics to PR fields
			originalPR, err := github.FindOriginalPR(entry.LinkedPR[0], c)
			if err != nil {
				log.Printf("check if the PR field is correct in changelog: %s", entry.File.Name)
				continue
			}

			b.changelog.Entries[i].LinkedPR = []int{originalPR}
		}
	}

	data, err := yaml.Marshal(&b.changelog)
	if err != nil {
		return fmt.Errorf("cannot marshall changelog: %w", err)
	}

	outFile := path.Join(b.dest, b.filename)
	log.Printf("saving changelog in %s\n", outFile)
	return afero.WriteFile(b.fs, outFile, data, changelogFilePerm)
}
