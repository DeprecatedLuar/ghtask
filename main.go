package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type Label struct {
	Name string `json:"name"`
}

type Issue struct {
	Number    int     `json:"number"`
	Title     string  `json:"title"`
	Labels    []Label `json:"labels"`
	CreatedAt string  `json:"createdAt"`
}

func main() {
	cmd, args := detectCommand()

	switch cmd {
	case "g0":
		createIssue("P0", args)
	case "g1":
		createIssue("P1", args)
	case "g2":
		createIssue("P2", args)
	case "g3":
		createIssue("P3", args)
	case "list", "":
		listIssues(args)
	case "p0", "p1", "p2", "p3":
		filter := strings.ToUpper(cmd)
		listIssues([]string{filter})
	case "active":
		listIssues([]string{"active"})
	case "start", "activate":
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Error: issue number required")
			fmt.Fprintln(os.Stderr, "Usage: gt start <issue-number>")
			os.Exit(1)
		}
		startIssue(args[0])
	case "done":
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Error: issue number required")
			fmt.Fprintln(os.Stderr, "Usage: gt done <issue-number>")
			os.Exit(1)
		}
		closeIssue(args[0])
	case "help", "--help", "-h":
		showHelp()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		showHelp()
		os.Exit(1)
	}
}

func detectCommand() (string, []string) {
	binary := filepath.Base(os.Args[0])

	// Linux/Mac: symlink name is the command
	if binary == "g0" || binary == "g1" || binary == "g2" || binary == "g3" {
		return binary, os.Args[1:]
	}

	// Windows: command is first argument
	if len(os.Args) > 1 {
		firstArg := os.Args[1]
		if firstArg == "g0" || firstArg == "g1" || firstArg == "g2" || firstArg == "g3" {
			return firstArg, os.Args[2:]
		}
	}

	// Main command
	if len(os.Args) > 1 {
		return os.Args[1], os.Args[2:]
	}

	// Just "gt" = list all
	return "list", []string{}
}

func getRepoFromGit() (string, error) {
	// Check for environment variable override
	if repo := os.Getenv("GT_REPO"); repo != "" {
		return repo, nil
	}

	// Get git remote URL
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not in a git repository or no origin remote")
	}

	remoteURL := strings.TrimSpace(string(output))
	return parseGitHubRepo(remoteURL)
}

func parseGitHubRepo(url string) (string, error) {
	// Remove .git suffix
	url = strings.TrimSuffix(url, ".git")

	// SSH format: git@github.com:owner/repo
	if strings.HasPrefix(url, "git@github.com:") {
		return strings.TrimPrefix(url, "git@github.com:"), nil
	}

	// HTTPS format: https://github.com/owner/repo
	if strings.Contains(url, "github.com/") {
		parts := strings.Split(url, "github.com/")
		if len(parts) == 2 {
			return parts[1], nil
		}
	}

	return "", fmt.Errorf("could not parse GitHub repo from: %s", url)
}

func createIssue(priority string, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: issue title required")
		fmt.Fprintf(os.Stderr, "Usage: %s <title>\n", strings.ToLower(priority))
		os.Exit(1)
	}

	repo, err := getRepoFromGit()
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

func listIssues(filters []string) {
	repo, err := getRepoFromGit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Make sure you're in a git repository with a GitHub remote")
		os.Exit(1)
	}

	cmd := exec.Command("gh", "issue", "list",
		"--repo", repo,
		"--state", "open",
		"--json", "number,title,labels,createdAt",
		"--limit", "1000")

	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing issues: %v\n", err)
		os.Exit(1)
	}

	var issues []Issue
	if err := json.Unmarshal(output, &issues); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing issues: %v\n", err)
		os.Exit(1)
	}

	// Filter issues
	filtered := filterIssues(issues, filters)

	// Sort by priority, then by creation date
	sortIssues(filtered)

	// Print issues
	if len(filtered) == 0 {
		fmt.Println("No issues found")
		return
	}

	for i, issue := range filtered {
		printIssue(issue, i)
	}
}

func filterIssues(issues []Issue, filters []string) []Issue {
	if len(filters) == 0 {
		return issues
	}

	var filtered []Issue
	for _, issue := range issues {
		if matchesFilters(issue, filters) {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

func matchesFilters(issue Issue, filters []string) bool {
	for _, filter := range filters {
		filter = strings.ToLower(filter)

		if filter == "active" {
			if !hasLabel(issue, "active") {
				return false
			}
		} else if strings.HasPrefix(filter, "p") && len(filter) == 2 {
			// Priority filter (p0, p1, p2, p3)
			priority := strings.ToUpper(filter)
			if !hasLabel(issue, priority) {
				return false
			}
		}
	}
	return true
}

func hasLabel(issue Issue, labelName string) bool {
	labelName = strings.ToLower(labelName)
	for _, label := range issue.Labels {
		if strings.ToLower(label.Name) == labelName {
			return true
		}
	}
	return false
}

func extractPriority(issue Issue) string {
	for _, label := range issue.Labels {
		if strings.HasPrefix(label.Name, "P") && len(label.Name) == 2 {
			return label.Name
		}
	}
	return "P2"
}

func isActive(issue Issue) bool {
	return hasLabel(issue, "active")
}

func sortIssues(issues []Issue) {
	sort.Slice(issues, func(i, j int) bool {
		priI := extractPriority(issues[i])
		priJ := extractPriority(issues[j])

		// Lower number = higher priority
		if priI != priJ {
			return priI < priJ
		}

		// Same priority: sort by creation date (oldest first)
		return issues[i].CreatedAt < issues[j].CreatedAt
	})
}

func printIssue(issue Issue, index int) {
	priority := extractPriority(issue)
	active := isActive(issue)

	// Get colors
	textColor := getPriorityColor(priority)
	bgColor := getBackgroundColor(index, active)
	reset := "\033[0m"

	// Active tasks get black text on light background
	if active {
		textColor = "\033[38;5;0m"
	}

	fmt.Printf("%s%s#%-5d %-4s %s%s\n",
		bgColor,
		textColor,
		issue.Number,
		priority,
		issue.Title,
		reset)
}

func getPriorityColor(priority string) string {
	switch priority {
	case "P0":
		return "\033[38;5;196m" // Red
	case "P1":
		return "\033[38;5;208m" // Orange
	case "P2":
		return "\033[38;5;250m" // White
	case "P3":
		return "\033[38;5;240m" // Dark Gray
	default:
		return "\033[38;5;250m" // Default: White
	}
}

func getBackgroundColor(index int, active bool) string {
	if active {
		return "\033[48;5;250m" // Light/white background for active
	}

	// Alternating backgrounds
	if index%2 == 0 {
		return "\033[48;5;233m" // Dark gray
	}
	return "\033[48;5;232m" // Darker gray
}

func startIssue(issueNum string) {
	repo, err := getRepoFromGit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Make sure you're in a git repository with a GitHub remote")
		os.Exit(1)
	}

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

func closeIssue(issueNum string) {
	repo, err := getRepoFromGit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Make sure you're in a git repository with a GitHub remote")
		os.Exit(1)
	}

	// Get issue title before closing
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

	// Close the issue
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

func showHelp() {
	help := `gt - GitHub Tasks

USAGE:
  gt                    List all open issues
  gt p0/p1/p2/p3        Filter by priority
  gt active             Show only active tasks
  gt start <number>     Mark issue as active
  gt done <number>      Close issue

  g0 <title>            Create P0 (critical) issue
  g1 <title>            Create P1 (important) issue
  g2 <title>            Create P2 (normal) issue
  g3 <title>            Create P3 (low) issue

EXAMPLES:
  g1 fix authentication bug
  gt p1
  gt start 234
  gt done 567

SETUP:
  1. Run: gh auth login
  2. Navigate to a git repo with GitHub remote
  3. Create labels: gh label create "inbox" --color "d4c5f9"
                    gh label create "active" --color "0e8a16"
                    gh label create "P0" --color "d93f0b"
                    gh label create "P1" --color "ff9800"
                    gh label create "P2" --color "ffeb3b"
                    gh label create "P3" --color "cccccc"
`
	fmt.Print(help)
}
