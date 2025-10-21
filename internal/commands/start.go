package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/DeprecatedLuar/ghtask/internal/github"
)

func StartIssue(issueNum string) {
	repo, err := github.GetRepoFromGit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Make sure you're in a git repository with a GitHub remote")
		os.Exit(1)
	}

	cmd := exec.Command("gh", "issue", "edit", issueNum,
		"--repo", repo,
		"--add-label", "active")

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error activating issue: %v\n", err)
		fmt.Fprintf(os.Stderr, "%s\n", output)
		os.Exit(1)
	}

	fmt.Printf("✓ Activated #%s\n", issueNum)
}
