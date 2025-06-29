package parsers

import (
	"path/filepath"

	"github.com/oss-health/background-worker/internal/dependency"
)

type Parser struct {
	Pattern   string
	Ecosystem string
	Parse     ParserFunc
}

type ParserFunc func(path string) ([]dependency.DependencyVersionPair, error)

var DependencyParsers []Parser

func RegisterParser(pattern, ecosystem string, fn ParserFunc) {
	DependencyParsers = append(DependencyParsers, Parser{
		Pattern:   pattern,
		Ecosystem: ecosystem,
		Parse:     fn,
	})
}

func GetParserForFile(path string) *Parser {
	base := filepath.Base(path)

	for _, p := range DependencyParsers {
		if p.Pattern == base {
			return &p
		}
	}
	return nil
}
