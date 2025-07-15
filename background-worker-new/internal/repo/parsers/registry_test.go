package parsers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oss-health/background-worker/internal/repo/parsers"
)

func TestParserRegistry_GetParsers_RequirementsTxt(t *testing.T) {
	registry := parsers.NewParserRegistry()

	p, ok := registry.GetParsers("requirements.txt")
	assert.True(t, ok)
	assert.Len(t, p, 1)
	assert.IsType(t, parsers.RequirementsParser{}, p[0])
}

func TestParserRegistry_GetParsers_PackageJson(t *testing.T) {
	registry := parsers.NewParserRegistry()

	p, ok := registry.GetParsers("package.json")
	assert.True(t, ok)
	assert.Len(t, p, 1)
	assert.IsType(t, parsers.NpmParser{}, p[0])
}

func TestParserRegistry_GetParsers_PyProjectToml(t *testing.T) {
	registry := parsers.NewParserRegistry()

	p, ok := registry.GetParsers("pyproject.toml")
	assert.True(t, ok)
	assert.Len(t, p, 2)

	foundPoetry := false
	foundPep621 := false

	for _, p := range p {
		switch p.(type) {
		case parsers.PoetryPyProjectParser:
			foundPoetry = true
		case parsers.Pep621PyProjectParser:
			foundPep621 = true
		}
	}

	assert.True(t, foundPoetry)
	assert.True(t, foundPep621)
}

func TestParserRegistry_GetParsers_UnsupportedFile(t *testing.T) {
	registry := parsers.NewParserRegistry()

	parsers, ok := registry.GetParsers("setup.cfg")
	assert.False(t, ok)
	assert.Empty(t, parsers)
}
