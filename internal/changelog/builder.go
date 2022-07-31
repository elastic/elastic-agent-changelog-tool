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
	"os/exec"
	"path"
	"path/filepath"
	"strings"

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

func (b Builder) Build(owner, repo string) error {
	log.Printf("building changelog for version: %s\n", b.changelog.Version)
	log.Printf("collecting fragments from %s\n", b.src)

	var files []string
	err := afero.Walk(b.fs, b.src, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == "fixtures" {
			return filepath.SkipDir
		}

		return collectFragment(b.fs, path, info, err, &files)
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
	graphqlClient := github.NewGraphQLClient(hc)

	for i, entry := range b.changelog.Entries {
		// Filling empty PR fields
		if len(entry.LinkedPR) == 0 {

			commitHash, err := GetLatestCommitHash(entry.File.Name)
			if err != nil {
				log.Printf("cannot find commit hash, fill the PR field in changelog: %s", entry.File.Name)
				continue
			}

			prIDs, err := FillEmptyPRField(commitHash, owner, repo, c)
			if err != nil {
				log.Printf("fill the PR field in changelog: %s", entry.File.Name)
				continue
			}

			if len(prIDs) > 1 {
				log.Printf("multiple PRs found for %s, please remove all but one of them", entry.File.Name)
			}

			b.changelog.Entries[i].LinkedPR = prIDs
		} else {
			// Applying heuristics to PR fields
			originalPR, err := FindOriginalPR(entry.LinkedPR[0], owner, repo, c)
			if err != nil {
				log.Printf("check if the PR field is correct in changelog: %s", entry.File.Name)
				continue
			}

			b.changelog.Entries[i].LinkedPR = []int{originalPR}
		}

		if len(entry.LinkedIssue) == 0 && len(entry.LinkedPR) > 0 {
			linkedIssues := []int{}

			for _, pr := range entry.LinkedPR {
				tempIssues, err := FindIssues(graphqlClient, context.Background(), owner, repo, pr, 50)
				if err != nil {
					log.Printf("could not find linked issues for pr id: %d", entry.LinkedPR)
					continue
				}

				linkedIssues = append(linkedIssues, tempIssues...)
			}

			b.changelog.Entries[i].LinkedIssue = linkedIssues
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

func collectFragment(fs afero.Fs, path string, info os.FileInfo, err error, files *[]string) error {
	if info, err := fs.Stat(path); err == nil && !info.IsDir() {
		if filepath.Ext(path) == ".yaml" {
			*files = append(*files, path)
		} else {
			log.Printf("skipping %s (not a YAML file)", path)
		}
	} else {
		return err
	}

	return nil
}

func GetLatestCommitHash(fileName string) (string, error) {
	response, err := exec.Command("git", "log", "--diff-filter=A", "--format=%H", "changelog/fragments/"+fileName).Output()
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(string(response), "\n", ""), nil
}

func FindIssues(graphqlClient *github.ClientGraphQL, ctx context.Context, owner, name string, prID, issuesLen int) ([]int, error) {
	issues, err := graphqlClient.PR.FindIssues(ctx, owner, name, prID, issuesLen)
	if err != nil {
		return nil, err
	}

	return issues, nil
}

func FillEmptyPRField(commitHash, owner, repo string, c *github.Client) ([]int, error) {
	pr, err := github.FindPR(context.Background(), c, owner, repo, commitHash)
	if err != nil {
		return []int{}, err
	}

	var prIDs []int

	for _, item := range pr.Items {
		prIDs = append(prIDs, item.PullRequestID)
	}

	return prIDs, nil
}

func FindOriginalPR(linkedPR int, owner, repo string, c *github.Client) (int, error) {
	pr, _, err := c.PullRequests.Get(context.Background(), owner, repo, linkedPR)
	if err != nil {
		return 0, err
	}

	prID, err := github.TestStrategies(pr, &github.BackportPRNumber{}, &github.PRNumber{})
	if err != nil {
		return 0, err
	}

	return prID, nil
}
