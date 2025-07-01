package parsers

import (
	"log"

	"github.com/oss-health/background-worker/internal/dependency"
)

type ParserFunc func(path string) ([]dependency.DependencyVersionPair, error)

type RegisteredParser struct {
	Pattern string
	Parser  ParserFunc
}

var registeredParsers []RegisteredParser

func RegisterParser(pattern, ecosystem string, fn ParserFunc) {
	log.Printf("registering %s", ecosystem)
	registeredParsers = append(registeredParsers, RegisteredParser{
		Pattern: pattern,
		Parser:  fn,
	})
}

func GetRegisteredParsers() []RegisteredParser {
	return registeredParsers
}
