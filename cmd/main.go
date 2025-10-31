// Package main implements gt (GitHub Tasks), a lightweight CLI tool for managing
// GitHub Issues with a dstask-inspired workflow. It provides fast task creation
// via priority shortcuts (gt0-gt3) and uses GitHub Issues as the single source of truth.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/DeprecatedLuar/ghtask/internal"
	"github.com/DeprecatedLuar/ghtask/internal/commands"
)

func main() {
	internal.HealShortcuts()

	cmd, args := detectCommand()

	switch cmd {
	case "gt0", "gt1", "gt2", "gt3", "create-default":
		openEditor, remainingArgs := commands.ParseBodyFlag(args)
		commands.CreateIssue(remainingArgs, cmd, openEditor)
	case "list", "":
		commands.ListIssues(args)
	case "p0", "p1", "p2", "p3":
		filter := strings.ToUpper(cmd)
		commands.ListIssues(append([]string{filter}, args...))
	case "active":
		commands.ListIssues(append([]string{"active"}, args...))
	case "start", "activate":
		commands.StartIssue(args)
	case "pause", "stop":
		commands.PauseIssue(args)
	case "done":
		commands.CloseIssue(args)
	case "rm", "delete":
		commands.DeleteIssue(args)
	case "setup":
		commands.SetupRepo()
	case "help", "--help", "-h":
		commands.ShowHelp()
	case "view":
		commands.ViewIssue(args)
	case "edit":
		commands.EditIssue(args)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		commands.ShowHelp()
		os.Exit(1)
	}
}

func detectCommand() (string, []string) {
	binary := filepath.Base(os.Args[0])

	if binary == "gt0" || binary == "gt1" || binary == "gt2" || binary == "gt3" {
		return binary, os.Args[1:]
	}

	if len(os.Args) > 1 {
		firstArg := os.Args[1]
		if firstArg == "gt0" || firstArg == "gt1" || firstArg == "gt2" || firstArg == "gt3" {
			return firstArg, os.Args[2:]
		}
	}

	if len(os.Args) == 1 {
		return "list", []string{}
	}

	firstArg := os.Args[1]

	knownCommands := []string{"list", "p0", "p1", "p2", "p3", "active", "start", "activate", "pause", "stop", "done", "rm", "delete", "setup", "help", "--help", "-h"}
	if slices.Contains(knownCommands, firstArg) {
		return firstArg, os.Args[2:]
	}

	if strings.HasPrefix(firstArg, "-") {
		return "list", os.Args[1:]
	}

	if _, err := strconv.Atoi(firstArg); err == nil {
		hasEdit := false
		for _, arg := range os.Args[2:] {
			if arg == "-e" || arg == "--edit" {
				hasEdit = true
				break
			}
		}
		if hasEdit {
			return "edit", os.Args[1:]
		}
		return "view", os.Args[1:]
	}

	return "create-default", os.Args[1:]
}

