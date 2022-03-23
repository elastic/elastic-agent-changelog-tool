package fragment

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

// File contains information of a Fragment file on disk.
type File struct {
	fs afero.Fs

	Content   []byte
	Fragment  Fragment
	Name      string
	Timestamp int64
	Title     string
}

// Load fills a File struct with data from a file on disk.
func Load(fs afero.Fs, path string) (File, error) {
	f := File{
		fs: fs,

		Name:     filepath.Base(path),
		Fragment: Fragment{},
	}

	ts, err := gettimestamp(f.Name)
	if err != nil {
		return f, fmt.Errorf("cannot extract timestamp from file name: %w", err)
	}

	f.Timestamp = ts
	f.Title = gettitle(f.Name)

	content, err := getcontent(f.fs, path)
	if err != nil {
		return f, err
	}

	f.Content = content

	if err := yaml.Unmarshal(content, &f.Fragment); err != nil {
		return f, fmt.Errorf("cannot unmarshal file: %w", err)
	}

	return f, nil
}

// Checksum computes SHA1 for file content and returns it as hex encoded string.
func (f File) Checksum() string {
	h := sha1.New()
	h.Write(f.Content)
	checksum := hex.EncodeToString([]byte(h.Sum(nil)))

	return checksum
}

// gettimestamp extracts timestamp from file name.
func gettimestamp(path string) (int64, error) {
	basename := filepath.Base(path)
	split := strings.Split(basename, "-")
	i, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot parse string to int: %w", err)
	}

	return int64(i), nil
}

// gettitle extracts fragment title from file name.
func gettitle(path string) string {
	basename := filepath.Base(path)
	split := strings.Split(basename, "-")
	title := strings.TrimSuffix(strings.Join(split[1:], "-"), filepath.Ext(basename))
	return title
}

// getcontent reads fragment file content.
func getcontent(fs afero.Fs, path string) ([]byte, error) {
	content, err := afero.ReadFile(fs, path)
	if err != nil {
		return []byte{}, fmt.Errorf("cannot read file: %w", err)
	}

	return content, nil
}
