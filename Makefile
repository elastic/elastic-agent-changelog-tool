.PHONY: build

build:
	go build -o elastic-agent-changelog-tool

test:
	go test -v ./...
