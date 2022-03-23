package changelog

type Changelog struct {
	Version string  `yaml:"version"`
	Entries []Entry `yaml:"entries"`
}
