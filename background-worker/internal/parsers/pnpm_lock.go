package parsers

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func init() {
	RegisterParser("pnpm-lock.yaml", "npm", ParsePnpmLock)
}

func ParsePnpmLock(path string) ([]Dependency, error) {
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

	if err := yaml.NewDecoder(file).Decode(&raw); err != nil {
		return nil, err
	}

	packages, ok := raw["packages"].(map[string]any)

	if !ok {
		return nil, nil
	}

	deps := []Dependency{}

	for pkgRef := range packages {
		ref := pkgRef

		if strings.Contains(ref, "node_modules") {
			continue
		}

		ref = strings.TrimPrefix(ref, "/")

		if strings.Count(ref, "@") >= 1 {
			parts := strings.Split(ref, "@")

			if len(parts) >= 2 {
				version := parts[len(parts)-1]
				name := strings.Join(parts[:len(parts)-1], "@")

				if name != "" {
					deps = append(deps, Dependency{
						Name:      name,
						Version:   version,
						Ecosystem: "npm",
					})
				}
			}
		}
	}

	return deps, nil
}
