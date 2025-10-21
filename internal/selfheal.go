// Package internal provides infrastructure utilities for the gt binary,
// including automatic shortcut creation for priority commands.
package internal

import (
	"os"
	"path/filepath"
	"runtime"
)

// HealShortcuts creates priority shortcuts (g0-g3) next to the gt binary.
// On Linux/Mac: creates symlinks
// On Windows: creates .bat files
// Silently succeeds if shortcuts already exist.
func HealShortcuts() {
	binaryPath, err := os.Executable()
	if err != nil {
		return
	}

	binaryDir := filepath.Dir(binaryPath)
	binaryName := filepath.Base(binaryPath)

	if runtime.GOOS == "windows" {
		createWindowsBatch(binaryDir, binaryName)
	} else {
		createLinuxSymlinks(binaryDir, binaryName)
	}
}

func createLinuxSymlinks(dir, binaryName string) {
	shortcuts := []string{"g0", "g1", "g2", "g3"}

	for _, shortcut := range shortcuts {
		linkPath := filepath.Join(dir, shortcut)

		if _, err := os.Lstat(linkPath); err == nil {
			continue
		}

		_ = os.Symlink(binaryName, linkPath)
	}
}

func createWindowsBatch(dir, binaryName string) {
	shortcuts := []string{"g0", "g1", "g2", "g3"}

	for _, shortcut := range shortcuts {
		batPath := filepath.Join(dir, shortcut+".bat")

		if _, err := os.Stat(batPath); err == nil {
			continue
		}

		content := "@echo off\r\n" + binaryName + " " + shortcut + " %*\r\n"
		_ = os.WriteFile(batPath, []byte(content), 0755)
	}
}
