package parsers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/pelletier/go-toml/v2"
)

type Pep621PyProjectParser struct{}

func (p Pep621PyProjectParser) CanParse(filename string) bool {
	return filename == "pyproject.toml"
}

func (p Pep621PyProjectParser) Parse(fileContent []byte) ([]Dependency, error) {
	var data map[string]any
	if err := toml.Unmarshal(fileContent, &data); err != nil {
		return nil, fmt.Errorf("failed to parse pyproject.toml: %w", err)
	}

	projectSection, ok := data["project"].(map[string]any)
	if !ok {
		return nil, nil // Skip if project section is not found
	}

	dependenciesSection, ok := projectSection["dependencies"].([]any)
	if !ok {
		return nil, nil // Skip if dependencies section is not found
	}

	reBracket := regexp.MustCompile(`\[.*?\]`)
	rePrefix := regexp.MustCompile(`^[~^<>=!]+`)

	var dependencies []Dependency

	for _, dep := range dependenciesSection {
		depstr, ok := dep.(string)
		if !ok {
			continue
		}

		name, version := parsePythonDependency(depstr, reBracket, rePrefix)

		parsedVersion, err := semver.NewVersion(version)
		if err != nil {
			version = "unknown"
		} else {
			version = parsedVersion.String()
		}

		dependency := Dependency{
			Name:      strings.ToLower(name),
			Version:   version,
			Ecosystem: "pypi",
		}
		dependencies = append(dependencies, dependency)
	}

	return dependencies, nil
}

func parsePythonDependency(dep string, reBrackets, rePrefix *regexp.Regexp) (string, string) {
	name := dep
	version := "unknown"

	// Format: "package (>=1.0,<2.0)"
	if strings.Contains(dep, " (") && strings.HasSuffix(dep, ")") {
		parts := strings.SplitN(dep[:len(dep)-1], " (", 2)
		if len(parts) == 2 {
			name = parts[0]
			if versionParts := strings.Split(parts[1], ","); len(versionParts) > 0 {
				version = versionParts[0]
			}
		}
	}

	name = strings.TrimSpace(reBrackets.ReplaceAllString(name, ""))
	version = strings.TrimSpace(rePrefix.ReplaceAllString(version, ""))
	return name, version
}
