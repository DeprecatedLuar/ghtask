// Package internal contains shared types and utilities used across gt commands.
package internal

import "strings"

// ExtractPriority extracts the priority label (P0-P3) from an issue's labels
func ExtractPriority(issue Issue) string {
	for _, label := range issue.Labels {
		if strings.HasPrefix(label.Name, "P") && len(label.Name) == 2 {
			return label.Name
		}
	}
	return "P2"
}

// GetPriorityColor returns the ANSI color code for a given priority level
func GetPriorityColor(priority string) string {
	switch priority {
	case "P0":
		return "\033[38;2;215;0;0m" // #d70000
	case "P1":
		return "\033[38;5;208m"
	case "P2":
		return "\033[38;5;250m"
	case "P3":
		return "\033[38;2;88;88;88m" // #585858
	default:
		return "\033[38;5;250m"
	}
}

// GetBackgroundColor returns the ANSI background color for list display
func GetBackgroundColor(index int, active bool) string {
	if active {
		return "\033[48;5;250m"
	}

	if index%2 == 0 {
		return "\033[48;2;18;18;18m" // #121212
	}
	return ""
}
