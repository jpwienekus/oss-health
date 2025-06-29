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
		return nil, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close file: %v\n", err)
		}
	}()

	var raw map[string]any

	if err := json.NewDecoder(file).Decode(&raw); err != nil {
		return nil, err
	}

	deps := []dependency.DependencyVersionPair{}

	if dependencies, ok := raw["dependencies"].(map[string]any); ok {
		for name, val := range dependencies {
			if info, ok := val.(map[string]any); ok {
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
		}
	}

	return deps, nil
}
