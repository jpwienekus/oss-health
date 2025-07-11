package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oss-health/background-worker/parser"
)

func TestNpmParser_CanParse(t *testing.T) {
	parser := parser.NpmParser{}

	assert.True(t, parser.CanParse("package.json"))
	assert.True(t, parser.CanParse("PACKAGE.JSON")) // case-insensitive
	assert.False(t, parser.CanParse("package-lock.json"))
	assert.False(t, parser.CanParse("requirements.txt"))
}

func TestNpmParser_Parse_WithDependencies(t *testing.T) {
	p := parser.NpmParser{}

	json := `{
		"dependencies": {
			"react": "^18.2.0",
			"lodash": "4.17.21"
		}
	}`

	deps, err := p.Parse([]byte(json))
	assert.NoError(t, err)
	assert.ElementsMatch(t, []parser.Dependency{
		{Name: "react", Version: "^18.2.0", Ecosystem: "npm"},
		{Name: "lodash", Version: "4.17.21", Ecosystem: "npm"},
	}, deps)
}

func TestNpmParser_Parse_WithDevDependencies(t *testing.T) {
	p := parser.NpmParser{}

	json := `{
		"devDependencies": {
			"jest": "^29.0.0"
		}
	}`

	deps, err := p.Parse([]byte(json))
	assert.NoError(t, err)
	assert.Equal(t, []parser.Dependency{
		{Name: "jest", Version: "^29.0.0", Ecosystem: "npm"},
	}, deps)
}

func TestNpmParser_Parse_WithBothDeps(t *testing.T) {
	p := parser.NpmParser{}

	json := `{
		"dependencies": {
			"express": "^4.18.1"
		},
		"devDependencies": {
			"mocha": "^10.0.0"
		}
	}`

	deps, err := p.Parse([]byte(json))
	assert.NoError(t, err)
	assert.ElementsMatch(t, []parser.Dependency{
		{Name: "express", Version: "^4.18.1", Ecosystem: "npm"},
		{Name: "mocha", Version: "^10.0.0", Ecosystem: "npm"},
	}, deps)
}

func TestNpmParser_Parse_Empty(t *testing.T) {
	p := parser.NpmParser{}

	json := `{}`

	deps, err := p.Parse([]byte(json))
	assert.NoError(t, err)
	assert.Empty(t, deps)
}

func TestNpmParser_Parse_InvalidJSON(t *testing.T) {
	p := parser.NpmParser{}

	json := `{"dependencies": "not-a-map"}`

	_, err := p.Parse([]byte(json))
	assert.Error(t, err)
}
