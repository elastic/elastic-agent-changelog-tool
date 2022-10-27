# How to use this tool

This tool provides functionalities for different personas:
- the customer or user needs to view the changelog for unreleased development (if they are interested in released changes they must refer to the public release notes)
- the developer needs to provide changelog fragments for changeset and may want to review the unreleased changelog
- the maintainer needs to verify the proper changelog fragment is provided in PRs and may want to review the unreleased changelog
- the release manager needs to build public release notes from the unreleased changelog.

This repository uses the tool to keep the changelog, so you can seeing it in action (look at `changelog` folder).

## I'm a customer/user

To be done.

## I'm a developer

As a developer you are responsible for adding one or multiple **changelog fragments** in your Pull Requests.
You may also be interested in  viewing the changelog for unreleased changes.

### Adding a changelog fragment

To create a new fragment:

```
$ elastic-agent-changelog-tool new a-changeset-filename
```

A new file will be created in the changelog fragments folder (default to `changelog/fragments`).

You must edit the created file compiling all uncommented fields and optionally uncomment and compile commented fields. Guidance is provided through comments in the file.

The fragment is created from the template available in [`../internal/changelog/fragment/template.yaml`](../internal/changelog/fragment/template.yaml).

### Previewing the changelog

From the repository root folder build the consolidated changelog with:

```
$ elastic-agent-changelog-tool build --version=next --owner <owner> --repo <repo>
```

the render the consolidated changelog with:

```
$ elastic-agent-changelog-tool render --version=next
```

An example is [`../changelog/0.1.0.yaml`](../changelog/0.1.0.yaml).

### My PR does not need a changelog

Repositories using this tool are expected to leverage the `pr-has-fragment` command to validate presence of a Changelog Fragment in a Pull Requests. The command can be configured to validate PRs without a Changelog Fragment but labelled with specific labels.

At this moment adding a label named `skip-changelog` or `backport` to a PR will skip the validation, allowing the labelled PR not to contain a Changelog Fragment.

## I'm a maintainer

As a maintainer you are responsible for ensureing **changelog fragments** are present in Pull Requests.
You may also be interested in  viewing the changelog for unreleased changes.

### Ensuring Changelog Fragment presence in PRs

`elastic-agent-chagelog-tool` has a dedicated command for this: `pr-has-fragment`.
Given a PR number this command checks for presence of an added Changelog Fragment.

For futrher details look at command usage: `elastic-agent-changelog-tool pr-has-fragment --help`

### Previewing the changelog

From the repository root folder build the consolidated changelog with:

```
$ elastic-agent-changelog-tool build --version=next --owner <owner> --repo <repo>
```

the render the consolidated changelog with:

```
$ $ elastic-agent-changelog-tool render --version=next
```

An example is [`../changelog/0.1.0.yaml`](../changelog/0.1.0.yaml).

### Including the changelog in a pre-release version

Pre-release versions are development releases, it's worth including the changelog but as they are not stable releases, you should not remove the fragments.
The side effect is that the changelog will include all entries from latest stable release up to current pre-release, but in the end this correctly reflects the content of the provided pre-release.

1. Create consolidated changelog with `$ elastic-agent-changelog-tool build --version <version> --owner <owner> --repo <repo>`;
* This will create `./changelog/x.y.z.yaml`;
2. Create rendered changelog with `$ elastic-agent-changelog-tool render --version <version>`;
* This will generate an asciidoc file in the `changelog/` directory;
3. Use the rendered changelog.

**Note**: we do not remove fragments, as they will be needed for the stable release version changelog.

## I'm the release manager

> NOTE: This instructions are for the 0.2.0 milestone (when ready). As more features and automation is added the process will become simpler.

### Preparing the changelog

* Wait for the last BC of the release. If another BC is generated after that or a patch version for a previous minor is released, you might need to restart the process.
* Create a branch **from the commit of the BC**.
* From the root folder of the repository run:

```
$ elastic-agent-changelog-tool build --version x.y.z --owner <owner> --repo <repo>
```
* Where:
  * `x.y.z` is the version to release.
  * `owner` is the user / organization the repository to use belongs to. The default value is `elastic`.
  * `repo` is the name of the repository containing the issues / PRs, etc. The default value is `elastic-agent`.
* This will create `./changelog/x.y.z.yaml`.
* From the root of the repository run:
```
$ elastic-agent-changelog-tool cleanup
```
* Commit the previous changes (consolidated changelod and removed files)
* From the root folder of the repository run:
```
$ elastic-agent-changelog-tool render --version x.y.z
```
* This will generate an asciidoc fragment in the `changelog/` directory.
* Integrate the generated fragment into the changelog. If the changelog is stored in the same repository, commit the changes in this same branch.
* Create a PR with the changes to the `x.y` branch.


### On Release Day

Once the release is given the final go on release day:
* Merge the PR created in the previous section.
* Forward-port it to the other active branches, up to `main`.
