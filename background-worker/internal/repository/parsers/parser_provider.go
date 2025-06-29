package parsers

import (
	"path/filepath"

	"github.com/oss-health/background-worker/internal/dependency"
	"github.com/oss-health/background-worker/internal/repository/interfaces"
)

type adapter struct {
	fn ParserFunc
}

func (a *adapter) Parse(path string) ([]dependency.DependencyVersionPair, error) {
	return a.fn(path)
}

type ParserProviderImpl struct{}

func (p *ParserProviderImpl) GetParser(path string) interfaces.Parser {
	base := filepath.Base(path)

	for _, entry := range GetRegisteredParsers() {
		if entry.Pattern == base {
			return &adapter{fn: entry.Parser}
		}
	}
	return nil
}
