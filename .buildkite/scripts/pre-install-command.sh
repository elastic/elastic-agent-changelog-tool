#!/bin/bash
set -euo pipefail

source .buildkite/scripts/tooling.sh

add_bin_path(){
    mkdir -p "${WORKSPACE}/bin"
    export PATH="${WORKSPACE}/bin:${PATH}"
}

with_go_junit_report() {
    go install github.com/jstemmer/go-junit-report/v2@latest
}

with_goreleaser() {
    go install github.com/goreleaser/goreleaser@v1.6.3
}

with_go() {
    go_version=$1
    url=$(get_gvm_link "${GVM}")
    retry 5 curl -sL -o "${WORKSPACE}/bin/gvm" "${url}"
    chmod +x "${WORKSPACE}/bin/gvm"
    ls ${WORKSPACE}/bin/ -l
    eval "$(gvm $go_version)"
    echo $GOPATH
    export PATH="${PATH}:${GOPATH}:${GOBIN}"
    go version
}

# for gvm link
get_gvm_link() {
    gvm_version=$1
    platform_type="$(uname)"
    platform_type_lowercase="${platform_type,,}"
    arch_type="$(uname -m)"
    [[ ${arch_type} == "aarch64" ]] && arch_type="arm64" # gvm do not have 'aarch64' name for archetecture type
    [[ ${arch_type} == "x86_64" ]] && arch_type="amd64"
    echo "https://github.com/andrewkroh/gvm/releases/download/${gvm_version}/gvm-${platform_type_lowercase}-${arch_type}"
}

# Required env variables:
#   WORKSPACE
WORKSPACE=${WORKSPACE:-"$(pwd)"}
