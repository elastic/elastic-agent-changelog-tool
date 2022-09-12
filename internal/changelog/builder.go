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
	"strconv"
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

	log.Println("Verifying entries:")

	for i, entry := range b.changelog.Entries {
		// Filling empty PR fields
		if len(entry.LinkedPR) == 0 {

			commitHash, err := GetLatestCommitHash(entry.File.Name)
			if err != nil {
				log.Printf("%s: cannot find commit hash, fill the PR field in changelog", entry.File.Name)
				continue
			}

			prIDs, err := FillEmptyPRField(commitHash, owner, repo, c)
			if err != nil {
				log.Printf("%s: fill the PR field in changelog", entry.File.Name)
				continue
			}

			if len(prIDs) > 1 {
				log.Printf("%s: multiple PRs found, please remove all but one of them", entry.File.Name)
			}

			b.changelog.Entries[i].LinkedPR = prIDs
		} else {
			// Applying heuristics to PR fields
			originalPR, err := FindOriginalPR(entry.LinkedPR[0], owner, repo, c)
			if err != nil {
				log.Printf("%s: check if the PR field is correct in changelog: %s", entry.File.Name, err.Error())
				continue
			}

			b.changelog.Entries[i].LinkedPR = []string{originalPR}
		}

		if len(entry.LinkedIssue) == 0 && len(b.changelog.Entries[i].LinkedPR) > 0 {
			linkedIssues := []string{}

			for _, prURL := range b.changelog.Entries[i].LinkedPR {
				tempIssues, err := FindIssues(graphqlClient, context.Background(), owner, repo, prURL, 50)
				if err != nil {
					log.Printf("%s: could not find linked issues for pr: %s: %s", entry.File.Name, entry.LinkedPR, err.Error())
					continue
				}

				linkedIssues = append(linkedIssues, tempIssues...)
				if len(linkedIssues) > 1 {
					log.Printf("%s: multiple issues found, please remove all but one of them", entry.File.Name)
				}
			}

			b.changelog.Entries[i].LinkedIssue = linkedIssues
		} else if len(entry.LinkedIssue) == 1 {
			_, err := ExtactEventNumber("issue", entry.LinkedIssue[0])
			if err != nil {
				log.Printf("%s: check if the issue field is correct in changelog: %s", entry.File.Name, err.Error())
			}
		}

	}

	data, err := yaml.Marshal(&b.changelog)
	if err != nil {
		return fmt.Errorf("cannot marshall changelog: %w", err)
	}

	outFile := path.Join(b.dest, b.filename)
	log.Printf("saving changelog in: %s\n", outFile)
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

func ExtactEventNumber(linkType, eventURL string) (string, error) {
	urlParts := strings.Split(eventURL, "/")

	if len(urlParts) < 1 {
		return "", fmt.Errorf("cant get event number")
	}

	switch linkType {
	// maybe use regex to validate instead of a simple string check
	case "pr":
		if !strings.Contains(eventURL, "pull") {
			return "", fmt.Errorf("link is invalid for pr")
		}
	case "issue":
		if !strings.Contains(eventURL, "issues") {
			return "", fmt.Errorf("link is invalid for issue")
		}
	}

	return urlParts[len(urlParts)-1], nil
}

func CreateEventLink(linkType, owner, repo, eventID string) string {
	switch linkType {
	case "issue":
		return fmt.Sprintf("https://github.com/%s/%s/issues/%s", owner, repo, eventID)
	case "pr":
		return fmt.Sprintf("https://github.com/%s/%s/pull/%s", owner, repo, eventID)
	default:
		return ""
	}
}

func GetLatestCommitHash(fileName string) (string, error) {
	response, err := exec.Command("git", "log", "--diff-filter=A", "--format=%H", "changelog/fragments/"+fileName).Output()
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(string(response), "\n", ""), nil
}

func FindIssues(graphqlClient *github.ClientGraphQL, ctx context.Context, owner, name string, prURL string, issuesLen int) ([]string, error) {
	prID, err := ExtactEventNumber("pr", prURL)
	if err != nil {
		return nil, err
	}

	prIDInt, _ := strconv.Atoi(prID)

	issues, err := graphqlClient.PR.FindIssues(ctx, owner, name, prIDInt, issuesLen)
	if err != nil {
		return nil, err
	}

	issueLinks := make([]string, len(issues))

	for i, issue := range issues {
		issueLinks[i] = CreateEventLink("issue", owner, name, issue)
	}

	return issueLinks, nil
}

func FillEmptyPRField(commitHash, owner, repo string, c *github.Client) ([]string, error) {
	pr, err := github.FindPR(context.Background(), c, owner, repo, commitHash)
	if err != nil {
		return []string{}, err
	}

	prLinks := []string{}

	for _, item := range pr.Items {
		prLinks = append(prLinks, CreateEventLink("pr", owner, repo, strconv.Itoa(item.PullRequestID)))
	}

	return prLinks, nil
}

func FindOriginalPR(prURL string, owner, repo string, c *github.Client) (string, error) {
	linkedPR, err := ExtactEventNumber("pr", prURL)
	if err != nil {
		return "", err
	}

	linkedPRString, _ := strconv.Atoi(linkedPR)

	pr, _, err := c.PullRequests.Get(context.Background(), owner, repo, linkedPRString)
	if err != nil {
		return "", err
	}

	prID, err := github.TestStrategies(pr, &github.BackportPRNumber{}, &github.PRNumber{})
	if err != nil {
		return "", err
	}

	prLink := CreateEventLink("pr", owner, repo, strconv.Itoa(prID))
	return prLink, nil
}
