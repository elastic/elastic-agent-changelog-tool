package changelog

import "github.com/elastic/elastic-agent-changelog-tool/internal/changelog/fragment"

type FragmentFileInfo struct {
	Name     string `yaml:"name"`
	Checksum string `yaml:"checksum"`
}

type Entry struct {
	Summary     string           `yaml:"summary"`
	Description string           `yaml:"description"`
	Kind        Kind             `yaml:"kind"`
	LinkedPR    int              `yaml:"pr"`
	LinkedIssue int              `yaml:"issue"`
	Timestamp   int64            `yaml:"timestamp"`
	File        FragmentFileInfo `yaml:"file"`
}

const (
	// NOTE: Kind values should be aligned with supported type from doc.elastic.co
	kindBreakingChange = "breaking-change"
	kindDeprecation    = "deprecation"
	kindBugfix         = "bug-fix"
	kindEnhancement    = "enhancement"
	kindFeature        = "feature"
	kindKnownIssue     = "known-issue"
	kindSecurity       = "security"
	kindUpgrade        = "upgrade"
	kindOther          = "other"
)

// EntriesFromFragment returns one or more entries based on the fragment File.
// A single Fragment can contain multiple Changelog entries.
func EntryFromFragment(f fragment.File) Entry {
	e := Entry{
		Summary:     f.Fragment.Summary,
		Description: f.Fragment.Description,
		Kind:        kind2kind(f),
		LinkedPR:    f.Fragment.Pr,
		LinkedIssue: f.Fragment.Issue,
		Timestamp:   f.Timestamp,
		File: FragmentFileInfo{
			Name:     f.Name,
			Checksum: f.Checksum(),
		},
	}

	return e
}

func kind2kind(f fragment.File) Kind {
	switch f.Fragment.Kind {
	case string(BreakingChange):
		return BreakingChange
	case string(BugFix):
		return BugFix
	case string(Deprecation):
		return Deprecation
	case string(Enhancement):
		return Enhancement
	case string(Feature):
		return Feature
	case string(KnownIssue):
		return KnownIssue
	case string(Security):
		return Security
	case string(Upgrade):
		return Upgrade
	case string(Other):
		return Other
	default:
		return Unknown
	}
}