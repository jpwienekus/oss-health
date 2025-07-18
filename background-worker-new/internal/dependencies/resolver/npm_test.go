package resolver_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/dependencies/resolver"
)

func TestCanResolve(t *testing.T) {
	r := &resolver.NPMResolver{}

	tests := []struct {
		name       string
		dependency db.Dependency
		expected   bool
	}{
		{"Valid NPM", db.Dependency{Ecosystem: "npm"}, true},
		{"Non-NPM", db.Dependency{Ecosystem: "pypi"}, false},
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

func TestResolve_Success(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"repository": map[string]interface{}{
				"url": "git+https://github.com/user/repo.git",
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	originalURL := resolver.NPMRegistryURL
	defer func() { resolver.NPMRegistryURL = originalURL }()
	resolver.NPMRegistryURL = server.URL

	r := &resolver.NPMResolver{}
	dep := db.Dependency{Name: "test-pkg", Ecosystem: "npm"}

	url, err := r.Resolve(context.Background(), dep)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := "https://github.com/user/repo"
	if url != expected {
		t.Errorf("expected %s, got %s", expected, url)
	}
}

func TestResolve_Non200Response(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	originalURL := resolver.NPMRegistryURL
	defer func() { resolver.NPMRegistryURL = originalURL }()
	resolver.NPMRegistryURL = server.URL

	r := &resolver.NPMResolver{}
	_, err := r.Resolve(context.Background(), db.Dependency{Name: "not-found", Ecosystem: "npm"})
	if err == nil || !strings.Contains(err.Error(), "received status code") {
		t.Errorf("expected error on non-200 response, got %v", err)
	}
}

func TestResolve_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{invalid json}"))
	}))
	defer server.Close()

	originalURL := resolver.NPMRegistryURL
	defer func() { resolver.NPMRegistryURL = originalURL }()
	resolver.NPMRegistryURL = server.URL

	r := &resolver.NPMResolver{}
	_, err := r.Resolve(context.Background(), db.Dependency{Name: "bad-json", Ecosystem: "npm"})
	if err == nil {
		t.Errorf("expected JSON decoding error, got nil")
	}
}
