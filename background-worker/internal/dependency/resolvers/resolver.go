package resolvers

import (
	"context"
	"net/http"
	"time"
)

var defaultHTTPClient = &http.Client{Timeout: 10 * time.Second}
var Resolvers = map[string]func(ctx context.Context, name string) (string, error){
	"npm":  GetNpmRepoURL(defaultHTTPClient, "https://registry.npmjs.org"),
	"pypi": GetPypiRepoURL(defaultHTTPClient, "https://pypi.org/pypi"),
}
