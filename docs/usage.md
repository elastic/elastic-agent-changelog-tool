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

### Viewing the changelog

There is no way to view the changelog at the moment. You can view the **consolidated changelog** by running:

```
$ elastic-agent-changelog-tool build --version=8.2.1
```

from the repository root folder.

An example is [`../changelog/0.1.0.yaml`](../changelog/0.1.0.yaml).

## I'm a maintainer

To be done.

## I'm the release manager

To be done.
