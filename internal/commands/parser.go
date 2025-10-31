package commands

import (
	"fmt"
	"strconv"
)

// ParseVerboseFlag extracts verbose flag from args and returns (verbose, remainingArgs)
func ParseVerboseFlag(args []string) (bool, []string) {
	verbose := false
	remaining := []string{}

	for _, arg := range args {
		if arg == "-v" || arg == "--verbose" {
			verbose = true
		} else {
			remaining = append(remaining, arg)
		}
	}

	return verbose, remaining
}

// ParseIssueNumber extracts and validates issue number from args
// Returns the issue number string or error if invalid/missing
func ParseIssueNumber(args []string, commandName string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("issue number required\nUsage: gt %s <issue-number>", commandName)
	}

	issueNum := args[0]

	// Validate it's a number
	if _, err := strconv.Atoi(issueNum); err != nil {
		return "", fmt.Errorf("invalid issue number: %s", issueNum)
	}

	return issueNum, nil
}

// ParseBodyFlag extracts --body flag from args and returns (hasBodyFlag, remainingArgs)
func ParseBodyFlag(args []string) (bool, []string) {
	hasBody := false
	remaining := []string{}

	for _, arg := range args {
		if arg == "--body" || arg == "-b" {
			hasBody = true
		} else {
			remaining = append(remaining, arg)
		}
	}

	return hasBody, remaining
}

// ParseEditFlag extracts -e flag and field (body/title) from args
// Returns (hasEditFlag, field, remainingArgs, error)
func ParseEditFlag(args []string) (bool, string, []string, error) {
	hasEdit := false
	field := ""
	remaining := []string{}

	for i := 0; i < len(args); i++ {
		if args[i] == "-e" || args[i] == "--edit" {
			hasEdit = true
			if i+1 < len(args) {
				field = args[i+1]
				if field != "body" && field != "title" {
					return false, "", nil, fmt.Errorf("invalid field: %s (must be 'body' or 'title')", field)
				}
				i++ // skip the field argument
			} else {
				return false, "", nil, fmt.Errorf("-e flag requires field argument (body or title)")
			}
		} else {
			remaining = append(remaining, args[i])
		}
	}

	return hasEdit, field, remaining, nil
}

// ParsePriorityFromCommand converts command name to priority label
func ParsePriorityFromCommand(cmd string) string {
	switch cmd {
	case "gt0":
		return "P0"
	case "gt1":
		return "P1"
	case "gt2":
		return "P2"
	case "gt3":
		return "P3"
	default:
		return "P2" // default priority
	}
}
