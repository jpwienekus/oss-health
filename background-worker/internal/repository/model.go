package repository

type Repository struct {
	ID            int
	URL           string
	LastScannedAt *string
	ScanStatus    string
}
