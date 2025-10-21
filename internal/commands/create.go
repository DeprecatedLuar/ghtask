package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/DeprecatedLuar/ghtask/internal/github"
)

func CreateIssue(priority string, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: issue title required")
		fmt.Fprintf(os.Stderr, "Usage: %s <title>\n", strings.ToLower(priority))
		os.Exit(1)
	}

	repo, err := github.GetRepoFromGit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Make sure you're in a git repository with a GitHub remote")
		os.Exit(1)
	}

	title := strings.Join(args, " ")
	labels := "inbox," + priority

	cmd := exec.Command("gh", "issue", "create",
		"--repo", repo,
		"--title", title,
		"--label", labels,
		"--body", "")

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating issue: %v\n", err)
		fmt.Fprintf(os.Stderr, "%s\n", output)
		os.Exit(1)
	}

	fmt.Printf("Created: %s\n", strings.TrimSpace(string(output)))
}
