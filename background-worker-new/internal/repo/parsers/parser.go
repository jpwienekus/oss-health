package parsers


// TODO: rename to DependencyVersion
type Dependency struct {
	Name    string
	Version string
	Ecosystem string
}

type Parser interface {
	Parse(fileContent []byte) ([]Dependency, error)
	CanParse(filename string) bool
}
