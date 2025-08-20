package dependency

import "time"

type Dependency struct {
	ID                 int64
	Name               string
	Ecosystem          string
	GithubURL          *string
	GithubURLResolved  bool
	GithubURLCheckedAt *time.Time
}

type DependencyVersionPair struct {
	Name      string
	Version   string
	Ecosystem string
}

type DependencyVersionResult struct {
	VersionID int
	Name      string
	Version   string
	Ecosystem string
}
