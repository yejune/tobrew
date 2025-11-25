package cmd

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yejune/tobrew/internal/config"
	"github.com/yejune/tobrew/internal/formula"
	"github.com/yejune/tobrew/internal/github"
	"github.com/yejune/tobrew/internal/version"
)

var (
	majorFlag bool
	minorFlag bool
	patchFlag bool
)

func ReleaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Create a complete release and update Homebrew tap",
		Long: `Create a complete release with automatic version bumping:

By default, increments patch version (v1.0.0 → v1.0.1)

Examples:
  tobrew release              # Patch: v1.0.0 → v1.0.1
  tobrew release --minor      # Minor: v1.0.1 → v1.1.0
  tobrew release --major      # Major: v1.1.0 → v2.0.0

The version is automatically managed in tobrew.lock file.

Process:
  1. Load current version from tobrew.lock
  2. Bump version according to flags
  3. Build the project
  4. Create and push git tag
  5. Download release tarball and calculate SHA256
  6. Generate Homebrew formula
  7. Update homebrew-tap repository
  8. Save new version to tobrew.lock`,
		RunE: runRelease,
	}

	cmd.Flags().BoolVar(&majorFlag, "major", false, "Increment major version (v1.0.0 → v2.0.0)")
	cmd.Flags().BoolVar(&minorFlag, "minor", false, "Increment minor version (v1.0.0 → v1.1.0)")
	cmd.Flags().BoolVar(&patchFlag, "patch", false, "Increment patch version (v1.0.0 → v1.0.1) - default")

	return cmd
}

func runRelease(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.Load("")
	if err != nil {
		return err
	}

	// Load lock file
	lock, err := version.LoadLock()
	if err != nil {
		return fmt.Errorf("failed to load version: %w", err)
	}

	// Determine bump type
	bumpType := version.BumpPatch // default
	if majorFlag {
		bumpType = version.BumpMajor
	} else if minorFlag {
		bumpType = version.BumpMinor
	}

	// Check for conflicting flags
	flagCount := 0
	if majorFlag {
		flagCount++
	}
	if minorFlag {
		flagCount++
	}
	if patchFlag {
		flagCount++
	}
	if flagCount > 1 {
		return fmt.Errorf("cannot use multiple version bump flags together")
	}

	// Bump version
	newVersion, err := lock.Bump(bumpType)
	if err != nil {
		return fmt.Errorf("failed to bump version: %w", err)
	}

	fmt.Printf("🚀 Starting release process for %s\n", cfg.Name)
	fmt.Printf("   Current version: %s\n", lock.Version)
	fmt.Printf("   New version:     %s\n\n", newVersion)

	// Confirm
	fmt.Print("Continue? (Y/n): ")
	var response string
	fmt.Scanln(&response)
	if response != "" && strings.ToLower(response) != "y" {
		return fmt.Errorf("release cancelled")
	}

	// Step 1: Build
	fmt.Println("\n📦 Building project...")
	if err := buildProject(cfg); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}
	fmt.Println("✓ Build successful")

	// Step 2: Git tag
	fmt.Printf("\n🏷️  Creating git tag %s...\n", newVersion)
	if err := createGitTag(newVersion); err != nil {
		return fmt.Errorf("git tag failed: %w", err)
	}
	fmt.Println("✓ Git tag created and pushed")

	// Wait for GitHub to process the tag
	fmt.Println("\n⏳ Waiting for GitHub to process the release...")
	time.Sleep(5 * time.Second)

	// Step 3: Download and calculate SHA256
	fmt.Println("\n🔐 Calculating SHA256 checksum...")
	tarballURL := cfg.GetTarballURL(newVersion)
	sha256sum, err := downloadAndHash(tarballURL)
	if err != nil {
		return fmt.Errorf("failed to download/hash tarball: %w", err)
	}
	fmt.Printf("✓ SHA256: %s\n", sha256sum)

	// Update lock file with SHA256
	lock.UpdateSHA256(sha256sum)

	// Step 4: Generate formula
	fmt.Println("\n📝 Generating Homebrew formula...")
	formulaContent, err := formula.Generate(cfg, newVersion, sha256sum)
	if err != nil {
		return fmt.Errorf("formula generation failed: %w", err)
	}

	// Save formula locally for review
	formulaFile := cfg.Name + ".rb"
	if err := os.WriteFile(formulaFile, []byte(formulaContent), 0644); err != nil {
		return fmt.Errorf("failed to write formula: %w", err)
	}
	fmt.Printf("✓ Formula generated: %s\n", formulaFile)

	// Step 5: Update homebrew-tap
	fmt.Println("\n🍺 Updating homebrew-tap repository...")
	if err := github.UpdateTap(cfg, formulaContent, newVersion); err != nil {
		return fmt.Errorf("tap update failed: %w", err)
	}
	fmt.Println("✓ Homebrew tap updated")

	// Step 6: Save lock file
	fmt.Println("\n💾 Saving version lock file...")
	if err := lock.Save(); err != nil {
		return fmt.Errorf("failed to save lock file: %w", err)
	}
	fmt.Println("✓ Version saved to tobrew.lock")

	// Success!
	fmt.Println("\n✅ Release complete!")
	fmt.Println()
	fmt.Printf("Version:  %s\n", newVersion)
	fmt.Printf("Released: %s\n", lock.LastRelease.Format(time.RFC3339))
	fmt.Println()
	fmt.Printf("Users can now install with:\n")
	fmt.Printf("  brew install %s/tap/%s\n", cfg.GitHub.User, cfg.Name)
	fmt.Println()
	fmt.Printf("Or upgrade with:\n")
	fmt.Printf("  brew upgrade %s\n", cfg.Name)

	return nil
}

func buildProject(cfg *config.Config) error {
	if cfg.Build.Command == "" {
		return fmt.Errorf("build.command not specified in config")
	}

	// Simple template replacement
	cmdStr := strings.ReplaceAll(cfg.Build.Command, "{{.Name}}", cfg.Name)

	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func createGitTag(version string) error {
	// Create annotated tag
	tagCmd := exec.Command("git", "tag", "-a", version, "-m", "Release "+version)
	tagCmd.Stdout = os.Stdout
	tagCmd.Stderr = os.Stderr
	if err := tagCmd.Run(); err != nil {
		return err
	}

	// Push tag
	pushCmd := exec.Command("git", "push", "origin", version)
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	return pushCmd.Run()
}

func downloadAndHash(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	h := sha256.New()
	if _, err := io.Copy(h, resp.Body); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
