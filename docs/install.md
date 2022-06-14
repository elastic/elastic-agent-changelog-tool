# Installation

## Install from release binary

Download the latest release from the [Releases](https://github.com/elastic/elastic-agent-changelog-tool/releases/latest) page.

On macOS, use `xattr -r -d com.apple.quarantine elastic-agent-changelog-tool` after downloading to allow the binary to run.

## Install with go install

If you have a `go` development environment setup, you can use `go install`:

```
$ go install github.com/elastic/elastic-agent-changelog-tool@latest
```

_Please make sure that you've correctly [setup environment variables](https://golang.org/doc/gopath_code.html#GOPATH) -
`$GOPATH` and `$PATH`, and `elastic-agent-changelog-tool` should be accessible from your `$PATH`._

## Install from source code

If you have a `go` development environment, you can clone this repo and build the source files with:

```
$ git clone https://github.com/elastic/elastic-agent-changelog-tool.git
$ cd elastic-agent-changelog-tool
$ make build
```

`elastic-agent-changelog-tool` binary will be created in the root folder.
Add it to your path.
