package resolvers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/oss-health/background-worker/internal/dependency/resolvers"
)

func TestNpmResolver_Success(t *testing.T) {
	expected := "https://github.com/example/repo"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"repository": map[string]string{
				"url": "git+https://github.com/example/repo.git",
			},
		})
	}))
	defer ts.Close()

	resolver := resolvers.GetNpmRepoURL(&http.Client{Timeout: 2 * time.Second}, ts.URL)
	url, err := resolver(context.Background(), "example")

	require.NoError(t, err)
	require.Equal(t, expected, url)
}

func TestPypiResolver_SourcePriority(t *testing.T) {
	expected := "https://github.com/example/project"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"info": map[string]any{
				"project_urls": map[string]string{
					"Source": expected,
				},
				"home_page": "https://fallback.com",
			},
		})
	}))
	defer ts.Close()

	resolver := resolvers.GetPypiRepoURL(&http.Client{Timeout: 2 * time.Second}, ts.URL)
	url, err := resolver(context.Background(), "example")

	require.NoError(t, err)
	require.Equal(t, expected, url)
}
