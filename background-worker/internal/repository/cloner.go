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
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	cmd := exec.Command("git", "clone", "--depth", "1", url, tempDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("git clone failed: %w", err)
	}

	return tempDir, nil
}
