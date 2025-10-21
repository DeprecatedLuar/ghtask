package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/DeprecatedLuar/ghtask/internal/github"
)

func SetupRepo() {
	repo, err := github.GetRepoFromGit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Make sure you're in a git repository with a GitHub remote")
		os.Exit(1)
	}

	fmt.Printf("Setting up labels for %s...\n", repo)

	requiredLabels := []struct {
		name  string
		color string
		desc  string
	}{
		{"inbox", "d4c5f9", "Newly created tasks"},
		{"active", "0e8a16", "Currently working on"},
		{"P0", "d93f0b", "Critical priority"},
		{"P1", "ff9800", "Important priority"},
		{"P2", "ffeb3b", "Normal priority"},
		{"P3", "cccccc", "Low priority"},
	}

	existingLabels := getExistingLabels(repo)
	created := 0
	skipped := 0

	for _, label := range requiredLabels {
		if labelExists(existingLabels, label.name) {
			fmt.Printf("  ✓ %s (already exists)\n", label.name)
			skipped++
		} else {
			if createLabel(repo, label.name, label.color, label.desc) {
				fmt.Printf("  ✓ %s (created)\n", label.name)
				created++
			} else {
				fmt.Printf("  ✗ %s (failed)\n", label.name)
			}
		}
	}

	fmt.Printf("\nSetup complete: %d created, %d already existed\n", created, skipped)
}

func getExistingLabels(repo string) []string {
	cmd := exec.Command("gh", "label", "list", "--repo", repo, "--json", "name")
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	var labels []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(output, &labels); err != nil {
		return []string{}
	}

	names := make([]string, len(labels))
	for i, label := range labels {
		names[i] = label.Name
	}
	return names
}

func labelExists(labels []string, name string) bool {
	for _, label := range labels {
		if strings.EqualFold(label, name) {
			return true
		}
	}
	return false
}

func createLabel(repo, name, color, description string) bool {
	cmd := exec.Command("gh", "label", "create", name,
		"--repo", repo,
		"--color", color,
		"--description", description)

	err := cmd.Run()
	return err == nil
}
