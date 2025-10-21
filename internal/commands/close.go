package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/DeprecatedLuar/ghtask/internal/github"
)

func CloseIssue(issueNum string) {
	repo, err := github.GetRepoFromGit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Make sure you're in a git repository with a GitHub remote")
		os.Exit(1)
	}

	cmd := exec.Command("gh", "issue", "view", issueNum,
		"--repo", repo,
		"--json", "title")

	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error viewing issue: %v\n", err)
		os.Exit(1)
	}

	var issueData struct {
		Title string `json:"title"`
	}
	if err := json.Unmarshal(output, &issueData); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing issue: %v\n", err)
		os.Exit(1)
	}

	cmd = exec.Command("gh", "issue", "close", issueNum,
		"--repo", repo)

	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error closing issue: %v\n", err)
		fmt.Fprintf(os.Stderr, "%s\n", output)
		os.Exit(1)
	}

	fmt.Printf("✓ Closed #%s: %s\n", issueNum, issueData.Title)
}
