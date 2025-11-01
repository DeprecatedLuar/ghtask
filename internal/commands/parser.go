package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

// ParseBodyFlag extracts --body flag and optional inline value from args
// Returns (hasBodyFlag, inlineValue, remainingArgs)
// Example: gt1 "title" --body "text" → returns (true, "text", ["title"])
func ParseBodyFlag(args []string) (bool, string, []string) {
	hasBody := false
	inlineValue := ""
	remaining := []string{}

	for i := 0; i < len(args); i++ {
		if args[i] == "--body" || args[i] == "-b" {
			hasBody = true
			// Collect everything after --body as inline value
			if i+1 < len(args) {
				inlineValue = strings.Join(args[i+1:], " ")
				break // we've consumed the rest
			}
		} else if !hasBody {
			remaining = append(remaining, args[i])
		}
	}

	return hasBody, inlineValue, remaining
}

// ParseEditFlag extracts -e flag, field (body/title), and optional inline value from args
// Returns (hasEditFlag, field, inlineValue, remainingArgs, error)
// Example: gt 123 -e body "new text" → returns (true, "body", "new text", ["123"], nil)
func ParseEditFlag(args []string) (bool, string, string, []string, error) {
	hasEdit := false
	field := ""
	inlineValue := ""
	remaining := []string{}

	for i := 0; i < len(args); i++ {
		if args[i] == "-e" || args[i] == "--edit" {
			hasEdit = true
			if i+1 < len(args) {
				field = args[i+1]
				if field != "body" && field != "title" {
					return false, "", "", nil, fmt.Errorf("invalid field: %s (must be 'body' or 'title')", field)
				}
				i++ // skip the field argument

				// Collect everything after field as inline value
				if i+1 < len(args) {
					inlineValue = strings.Join(args[i+1:], " ")
					break // we've consumed the rest
				}
			} else {
				return false, "", "", nil, fmt.Errorf("-e flag requires field argument (body or title)")
			}
		} else {
			remaining = append(remaining, args[i])
		}
	}

	return hasEdit, field, inlineValue, remaining, nil
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

// GetContentFromInput determines content source in priority order:
// 1. Stdin (pipes, automation) - if available
// 2. Inline value - if provided
// 3. Editor - if flag was set
// Returns content string or error
func GetContentFromInput(flagProvided bool, inlineValue string, fieldName string) (string, error) {
	stat, _ := os.Stdin.Stat()
	isPiped := (stat.Mode() & os.ModeCharDevice) == 0

	if isPiped {
		return readFromStdin()
	}

	if inlineValue != "" {
		return inlineValue, nil
	}

	if flagProvided {
		return openEditorForContent(fieldName)
	}

	return "", nil
}

func readFromStdin() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var content strings.Builder

	for scanner.Scan() {
		content.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading from stdin: %w", err)
	}

	result := strings.TrimSpace(content.String())
	if result == "" {
		return "", fmt.Errorf("no content provided via stdin")
	}

	fmt.Fprintf(os.Stderr, "✓ Read %d bytes from stdin\n", len(result))
	return result, nil
}

func openEditorForContent(fieldName string) (string, error) {
	editor := GetEditor()

	tmpFile, err := os.CreateTemp("", "ghtask-"+fieldName+"-*.md")
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
