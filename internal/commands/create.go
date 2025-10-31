package commands

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/DeprecatedLuar/ghtask/internal/github"
)

func CreateIssue(args []string, cmd string, openEditor bool) {
	priority := ParsePriorityFromCommand(cmd)

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: issue title required")
		fmt.Fprintf(os.Stderr, "Usage: %s <title> [--body]\n", cmd)
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

	body := ""
	if openEditor {
		var err error
		body, err = openEditorForBody()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening editor: %v\n", err)
			os.Exit(1)
		}
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

func openEditorForBody() (string, error) {
	editor := GetEditor()

	tmpFile, err := os.CreateTemp("", "ghtask-body-*.md")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()

	defer os.Remove(tmpPath)

	cmd := exec.Command(editor, tmpPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor exited with error: %w", err)
	}

	content, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", fmt.Errorf("failed to read temp file: %w", err)
	}

	return strings.TrimSpace(string(content)), nil
}

func GetEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}

	if runtime.GOOS == "windows" {
		return "notepad.exe"
	}

	for _, editor := range []string{"vim", "nano", "vi"} {
		if _, err := exec.LookPath(editor); err == nil {
			return editor
		}
	}

	return "vi"
}
