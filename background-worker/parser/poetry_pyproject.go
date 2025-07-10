package parser

import (
	"fmt"

	"github.com/pelletier/go-toml/v2"
)

type PoetryPyProjectParser struct{}

func (p PoetryPyProjectParser) CanParse(filename string) bool {
	return filename == "pyproject.toml"
}

func (p PoetryPyProjectParser) Parse(fileContent []byte) ([]Dependency, error) {
	var dependencies []Dependency

	var data map[string]any
	if err := toml.Unmarshal(fileContent, &data); err != nil {
		return nil, fmt.Errorf("failed to parse pyproject.toml: %w", err)
	}

	toolSection, ok := data["tool"].(map[string]any)
	if !ok {
		return nil, nil // Skip if tool section is not found
	}

	poetrySection, ok := toolSection["poetry"].(map[string]any)
	groupSection, _ := poetrySection["group"].(map[string]any)
	devSection, _ := groupSection["dev"].(map[string]any)

	if !ok {
		return nil, nil // Skip if poetry section is not found
	}

	mainDependencies, _ := poetrySection["dependencies"].(map[string]any)
	devDependencies, _ := devSection["dependencies"].(map[string]any)

	for name, version := range mainDependencies {
		dependency := Dependency{
			Name:      name,
			Version:   version.(string),
			Ecosystem: "pypi",
		}
		dependencies = append(dependencies, dependency)
	}

	for name, version := range devDependencies {
		dependency := Dependency{
			Name:      name,
			Version:   version.(string),
			Ecosystem: "pypi",
		}
		dependencies = append(dependencies, dependency)
	}

	return dependencies, nil
}
