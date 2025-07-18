package parsers_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oss-health/background-worker/internal/repo/parsers"
)

func TestPoetryPyProjectParser_CanParse(t *testing.T) {
	p := parsers.PoetryPyProjectParser{}

	assert.True(t, p.CanParse("pyproject.toml"))
	assert.False(t, p.CanParse("setup.cfg"))
	assert.False(t, p.CanParse("package.json"))
}

func TestPoetryPyProjectParser_Parse_MainDependenciesOnly(t *testing.T) {
	p := parsers.PoetryPyProjectParser{}

	tomlData := `
[tool.poetry]
name = "example"
version = "0.1.0"

[tool.poetry.dependencies]
python = "^3.10"
requests = "^2.31.0"
`

	deps, err := p.Parse([]byte(tomlData))
	assert.NoError(t, err)
	assert.ElementsMatch(t, []parsers.Dependency{
		{Name: "python", Version: "3.10.0", Ecosystem: "pypi"},
		{Name: "requests", Version: "2.31.0", Ecosystem: "pypi"},
	}, deps)
}

func TestPoetryPyProjectParser_Parse_WithDevDependencies(t *testing.T) {
	p := parsers.PoetryPyProjectParser{}

	tomlData := `
[tool.poetry]
name = "example"
version = "0.1.0"

[tool.poetry.dependencies]
flask = "^2.3.0"

[tool.poetry.group.dev.dependencies]
pytest = "^7.4.0"
black = "^23.7.0"
`

	deps, err := p.Parse([]byte(tomlData))
	assert.NoError(t, err)
	assert.ElementsMatch(t, []parsers.Dependency{
		{Name: "flask", Version: "2.3.0", Ecosystem: "pypi"},
		{Name: "pytest", Version: "7.4.0", Ecosystem: "pypi"},
		{Name: "black", Version: "23.7.0", Ecosystem: "pypi"},
	}, deps)
}

func TestPoetryPyProjectParser_Parse_NoPoetrySection(t *testing.T) {
	p := parsers.PoetryPyProjectParser{}

	tomlData := `
[tool.setuptools]
name = "example"
version = "0.1.0"
`

	deps, err := p.Parse([]byte(tomlData))
	assert.NoError(t, err)
	assert.Nil(t, deps)
}

func TestPoetryPyProjectParser_Parse_NoToolSection(t *testing.T) {
	p := parsers.PoetryPyProjectParser{}

	tomlData := `
[build-system]
requires = ["poetry-core"]
`

	deps, err := p.Parse([]byte(tomlData))
	assert.NoError(t, err)
	assert.Nil(t, deps)
}

func TestPoetryPyProjectParser_Parse_InvalidTOML(t *testing.T) {
	p := parsers.PoetryPyProjectParser{}

	invalid := `invalid = = = toml`

	deps, err := p.Parse([]byte(invalid))
	assert.Error(t, err)
	assert.Nil(t, deps)
}

func TestPoetryPyProjectParser_Parse_FieldsAreLowercased(t *testing.T) {
	p := parsers.PoetryPyProjectParser{}

	tomlData := `
[tool.poetry]
name = "Example"
version = "0.1.0"

[tool.poetry.dependencies]
Flask = "^2.3.0"

[tool.poetry.group.Dev.dependencies]
PyTest = "^7.4.0"
BLACK = "^23.7.0"
`

	deps, err := p.Parse([]byte(tomlData))
	assert.NoError(t, err)

	for _, dep := range deps {
		assert.Equal(t, dep.Name, strings.ToLower(dep.Name), "name should be lowercase")
		assert.Equal(t, dep.Version, strings.ToLower(dep.Version), "version should be lowercase")
		assert.Equal(t, dep.Ecosystem, strings.ToLower(dep.Ecosystem), "ecosystem should be lowercase")
	}
}
