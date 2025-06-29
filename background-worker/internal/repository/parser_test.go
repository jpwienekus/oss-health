package repository_test

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/oss-health/background-worker/internal/repository"
	"github.com/oss-health/background-worker/internal/repository/parsers"
)

func writeTempFile(t *testing.T, filename, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	fullPath := filepath.Join(tmpDir, filename)
	err := os.WriteFile(fullPath, []byte(content), 0644)

	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	return fullPath
}

func findDep(deps []repository.DependencyParsed, name, version, ecosystem string) bool {
	for _, dep := range deps {
		if dep.Name == name && dep.Version == version && dep.Ecosystem == ecosystem {
			return true
		}
	}
	return false
}

func TestParsePackageLock(t *testing.T) {
	path := writeTempFile(t, "package-lock.json", `{"dependencies": {"lodash": {"version": "4.17.21"}}}`)

	deps, err := parsers.ParsePackageLock(path)
	log.Print(deps)

	if err != nil {
		t.Fatal(err)
	}

	if len(deps) != 1 {
		t.Fatalf("Expected 1 dep, got %d", len(deps))
	}

	dep := deps[0]

	if dep.Name != "lodash" || dep.Version != "4.17.21" || dep.Ecosystem != "npm" {
		t.Errorf("Unexpected result: %+v", dep)
	}
}

func TestParsePnpmLock(t *testing.T) {
	content := `
packages:
  axios@1.0.0:
    resolution: {integrity: sha512}
  node_modules/ignored@1.0.0:
    resolution: {integrity: sha512}
`
	path := writeTempFile(t, "pnpm-lock.yaml", content)

	deps, err := parsers.ParsePnpmLock(path)

	if err != nil {
		t.Fatal(err)
	}

	if len(deps) != 1 {
		t.Fatalf("Expected 1 dep, got %d", len(deps))
	}

	if deps[0].Name != "axios" || deps[0].Version != "1.0.0" || deps[0].Ecosystem != "npm" {
		t.Errorf("Unexpected result: %+v", deps[0])
	}
}

func TestParsePep621Pyproject(t *testing.T) {
	content := `
[project]
name = "example"
version = "0.1.0"
dependencies = [
    "requests (>=2.25)",
    "httpx (==0.27.0)",
    "custom-lib",
    "fastapi[standard] (==0.115.13)",
]
`
	path := writeTempFile(t, "pyproject.toml", content)
	deps, err := parsers.ParsePep621Pyproject(path)

	if err != nil {
		t.Fatal(err)
	}

	cases := []repository.DependencyParsed{
		{"requests", "2.25", "PyPI"},
		{"httpx", "0.27.0", "PyPI"},
		{"custom-lib", "unknown", "PyPI"},
		{"fastapi", "0.115.13", "PyPI"},
	}

	for _, expected := range cases {
		if !findDep(deps, expected.Name, expected.Version, expected.Ecosystem) {
			t.Errorf("Missing dependency: %+v", expected)
		}
	}
}

func TestParsePoetryPyproject(t *testing.T) {
	content := `
[tool.poetry]
name = "example"
version = "0.1.0"

[tool.poetry.dependencies]
python = "^3.10"
requests = "^2.31.0"
httpx = { version = "^0.27.0", extras = ["http2"] }
custom = {}

[tool.poetry.group.dev.dependencies]
pytest = "^8.0.0"
black = { version = "^24.3.0" }
mypy = { some_other_field = "irrelevant" }
`
	path := writeTempFile(t, "pyproject.toml", content)
	deps, err := parsers.ParsePoetryPyproject(path)

	if err != nil {
		t.Fatal(err)
	}

	expected := []repository.DependencyParsed{
		{"requests", "2.31.0", "PyPI"},
		{"httpx", "0.27.0", "PyPI"},
		{"custom", "unknown", "PyPI"},
		{"pytest", "8.0.0", "PyPI"},
		{"black", "24.3.0", "PyPI"},
		{"mypy", "unknown", "PyPI"},
	}

	for _, exp := range expected {
		if !findDep(deps, exp.Name, exp.Version, exp.Ecosystem) {
			t.Errorf("Missing dependency: %+v", exp)
		}
	}

	for _, dep := range deps {
		if strings.ToLower(dep.Name) == "python" {
			t.Error("Should not include 'python' dependency")
		}
	}
}

func TestParseRequirementsTxt(t *testing.T) {
	content := `
requests==2.25.1
numpy
# comment
`
	path := writeTempFile(t, "requirements.txt", content)
	deps, err := parsers.ParseRequirementsTxt(path)

	if err != nil {
		t.Fatal(err)
	}

	expected := []repository.DependencyParsed{
		{"requests", "2.25.1", "PyPI"},
		{"numpy", "unknown", "PyPI"},
	}

	if len(deps) != 2 {
		t.Fatalf("Expected 2 dependencies, got %d", len(deps))
	}

	for _, exp := range expected {
		if !findDep(deps, exp.Name, exp.Version, exp.Ecosystem) {
			t.Errorf("Missing dependency: %+v", exp)
		}
	}
}
