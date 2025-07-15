package parsers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oss-health/background-worker/internal/repo/parsers"
)

func TestRequirementsParser_CanParse(t *testing.T) {
	p := parsers.RequirementsParser{}

	assert.True(t, p.CanParse("requirements.txt"))
	assert.True(t, p.CanParse("REQUIREMENTS.TXT")) // case-insensitive
	assert.False(t, p.CanParse("setup.py"))
	assert.False(t, p.CanParse("pyproject.toml"))
}

func TestRequirementsParser_Parse_BasicDependencies(t *testing.T) {
	p := parsers.RequirementsParser{}

	data := `
requests==2.31.0
flask==2.3.3
`

	deps, err := p.Parse([]byte(data))
	assert.NoError(t, err)
	assert.ElementsMatch(t, []parsers.Dependency{
		{Name: "requests", Version: "2.31.0", Ecosystem: "pypi"},
		{Name: "flask", Version: "2.3.3", Ecosystem: "pypi"},
	}, deps)
}

func TestRequirementsParser_Parse_IgnoresCommentsAndBlankLines(t *testing.T) {
	p := parsers.RequirementsParser{}

	data := `
# This is a comment
requests==2.31.0

# Another comment
flask==2.3.3
`

	deps, err := p.Parse([]byte(data))
	assert.NoError(t, err)
	assert.ElementsMatch(t, []parsers.Dependency{
		{Name: "requests", Version: "2.31.0", Ecosystem: "pypi"},
		{Name: "flask", Version: "2.3.3", Ecosystem: "pypi"},
	}, deps)
}

func TestRequirementsParser_Parse_InvalidFormatLines(t *testing.T) {
	p := parsers.RequirementsParser{}

	data := `
requests>=2.31.0
flask
numpy==1.25.0
`

	deps, err := p.Parse([]byte(data))
	assert.NoError(t, err)
	assert.Equal(t, []parsers.Dependency{
		{Name: "numpy", Version: "1.25.0", Ecosystem: "pypi"},
	}, deps)
}

func TestRequirementsParser_Parse_ExtraSpaces(t *testing.T) {
	p := parsers.RequirementsParser{}

	data := `
  requests  ==  2.31.0
  flask==  2.3.3 
`

	deps, err := p.Parse([]byte(data))
	assert.NoError(t, err)
	assert.ElementsMatch(t, []parsers.Dependency{
		{Name: "requests", Version: "2.31.0", Ecosystem: "pypi"},
		{Name: "flask", Version: "2.3.3", Ecosystem: "pypi"},
	}, deps)
}
