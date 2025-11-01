// Package internal provides editor utilities for creating and editing content
// in temporary files using the user's preferred editor.
package internal

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// GetEditor returns the user's preferred editor.
// Priority: $EDITOR env var > platform defaults (vim/nano/vi on Unix, notepad on Windows)
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

// OpenEditorWithContent opens the user's editor with optional initial content.
// If content is empty, creates an empty temp file.
// Returns the edited content after the editor exits.
func OpenEditorWithContent(initialContent, fieldName string) (string, error) {
	editor := GetEditor()

	tmpFile, err := os.CreateTemp("", "ghtask-"+fieldName+"-*.md")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	if initialContent != "" {
		if _, err := tmpFile.WriteString(initialContent); err != nil {
			tmpFile.Close()
			os.Remove(tmpPath)
			return "", fmt.Errorf("failed to write to temp file: %w", err)
		}
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
