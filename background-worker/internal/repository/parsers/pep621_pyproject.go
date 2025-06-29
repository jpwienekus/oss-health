package parsers

import (
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
		return nil, err
	}

	var data map[string]any

	if err := toml.Unmarshal(content, &data); err != nil {
		return nil, err
	}

	project, _ := data["project"].(map[string]any)
	rawDeps, _ := project["dependencies"].([]any)
	deps := []dependency.DependencyVersionPair{}

	reBracket := regexp.MustCompile(`\[.*?\]`)
	rePrefix := regexp.MustCompile(`^[~^<>=!]+`)

	for _, item := range rawDeps {
		if depStr, ok := item.(string); ok {
			var name, version string

			if strings.Contains(depStr, " (") && strings.HasSuffix(depStr, ")") {
				parts := strings.SplitN(depStr[:len(depStr)-1], " (", 2)
				name = parts[0]
				version = parts[1]
			} else {
				name = depStr
				version = "unknown"
			}

			name = strings.TrimSpace(reBracket.ReplaceAllString(name, ""))
			version = strings.TrimSpace(rePrefix.ReplaceAllString(version, ""))
			deps = append(deps, dependency.DependencyVersionPair{
				Name:      name,
				Version:   version,
				Ecosystem: "PyPI",
			})
		}
	}

	return deps, nil
}
