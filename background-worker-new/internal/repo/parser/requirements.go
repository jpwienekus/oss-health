package parser

import "strings"

type RequirementsParser struct {}

func (r RequirementsParser) CanParse(filename string) bool {
	return strings.EqualFold(filename, "requirements.txt")
}

func (r RequirementsParser) Parse(fileContent []byte) ([]Dependency, error) {
	lines := strings.Split(string(fileContent), "\n")
	var dependencies []Dependency

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "==")
		if len(parts) != 2 {
			continue
		}

		dependency := Dependency{
			Name:      strings.TrimSpace(parts[0]),
			Version:   strings.TrimSpace(parts[1]),
			Ecosystem: "pypi",
		}
		dependencies = append(dependencies, dependency)
	}

	return dependencies, nil
}
