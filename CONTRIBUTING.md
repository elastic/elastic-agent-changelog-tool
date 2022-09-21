# Contributing

Thanks for considering contributing to this project!

## Opening issues

If you find a bug, please feel free to [open an issue](https://github.com/elastic/elastic-agent-changelog-tool/issues).

## Fixing bugs

We love pull requests. Here’s a quick guide:

1. [Fork this repository](https://github.com/elastic/elastic-agent-changelog-tool/fork) and then clone it locally:

  ```bash
  $ git clone https://github.com/elastic/elastic-agent-changelog-tool
  $ cd elastic-agent-changelog-tool
  $ make build # ensure building the binary works
  $ make test # ensure tests are green
  ```

2. Create a topic branch for your changes:

  ```bash
  git checkout -b fix-for-that-thing
  ```
3. Commit a failing test for the bug:

  ```bash
  git commit -am "Adds a failing test to demonstrate that thing"
  ```

4. Commit a fix that makes the test pass:

  ```bash
  git commit -am "Adds a fix for that thing!"
  ```

5. Run the tests:

  ```bash
  make test
  ```

6. If everything looks good, push to your fork:

  ```bash
  git push origin fix-for-that-thing
  ```

7. [Submit a pull request.](https://help.github.com/articles/creating-a-pull-request)

## Adding new features

[Open an issue](https://github.com/elastic/elastic-agent-changelog-tool/issues) and let’s design it together.

## Releasing a new version

Versions should use the format `v<semver>`, like `v0.1.0` or `v1.2.3`. 
To release a **pre-release** version see below.

1. Create a dedicated branch;
2. Create and commit consolidated changelog with `elastic-agent-changelog-tool build --version <version> --repo elastic-agent-changelog-tool`;
3. Remove all previous fragments from git with `git rm changelog/fragments/*`;
4. Create a PR with changes;
5. Merge it;
6. Pull `main` and tag it with `<version>`;
7. Push the tag, the release process will kick off automatically.

### Releasing a pre-release version

Versions should use the format `v<semver>-<prerelease info>`, like `v0.1.0-beta.1` or `v1.2.3-rc.2`. 

1. Create a dedicated branch;
2. Create and commit consolidated changelog with `elastic-agent-changelog-tool build --version <version> --repo elastic-agent-changelog-tool`;
3. Create a PR with changes;
4. Merge it;
5. Pull `main` and tag it with `<version>`;
6. Push the tag, the release process will kick off automatically.

**Note**: in this version we do not remove fragments, as they will be needed for the stable release version changelog.