package dependency

import (
	"context"
	"net/http"
	"time"

	"github.com/oss-health/background-worker/internal/dependency/resolvers"
)

var defaultHTTPClient = &http.Client{Timeout: 10 * time.Second}
var Resolvers = map[string]func(ctx context.Context, name string) (string, error){
	"npm":  resolvers.GetNpmRepoURL(defaultHTTPClient, "https://registry.npmjs.org"),
	"pypi": resolvers.GetPypiRepoURL(defaultHTTPClient, "https://pypi.org/pypi"),
}
