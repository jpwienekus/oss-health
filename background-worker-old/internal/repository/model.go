package repository

import "time"

type Repository struct {
	ID            int
	URL           string
	GithubId      int
	LastScannedAt *time.Time
	ScanStatus    string
}
