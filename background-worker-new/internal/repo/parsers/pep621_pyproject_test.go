package parsers_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oss-health/background-worker/internal/repo/parsers"
)

func TestPep621PyProjectParser_CanParse(t *testing.T) {
	p := parsers.Pep621PyProjectParser{}

	assert.True(t, p.CanParse("pyproject.toml"))
	assert.False(t, p.CanParse("setup.py"))
	assert.False(t, p.CanParse("requirements.txt"))
}

func TestPep621PyProjectParser_Parse_WithDependencies(t *testing.T) {
	p := parsers.Pep621PyProjectParser{}

	tomlData := `
[project]
dependencies = [
  "requests (>=2.0,<3.0)",
  "flask",
  "numpy (^1.25.0)"
]`

	deps, err := p.Parse([]byte(tomlData))
	assert.NoError(t, err)
	assert.ElementsMatch(t, []parsers.Dependency{
		{Name: "requests", Version: "2.0.0", Ecosystem: "pypi"},
		{Name: "flask", Version: "unknown", Ecosystem: "pypi"},
		{Name: "numpy", Version: "1.25.0", Ecosystem: "pypi"},
	}, deps)
}

func TestPep621PyProjectParser_Parse_NoProjectSection(t *testing.T) {
	p := parsers.Pep621PyProjectParser{}

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
	p := parsers.Pep621PyProjectParser{}

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
	p := parsers.Pep621PyProjectParser{}

	tomlData := `
[project]
dependencies = [123, true, "pandas (>=1.3.0)"]
`

	deps, err := p.Parse([]byte(tomlData))
	assert.NoError(t, err)
	assert.Equal(t, []parsers.Dependency{
		{Name: "pandas", Version: "1.3.0", Ecosystem: "pypi"},
	}, deps)
}

func TestPep621PyProjectParser_Parse_ComplexVersions(t *testing.T) {
	p := parsers.Pep621PyProjectParser{}

	tomlData := `
[project]
dependencies = [
  "example (>=0.1.0,!=0.1.5,<1.0.0)",
  "toolkit"
]
`

	deps, err := p.Parse([]byte(tomlData))
	assert.NoError(t, err)

	assert.ElementsMatch(t, []parsers.Dependency{
		{Name: "example", Version: "0.1.0", Ecosystem: "pypi"},
		{Name: "toolkit", Version: "unknown", Ecosystem: "pypi"},
	}, deps)
}

func TestPep621PyProjectParser_Parse_FieldsAreLowercased(t *testing.T) {
	p := parsers.Pep621PyProjectParser{}

	tomlData := `
[project]
dependencies = [
  "Requests (>=2.28.0)",
  "NumPy (^1.26.0)",
  "Flask"
]
`

	deps, err := p.Parse([]byte(tomlData))
	assert.NoError(t, err)

	for _, dep := range deps {
		assert.Equal(t, dep.Name, strings.ToLower(dep.Name), "name should be lowercase")
		assert.Equal(t, dep.Version, strings.ToLower(dep.Version), "version should be lowercase")
		assert.Equal(t, dep.Ecosystem, strings.ToLower(dep.Ecosystem), "ecosystem should be lowercase")
	}
}
