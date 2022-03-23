package changelog_test

import (
	"path"
	"testing"

	"github.com/elastic/elastic-agent-changelog-tool/internal/changelog"
	"github.com/elastic/elastic-agent-changelog-tool/internal/changelog/fragment"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func loadEntries(t *testing.T, fixture string) (fragment.File, changelog.Entry) {
	fs := afero.NewCopyOnWriteFs(afero.NewOsFs(), afero.NewMemMapFs())

	viper.Set("fragment_location", "testdata")

	f1, err := fragment.Load(fs, path.Join(viper.GetString("fragment_location"), fixture))
	require.Nil(t, err)

	return f1, changelog.EntryFromFragment(f1)
}

func TestEntriesFromFragment_breaking(t *testing.T) {
	fixture := "1648040928-breaking-change.yaml"
	f, e := loadEntries(t, fixture)
	require.Equal(t, e.File.Name, fixture)
	require.Equal(t, e.File.Checksum, f.Checksum())
	require.Equal(t, e.Timestamp, f.Timestamp)
	require.Equal(t, e.Kind, changelog.BreakingChange)
	require.Equal(t, e.Summary, f.Fragment.Summary)
}

func TestEntriesFromFragment_deprecation(t *testing.T) {
	fixture := "1648040928-bug-fix.yaml"
	f, e := loadEntries(t, fixture)
	require.Equal(t, e.File.Name, fixture)
	require.Equal(t, e.File.Checksum, f.Checksum())
	require.Equal(t, e.Timestamp, f.Timestamp)
	require.Equal(t, e.Kind, changelog.BugFix)
	require.Equal(t, e.Summary, f.Fragment.Summary)
}

func TestEntriesFromFragment_bugfix(t *testing.T) {
	fixture := "1648040928-deprecation.yaml"
	f, e := loadEntries(t, fixture)
	require.Equal(t, e.File.Name, fixture)
	require.Equal(t, e.File.Checksum, f.Checksum())
	require.Equal(t, e.Timestamp, f.Timestamp)
	require.Equal(t, e.Kind, changelog.Deprecation)
	require.Equal(t, e.Summary, f.Fragment.Summary)
}

func TestEntriesFromFragment_enhancement(t *testing.T) {
	fixture := "1648040928-enhancement.yaml"
	f, e := loadEntries(t, fixture)
	require.Equal(t, e.File.Name, fixture)
	require.Equal(t, e.File.Checksum, f.Checksum())
	require.Equal(t, e.Timestamp, f.Timestamp)
	require.Equal(t, e.Kind, changelog.Enhancement)
	require.Equal(t, e.Summary, f.Fragment.Summary)
}

func TestEntriesFromFragment_feature(t *testing.T) {
	fixture := "1648040928-feature.yaml"
	f, e := loadEntries(t, fixture)
	require.Equal(t, e.File.Name, fixture)
	require.Equal(t, e.File.Checksum, f.Checksum())
	require.Equal(t, e.Timestamp, f.Timestamp)
	require.Equal(t, e.Kind, changelog.Feature)
	require.Equal(t, e.Summary, f.Fragment.Summary)
}

func TestEntriesFromFragment_knownissue(t *testing.T) {
	fixture := "1648040928-known-issue.yaml"
	f, e := loadEntries(t, fixture)
	require.Equal(t, e.File.Name, fixture)
	require.Equal(t, e.File.Checksum, f.Checksum())
	require.Equal(t, e.Timestamp, f.Timestamp)
	require.Equal(t, e.Kind, changelog.KnownIssue)
	require.Equal(t, e.Summary, f.Fragment.Summary)

}
func TestEntriesFromFragment_security(t *testing.T) {
	fixture := "1648040928-security.yaml"
	f, e := loadEntries(t, fixture)
	require.Equal(t, e.File.Name, fixture)
	require.Equal(t, e.File.Checksum, f.Checksum())
	require.Equal(t, e.Timestamp, f.Timestamp)
	require.Equal(t, e.Kind, changelog.Security)
	require.Equal(t, e.Summary, f.Fragment.Summary)

}
func TestEntriesFromFragment_upgrade(t *testing.T) {
	fixture := "1648040928-upgrade.yaml"
	f, e := loadEntries(t, fixture)
	require.Equal(t, e.File.Name, fixture)
	require.Equal(t, e.File.Checksum, f.Checksum())
	require.Equal(t, e.Timestamp, f.Timestamp)
	require.Equal(t, e.Kind, changelog.Upgrade)
	require.Equal(t, e.Summary, f.Fragment.Summary)

}
func TestEntriesFromFragment_other(t *testing.T) {
	fixture := "1648040928-other.yaml"
	f, e := loadEntries(t, fixture)
	require.Equal(t, e.File.Name, fixture)
	require.Equal(t, e.File.Checksum, f.Checksum())
	require.Equal(t, e.Timestamp, f.Timestamp)
	require.Equal(t, e.Kind, changelog.Other)
	require.Equal(t, e.Summary, f.Fragment.Summary)
}

func TestEntriesFromFragment_unknown(t *testing.T) {
	fixture := "1648040928-unknown.yaml"
	f, e := loadEntries(t, fixture)
	require.Equal(t, e.File.Name, fixture)
	require.Equal(t, e.File.Checksum, f.Checksum())
	require.Equal(t, e.Timestamp, f.Timestamp)
	require.Equal(t, e.Kind, changelog.Unknown)
	require.Equal(t, e.Summary, f.Fragment.Summary)
}
