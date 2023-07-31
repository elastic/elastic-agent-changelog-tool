#!/bin/bash

set -euo pipefail

echo "--- Pre install"
source .buildkite/scripts/pre-install-command.sh
add_bin_path
with_go "${GO_VERSION}"
with_goreleaser

echo "--- Snapshot"
# goreleaser release --snapshot
git rev-parse --is-shallow-repository
git --no-pager tag
