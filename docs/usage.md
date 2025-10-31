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

These steps require [GitHub Authentication](./github-authentication.md).

From the repository root folder build the consolidated changelog with:

```
$ elastic-agent-changelog-tool build --version=next --owner <owner> --repo <repo>
```

then render the consolidated changelog with:
```
$ elastic-agent-changelog-tool render --version=next --file_type <asciidoc|markdown>
```

An example is [`../changelog/0.1.0.yaml`](../changelog/0.1.0.yaml).

### My PR does not need a changelog

Repositories using this tool are expected to leverage the `pr-has-fragment` command to validate presence of a Changelog Fragment in a Pull Requests. The command can be configured to validate PRs without a Changelog Fragment but labelled with specific labels.

At this moment adding a label named `skip-changelog` or `backport` to a PR will skip the validation, allowing the labelled PR not to contain a Changelog Fragment.

## I'm a maintainer

As a maintainer you are responsible for ensureing **changelog fragments** are present in Pull Requests.
You may also be interested in  viewing the changelog for unreleased changes.

### Ensuring Changelog Fragment presence in PRs

These steps require [GitHub Authentication](./github-authentication.md).

`elastic-agent-chagelog-tool` has a dedicated command for this: `pr-has-fragment`.
Given a PR number this command checks for presence of an added Changelog Fragment.

For futrher details look at command usage: `elastic-agent-changelog-tool pr-has-fragment --help`

### Previewing the changelog

These steps require [GitHub Authentication](./github-authentication.md).

From the repository root folder build the consolidated changelog with:

```
$ elastic-agent-changelog-tool build --version=next --owner <owner> --repo <repo>
```

then render the consolidated changelog with:
```
$ elastic-agent-changelog-tool render --version=next --file_type <asciidoc|markdown>
```

An example is [`../changelog/0.1.0.yaml`](../changelog/0.1.0.yaml).

### Including the changelog in a pre-release version

These steps require [GitHub Authentication](./github-authentication.md).

Pre-release versions are development releases, it's worth including the changelog but as they are not stable releases, you should not remove the fragments.
The side effect is that the changelog will include all entries from latest stable release up to current pre-release, but in the end this correctly reflects the content of the provided pre-release.

1. Create consolidated changelog with `$ elastic-agent-changelog-tool build --version <version> --owner <owner> --repo <repo>`;
* This will create `./changelog/x.y.z.yaml`;
2. Create rendered changelog with `$ elastic-agent-changelog-tool render --version <version> --file_type <asciidoc|markdown>`;

    Depending on the specified `file_type`, this will generate the following files:
    * `markdown`:
      * Release notes: `./changelog/<version>/index.md`
      * Breaking changes: `./changelog/<version>/breaking-changes.md`
      * Deprecations: `./changelog/<version>/deprecations.md`
    * `asciidoc`: `changelog/<version>.asciidoc`
3. Use the rendered changelog.

**Note**: we do not remove fragments, as they will be needed for the stable release version changelog.

## I'm the release manager

### Preparing the changelog

These steps require [GitHub Authentication](./github-authentication.md).

1. Wait for the last BC of the release. If another BC is generated after that or a patch version for a previous minor is released, you might need to restart the process.
1. Create a branch **from the commit of the BC**.
1. From the root folder of the repository run:

    ```
    $ elastic-agent-changelog-tool build --version x.y.z --owner <owner> --repo <repo>
    ```

    Where:

    * `x.y.z` is the version to release.
    * `owner` is the user / organization the repository to use belongs to. The default value is `elastic`.
    * `repo` is the name of the repository containing the issues and PRs. The default value is `elastic-agent`.

    >NOTE: If the repo you're targeting has an `owner` and `repo` defined in a `config.changelog.yaml` file, you do not need to specify them when running the command.

    This will create `./changelog/x.y.z.yaml`.

    >NOTE: If there are no changes included in the release (there are no changelog fragment files), the changelog file will be generated with an empty array of entries.
1. From the root of the repository run:
    ```
    $ elastic-agent-changelog-tool cleanup
    ```
1. Commit the previous changes (consolidated changelog and removed files)
1. From the root folder of the repository run:
    ```
    $ elastic-agent-changelog-tool render --version x.y.z --file_type <asciidoc|markdown>
    ```

    >IMPORTANT: Use `file_type` `markdown` for 9.x versions and `asciidoc` for 8.x versions.
    >If the repo you're targeting has a `file_type` defined in a `config.changelog.yaml` file, you do not need to specify it when running the command.

    The files that are generated depend on the specified `file_type`. The destination directory depends on the `rendered_changelog_destination` defined in the the repo's `config.changelog.yaml`. If no `rendered_changelog_destination` is specified, it will be added to the `changelog` directory.

    * `markdown`: These files will be created:
      * Release notes: `<rendered_changelog_destination>/<version>/index.md`
      * Breaking changes: `<rendered_changelog_destination>/<version>/breaking-changes.md`
      * Deprecations: `<rendered_changelog_destination>/<version>/deprecations.md`

      If the `rendered_changelog_destination` is set to `release-notes/_snippets`, the related `_snippets` files will automatically be updated.

    * `asciidoc`: There will be one file created, `<rendered_changelog_destination>/<version>.asciidoc`, and you will need to integrate the generated content into the changelog.

1. If the changelog is stored in the same repository, commit the changes in this same branch.
1. Create a PR with the changes to the `x.y` branch.


### On Release Day

Once the release is given the final go on release day:
* Merge the PR created in the previous section.
* Forward-port it to the other active branches, up to `main`.
