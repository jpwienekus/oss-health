package models

type DependencyVersion struct {
	ID                  int
	Version             string
	DependencyId        int
	DependencyName      string
	DependencyEcosystem string
}
