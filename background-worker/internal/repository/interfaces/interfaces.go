package interfaces

import (
	"github.com/oss-health/background-worker/internal/dependency"
)

type Parser interface {
	Parse(path string) ([]dependency.DependencyVersionPair, error)
}

type ParserProvider interface {
	GetParser(path string) Parser
}
