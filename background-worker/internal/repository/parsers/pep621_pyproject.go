package parsers

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/pelletier/go-toml/v2"

	"github.com/oss-health/background-worker/internal/dependency"
)

func init() {
	RegisterParser("pyproject.toml", "pypi", ParsePep621Pyproject)
}

func ParsePep621Pyproject(path string) ([]dependency.DependencyVersionPair, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read pyproject.toml at %q: %w", path, err)
	}

	var data map[string]any
	if err := toml.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("unmarshal TOML in %q: %w", path, err)
	}

	projectSection, ok := data["project"].(map[string]any)
	if !ok {
		// not an error, just missing 'project' section
		return nil, nil
	}

	rawDeps, ok := projectSection["dependencies"].([]any)
	if !ok {
		return nil, nil // no dependencies defined
	}

	var deps []dependency.DependencyVersionPair
	reBracket := regexp.MustCompile(`\[.*?\]`)
	rePrefix := regexp.MustCompile(`^[~^<>=!]+`)

	// project, _ := data["project"].(map[string]any)
	// rawDeps, _ := project["dependencies"].([]any)
	// deps := []dependency.DependencyVersionPair{}
	//
	// reBracket := regexp.MustCompile(`\[.*?\]`)
	// rePrefix := regexp.MustCompile(`^[~^<>=!]+`)

	for _, item := range rawDeps {
		depStr, ok := item.(string)
		if !ok {
			continue
		}

		name := depStr
		version := "unknown"

		if strings.Contains(depStr, " (") && strings.HasSuffix(depStr, ")") {
			parts := strings.SplitN(depStr[:len(depStr)-1], " (", 2)
			if len(parts) == 2 {
				name = parts[0]
				version = parts[1]
			}
		}

		name = strings.TrimSpace(reBracket.ReplaceAllString(name, ""))
		version = strings.TrimSpace(rePrefix.ReplaceAllString(version, ""))

		deps = append(deps, dependency.DependencyVersionPair{
			Name:      name,
			Version:   version,
			Ecosystem: "PyPI",
		})
	}

	return deps, nil
}
