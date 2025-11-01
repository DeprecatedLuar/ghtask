package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/DeprecatedLuar/ghtask/internal"
)

func CreateIssue(args []string, cmd string, hasBody bool, bodyValue string) {
	priority := ParsePriorityFromCommand(cmd)

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: issue title required")
		fmt.Fprintf(os.Stderr, "Usage: %s <title> [--body [inline-text]]\n", cmd)
		os.Exit(1)
	}

	repo := internal.GetRepoOrDie()

	title := strings.Join(args, " ")
	labels := "inbox," + priority

	body, err := GetContentFromInput(hasBody, bodyValue, "body")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting body: %v\n", err)
		os.Exit(1)
	}

	ghCmd := exec.Command("gh", "issue", "create",
		"--repo", repo,
		"--title", title,
		"--label", labels,
		"--body", body)

	output, err := ghCmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating issue: %v\n", err)
		fmt.Fprintf(os.Stderr, "%s\n", output)
		os.Exit(1)
	}

	fmt.Printf("Created: %s\n", strings.TrimSpace(string(output)))
}
