FROM docker.io/library/alpine:3.15.6

COPY elastic-agent-changelog-tool /usr/bin/elastic-agent-changelog-tool

ENTRYPOINT ["/usr/bin/elastic-agent-changelog-tool"]
