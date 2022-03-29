MODULE = github.com/elastic/elastic-agent-changelog-tool
VERSION_IMPORT_PATH = $(MODULE)/internal/version
VERSION_COMMIT_HASH = `git describe --always --long --dirty`
VERSION_BUILD_TIME = `date +%s`
VERSION_TAG = `(git describe --exact-match --tags 2>/dev/null || echo '') | tr -d '\n'`
VERSION_LDFLAGS = -X $(VERSION_IMPORT_PATH).CommitHash=$(VERSION_COMMIT_HASH) -X $(VERSION_IMPORT_PATH).BuildTime=$(VERSION_BUILD_TIME) -X $(VERSION_IMPORT_PATH).Tag=$(VERSION_TAG)

.PHONY: build

build:
	go build -ldflags "$(VERSION_LDFLAGS)" -o elastic-agent-changelog-tool

licenser:
	go run github.com/elastic/go-licenser -license Elasticv2

test:
	go test -v ./...
