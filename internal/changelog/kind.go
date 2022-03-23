package changelog

type Kind string

const (
	// NOTE: Kind values should be aligned with supported type from doc.elastic.co
	BreakingChange Kind = "breaking-change"
	Deprecation    Kind = "deprecation"
	BugFix         Kind = "bug-fix"
	Enhancement    Kind = "enhancement"
	Feature        Kind = "feature"
	KnownIssue     Kind = "known-issue"
	Security       Kind = "security"
	Upgrade        Kind = "upgrade"
	Other          Kind = "other"
	// Unknown kind is used if no matching kind is found
	Unknown Kind = "unknown"
)
