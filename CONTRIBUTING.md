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
