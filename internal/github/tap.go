package github

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/yejune/tobrew/internal/config"
)

// UpdateTap updates the homebrew-tap repository with the formula
func UpdateTap(cfg *config.Config, formulaContent string, version string) error {
	commitMsg := fmt.Sprintf("Update %s to %s", cfg.Name, version)
	return UpdateTapWithMessage(cfg, formulaContent, commitMsg)
}

// UpdateTapWithMessage updates tap with custom commit message
func UpdateTapWithMessage(cfg *config.Config, formulaContent string, commitMsg string) error {
	// Create temporary directory
	tmpDir := filepath.Join(os.TempDir(), "homebrew-tap-"+cfg.GitHub.TapRepo)

	// Clean up old tmp dir if exists
	os.RemoveAll(tmpDir)

	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return fmt.Errorf("failed to create tmp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize git repo
	if err := runCmd(tmpDir, "git", "init"); err != nil {
		return err
	}

	// Write formula
	formulaFile := filepath.Join(tmpDir, cfg.Name+".rb")
	if err := os.WriteFile(formulaFile, []byte(formulaContent), 0644); err != nil {
		return fmt.Errorf("failed to write formula: %w", err)
	}

	// Git add and commit
	if err := runCmd(tmpDir, "git", "add", cfg.Name+".rb"); err != nil {
		return err
	}

	if err := runCmd(tmpDir, "git", "commit", "-m", commitMsg); err != nil {
		return err
	}

	// Add remote and push
	tapURL := cfg.GetTapRepoURL()
	if err := runCmd(tmpDir, "git", "remote", "add", "origin", tapURL); err != nil {
		return err
	}

	if err := runCmd(tmpDir, "git", "branch", "-M", "main"); err != nil {
		return err
	}

	if err := runCmd(tmpDir, "git", "push", "-f", "origin", "main"); err != nil {
		return err
	}

	return nil
}

// runCmd executes a command in a specific directory
func runCmd(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
