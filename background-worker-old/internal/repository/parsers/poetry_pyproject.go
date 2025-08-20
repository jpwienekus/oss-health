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
	RegisterParser("pyproject.toml", "pypi", ParsePoetryPyproject)
}

func ParsePoetryPyproject(path string) ([]dependency.DependencyVersionPair, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read pyproject.toml at %q: %w", path, err)
	}

	var data map[string]any
	if err := toml.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("unmarshal TOML in %q: %w", path, err)
	}

	getToolSection := func(keys ...string) map[string]any {
		section := data

		for _, key := range keys {
			val, ok := section[key].(map[string]any)
			if !ok {
				return nil
			}

			section = val
		}

		return section
	}

	mainDeps := getToolSection("tool", "poetry", "dependencies")
	devDeps := getToolSection("tool", "poetry", "group", "dev", "dependencies")

	mergedDeps := mergeMaps(mainDeps, devDeps)
	var deps []dependency.DependencyVersionPair

	for name, version := range mergedDeps {
		if strings.ToLower(name) == "python" {
			continue
		}

		name, versionStr := normalizeNameVersion(name, version)
		deps = append(deps, dependency.DependencyVersionPair{
			Name:      name,
			Version:   versionStr,
			Ecosystem: "PyPI",
		})
	}

	return deps, nil
}

func mergeMaps(maps ...map[string]any) map[string]any {
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
