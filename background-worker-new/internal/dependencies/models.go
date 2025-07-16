package dependencies

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
