package cmd

import (
	"fmt"
	"io"
	"net/http"

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
	const guideURL = "https://raw.githubusercontent.com/yejune/tobrew/main/SELFUPDATE.md"

	fmt.Println("📖 Self-Update Implementation Guide")
	fmt.Println("=====================================")
	fmt.Println()
	fmt.Println("Fetching latest guide from GitHub...")
	fmt.Println()

	// Fetch the guide from GitHub
	resp, err := http.Get(guideURL)
	if err != nil {
		fmt.Printf("⚠️  Could not fetch guide: %v\n", err)
		fmt.Println()
		fmt.Println("View online at:")
		fmt.Println("https://github.com/yejune/tobrew/blob/main/SELFUPDATE.md")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("⚠️  Could not fetch guide (status %d)\n", resp.StatusCode)
		fmt.Println()
		fmt.Println("View online at:")
		fmt.Println("https://github.com/yejune/tobrew/blob/main/SELFUPDATE.md")
		return
	}

	// Read and display the guide
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("⚠️  Could not read guide: %v\n", err)
		return
	}

	fmt.Println(string(content))
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("💡 This guide is always up-to-date from GitHub")
	fmt.Println("   View online: https://github.com/yejune/tobrew/blob/main/SELFUPDATE.md")
}
