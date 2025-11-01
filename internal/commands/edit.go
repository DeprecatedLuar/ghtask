package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/DeprecatedLuar/ghtask/internal/github"
)

func EditIssue(args []string) {
	hasEdit, field, inlineValue, remainingArgs, err := ParseEditFlag(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if !hasEdit {
		fmt.Fprintln(os.Stderr, "Error: -e flag required")
		fmt.Fprintln(os.Stderr, "Usage: gt <issue-number> -e <body|title> [inline-value]")
		os.Exit(1)
	}

	issueNum, err := ParseIssueNumber(remainingArgs, "edit")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	repo, err := github.GetRepoFromGit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Make sure you're in a git repository with a GitHub remote")
		os.Exit(1)
	}

	var updateFlag string
	switch field {
	case "body":
		updateFlag = "--body"
	case "title":
		updateFlag = "--title"
	}

	var newContent string
	stat, _ := os.Stdin.Stat()
	isPiped := (stat.Mode() & os.ModeCharDevice) == 0

	if isPiped || inlineValue != "" {
		newContent, err = GetContentFromInput(true, inlineValue, field)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting content: %v\n", err)
			os.Exit(1)
		}
	} else {
		cmd := exec.Command("gh", "issue", "view", issueNum,
			"--repo", repo,
			"--json", "title,body")

		output, err := cmd.Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching issue: %v\n", err)
			os.Exit(1)
		}

		var issueData struct {
			Title string `json:"title"`
			Body  string `json:"body"`
		}

		if err := json.Unmarshal(output, &issueData); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing issue: %v\n", err)
			os.Exit(1)
		}

		var currentContent string
		switch field {
		case "body":
			currentContent = issueData.Body
		case "title":
			currentContent = issueData.Title
		}

		newContent, err = openEditorWithContent(currentContent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening editor: %v\n", err)
			os.Exit(1)
		}

		if strings.TrimSpace(newContent) == strings.TrimSpace(currentContent) {
			fmt.Println("No changes made")
			return
		}
	}

	updateCmd := exec.Command("gh", "issue", "edit", issueNum,
		"--repo", repo,
		updateFlag, newContent)

	if err := updateCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating issue: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Updated %s for issue #%s\n", field, issueNum)
}

func openEditorWithContent(content string) (string, error) {
	editor := GetEditor()

	tmpFile, err := os.CreateTemp("", "ghtask-edit-*.md")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	if _, err := tmpFile.WriteString(content); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}
	tmpFile.Close()

	defer os.Remove(tmpPath)

	cmd := exec.Command(editor, tmpPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor exited with error: %w", err)
	}

	editedContent, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", fmt.Errorf("failed to read temp file: %w", err)
	}

	return strings.TrimSpace(string(editedContent)), nil
}
