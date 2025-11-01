// Package internal contains shared helper functions used across gt commands.
package internal

import (
	"fmt"
	"os"

	"github.com/DeprecatedLuar/ghtask/internal/github"
)

// GetRepoOrDie retrieves the GitHub repository from git remote or GT_REPO env var.
// Exits with error message if not in a git repository or remote is not GitHub.
func GetRepoOrDie() string {
	repo, err := github.GetRepoFromGit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Make sure you're in a git repository with a GitHub remote")
		os.Exit(1)
	}
	return repo
}
