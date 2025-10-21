// Package main implements gt (GitHub Tasks), a lightweight CLI tool for managing
// GitHub Issues with a dstask-inspired workflow. It provides fast task creation
// via priority shortcuts (g0-g3) and uses GitHub Issues as the single source of truth.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/DeprecatedLuar/ghtask/internal/commands"
)

func main() {
	cmd, args := detectCommand()

	switch cmd {
	case "g0":
		commands.CreateIssue("P0", args)
	case "g1":
		commands.CreateIssue("P1", args)
	case "g2":
		commands.CreateIssue("P2", args)
	case "g3":
		commands.CreateIssue("P3", args)
	case "create-default":
		commands.CreateIssue("P2", args)
	case "list", "":
		commands.ListIssues(args, false)
	case "p0", "p1", "p2", "p3":
		filter := strings.ToUpper(cmd)
		commands.ListIssues([]string{filter}, false)
	case "active":
		commands.ListIssues([]string{"active"}, false)
	case "start", "activate":
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Error: issue number required")
			fmt.Fprintln(os.Stderr, "Usage: gt start <issue-number>")
			os.Exit(1)
		}
		commands.StartIssue(args[0])
	case "done":
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Error: issue number required")
			fmt.Fprintln(os.Stderr, "Usage: gt done <issue-number>")
			os.Exit(1)
		}
		commands.CloseIssue(args[0])
	case "rm", "delete":
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Error: issue number required")
			fmt.Fprintln(os.Stderr, "Usage: gt rm <issue-number>")
			os.Exit(1)
		}
		commands.DeleteIssue(args[0])
	case "setup":
		commands.SetupRepo()
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

	if binary == "g0" || binary == "g1" || binary == "g2" || binary == "g3" {
		return binary, os.Args[1:]
	}

	if len(os.Args) > 1 {
		firstArg := os.Args[1]
		if firstArg == "g0" || firstArg == "g1" || firstArg == "g2" || firstArg == "g3" {
			return firstArg, os.Args[2:]
		}
	}

	if len(os.Args) == 1 {
		return "list", []string{}
	}

	firstArg := os.Args[1]

	knownCommands := []string{"list", "p0", "p1", "p2", "p3", "active", "start", "activate", "done", "rm", "delete", "setup", "help", "--help", "-h"}
	if slices.Contains(knownCommands, firstArg) {
		return firstArg, os.Args[2:]
	}

	if strings.HasPrefix(firstArg, "-") {
		return "list", os.Args[1:]
	}

	return "create-default", os.Args[1:]
}

func showHelp() {
	help := `gt - GitHub Tasks

USAGE:
  gt [-v]               List all open issues
  gt <title>            Create P2 (normal) issue (default)
  gt p0/p1/p2/p3 [-v]   Filter by priority
  gt active [-v]        Show only active tasks
  gt start <number>     Mark issue as active
  gt done <number>      Close issue
  gt rm <number>        Delete issue (permanent)
  gt setup              Create required labels in repo

  g0 <title>            Create P0 (critical) issue
  g1 <title>            Create P1 (important) issue
  g2 <title>            Create P2 (normal) issue
  g3 <title>            Create P3 (low) issue

FLAGS:
  -v, --verbose         Show priority labels in output

EXAMPLES:
  gt setup                      # Setup labels for this repo
  gt                            # List all tasks (colors only)
  gt -v                         # List all tasks (with priority text)
  gt refactor legacy code       # Create default P2 task
  g1 fix authentication bug     # Create P1 task
  gt p1 -v                      # List P1 tasks with priority labels
  gt start 234                  # Mark #234 as active
  gt done 567                   # Close #567
  gt rm 890                     # Delete #890 (permanent)

WORKFLOW:
  g2 <title>    - Creates a P2 issue
  gt p2         - Lists existing P2 issues

SETUP:
  1. Run: gh auth login
  2. Navigate to a git repo with GitHub remote
  3. Run: gt setup (creates all required labels)
`
	fmt.Print(help)
}
