.PHONY: build

build:
	go build -o elastic-agent-changelog-tool

licenser:
	go run github.com/elastic/go-licenser -license Elasticv2

test:
	go test -v ./...
