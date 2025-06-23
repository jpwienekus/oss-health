package repository

import "time"

type Repository struct {
	ID            int
	URL           string
	LastScannedAt *time.Time
	ScanStatus    time.Time
}
