package dependencies

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/oss-health/background-worker/internal/db"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewOSVClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    "https://api.osv.dev/v1",
	}
}

type osvQuery struct {
	Package struct {
		Name      string `json:"name"`
		Ecosystem string `json:"ecosystem"`
	} `json:"package"`
	Versions []string `json:"versions"`
}

type osvBatchRequest struct {
	Query []osvQuery `json:"query"`
}

type osvBatchResponse struct {
	Results []struct {
		Vulnerabilities []struct {
			ID					string `json:"id"`
			Summary			string `json:"summary"`
			Details			string `json:"details"`
		} `json:"vulnerabilities"`
	} `json:"results"`
}

func (c *Client) QueryBatch(ctx context.Context, dependencyVersions []db.DependencyVersion) (map[string][]string, error) {
	queries := make([]osvQuery, len(dependencyVersions))

	for i, dv := range dependencyVersions {
		queries[i] = osvQuery{
			Package: struct {
				Name      string `json:"name"`
				Ecosystem string `json:"ecosystem"`
			}{
				Name:      dv.DependencyName,
				Ecosystem: dv.DependencyEcosystem,
			},
			Versions: []string{dv.Version},
		}
	}

	requestBody, err := json.Marshal(osvBatchRequest{Query: queries})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OSV batch request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/querybatch", bytes.NewBuffer(requestBody))
	log.Printf("OSV Batch Request: %s", request.URL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create OSV batch request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to execute OSV batch request: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OSV batch request failed with status: %s", response.Status)
	}

	var osvResponse osvBatchResponse
	err = json.NewDecoder(response.Body).Decode(&osvResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode OSV batch response: %w", err)
	}

	results := make(map[string][]string)
	for i, result := range osvResponse.Results {
		key := fmt.Sprintf("%s:%s", dependencyVersions[i].DependencyName, dependencyVersions[i].Version)
		for _, vuln := range result.Vulnerabilities {
			results[key] = append(results[key], vuln.ID)
		}
	}

	return results, nil
}
