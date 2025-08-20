package repository

import (
	"fmt"
	"os"
	"os/exec"
)

type Cloner interface {
	CloneRepository(url string) (string, error)
}

type GitCloner struct{}

func (g *GitCloner) CloneRepository(url string) (string, error) {
	tempDir, err := os.MkdirTemp("", "repo-clone-*")
	if err != nil {
		return "", fmt.Errorf("create temp dir for cloning %q: %w", url, err)
	}

	cmd := exec.Command("git", "clone", "--depth", "1", url, tempDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if removeErr := os.RemoveAll(tempDir); removeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to clean up temp dir %s after clone failure: %v\n", tempDir, removeErr)
		}
		return "", fmt.Errorf("git clone %q into %s failed: %w", url, tempDir, err)
	}

	return tempDir, nil
}
