// Package internal contains shared types and utilities used across gt commands.
package internal

import (
	"fmt"
	"strings"
)

const (
	// Priority label format
	priorityLabelLength = 2 // Expected length of priority labels (P0, P1, P2, P3)

	// Priority colors (ANSI 256-color and RGB codes)
	colorP0Red       = 215 // P0 critical - red component
	colorP0Green     = 0   // P0 critical - green component
	colorP0Blue      = 0   // P0 critical - blue component
	colorP1Orange    = 208 // P1 important - orange
	colorP2Gray      = 250 // P2 normal - gray/white
	colorP3Gray      = 88  // P3 low - dark gray component
	colorP3GrayGreen = 88  // P3 low - dark gray green component
	colorP3GrayBlue  = 88  // P3 low - dark gray blue component

	// Background colors
	colorActiveBg      = 250 // Active issue background - light gray
	colorAltRowBgRed   = 18  // Alternating row background - dark gray red component
	colorAltRowBgGreen = 18  // Alternating row background - dark gray green component
	colorAltRowBgBlue  = 18  // Alternating row background - dark gray blue component

	// Row alternation
	evenRowModulo = 2 // Modulo value for alternating row backgrounds
)

// ExtractPriority extracts the priority label (P0-P3) from an issue's labels
func ExtractPriority(issue Issue) string {
	for _, label := range issue.Labels {
		if strings.HasPrefix(label.Name, "P") && len(label.Name) == priorityLabelLength {
			return label.Name
		}
	}
	return "P2"
}

// GetPriorityColor returns the ANSI color code for a given priority level
func GetPriorityColor(priority string) string {
	switch priority {
	case "P0":
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", colorP0Red, colorP0Green, colorP0Blue) // Critical - red
	case "P1":
		return fmt.Sprintf("\033[38;5;%dm", colorP1Orange) // Important - orange
	case "P2":
		return fmt.Sprintf("\033[38;5;%dm", colorP2Gray) // Normal - gray/white
	case "P3":
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", colorP3Gray, colorP3GrayGreen, colorP3GrayBlue) // Low - dark gray
	default:
		return fmt.Sprintf("\033[38;5;%dm", colorP2Gray) // Default to P2 color
	}
}

// GetBackgroundColor returns the ANSI background color for list display
func GetBackgroundColor(index int, active bool) string {
	if active {
		return fmt.Sprintf("\033[48;5;%dm", colorActiveBg) // Light gray background for active issues
	}

	if index%evenRowModulo == 0 {
		return fmt.Sprintf("\033[48;2;%d;%d;%dm", colorAltRowBgRed, colorAltRowBgGreen, colorAltRowBgBlue) // Dark gray for alternating rows
	}
	return "" // No background for odd rows
}
