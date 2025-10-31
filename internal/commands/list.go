package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/DeprecatedLuar/ghtask/internal"
	"github.com/DeprecatedLuar/ghtask/internal/github"
	"golang.org/x/term"
)

func ListIssues(args []string) {
	verbose, filters := ParseVerboseFlag(args)

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

	var issues []internal.Issue
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

func filterIssues(issues []internal.Issue, filters []string) []internal.Issue {
	if len(filters) == 0 {
		return issues
	}

	var filtered []internal.Issue
	for _, issue := range issues {
		if matchesFilters(issue, filters) {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

func matchesFilters(issue internal.Issue, filters []string) bool {
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

func hasLabel(issue internal.Issue, labelName string) bool {
	labelName = strings.ToLower(labelName)
	for _, label := range issue.Labels {
		if strings.ToLower(label.Name) == labelName {
			return true
		}
	}
	return false
}

func isActive(issue internal.Issue) bool {
	return hasLabel(issue, "active")
}

func sortIssues(issues []internal.Issue) {
	sort.Slice(issues, func(i, j int) bool {
		priI := internal.ExtractPriority(issues[i])
		priJ := internal.ExtractPriority(issues[j])

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

func printIssue(issue internal.Issue, index int, verbose bool) {
	priority := internal.ExtractPriority(issue)
	active := isActive(issue)

	textColor := internal.GetPriorityColor(priority)
	bgColor := internal.GetBackgroundColor(index, active)
	reset := "\033[0m"

	if active {
		textColor = "\033[38;5;0m"
	}

	title := issue.Title
	isTerminal := term.IsTerminal(int(os.Stdout.Fd()))

	if !verbose && isTerminal {
		termWidth := getTerminalWidth()
		reservedSpace := 7

		availableWidth := termWidth - reservedSpace
		if availableWidth < 40 {
			availableWidth = 40
		}

		title = truncateTitle(title, availableWidth)
	}

	var content string
	if verbose {
		paddedNum := fmt.Sprintf("%03d", issue.Number)
		grayLeadingZeros := formatLeadingZeros(paddedNum, textColor)
		content = fmt.Sprintf("[%s-%s]   %s", grayLeadingZeros, priority, title)
	} else {
		content = fmt.Sprintf("%-5d %s", issue.Number, title)
	}

	padding := ""
	if isTerminal {
		termWidth := getTerminalWidth()
		visibleLen := getVisibleLength(content)
		if visibleLen < termWidth {
			padding = strings.Repeat(" ", termWidth-visibleLen)
		}
	}

	fmt.Printf("%s%s%s%s%s\n", bgColor, textColor, content, padding, reset)
}

func getVisibleLength(s string) int {
	visible := 0
	inEscape := false

	for i := 0; i < len(s); i++ {
		if s[i] == '\033' {
			inEscape = true
		} else if inEscape && s[i] == 'm' {
			inEscape = false
		} else if !inEscape {
			visible++
		}
	}

	return visible
}

func formatLeadingZeros(paddedNum string, mainColor string) string {
	gray := "\033[38;5;235m"

	isTerminal := term.IsTerminal(int(os.Stdout.Fd()))
	if !isTerminal {
		return paddedNum
	}

	firstNonZero := 0
	for i, ch := range paddedNum {
		if ch != '0' {
			firstNonZero = i
			break
		}
	}

	if firstNonZero == 0 {
		return paddedNum
	}

	leadingZeros := paddedNum[:firstNonZero]
	significantDigits := paddedNum[firstNonZero:]

	return gray + leadingZeros + mainColor + significantDigits
}
