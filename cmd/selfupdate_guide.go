package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func SelfUpdateGuideCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "selfupdate-guide",
		Short: "Show guide for implementing self-update in your CLI tool",
		Long: `Display information about implementing self-update functionality in your own CLI tools.

This guide explains:
  - How to add self-update commands to your Go CLI applications
  - Detecting Homebrew vs direct installations
  - Downloading and replacing binaries safely
  - Best practices and platform considerations

For detailed implementation guide, see SELFUPDATE.md in the tobrew repository:
https://github.com/yejune/tobrew/blob/main/SELFUPDATE.md`,
		Run: runSelfUpdateGuide,
	}

	return cmd
}

func runSelfUpdateGuide(cmd *cobra.Command, args []string) {
	fmt.Println(`
╭─────────────────────────────────────────────────────────────╮
│                                                             │
│         📖  Self-Update Implementation Guide                │
│                                                             │
╰─────────────────────────────────────────────────────────────╯

This guide helps you add a self-update command to your own CLI tool
that works seamlessly with both Homebrew and direct installations.

📚 Full Documentation:
   SELFUPDATE.md in the tobrew repository

🔗 Online:
   https://github.com/yejune/tobrew/blob/main/SELFUPDATE.md

📦 What you'll learn:
   • Detecting installation method (Homebrew vs direct)
   • Using 'brew upgrade' for Homebrew installations
   • Downloading and replacing binaries from GitHub
   • Platform-specific considerations (macOS, Linux)
   • Safe binary replacement with backups
   • Error handling and user feedback

💡 Quick Start:
   1. Add a 'selfupdate' command to your CLI
   2. Detect if installed via Homebrew (check path)
   3. If Homebrew: run 'brew upgrade your-app'
   4. If direct: download from GitHub releases and replace binary

⚙️  Example Implementation:
   See SELFUPDATE.md for complete Go code examples including:
   - Command setup with cobra
   - GitHub releases API integration
   - Safe binary replacement logic
   - Platform detection

🛠️  tobrew's Implementation:
   This tool uses the same pattern! Check out:
   • cmd/selfupdate.go - Our implementation
   • 'tobrew self-update' - Try it yourself

For detailed code examples, best practices, and troubleshooting,
please refer to SELFUPDATE.md in the repository.
`)
}
