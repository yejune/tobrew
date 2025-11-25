package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func InstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install tobrew to /usr/local/bin",
		Long: `Install tobrew binary to /usr/local/bin.

This makes tobrew available system-wide.

Example:
  tobrew install`,
		RunE: runInstall,
	}

	return cmd
}

func runInstall(cmd *cobra.Command, args []string) error {
	// Get current executable path
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Resolve symlinks
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	targetPath := "/usr/local/bin/tobrew"

	fmt.Println("📦 Installing tobrew...")
	fmt.Printf("From: %s\n", exePath)
	fmt.Printf("To:   %s\n\n", targetPath)

	// Check if already installed
	if exePath == targetPath {
		fmt.Println("✓ tobrew is already installed at /usr/local/bin/tobrew")
		return nil
	}

	// Copy binary with sudo
	fmt.Println("Copying binary (sudo required)...")
	copyCmd := exec.Command("sudo", "cp", exePath, targetPath)
	copyCmd.Stdin = os.Stdin
	copyCmd.Stdout = os.Stdout
	copyCmd.Stderr = os.Stderr

	if err := copyCmd.Run(); err != nil {
		return fmt.Errorf("failed to copy binary: %w", err)
	}

	// Set executable permissions
	chmodCmd := exec.Command("sudo", "chmod", "+x", targetPath)
	if err := chmodCmd.Run(); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	fmt.Println()
	fmt.Println("✅ Installation complete!")
	fmt.Printf("tobrew is now available at: %s\n", targetPath)

	return nil
}
