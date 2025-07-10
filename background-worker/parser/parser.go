package parser


type Dependency struct {
	Name    string
	Version string
	Ecosystem string
}

type Parser interface {
	Parse(fileContent []byte) ([]Dependency, error)
	CanParse(filename string) bool
}
