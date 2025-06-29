package parsers

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/oss-health/background-worker/internal/repository"
)


func init() {
	repository.RegisterParser("requirements.txt", "pypi", ParseRequirementsTxt)
	repository.RegisterParser("requirements-*.txt", "pypi", ParseRequirementsTxt)
	repository.RegisterParser("requirements/*.txt", "pypi", ParseRequirementsTxt)
}

func ParseRequirementsTxt(path string) ([]repository.DependencyParsed, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close file: %v\n", err)
		}
	}()

	deps := []repository.DependencyParsed{}
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
			name = parts[0]
			version = parts[1]
		}

		deps = append(deps, repository.DependencyParsed{name, version, "PyPI"})
	}

	return deps, scanner.Err()
}
