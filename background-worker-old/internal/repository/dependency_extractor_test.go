package repository_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/oss-health/background-worker/internal/dependency"
	"github.com/oss-health/background-worker/internal/repository"
	"github.com/oss-health/background-worker/internal/repository/interfaces"
)

type mockParser struct{}

func (m *mockParser) Parse(path string) ([]dependency.DependencyVersionPair, error) {
	return []dependency.DependencyVersionPair{
		{Name: "mocklib", Version: "1.2.3", Ecosystem: "mock"},
	}, nil
}

type mockProvider struct{}

func (m *mockProvider) GetParser(path string) interfaces.Parser {
	if strings.HasSuffix(path, "mockfile.txt") {
		return &mockParser{}
	}
	return nil
}

func TestExtractDependencies(t *testing.T) {
	tempDir := t.TempDir()

	testFile := filepath.Join(tempDir, "mockfile.txt")

	if err := os.WriteFile(testFile, []byte("irrelevant content"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	extractor := &repository.DependencyExtractor{
		Provider: &mockProvider{},
	}

	deps, err := extractor.ExtractDependencies(tempDir)

	if err != nil {
		t.Fatalf("ExtractDependencies returned error: %v", err)
	}

	if len(deps) != 1 {
		t.Fatalf("expected 1 dependency, got %d", len(deps))
	}

	dep := deps[0]

	if dep.Name != "mocklib" || dep.Version != "1.2.3" || dep.Ecosystem != "mock" {
		t.Errorf("unexpected dependency: %+v", dep)
	}
}

type dupMockParser struct{}

func (m *dupMockParser) Parse(path string) ([]dependency.DependencyVersionPair, error) {
	return []dependency.DependencyVersionPair{
		{Name: "dup-lib", Version: "1.0.0", Ecosystem: "mock"},
	}, nil
}

type dupMockProvider struct{}

func (m *dupMockProvider) GetParser(path string) interfaces.Parser {
	if strings.HasSuffix(path, "mockfile1.txt") || strings.HasSuffix(path, "mockfile2.txt") {
		return &dupMockParser{}
	}
	return nil
}
func TestExtractDependencies_Deduplication(t *testing.T) {
	tempDir := t.TempDir()

	file1 := filepath.Join(tempDir, "mockfile1.txt")
	file2 := filepath.Join(tempDir, "mockfile2.txt")

	for _, f := range []string{file1, file2} {
		if err := os.WriteFile(f, []byte("irrelevant content"), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}
	}

	extractor := &repository.DependencyExtractor{
		Provider: &dupMockProvider{},
	}

	deps, err := extractor.ExtractDependencies(tempDir)

	if err != nil {
		t.Fatalf("ExtractDependencies returned error: %v", err)
	}

	if len(deps) != 1 {
		t.Fatalf("expected 1 unique dependency after deduplication, got %d", len(deps))
	}

	dep := deps[0]
	if dep.Name != "dup-lib" || dep.Version != "1.0.0" || dep.Ecosystem != "mock" {
		t.Errorf("unexpected dependency: %+v", dep)
	}
}
