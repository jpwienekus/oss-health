package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oss-health/background-worker/parser"
)

func TestPep621PyProjectParser_CanParse(t *testing.T) {
	p := parser.Pep621PyProjectParser{}

	assert.True(t, p.CanParse("pyproject.toml"))
	assert.False(t, p.CanParse("setup.py"))
	assert.False(t, p.CanParse("requirements.txt"))
}

func TestPep621PyProjectParser_Parse_WithDependencies(t *testing.T) {
	p := parser.Pep621PyProjectParser{}

	tomlData := `
[project]
dependencies = [
  "requests (>=2.0,<3.0)",
  "flask",
  "numpy (^1.25.0)"
]`

	deps, err := p.Parse([]byte(tomlData))
	assert.NoError(t, err)
	assert.ElementsMatch(t, []parser.Dependency{
		{Name: "requests", Version: "2.0", Ecosystem: "pypi"},
		{Name: "flask", Version: "unknown", Ecosystem: "pypi"},
		{Name: "numpy", Version: "1.25.0", Ecosystem: "pypi"},
	}, deps)
}

func TestPep621PyProjectParser_Parse_NoProjectSection(t *testing.T) {
	p := parser.Pep621PyProjectParser{}

	tomlData := `
[tool.poetry]
name = "example"
version = "0.1.0"
`

	deps, err := p.Parse([]byte(tomlData))
	assert.NoError(t, err)
	assert.Nil(t, deps)
}

func TestPep621PyProjectParser_Parse_NoDependenciesSection(t *testing.T) {
	p := parser.Pep621PyProjectParser{}

	tomlData := `
[project]
name = "example"
version = "0.1.0"
`

	deps, err := p.Parse([]byte(tomlData))
	assert.NoError(t, err)
	assert.Nil(t, deps)
}

func TestPep621PyProjectParser_Parse_InvalidDependencyType(t *testing.T) {
	p := parser.Pep621PyProjectParser{}

	tomlData := `
[project]
dependencies = [123, true, "pandas (>=1.3.0)"]
`

	deps, err := p.Parse([]byte(tomlData))
	assert.NoError(t, err)
	assert.Equal(t, []parser.Dependency{
		{Name: "pandas", Version: "1.3.0", Ecosystem: "pypi"},
	}, deps)
}

func TestPep621PyProjectParser_Parse_ComplexVersions(t *testing.T) {
	p := parser.Pep621PyProjectParser{}

	tomlData := `
[project]
dependencies = [
  "example (>=0.1.0,!=0.1.5,<1.0.0)",
  "toolkit"
]
`

	deps, err := p.Parse([]byte(tomlData))
	assert.NoError(t, err)

	assert.ElementsMatch(t, []parser.Dependency{
		{Name: "example", Version: "0.1.0", Ecosystem: "pypi"},
		{Name: "toolkit", Version: "unknown", Ecosystem: "pypi"},
	}, deps)
}
