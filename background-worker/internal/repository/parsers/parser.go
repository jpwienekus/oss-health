package parsers

import (
	"github.com/oss-health/background-worker/internal/dependency"
)

type ParserFunc func(path string) ([]dependency.DependencyVersionPair, error)

type RegisteredParser struct {
	Pattern string
	Parser  ParserFunc
}

var registeredParsers []RegisteredParser

func RegisterParser(pattern, ecosystem string, fn ParserFunc) {
	registeredParsers = append(registeredParsers, RegisteredParser{
		Pattern: pattern,
		Parser:  fn,
	})
}

func GetRegisteredParsers() []RegisteredParser {
	return registeredParsers
}
