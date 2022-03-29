// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

/*
gitreporoot package implement functionalities to find the Git repository root folder.

The package uses a shell execution of git CLI to get the repo root, as this functionality is still
not supported in https://github.com/go-git/go-git
*/
package gitreporoot

import (
	"fmt"
	"os/exec"
	"strings"
)

// Find uses git via shell to locate the top level directory
func Find() (string, error) {
	p, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("cannot find working tree top level path: %w", err)
	}

	return strings.TrimSpace(string(p)), nil
}
