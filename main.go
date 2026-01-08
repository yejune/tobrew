package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yejune/tobrew/cmd"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:   "tobrew",
		Short: "Automate Homebrew tap releases for Go projects",
		Long: `tobrew - Automated Homebrew Release Tool

A CLI tool to automate the entire Homebrew tap release process:
  - Automatic semantic versioning (tobrew.lock)
  - Build and create GitHub releases
  - Calculate SHA256 checksums automatically
  - Generate/update Homebrew formulas
  - Manage homebrew-tap repository

Simple workflow:
  1. tobrew init              # Create config (once)
  2. tobrew release           # Release with patch bump (v1.0.0 → v1.0.1)
  3. tobrew release --minor   # Release with minor bump (v1.0.1 → v1.1.0)
  4. tobrew release --major   # Release with major bump (v1.1.0 → v2.0.0)

Building a CLI tool with tobrew?
  Run 'tobrew selfupdate-guide' to learn how to add self-update functionality
  to your own CLI applications. See also: SELFUPDATE.md`,
		Version:              version,
		PersistentPreRun:     checkMultipleInstallations,
		SilenceUsage:         true,
	}

	rootCmd.AddCommand(cmd.InitCmd())
	rootCmd.AddCommand(cmd.ReleaseCmd())
	rootCmd.AddCommand(cmd.SyncCmd())
	rootCmd.AddCommand(cmd.InstallCmd())
	rootCmd.AddCommand(cmd.SelfUpdateCmd())
	rootCmd.AddCommand(cmd.SelfUpdateGuideCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// checkMultipleInstallations warns if multiple tobrew installations exist
func checkMultipleInstallations(cmd *cobra.Command, args []string) {
	// Get current executable path
	exePath, err := os.Executable()
	if err != nil {
		return
	}

	// Resolve symlinks
	currentPath, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		currentPath = exePath
	}

	// Use 'which -a' to find all tobrew installations in PATH
	whichCmd := exec.Command("which", "-a", "tobrew")
	output, err := whichCmd.Output()
	if err != nil {
		return // Silently ignore if which command fails
	}

	paths := strings.Split(strings.TrimSpace(string(output)), "\n")

	// Resolve symlinks for all found paths
	var otherPaths []string
	for _, path := range paths {
		if path == "" {
			continue
		}
		resolved, err := filepath.EvalSymlinks(path)
		if err != nil {
			resolved = path
		}
		// Skip if it's the same as current
		if resolved != currentPath {
			otherPaths = append(otherPaths, resolved)
		}
	}

	if len(otherPaths) > 0 {
		fmt.Println("⚠️  Warning: Multiple tobrew installations detected!")
		fmt.Printf("   Currently using: %s\n", currentPath)
		for _, path := range otherPaths {
			installType := "direct"
			if isHomebrewPath(path) {
				installType = "Homebrew"
			}
			fmt.Printf("   Also found (%s): %s\n", installType, path)
		}
		fmt.Println()
		fmt.Println("   Consider removing duplicate installations to avoid confusion.")
		fmt.Println("   Run 'which -a tobrew' to see all installations in your PATH.")
		fmt.Println()
	}
}

// isHomebrewPath checks if a path is a Homebrew installation
func isHomebrewPath(path string) bool {
	return strings.Contains(path, "/Cellar/") ||
		strings.Contains(path, "/opt/homebrew/") ||
		strings.Contains(path, "/usr/local/Cellar/") ||
		strings.Contains(path, "homebrew")
}
