# Self-Update Feature Implementation Guide

This guide helps you add a `self-update` command to your CLI tool that works seamlessly with both Homebrew installations and direct binary installations.

## Why Self-Update?

A self-update feature provides:
- **Convenience**: Users can update with a single command
- **Flexibility**: Works with both Homebrew and direct installations
- **User Experience**: No need to remember different update methods

## How It Works

Your self-update command should:

1. **Detect installation method**: Check if installed via Homebrew
2. **Use appropriate update method**:
   - Homebrew installation → Run `brew upgrade your-app`
   - Direct installation → Download and replace binary from GitHub

## Implementation (Go Example)

### 1. Add Self-Update Command

Create a new command file (e.g., `cmd/selfupdate.go`):

```go
package cmd

import (
  "fmt"
  "os"
  "os/exec"
  "runtime"
  "strings"

  "github.com/spf13/cobra"
  "your-project/internal/update"
)

var selfUpdateCmd = &cobra.Command{
  Use:     "self-update",
  Aliases: []string{"selfupdate"},  // Support both hyphenated and non-hyphenated
  Short:   "Update your-app to the latest version",
  Long: `Check for and install the latest version of your-app.

If installed via Homebrew, it will use 'brew upgrade'.
Otherwise, it will download and replace the binary directly.

Examples:
  your-app self-update`,
  RunE: func(cmd *cobra.Command, args []string) error {
    // Check if installed via Homebrew
    execPath, err := os.Executable()
    if err == nil && isHomebrewInstall(execPath) {
      fmt.Println("Detected Homebrew installation")
      fmt.Println("Running: brew upgrade your-app")

      brewCmd := exec.Command("brew", "upgrade", "your-app")
      brewCmd.Stdout = os.Stdout
      brewCmd.Stderr = os.Stderr
      return brewCmd.Run()
    }

    // Direct binary installation
    fmt.Println("Checking for updates...")
    repo := "username/your-app"
    currentVersion := Version // Your app's current version

    if err := update.PerformUpdate(repo, currentVersion, runtime.GOOS, runtime.GOARCH); err != nil {
      return fmt.Errorf("update failed: %w", err)
    }

    fmt.Println("Update completed successfully!")
    return nil
  },
}

func isHomebrewInstall(execPath string) bool {
  // Check common Homebrew installation paths
  return strings.Contains(execPath, "/Cellar/") ||
         strings.Contains(execPath, "/opt/homebrew/") ||
         strings.Contains(execPath, "homebrew")
}

func init() {
  rootCmd.AddCommand(selfUpdateCmd)
}
```

### 2. Create Update Package

Create `internal/update/update.go`:

```go
package update

import (
  "archive/tar"
  "compress/gzip"
  "encoding/json"
  "fmt"
  "io"
  "net/http"
  "os"
  "path/filepath"
  "runtime"
  "strings"
)

type GitHubRelease struct {
  TagName string `json:"tag_name"`
  Assets  []struct {
    Name               string `json:"name"`
    BrowserDownloadURL string `json:"browser_download_url"`
  } `json:"assets"`
}

// PerformUpdate downloads and installs the latest version
func PerformUpdate(repo, currentVersion, goos, goarch string) error {
  // Get latest release info from GitHub
  latestRelease, err := getLatestRelease(repo)
  if err != nil {
    return fmt.Errorf("failed to get latest release: %w", err)
  }

  fmt.Printf("Current version: %s\n", currentVersion)
  fmt.Printf("Latest version:  %s\n", latestRelease.TagName)

  // Check if already up to date
  if currentVersion == latestRelease.TagName {
    fmt.Println("You are already running the latest version!")
    return nil
  }

  // Find appropriate asset for current platform
  assetURL, assetName := findAsset(latestRelease, goos, goarch)
  if assetURL == "" {
    return fmt.Errorf("no release found for %s/%s", goos, goarch)
  }

  fmt.Printf("Downloading %s...\n", assetName)

  // Download and extract
  if err := downloadAndInstall(assetURL, assetName); err != nil {
    return fmt.Errorf("failed to install: %w", err)
  }

  return nil
}

func getLatestRelease(repo string) (*GitHubRelease, error) {
  url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

  resp, err := http.Get(url)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
  }

  var release GitHubRelease
  if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
    return nil, err
  }

  return &release, nil
}

func findAsset(release *GitHubRelease, goos, goarch string) (string, string) {
  // Look for asset matching OS and architecture
  // Example: your-app_1.0.0_darwin_amd64.tar.gz
  for _, asset := range release.Assets {
    name := strings.ToLower(asset.Name)
    if strings.Contains(name, goos) && strings.Contains(name, goarch) {
      return asset.BrowserDownloadURL, asset.Name
    }
  }
  return "", ""
}

func downloadAndInstall(url, filename string) error {
  // Download file
  resp, err := http.Get(url)
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  // Create temp file
  tmpFile, err := os.CreateTemp("", "update-*.tar.gz")
  if err != nil {
    return err
  }
  defer os.Remove(tmpFile.Name())
  defer tmpFile.Close()

  // Write to temp file
  if _, err := io.Copy(tmpFile, resp.Body); err != nil {
    return err
  }

  // Reset file pointer
  if _, err := tmpFile.Seek(0, 0); err != nil {
    return err
  }

  // Extract binary
  binary, err := extractBinary(tmpFile)
  if err != nil {
    return err
  }
  defer os.Remove(binary)

  // Replace current executable
  execPath, err := os.Executable()
  if err != nil {
    return err
  }

  // Make backup
  backup := execPath + ".backup"
  if err := os.Rename(execPath, backup); err != nil {
    return fmt.Errorf("failed to backup current binary: %w", err)
  }

  // Copy new binary
  if err := copyFile(binary, execPath); err != nil {
    // Restore backup on failure
    os.Rename(backup, execPath)
    return fmt.Errorf("failed to install new binary: %w", err)
  }

  // Make executable
  if err := os.Chmod(execPath, 0755); err != nil {
    return err
  }

  // Remove backup
  os.Remove(backup)

  fmt.Println("Successfully updated!")
  return nil
}

func extractBinary(tarGz *os.File) (string, error) {
  gzr, err := gzip.NewReader(tarGz)
  if err != nil {
    return "", err
  }
  defer gzr.Close()

  tr := tar.NewReader(gzr)

  for {
    header, err := tr.Next()
    if err == io.EOF {
      break
    }
    if err != nil {
      return "", err
    }

    // Look for the main binary (not .md, .txt, etc.)
    if header.Typeflag == tar.TypeReg && !strings.Contains(header.Name, ".") {
      tmpBinary, err := os.CreateTemp("", "binary-*")
      if err != nil {
        return "", err
      }
      defer tmpBinary.Close()

      if _, err := io.Copy(tmpBinary, tr); err != nil {
        return "", err
      }

      return tmpBinary.Name(), nil
    }
  }

  return "", fmt.Errorf("binary not found in archive")
}

func copyFile(src, dst string) error {
  source, err := os.Open(src)
  if err != nil {
    return err
  }
  defer source.Close()

  destination, err := os.Create(dst)
  if err != nil {
    return err
  }
  defer destination.Close()

  _, err = io.Copy(destination, source)
  return err
}
```

## Platform-Specific Considerations

### macOS Homebrew Paths

Homebrew installs to different locations:
- **Intel Macs**: `/usr/local/Cellar/your-app/`
- **Apple Silicon**: `/opt/homebrew/Cellar/your-app/`

Both can be detected by checking if the executable path contains `/Cellar/` or `/opt/homebrew/`.

### Linux Package Managers

For Linux, you might also want to detect:
```go
func isPackageManagerInstall(execPath string) bool {
  return strings.HasPrefix(execPath, "/usr/bin/") ||
         strings.HasPrefix(execPath, "/usr/local/bin/") ||
         strings.Contains(execPath, "linuxbrew")
}
```

## Testing

### Test Homebrew Installation

```bash
# Install via Homebrew
brew install username/tap/your-app

# Test self-update
your-app self-update
# Should output: "Detected Homebrew installation"
# Should run: "brew upgrade your-app"
```

### Test Direct Installation

```bash
# Install directly
curl -sSL https://your-domain.com/install.sh | sh

# Test self-update
your-app self-update
# Should download from GitHub and replace binary
```

## Best Practices

### 1. Version Detection and Build-Time Injection

To make self-update work properly, your CLI needs to know its current version. This is done by injecting the version at build time using Go's `-ldflags`.

#### In Your Code

```go
// main.go or cmd/root.go
package main

var Version = "dev" // Default value for development
```

#### When Building Manually

```bash
go build -ldflags "-X main.Version=v1.0.0" -o myapp
```

#### In Your Homebrew Formula

**This is crucial for Homebrew installations!** Your `tobrew.yaml` formula should include ldflags:

```yaml
formula:
  install: |
    ldflags = "-s -w -X main.Version=#{version}"
    system "go", "build", *std_go_args(ldflags:, output: bin/"myapp"), "."
```

Breaking it down:
- `-s -w`: Strip debug info and symbol table (reduces binary size)
- `-X main.Version=#{version}`: Inject Homebrew's version variable
- `*std_go_args(ldflags:, ...)`: Homebrew's standard Go build args with your ldflags
- `#{version}`: Ruby interpolation - Homebrew automatically provides this (e.g., "v1.2.3")

#### Different Package Paths

If your `Version` variable is in a different package:

```go
// cmd/version.go
package cmd

var Version = "dev"
```

Then in your formula:
```ruby
ldflags = "-s -w -X github.com/username/myapp/cmd.Version=#{version}"
system "go", "build", *std_go_args(ldflags:, output: bin/"myapp"), "."
```

The format is: `-X <package-path>.<variable>=<value>`

#### Verifying Version Injection

After Homebrew installation, users should see the correct version:
```bash
$ myapp --version
myapp v1.2.3
```

Not:
```bash
$ myapp --version
myapp dev  # ❌ This means ldflags didn't work
```

### 2. Safe Binary Replacement

Always:
- Create a backup before replacing
- Verify download integrity (checksums if available)
- Restore backup on failure
- Make new binary executable

### 3. User Feedback

Provide clear messages:
```go
fmt.Println("Detected Homebrew installation")
fmt.Println("Current version: v1.0.0")
fmt.Println("Latest version: v1.1.0")
fmt.Println("Downloading...")
fmt.Println("Successfully updated!")
```

### 4. Error Handling

Handle common errors:
- Network failures
- Permission errors (suggest `sudo` if needed)
- GitHub API rate limits
- Invalid/missing releases

### 5. Permissions

If the user lacks permissions to replace the binary:
```go
if err := os.Rename(execPath, backup); err != nil {
  if os.IsPermission(err) {
    return fmt.Errorf("permission denied. Try running with sudo: sudo %s self-update", os.Args[0])
  }
  return err
}
```

### 6. Detecting Multiple Installations

Users may have both Homebrew and direct installations, causing confusion. Detect and warn about this:

```go
// In main.go or root command
func checkMultipleInstallations(cmd *cobra.Command, args []string) {
    // Get current executable path
    exePath, err := os.Executable()
    if err != nil {
        return
    }

    currentPath, err := filepath.EvalSymlinks(exePath)
    if err != nil {
        currentPath = exePath
    }

    // Find all installations in PATH
    whichCmd := exec.Command("which", "-a", "your-app")
    output, err := whichCmd.Output()
    if err != nil {
        return
    }

    paths := strings.Split(strings.TrimSpace(string(output)), "\n")

    // Resolve symlinks and collect unique paths
    var resolvedPaths []string
    seen := make(map[string]bool)
    currentInPath := false

    for _, path := range paths {
        if path == "" {
            continue
        }
        resolved, err := filepath.EvalSymlinks(path)
        if err != nil {
            resolved = path
        }
        if !seen[resolved] {
            resolvedPaths = append(resolvedPaths, resolved)
            seen[resolved] = true
        }
        if resolved == currentPath {
            currentInPath = true
        }
    }

    // Only warn if:
    // 1. Multiple installations in PATH
    // 2. Current executable is in PATH (not running from dev directory)
    if len(resolvedPaths) > 1 && currentInPath {
        fmt.Println("⚠️  Warning: Multiple installations detected!")
        for _, path := range resolvedPaths {
            installType := "direct"
            if strings.Contains(path, "homebrew") || strings.Contains(path, "/Cellar/") {
                installType = "Homebrew"
            }
            fmt.Printf("   - %s (%s)\n", path, installType)
        }
        fmt.Println()
        fmt.Println("   Consider removing duplicates to avoid confusion.")
        fmt.Println("   Run 'which your-app' to see which one is currently active.")
    }
}

// In your root command
rootCmd := &cobra.Command{
    Use:              "your-app",
    PersistentPreRun: checkMultipleInstallations,
    // ...
}
```

**Why this matters:**
- User installs via Homebrew: `/opt/homebrew/bin/your-app`
- Later does direct install: `/usr/local/bin/your-app`
- PATH determines which runs: `which your-app` shows first match
- Updates may affect wrong installation
- Version confusion: `your-app --version` shows one, but Homebrew has different version

**Best practice:**
1. Check on every command run (use `PersistentPreRun`)
2. Only warn if multiple installations **in PATH** exist
3. Skip warning if running from dev directory (not in PATH)
4. List all installations with their types (Homebrew/direct)
5. Let user decide which to keep - don't auto-remove

## Alternative: Using goreleaser with Built-in Updates

If you're using [goreleaser](https://goreleaser.com/), you can use their built-in update feature:

```yaml
# .goreleaser.yml
brews:
  - name: your-app
    homepage: https://github.com/username/your-app
    repository:
      owner: username
      name: homebrew-tap

# Use minio/selfupdate package
# See: https://github.com/minio/selfupdate
```

## Resources

- [goreleaser](https://goreleaser.com/) - Release automation tool
- [minio/selfupdate](https://github.com/minio/selfupdate) - Self-update library for Go
- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook) - Homebrew formula guide

## Summary

A good self-update implementation:
1. ✅ Detects installation method (Homebrew vs direct)
2. ✅ Uses `brew upgrade` for Homebrew installations
3. ✅ Downloads and replaces binary for direct installations
4. ✅ Provides clear feedback to users
5. ✅ Handles errors gracefully
6. ✅ Creates backups before replacing binaries
7. ✅ Works across different platforms (macOS, Linux, Windows)

With this guide, your CLI tool users will have a seamless update experience regardless of how they installed your tool!
