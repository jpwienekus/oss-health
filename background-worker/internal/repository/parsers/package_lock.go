package parsers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/oss-health/background-worker/internal/repository"
)

func init() {
	repository.RegisterParser("package-lock.json", "npm", ParsePackageLock)
}

func ParsePackageLock(path string) ([]repository.DependencyParsed, error) {
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

	deps := []repository.DependencyParsed{}

	if dependencies, ok := raw["dependencies"].(map[string]any); ok {
		for name, val := range dependencies {
			if info, ok := val.(map[string]any); ok {
				version := "unknown"

				if v, ok := info["version"].(string); ok {
					version = v
				}

				deps = append(deps, repository.DependencyParsed{
					Name:      name,
					Version:   version,
					Ecosystem: "npm",
				})
			}
		}
	}

	return deps, nil
}
