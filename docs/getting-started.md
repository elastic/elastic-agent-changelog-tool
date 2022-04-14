# Getting started

This page describes all steps required for going from zero to viewing a consolidated changelog.

## 1. Prerequisites

Ensure you have `git` installed and available.

## 2. Install

Follow one of the installation ways described in [`./install.md`](./install.md).

## 3. Create a changelog fragment

From the root folder of the repository run:

```
$ elastic-agent-changelog-tool new "my test fragment"
```

This will create `./changelog/fragments/<timestamp>-my-test-fragment.yaml` with this content:

```yaml
# one of:
# - breaking-change: a change to previously-documented behavior
# - deprecation: functionality that is being removed in a later release
# - bug-fix: fixes a problem in a previous version
# - enhancement: extends functionality but does not break or fix existing behavior
# - feature: new functionality
# - known-issue: problems that we are aware of in a given version
# - security: impacts on the security of a product or a user’s deployment.
# - upgrade: important information for someone upgrading from a prior version
# - other: does not fit into any of the other categories
kind: feature
# Change summary; a 80ish characters long description of the change.
summary: 
# Long description; in case the summary is not enough to describe the change this field accomodate a description without length limits.
# description: 
# Affected component; a word indicating the component this changeset affects.
component:
# PR number; optional; the PR number that added the changeset.
# If not present is automatically filled by the tooling finding the PR where this changelog fragment has been added.
# NOTE: the tooling supports backports, so it's able to fill the original PR number instead of the backport PR number.
# Please provide it if you are adding a fragment for a different PR.
# pr: 1234
# Issue number; optional; the GitHub issue related to this changeset (either closes or is part of).
# If not present is automatically filled by the tooling with the issue linked to the PR number.
# issue: 1234
# Repository URL; optional; the repository URL related to this changeset and pr and issue numbers.
# If not present is automatically filled by the tooling based on the repository this file has been committed in.
# repository: https://github.com/elastic/elastic-agent-changelog-tool
```

Ensure `kind` is correct and fill the `summary` field with a brief description. You can ignore `component`, but you must set `pr`, `issue` and `repository` manually, the logic to fill them automatically is still work in progress.

Save and close the file.

## 4. Create the consolidated changelog

From the root folder of the repository run:

```
$ elastic-agent-changelog-tool build
```

This will create `./changelog/0.1.0.yaml` with content similar to this:

```yaml
version: 0.1.0
entries:
    - summary: Add Changelog Fragment creation
      description: ""
      kind: feature
      pr: 13
      issue: 21
      timestamp: 1649924282
      file:
        name: 1649924282-changelog-fragment-creation.yaml
        checksum: 10bf1bf67509a524e48f0795be75c278e29c0b47

```

There will be multiple entries, one for each files in `changelog/fragments`.

**NOTE:** the version is currently hard coded to `0.1.0` (the first version for `elastic-agent-changelog-tool`).
