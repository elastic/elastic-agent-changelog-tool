// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package changelog

import (
	"context"
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
		// Applying heuristics
		originalPR, err := FindOriginalPR(entry, c)
		if err == nil {
			b.changelog.Entries[i].LinkedPR = originalPR
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

func FindOriginalPR(entry Entry, c *github.Client) (int, error) {
	pr, _, err := c.PullRequests.Get(context.Background(), "elastic", "beats", entry.LinkedPR)
	if err != nil {
		return 0, err
	}

	prID, err := github.TestStrategies(pr, &github.BackportPRNumber{}, &github.PRNumber{})
	if err != nil {
		return 0, err
	}

	return prID, nil
}
