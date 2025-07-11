package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oss-health/background-worker/parser"
)

func TestParserRegistry_GetParsers_RequirementsTxt(t *testing.T) {
	registry := parser.NewParserRegistry()

	parsers, ok := registry.GetParsers("requirements.txt")
	assert.True(t, ok)
	assert.Len(t, parsers, 1)
	assert.IsType(t, parser.RequirementsParser{}, parsers[0])
}

func TestParserRegistry_GetParsers_PackageJson(t *testing.T) {
	registry := parser.NewParserRegistry()

	parsers, ok := registry.GetParsers("package.json")
	assert.True(t, ok)
	assert.Len(t, parsers, 1)
	assert.IsType(t, parser.NpmParser{}, parsers[0])
}

func TestParserRegistry_GetParsers_PyProjectToml(t *testing.T) {
	registry := parser.NewParserRegistry()

	parsers, ok := registry.GetParsers("pyproject.toml")
	assert.True(t, ok)
	assert.Len(t, parsers, 2)

	foundPoetry := false
	foundPep621 := false

	for _, p := range parsers {
		switch p.(type) {
		case parser.PoetryPyProjectParser:
			foundPoetry = true
		case parser.Pep621PyProjectParser:
			foundPep621 = true
		}
	}

	assert.True(t, foundPoetry)
	assert.True(t, foundPep621)
}

func TestParserRegistry_GetParsers_UnsupportedFile(t *testing.T) {
	registry := parser.NewParserRegistry()

	parsers, ok := registry.GetParsers("setup.cfg")
	assert.False(t, ok)
	assert.Empty(t, parsers)
}
