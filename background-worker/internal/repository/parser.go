package repository

import (
	"path/filepath"
)

type DependencyParsed struct {
	Name      string
	Version   string
	Ecosystem string
}

type Parser struct {
	Pattern   string
	Ecosystem string
	Parse     ParserFunc
}

type ParserFunc func(path string) ([]DependencyParsed, error)

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
