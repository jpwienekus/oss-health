package resolvers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GetPypiRepoURL(client *http.Client, baseURL string) func(ctx context.Context, name string) (string, error) {
	return func(ctx context.Context, name string) (string, error) {
		url := fmt.Sprintf("%s/%s/json", baseURL, name)
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
