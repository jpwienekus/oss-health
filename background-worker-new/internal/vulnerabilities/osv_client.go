package vulnerabilities

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

type osvPackage struct {
	Name      string `json:"name"`
	Ecosystem string `json:"ecosystem"`
}

type osvQuery struct {
	Version string     `json:"version"`
	Package osvPackage `json:"package"`
}

type osvBatchRequest struct {
	Queries []osvQuery `json:"queries"`
}

type osvBatchResponse struct {
	Results []osvQueryResult `json:"results"`
}

type osvQueryResult struct {
	Vulnerabilities []osvBatchResultVulnerability `json:"vulns"`
}

type osvBatchResultVulnerability struct {
	ID string `json:"id"`
}

type osvVulnerability struct {
	ID               string            `json:"id"`
	Summary          string            `json:"summary"`
	DatabaseSpecific DatabbaseSpecific `json:"database_specific"`
}

type osvAffectedVersion struct {
	Package osvPackage `json:"package"`
	Ranges  []osvRange `json:"ranges"`
}

type DatabbaseSpecific struct {
	Severity string `json:"severity"`
}

type osvRange struct {
	Type   string     `json:"type"`
	Events []osvEvent `json:"events"`
}

type osvEvent struct {
	Introduced string `json:"introduced,omitempty"`
	Fixed      string `json:"fixed,omitempty"`
}

type Vulnerability struct {
	ID       string
	Summary  string
	Severity string
}

func (c *Client) GetVulnerabilityData(ctx context.Context, dependencyVersions []db.DependencyVersion) (map[int][]Vulnerability, error) {
	queryResults, err := c.queryBatch(ctx, dependencyVersions)
	if err != nil {
		return nil, fmt.Errorf("failed to query OSV batch: %w", err)
	}

	detailsResults, err := c.getVulnerabilityDetails(ctx, queryResults)

	return detailsResults, nil
}

func (c *Client) queryBatch(ctx context.Context, dependencyVersions []db.DependencyVersion) (map[int][]string, error) {
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
			Version: dv.Version,
		}
	}

	requestBody, err := json.Marshal(osvBatchRequest{Queries: queries})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OSV batch request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/querybatch", bytes.NewBuffer(requestBody))
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
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("OSV batch request failed with status: %s, body: %s", response.Status, string(body))
	}

	var osvResponse osvBatchResponse
	err = json.NewDecoder(response.Body).Decode(&osvResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode OSV batch response: %w", err)
	}

	results := make(map[int][]string)
	for i, result := range osvResponse.Results {
		key := dependencyVersions[i].ID
		for _, vuln := range result.Vulnerabilities {
			results[key] = append(results[key], vuln.ID)
		}
	}

	return results, nil
}

func (c *Client) getVulnerabilityDetails(ctx context.Context, vulnerabilities map[int][]string) (map[int][]Vulnerability, error) {
	severityResults := make(map[int][]Vulnerability)

	for key, result := range vulnerabilities {
		for _, vulnID := range result {
			severityURL := fmt.Sprintf("%s/vulns/%s", c.baseURL, vulnID)
			severityRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, severityURL, nil)

			if err != nil {
				return nil, fmt.Errorf("failed to create OSV severity request for %s: %w", vulnID, err)
			}

			severityResponse, err := c.httpClient.Do(severityRequest)
			if err != nil {
				return nil, fmt.Errorf("failed to execute OSV severity request for %s: %w", vulnID, err)
			}

			defer severityResponse.Body.Close()

			if severityResponse.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(severityResponse.Body)
				return nil, fmt.Errorf("OSV severity request for %s failed with status: %s, body: %s", vulnID, severityResponse.Status, string(body))
			}

			vulnerabilityResponse := osvVulnerability{}
			err = json.NewDecoder(severityResponse.Body).Decode(&vulnerabilityResponse)
			if err != nil {
				return nil, fmt.Errorf("failed to decode OSV severity response for %s: %w", vulnID, err)
			}

			if vulnerabilityResponse.DatabaseSpecific.Severity != "" {
				severityResults[key] = append(severityResults[key], Vulnerability{
					ID:       vulnID,
					Summary:  vulnerabilityResponse.Summary,
					Severity: vulnerabilityResponse.DatabaseSpecific.Severity,
				})
			}
		}
	}

	return severityResults, nil
}
