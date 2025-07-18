package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/oss-health/background-worker/internal/db"
)

type NPMResolver struct{}

var NPMRegistryURL = "https://registry.npmjs.org"

func (r *NPMResolver) CanResolve(dependency db.Dependency) bool {
	return dependency.Ecosystem == "npm"
}

func (r *NPMResolver) Resolve(ctx context.Context, dependency db.Dependency) (string, error) {
	url := fmt.Sprintf("%s/%s", NPMRegistryURL, dependency.Name)

	request, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to resolve npm package %s: %w", dependency.Name, err)
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to resolve npm package %s: received status code %d", dependency.Name, response.StatusCode)
	}

	var data struct {
		Repository struct {
			URL string `json:"url"`
		} `json:"repository"`
	}

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return "", err
	}

	repoURL := data.Repository.URL
	repoURL = strings.ReplaceAll(repoURL, "git+", "")
	repoURL = strings.ReplaceAll(repoURL, ".git", "")
	repoURL = strings.ReplaceAll(repoURL, "git://", "https://")

	return repoURL, nil
}
