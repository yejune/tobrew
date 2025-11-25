package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yejune/tobrew/internal/version"
)

func SyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync lock file with remote tags",
		Long: `Sync the local tobrew.lock file with remote git tags.

This is useful when:
  - Lock file is out of sync with actual releases
  - Working on a different machine
  - Recovering from a failed release

The command will:
  1. Fetch remote tags
  2. Find the latest version tag
  3. Update tobrew.lock with the latest version`,
		RunE: runSync,
	}

	return cmd
}

func runSync(cmd *cobra.Command, args []string) error {
	// Load lock file
	lock, err := version.LoadLock()
	if err != nil {
		return fmt.Errorf("failed to load lock file: %w", err)
	}

	currentVersion := lock.Version
	fmt.Printf("📋 Current lock version: %s\n", currentVersion)

	// Get latest remote tag
	fmt.Println("🔄 Fetching remote tags...")
	latestTag, err := getLatestRemoteTag()
	if err != nil {
		return fmt.Errorf("failed to get remote tags: %w", err)
	}

	fmt.Printf("   Latest remote tag: %s\n", latestTag)

	// Compare and update
	if compareVersions(latestTag, currentVersion) > 0 {
		lock.Version = latestTag
		lock.UpdateFingerprint()
		if err := lock.Save(); err != nil {
			return fmt.Errorf("failed to save lock file: %w", err)
		}
		fmt.Printf("\n✅ Lock file updated: %s → %s\n", currentVersion, latestTag)
	} else if compareVersions(latestTag, currentVersion) == 0 {
		lock.UpdateFingerprint()
		if err := lock.Save(); err != nil {
			return fmt.Errorf("failed to save lock file: %w", err)
		}
		fmt.Println("\n✅ Already in sync")
	} else {
		fmt.Printf("\n⚠️  Lock file (%s) is ahead of remote (%s)\n", currentVersion, latestTag)
	}

	return nil
}
