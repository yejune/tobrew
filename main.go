package main

import (
	"fmt"
	"os"

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
  4. tobrew release --major   # Release with major bump (v1.1.0 → v2.0.0)`,
		Version: version,
	}

	rootCmd.AddCommand(cmd.InitCmd())
	rootCmd.AddCommand(cmd.ReleaseCmd())
	rootCmd.AddCommand(cmd.InstallCmd())
	rootCmd.AddCommand(cmd.SelfUpdateCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
