// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package fragment

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/spf13/afero"
)

// timestamp represent a function providing a timestamp.
// It's used to allow replacing the value with a known one during testing.
type timestamp func() time.Time

func NewCreator(fs afero.Fs, location string) FragmentCreator {
	return FragmentCreator{
		fs:        fs,
		location:  location,
		timestamp: time.Now,
	}
}

// TestNewCreator sets up a FragmentCreator configured to be used in testing.
func TestNewCreator() FragmentCreator {
	f := NewCreator(afero.NewMemMapFs(), "testdata")
	tz, _ := time.LoadLocation("MST")
	f.timestamp = func() time.Time { return time.Date(2006, time.January, 2, 15, 4, 5, 0, tz) }
	return f
}

type FragmentCreator struct {
	fs       afero.Fs
	location string
	// timestamp allow overriding value in tests
	timestamp timestamp
}

func (c FragmentCreator) Location() string {
	return c.location
}

// filename computes the filename for the changelog fragment to be created.
// To provide unique names the provided slug is prepended with current timestamp.
func (f FragmentCreator) filename(slug string) string {
	filename := fmt.Sprintf("%d-%s.yaml", f.timestamp().Unix(), sanitizeFilename(slug))
	return filename
}

var fragmentLocPerm = os.FileMode(0770)
var fragmentPerm = os.FileMode(0660)

// Create marshal changelog fragment and persist it to file.
func (c FragmentCreator) Create(slug string) error {
	template, err := Template()
	if err != nil {
		return err
	}
	data := bytes.Replace(template, []byte("summary:"), []byte("summary: "+slug), 1)

	if err := c.fs.MkdirAll(c.location, fragmentLocPerm); err != nil {
		return fmt.Errorf("cannot create fragment location folder: %v", err)
	}

	filePath := path.Join(c.location, c.filename(slug))
	if err := afero.WriteFile(c.fs, filePath, data, fragmentPerm); err != nil {
		return err
	}

	log.Print("created fragment ", filePath)
	return nil
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
