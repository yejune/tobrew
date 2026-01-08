package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

const (
	githubAPI     = "https://api.github.com/repos/yejune/tobrew/releases/latest"
	downloadURL   = "https://github.com/yejune/tobrew/releases/download/%s/tobrew-%s-%s"
)

func SelfUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "selfupdate",
		Short: "Update tobrew to the latest version",
		Long: `Update tobrew to the latest version.

If installed via Homebrew, it will use 'brew upgrade'.
Otherwise, it will download and replace the binary directly from GitHub releases.

Example:
  tobrew selfupdate`,
		RunE: runSelfUpdate,
	}

	return cmd
}

func runSelfUpdate(cmd *cobra.Command, args []string) error {
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

	// Check if installed via Homebrew
	if isHomebrewInstall(exePath) {
		fmt.Println("🍺 Detected Homebrew installation")
		fmt.Println("Running: brew upgrade tobrew")
		fmt.Println()

		brewCmd := exec.Command("brew", "upgrade", "yejune/tap/tobrew")
		brewCmd.Stdout = os.Stdout
		brewCmd.Stderr = os.Stderr
		brewCmd.Stdin = os.Stdin

		if err := brewCmd.Run(); err != nil {
			return fmt.Errorf("brew upgrade failed: %w", err)
		}

		fmt.Println()
		fmt.Println("✅ Update complete!")
		return nil
	}

	// Direct installation - download from GitHub
	fmt.Println("🔍 Checking for updates...")

	// Get latest version from GitHub
	latestVersion, err := getLatestVersion()
	if err != nil {
		return fmt.Errorf("failed to check latest version: %w", err)
	}

	fmt.Printf("Latest version: %s\n", latestVersion)
	fmt.Printf("Current path:   %s\n\n", exePath)

	// Determine platform
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// Build download URL
	downloadPath := fmt.Sprintf(downloadURL, latestVersion, goos, goarch)

	fmt.Printf("📦 Downloading %s...\n", latestVersion)

	// Download binary
	tmpFile := filepath.Join(os.TempDir(), "tobrew-new")
	if err := downloadFile(tmpFile, downloadPath); err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer os.Remove(tmpFile)

	// Make executable
	if err := os.Chmod(tmpFile, 0755); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	fmt.Println("✓ Download complete")
	fmt.Println()

	// Replace binary
	fmt.Println("📋 Installing update...")

	// If installed in /usr/local/bin, need sudo
	needsSudo := strings.HasPrefix(exePath, "/usr/local")

	var replaceCmd *exec.Cmd
	if needsSudo {
		fmt.Println("(sudo required)")
		replaceCmd = exec.Command("sudo", "mv", tmpFile, exePath)
	} else {
		replaceCmd = exec.Command("mv", tmpFile, exePath)
	}

	replaceCmd.Stdin = os.Stdin
	replaceCmd.Stdout = os.Stdout
	replaceCmd.Stderr = os.Stderr

	if err := replaceCmd.Run(); err != nil {
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	// Ensure executable
	if needsSudo {
		chmodCmd := exec.Command("sudo", "chmod", "+x", exePath)
		chmodCmd.Run()
	}

	fmt.Println()
	fmt.Println("✅ Update complete!")
	fmt.Printf("tobrew has been updated to %s\n", latestVersion)

	return nil
}

func getLatestVersion() (string, error) {
	resp, err := http.Get(githubAPI)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	// Simple parsing - just get tag_name from JSON
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	bodyStr := string(body)

	// Find "tag_name":"v1.2.3"
	tagStart := strings.Index(bodyStr, `"tag_name":"`)
	if tagStart == -1 {
		return "", fmt.Errorf("could not find tag_name in response")
	}

	tagStart += len(`"tag_name":"`)
	tagEnd := strings.Index(bodyStr[tagStart:], `"`)
	if tagEnd == -1 {
		return "", fmt.Errorf("could not parse tag_name")
	}

	return bodyStr[tagStart : tagStart+tagEnd], nil
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// isHomebrewInstall checks if the binary is installed via Homebrew
func isHomebrewInstall(execPath string) bool {
	// Check common Homebrew installation paths
	return strings.Contains(execPath, "/Cellar/") ||
		strings.Contains(execPath, "/opt/homebrew/") ||
		strings.Contains(execPath, "/usr/local/Cellar/") ||
		strings.Contains(execPath, "homebrew")
}
