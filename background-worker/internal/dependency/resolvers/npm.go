package resolvers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func GetNpmRepoURL(client *http.Client, baseURL string) func(ctx context.Context, name string) (string, error) {
	return func(ctx context.Context, name string) (string, error) {
		url := fmt.Sprintf("%s/%s", baseURL, name)
		req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("failed to close response body: %v", err)
			}
		}()

		if resp.StatusCode != 200 {
			return "", nil
		}

		var data struct {
			Repository struct {
				URL string `json:"url"`
			} `json:"repository"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return "", err
		}

		repoURL := data.Repository.URL
		repoURL = strings.ReplaceAll(repoURL, "git+", "")
		repoURL = strings.ReplaceAll(repoURL, ".git", "")
		repoURL = strings.ReplaceAll(repoURL, "git://", "https://")

		return repoURL, nil
	}
}
