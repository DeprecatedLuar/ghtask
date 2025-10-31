package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/DeprecatedLuar/ghtask/internal/github"
)

func PauseIssue(args []string) {
	issueNum, err := ParseIssueNumber(args, "pause")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	repo, err := github.GetRepoFromGit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Make sure you're in a git repository with a GitHub remote")
		os.Exit(1)
	}

	cmd := exec.Command("gh", "issue", "edit", issueNum,
		"--repo", repo,
		"--remove-label", "active")

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error pausing issue: %v\n", err)
		fmt.Fprintf(os.Stderr, "%s\n", output)
		os.Exit(1)
	}

	fmt.Printf("✓ Paused #%s\n", issueNum)
}
