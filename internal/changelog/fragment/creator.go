// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package fragment

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

// timestamp represent a function providing a timestamp.
// It's used to allow replacing the value with a known one during testing.
type timestamp func() int64

func NewCreator(fs afero.Fs) FragmentCreator {
	return FragmentCreator{
		fs:        fs,
		timestamp: time.Now().Unix,
	}
}

// TestNewCreator sets up a FragmentCreator configured to be used in testing.
func TestNewCreator() FragmentCreator {
	f := NewCreator(afero.NewMemMapFs())
	f.timestamp = func() int64 { return 1647345675 }
	return f
}

type FragmentCreator struct {
	fs afero.Fs
	// timestamp allow overriding value in tests
	timestamp timestamp
}

// filename computes the filename for the changelog fragment to be created.
// To provide unique names the provided slug is prepended with current timestamp.
func (f FragmentCreator) filename(slug string) string {
	filename := fmt.Sprintf("%d-%s.yaml", f.timestamp(), sanitizeFilename(slug))
	return filename
}

// Create marshal changelog fragment and persist it to file.
func (c FragmentCreator) Create(location, slug string) error {
	frg := Fragment{}

	data, err := yaml.Marshal(&frg)
	if err != nil {
		return err
	}

	return afero.WriteFile(c.fs, path.Join(location, c.filename(slug)), data, 0660)
}

// sanitizeFilename takes care of removing dangerous elements from a string so it can be safely
// used as a filename.
// NOTE: does not prevent command injection or ensure complete escaping of input
func sanitizeFilename(s string) string {
	s = strings.Replace(s, " ", "-", -1)
	s = strings.Replace(s, "/", "-", -1)
	s = strings.Replace(s, "\\", "-", -1)
	return s
}
