package resolver_test

import (
	"context"
	"testing"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/dependencies/resolver"
)

func TestPyPIResolver_CanResolve(t *testing.T) {
	r := &resolver.PyPIResolver{}

	tests := []struct {
		name       string
		dependency db.Dependency
		expected   bool
	}{
		{"Valid PyPI", db.Dependency{Ecosystem: "pypi"}, true},
		{"Invalid NPM", db.Dependency{Ecosystem: "npm"}, false},
		{"Empty Ecosystem", db.Dependency{Ecosystem: ""}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := r.CanResolve(tt.dependency)
			if actual != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}

func TestPyPIResolver_Resolve(t *testing.T) {
	r := &resolver.PyPIResolver{}

	dep := db.Dependency{Name: "requests", Ecosystem: "pypi"}
	expected := "https://pypi.org/project/requests/"

	url, err := r.Resolve(context.Background(), dep)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if url != expected {
		t.Errorf("expected %s, got %s", expected, url)
	}
}
