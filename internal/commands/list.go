package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/DeprecatedLuar/ghtask/internal/github"
	"golang.org/x/term"
)

func ListIssues(args []string, verbose bool) {
	filters := []string{}
	for _, arg := range args {
		if arg == "-v" || arg == "--verbose" {
			verbose = true
		} else {
			filters = append(filters, arg)
		}
	}

	repo, err := github.GetRepoFromGit()
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

	filtered := filterIssues(issues, filters)
	sortIssues(filtered)

	if len(filtered) == 0 {
		fmt.Println("No issues found")
		return
	}

	for i, issue := range filtered {
		printIssue(issue, i, verbose)
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

		if priI != priJ {
			return priI < priJ
		}

		return issues[i].CreatedAt < issues[j].CreatedAt
	})
}

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80
	}
	return width
}

func truncateTitle(title string, maxWidth int) string {
	if len(title) <= maxWidth {
		return title
	}
	if maxWidth < 1 {
		return ""
	}
	return title[:maxWidth-1] + ">"
}

func printIssue(issue Issue, index int, verbose bool) {
	priority := extractPriority(issue)
	active := isActive(issue)

	textColor := getPriorityColor(priority)
	bgColor := getBackgroundColor(index, active)
	reset := "\033[0m"

	if active {
		textColor = "\033[38;5;0m"
	}

	title := issue.Title
	if !verbose {
		termWidth := getTerminalWidth()
		reservedSpace := 7

		availableWidth := termWidth - reservedSpace
		if availableWidth < 40 {
			availableWidth = 40
		}

		title = truncateTitle(title, availableWidth)
	}

	if verbose {
		fmt.Printf("%s%s#%-5d %-4s %s%s\n",
			bgColor,
			textColor,
			issue.Number,
			priority,
			title,
			reset)
	} else {
		fmt.Printf("%s%s#%-5d %s%s\n",
			bgColor,
			textColor,
			issue.Number,
			title,
			reset)
	}
}

func getPriorityColor(priority string) string {
	switch priority {
	case "P0":
		return "\033[38;5;196m"
	case "P1":
		return "\033[38;5;208m"
	case "P2":
		return "\033[38;5;250m"
	case "P3":
		return "\033[38;5;240m"
	default:
		return "\033[38;5;250m"
	}
}

func getBackgroundColor(index int, active bool) string {
	if active {
		return "\033[48;5;250m"
	}

	if index%2 == 0 {
		return "\033[48;5;233m"
	}
	return "\033[48;5;232m"
}
