package vulnerabilities

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"sync"
	"time"

	"github.com/oss-health/background-worker/internal/db"
	"golang.org/x/time/rate"
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

func (c *Client) QueryBatch(ctx context.Context, dependencyVersions []db.DependencyVersion) (map[int][]string, error) {
	batchSize := 100
	maxRPS := 10
	batches := c.splitIntoBatches(dependencyVersions, batchSize)
	rateLimiter := rate.NewLimiter(rate.Limit(maxRPS), 1)

	var (
		mu      sync.Mutex
		results = make(map[int][]string)
		wg      sync.WaitGroup
		errCh   = make(chan error, len(batches))
	)

	for _, b := range batches {
		// Rate limit before launching
		if err := rateLimiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limiter wait failed: %w", err)
		}

		wg.Add(1)
		go func(depBatch []db.DependencyVersion) {
			log.Printf("Processing batch of %d dependencies\n", len(b))
			defer wg.Done()

			partialResults, err := c.querySingleBatch(ctx, depBatch)
			if err != nil {
				errCh <- err
				return
			}

			mu.Lock()
			maps.Copy(results, partialResults)
			// for k, v := range partialResults {
			// 	results[k] = v
			// }
			mu.Unlock()
		}(b)
	}

	wg.Wait()
	close(errCh)

	if len(errCh) > 0 {
		return nil, <-errCh // return first error
	}

	return results, nil
}

func (c *Client) GetVulnerabilityDetails(ctx context.Context, vulnerabilities map[int][]string) (map[int][]Vulnerability, error) {
	const maxRPS = 10
	rateLimiter := rate.NewLimiter(rate.Limit(maxRPS), 1)

	type job struct {
		DepVersionID int
		VulnID       string
	}

	var jobs []job
	for depID, vulnIDs := range vulnerabilities {
		for _, vulnID := range vulnIDs {
			jobs = append(jobs, job{DepVersionID: depID, VulnID: vulnID})
		}
	}

	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		results = make(map[int][]Vulnerability)
		errCh   = make(chan error, len(jobs))
	)

	for _, j := range jobs {
		if err := rateLimiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limiter wait failed: %w", err)
		}

		wg.Add(1)
		go func(j job) {
			defer wg.Done()

			url := fmt.Sprintf("%s/vulns/%s", c.baseURL, j.VulnID)
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				errCh <- fmt.Errorf("failed to create request for %s: %w", j.VulnID, err)
				return
			}

			resp, err := c.httpClient.Do(req)
			if err != nil {
				errCh <- fmt.Errorf("failed to fetch %s: %w", j.VulnID, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				errCh <- fmt.Errorf("OSV response for %s failed: %s - %s", j.VulnID, resp.Status, string(body))
				return
			}

			var osvResp osvVulnerability
			if err := json.NewDecoder(resp.Body).Decode(&osvResp); err != nil {
				errCh <- fmt.Errorf("failed to decode %s: %w", j.VulnID, err)
				return
			}

			if osvResp.DatabaseSpecific.Severity != "" {
				mu.Lock()
				results[j.DepVersionID] = append(results[j.DepVersionID], Vulnerability{
					ID:       j.VulnID,
					Summary:  osvResp.Summary,
					Severity: osvResp.DatabaseSpecific.Severity,
				})
				mu.Unlock()
			}
		}(j)
	}

	wg.Wait()
	close(errCh)

	if len(errCh) > 0 {
		return nil, <-errCh // return first error
	}

	return results, nil
}

func (c *Client) querySingleBatch(ctx context.Context, dependencyVersions []db.DependencyVersion) (map[int][]string, error) {
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/querybatch", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create OSV batch request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute OSV batch request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OSV batch request failed with status: %s, body: %s", resp.Status, string(body))
	}

	var osvResponse osvBatchResponse
	if err := json.NewDecoder(resp.Body).Decode(&osvResponse); err != nil {
		return nil, fmt.Errorf("failed to decode OSV batch response: %w", err)
	}

	results := make(map[int][]string)
	for i, result := range osvResponse.Results {
		id := dependencyVersions[i].ID
		for _, vuln := range result.Vulnerabilities {
			results[id] = append(results[id], vuln.ID)
		}
	}

	return results, nil
}

func (c *Client) splitIntoBatches(deps []db.DependencyVersion, size int) [][]db.DependencyVersion {
	var batches [][]db.DependencyVersion
	for size < len(deps) {
		deps, batches = deps[size:], append(batches, deps[0:size:size])
	}
	batches = append(batches, deps)
	return batches
}
