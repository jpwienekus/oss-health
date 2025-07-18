package parsers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
)

type NpmParser struct{}

func (n NpmParser) CanParse(filename string) bool {
	return strings.EqualFold(filename, "package.json")
}

func (n NpmParser) Parse(fileContent []byte) ([]Dependency, error) {
	var jsonData struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	if err := json.Unmarshal(fileContent, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	toDependencies := func(m map[string]string) []Dependency {
		dependencies := make([]Dependency, 0, len(m))
		for name, version := range m {
			cleanedVersion := strings.TrimLeft(version, "^~><= ")
			parsedVersion, err := semver.NewVersion(cleanedVersion)
			if err != nil {
				cleanedVersion = "unknown"
			} else {
				cleanedVersion = parsedVersion.String()
			}

			dependencies = append(dependencies, Dependency{
				Name:      strings.ToLower(name),
				Version:   cleanedVersion,
				Ecosystem: "npm",
			})
		}
		return dependencies
	}

	dependencies := append(
		toDependencies(jsonData.Dependencies),
		toDependencies(jsonData.DevDependencies)...,
	)

	return dependencies, nil
}
