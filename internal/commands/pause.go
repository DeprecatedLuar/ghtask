package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/DeprecatedLuar/ghtask/internal"
)

func PauseIssue(args []string) {
	issueNum, err := ParseIssueNumber(args, "pause")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	repo := internal.GetRepoOrDie()

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
