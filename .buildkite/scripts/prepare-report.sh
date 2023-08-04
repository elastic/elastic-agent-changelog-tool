#!/bin/bash

set -euo pipefail

echo "--- Pre install"
source .buildkite/scripts/pre-install-command.sh
add_bin_path
with_go_junit_report

# Create Junit report for junit annotation plugin
go-junit-report > junit-report.xml < tests-report.txt
exit $exit_code
