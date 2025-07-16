package db

import "time"

type Repository struct {
	ID            int
	URL           string
	GithubId      int
	LastScannedAt *time.Time
	ScanStatus    string
}

type Dependency struct {
	ID        int
	Name      string
	Ecosystem string
}

type DependencyVersion struct {
	ID                  int
	Version             string
	DependencyId        int
	DependencyName      string
	DependencyEcosystem string
}
