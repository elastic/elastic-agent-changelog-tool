// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package changelog

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidators(t *testing.T) {
	testcases := []struct {
		name          string
		entry         Entry
		validatorFunc func(entry Entry) error
		expectedErr   error
	}{
		// PRMultipleIDs
		{
			"pr multiple ids: 1 id",
			Entry{
				LinkedPR: []string{"1"},
			},
			validator_PRMultipleIDs,
			nil,
		},
		{
			"pr multiple ids: multiple ids",
			Entry{
				LinkedPR: []string{"1", "2"},
			},
			validator_PRMultipleIDs,
			fmt.Errorf("changelog entry: %s has multiple PR ids", ""),
		},
		// PRnoIDs
		{
			"pr multiple ids: error",
			Entry{
				LinkedPR: []string{},
			},
			validator_PRnoIDs,
			fmt.Errorf("changelog entry: %s has no PR id", ""),
		},
		// IssueNoIDs
		{
			"issue no ids: error",
			Entry{
				LinkedIssue: []string{},
			},
			validator_IssueNoIDs,
			fmt.Errorf("changelog entry: %s has no issue id", ""),
		},
		// ComponentValid
		{
			"component valid: beats",
			Entry{
				Component: "beats",
			},
			validator_componentValid([]string{"beats"}),
			nil,
		},
		{
			"component valid: not found in config",
			Entry{
				Component: "agent",
			},
			validator_componentValid([]string{"beats"}),
			fmt.Errorf("changelog entry: %s -> component [%s] not found in config: [%s]", "", "agent", "beats"),
		},
		{
			"component valid: no component",
			Entry{
				Component: "",
			},
			validator_componentValid([]string{"beats", "agent"}),
			fmt.Errorf("changelog entry: %s -> component cannot be assumed, choose it from config: %s", "", []string{"beats", "agent"}),
		},
		{
			"component valid: invalid component",
			Entry{
				Component: "invalid_component",
			},
			validator_componentValid([]string{"beats"}),
			fmt.Errorf("changelog entry: %s -> component [%s] not found in config: [%s]", "", "invalid_component", "beats"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.validatorFunc(tc.entry)
			require.Equal(t, err, tc.expectedErr)
		})
	}
}
