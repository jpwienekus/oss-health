package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
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
		defer resp.Body.Close()

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

func GetPypiRepoURL(client *http.Client, baseURL string) func(ctx context.Context, name string) (string, error) {
	return func(ctx context.Context, name string) (string, error) {
		url := fmt.Sprintf("%s/%s/json", baseURL, name)
		req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return "", nil
		}

		var data struct {
			Info struct {
				ProjectURLs map[string]string `json:"project_urls"`
				HomePage    string            `json:"home_page"`
			} `json:"info"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return "", err
		}

		if src, ok := data.Info.ProjectURLs["Source"]; ok {
			return src, nil
		}
		if homepage := data.Info.ProjectURLs["Homepage"]; homepage != "" {
			return homepage, nil
		}
		return data.Info.HomePage, nil
	}
}
