# Release process

This project uses [GoReleaser] to release a new version of the application.

Release publishing is automatically managed by the Jenkins CI ([Jenkinsfile]) and it's triggered by Git tags. 

Release artifacts are available in the [Releases] section. This project supports **reproducible builds**.

## Creating a new release

1. Fetch latest main from origin (remember to rebase the branch):

        git fetch origin
        git rebase origin/main

2. Create Git tag with release candidate:

        git tag v0.2.0 -m "Release version 0.2.0"

3. Push new tag to the upstream.

        git push origin v0.2.0 

4. The CI will run a new job for the just pushed tag and publish released artifacts.

## Reproducible builds

This project supports [reproducible builds] (builds that always produce the exact same binary given the same sources).

Support is achieved through:
- `Makefile`: using [source epoch] instead of build datetime
- `GoReleaser`(see [docs][gorel-docs]): using `trimpath`, `mod_timestamp` and `.CommitDate`


[GoReleaser]: https://goreleaser.com/
[Jenkinsfile]: https://github.com/elastic/elastic-agent-changelog-tool/blob/main/.ci/Jenkinsfile
[Releases]: https://github.com/elastic/elastic-package/releases
[reproducible builds]: https://reproducible-builds.org/
[source epoch]: https://reproducible-builds.org/docs/source-date-epoch/
[gorel-docs]: https://goreleaser.com/customization/build/#reproducible-builds