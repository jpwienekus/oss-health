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

	err = cmd.Run()

	if err != nil {
		removeErr := os.RemoveAll(tempDir)

		if removeErr != nil {
			fmt.Printf("warning: failed to remove temp dir %s: %v\n", tempDir, removeErr)
		}

		return "", fmt.Errorf("git clone failed: %w", err)
	}

	return tempDir, nil
}
