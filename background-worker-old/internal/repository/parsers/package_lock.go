package parsers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/oss-health/background-worker/internal/dependency"
)

func init() {
	RegisterParser("package-lock.json", "npm", ParsePackageLock)
}

func ParsePackageLock(path string) ([]dependency.DependencyVersionPair, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open package-lock.json at %q: %w", path, err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close file %q: %v\n", path, cerr)
		}
	}()

	var raw map[string]any
	if err := json.NewDecoder(file).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode JSON in %q: %w", path, err)
	}

	deps := []dependency.DependencyVersionPair{}
	dependencies, ok := raw["dependencies"].(map[string]any)
	if !ok {
		// not an error, just empty dependencies section
		return deps, nil
	}

	for name, val := range dependencies {
		info, ok := val.(map[string]any)
		if !ok {
			continue // skip invalid entries
		}

		version := "unknown"
		if v, ok := info["version"].(string); ok {
			version = v
		}

		deps = append(deps, dependency.DependencyVersionPair{
			Name:      name,
			Version:   version,
			Ecosystem: "npm",
		})
	}

	return deps, nil
}
