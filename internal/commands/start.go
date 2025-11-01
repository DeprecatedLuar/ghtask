package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/DeprecatedLuar/ghtask/internal"
)

func StartIssue(args []string) {
	issueNum, err := ParseIssueNumber(args, "start")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	repo := internal.GetRepoOrDie()

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
