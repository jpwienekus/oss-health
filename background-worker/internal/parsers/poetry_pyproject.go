package parsers

import (
	"os"
	"regexp"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

func init() {
	RegisterParser("pyproject.toml", "pypi", ParsePoetryPyproject)
}

func ParsePoetryPyproject(path string) ([]Dependency, error) {
	content, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var data map[string]any

	if err := toml.Unmarshal(content, &data); err != nil {
		return nil, err
	}

	getToolSection := func(keys ...string) map[string]any {
		section := data

		for _, key := range keys {
			if val, ok := section[key].(map[string]any); ok {
				section = val
			} else {
				return nil
			}
		}

		return section
	}

	deps := []Dependency{}
	mainDeps := getToolSection("tool", "poetry", "dependencies")
	devDeps := getToolSection("tool", "poetry", "group", "dev", "dependencies")

	for name, version := range mergeMaps(mainDeps, devDeps) {
		if strings.ToLower(name) == "python" {
			continue
		}

		name, versionStr := normalizeNameVersion(name, version)
		deps = append(deps, Dependency{name, versionStr, "PyPI"})
	}

	return deps, nil
}

func mergeMaps(maps ...map[string]any) map[string]interface{} {
	merged := make(map[string]any)
	for _, m := range maps {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}

func normalizeNameVersion(name string, version any) (string, string) {
	reBracket := regexp.MustCompile(`\[.*?\]`)
	rePrefix := regexp.MustCompile(`^[~^<>=!]+`)

	cleanName := strings.TrimSpace(reBracket.ReplaceAllString(name, ""))
	versionStr := "unknown"

	switch v := version.(type) {
	case string:
		versionStr = v
	case map[string]any:
		if ver, ok := v["version"].(string); ok {
			versionStr = ver
		}
	}

	return cleanName, strings.TrimSpace(rePrefix.ReplaceAllString(versionStr, ""))
}
