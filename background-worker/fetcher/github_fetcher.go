package fetcher

import (
	"fmt"
	"os/exec"
	"strings"
)

type GithubFetcher struct{}

var execCommand = exec.Command

func (g GithubFetcher) CanHandle(repoUrl string) bool {
	return strings.HasPrefix(repoUrl, "git@") || strings.Contains(repoUrl, ".git")

}

func (g GithubFetcher) Fetch(repoUrl string, destPath string) error {
	cmd := execCommand("git", "clone", "--depth=1", repoUrl, destPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to clone repository: %v, output: %s", err, output)
	}

	return nil
}

