// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package version

import (
	"strconv"
	"time"
)

var (
	// BuildTime is the build time of the binary (set externally with ldflags).
	BuildTime = "unknown"
	// CommitHash is the Git hash of the branch, used for version purposes (set externally with ldflags).
	CommitHash = "undefined"
	// Tag describes the semver version of the application (set externally with ldflags).
	Tag string
)

// BuildTimeFormatted method returns the build time in UTC preserving the RFC3339 format.
func BuildTimeFormatted() string {
	if BuildTime == "unknown" {
		return BuildTime
	}

	seconds, err := strconv.ParseInt(BuildTime, 10, 64)
	if err != nil {
		return "invalid"
	}

	// NOTE: time is returned in UTC to avoid timezone difference issues
	return time.Unix(seconds, 0).UTC().Format(time.RFC3339)
}
