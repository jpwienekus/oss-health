package parsers

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/oss-health/background-worker/internal/dependency"
)

func init() {
	RegisterParser("requirements.txt", "pypi", ParseRequirementsTxt)
	RegisterParser("requirements-*.txt", "pypi", ParseRequirementsTxt)
	RegisterParser("requirements/*.txt", "pypi", ParseRequirementsTxt)
}

func ParseRequirementsTxt(path string) ([]dependency.DependencyVersionPair, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("open requirements file %q: %w", path, err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close file %q: %v\n", path, cerr)
		}
	}()

	deps := []dependency.DependencyVersionPair{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		name := line
		version := "unknown"

		if strings.Contains(line, "==") {
			parts := strings.SplitN(line, "==", 2)

			if len(parts) == 2 {
				name = strings.TrimSpace(parts[0])
				version = strings.TrimSpace(parts[1])
			}
		}

		deps = append(deps, dependency.DependencyVersionPair{
			Name:      name,
			Version:   version,
			Ecosystem: "PyPI",
		})
	}

	return deps, scanner.Err()
}
