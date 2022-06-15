MODULE = github.com/elastic/elastic-agent-changelog-tool
VERSION_IMPORT_PATH = $(MODULE)/internal/version
VERSION_COMMIT_HASH = `git describe --always --long --dirty`
SOURCE_DATE_EPOCH = `git log -1 --pretty=%ct` # https://reproducible-builds.org/docs/source-date-epoch/
VERSION_TAG = `(git describe --exact-match --tags 2>/dev/null || echo '') | tr -d '\n'`
VERSION_LDFLAGS = -X $(VERSION_IMPORT_PATH).CommitHash=$(VERSION_COMMIT_HASH) -X $(VERSION_IMPORT_PATH).SourceDateEpoch=$(SOURCE_DATE_EPOCH) -X $(VERSION_IMPORT_PATH).Tag=$(VERSION_TAG)

.PHONY: build

# NOTE: this command is only for dev builds, releases are build using goreleaser
build:
	go build -ldflags "$(VERSION_LDFLAGS)" -o elastic-agent-changelog-tool

licenser:
	go-licenser -license Elasticv2

test:
	go test -v ./...
