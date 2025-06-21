package db

import "time"

type Dependency struct {
	ID                 int64
	Name               string
	Ecosystem          string
	GithubURL          *string
	GithubURLResolved  bool
	GithubURLCheckedAt *time.Time
}
