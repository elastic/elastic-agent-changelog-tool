// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package github

import (
	"errors"
	"fmt"
)

// EnsureAuthConfigured method ensures that GitHub auth token is available.
func EnsureAuthConfigured(tk AuthToken) (bool, error) {
	tkloc, err := TokenLocation()
	if err != nil {
		return false, fmt.Errorf("cannot determine token location: %w", err)
	}

	val, err := tk.AuthToken()
	if err != nil {
		return false, fmt.Errorf("GitHub authorization token is missing. Please use either environment variable %s or ~/%s: %w",
			envAuth, tkloc, err)
	}

	if val == "" {
		return false, errors.New("GitHub authorization token is empty. Make sure a value is provided")
	}

	return true, nil
}
