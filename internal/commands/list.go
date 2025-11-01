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

const (
	// Display formatting
	defaultTerminalWidth = 80  // Fallback terminal width when detection fails
	minTitleWidth        = 40  // Minimum width for title column
	issueNumReserved     = 7   // Reserved space for issue number column (#1234 + spacing)
	issueNumWidth        = 5   // Width for issue number formatting (%-5d)
	issueNumPadding      = 3   // Zero-padding width for verbose mode (03d)

	// List fetching
	maxIssueLimit = 1000 // Maximum number of issues to fetch from GitHub

	// Color codes
	colorBlackText = 0   // Black text for active issues
	colorGrayZeros = 235 // Gray color for leading zeros in verbose mode
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
		"--limit", fmt.Sprintf("%d", maxIssueLimit))

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
		return defaultTerminalWidth
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
		textColor = fmt.Sprintf("\033[38;5;%dm", colorBlackText)
	}

	title := issue.Title
	isTerminal := term.IsTerminal(int(os.Stdout.Fd()))

	if !verbose && isTerminal {
		termWidth := getTerminalWidth()

		availableWidth := termWidth - issueNumReserved
		if availableWidth < minTitleWidth {
			availableWidth = minTitleWidth
		}

		title = truncateTitle(title, availableWidth)
	}

	var content string
	if verbose {
		paddedNum := fmt.Sprintf("%0*d", issueNumPadding, issue.Number)
		grayLeadingZeros := formatLeadingZeros(paddedNum, textColor)
		content = fmt.Sprintf("[%s-%s]   %s", grayLeadingZeros, priority, title)
	} else {
		content = fmt.Sprintf("%-*d %s", issueNumWidth, issue.Number, title)
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
	gray := fmt.Sprintf("\033[38;5;%dm", colorGrayZeros)

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
