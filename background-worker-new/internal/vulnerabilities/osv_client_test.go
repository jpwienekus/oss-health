package vulnerabilities_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/oss-health/background-worker/internal/db"
	"github.com/oss-health/background-worker/internal/vulnerabilities"
	"github.com/stretchr/testify/assert"
)

// mockRoundTripper implements http.RoundTripper
type mockRoundTripper struct {
	responder func(*http.Request) *http.Response
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.responder(req), nil
}

func newMockClient(responder func(*http.Request) *http.Response) *vulnerabilities.OSVClient {
	httpClient := &http.Client{
		Transport: &mockRoundTripper{responder: responder},
		Timeout:   2 * time.Second,
	}

	return vulnerabilities.NewTestClient(httpClient, "https://api.osv.dev/v1")
}

func TestQuerySingleBatch_Success(t *testing.T) {
	mockResp := `{
		"results": [
			{ "vulns": [ { "id": "OSV-1" } ] },
			{ "vulns": [ { "id": "OSV-2" }, { "id": "OSV-3" } ] }
		]
	}`

	client := newMockClient(func(req *http.Request) *http.Response {
		assert.Equal(t, "https://api.osv.dev/v1/querybatch", req.URL.String())
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockResp)),
			Header:     make(http.Header),
		}
	})

	input := []db.DependencyVersion{
		{ID: 1, DependencyName: "lodash", DependencyEcosystem: "npm", Version: "1.0.0"},
		{ID: 2, DependencyName: "react", DependencyEcosystem: "npm", Version: "2.0.0"},
	}

	results, err := client.QueryBatch(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, []string{"OSV-1"}, results[1])
	assert.Equal(t, []string{"OSV-2", "OSV-3"}, results[2])
}

func TestGetVulnerabilityDetails_Success(t *testing.T) {
	mockResp := `{
		"id": "OSV-1",
		"summary": "Critical vuln",
		"database_specific": {
			"severity": "CRITICAL"
		}
	}`

	client := newMockClient(func(req *http.Request) *http.Response {
		assert.Contains(t, req.URL.Path, "/vulns/OSV-1")
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockResp)),
			Header:     make(http.Header),
		}
	})

	input := map[int][]string{
		123: {"OSV-1"},
	}

	results, err := client.GetVulnerabilityDetails(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, "OSV-1", results[123][0].ID)
	assert.Equal(t, "CRITICAL", results[123][0].Severity)
}

func TestQuerySingleBatch_Failure(t *testing.T) {
	client := newMockClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(bytes.NewBufferString("Internal Server Error")),
			Header:     make(http.Header),
		}
	})

	input := []db.DependencyVersion{
		{ID: 1, DependencyName: "axios", DependencyEcosystem: "npm", Version: "1.0.0"},
	}

	_, err := client.QueryBatch(context.Background(), input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "batch request failed")
}

func TestGetVulnerabilityDetails_BadJSON(t *testing.T) {
	client := newMockClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`{not json}`)),
			Header:     make(http.Header),
		}
	})

	input := map[int][]string{
		456: {"OSV-123"},
	}

	_, err := client.GetVulnerabilityDetails(context.Background(), input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode")
}
