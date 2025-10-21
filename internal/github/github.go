// Package github provides utilities for Git repository detection and GitHub
// remote URL parsing. It auto-detects the current repository from git remotes
// and supports both SSH and HTTPS URL formats.
package github

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetRepoFromGit() (string, error) {
	if repo := os.Getenv("GT_REPO"); repo != "" {
		return repo, nil
	}

	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not in a git repository or no origin remote")
	}

	remoteURL := strings.TrimSpace(string(output))
	return parseGitHubRepo(remoteURL)
}

func parseGitHubRepo(url string) (string, error) {
	url = strings.TrimSuffix(url, ".git")

	if strings.HasPrefix(url, "git@github.com:") {
		return strings.TrimPrefix(url, "git@github.com:"), nil
	}

	if strings.Contains(url, "github.com/") {
		parts := strings.Split(url, "github.com/")
		if len(parts) == 2 {
			return parts[1], nil
		}
	}

	return "", fmt.Errorf("could not parse GitHub repo from: %s", url)
}
