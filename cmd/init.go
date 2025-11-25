package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"github.com/yejune/tobrew/internal/config"
	"gopkg.in/yaml.v3"
)

var (
	formatFlag string
	outputFlag string
)

func InitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new tobrew configuration file",
		Long: `Create a new tobrew configuration file with default values.

Supports multiple formats:
  - yaml (default)
  - json
  - toml

Example:
  tobrew init
  tobrew init --format json
  tobrew init --format toml -o release.toml`,
		RunE: runInit,
	}

	cmd.Flags().StringVarP(&formatFlag, "format", "f", "yaml", "Config file format (yaml, json, toml)")
	cmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Output file path (default: tobrew.{format})")

	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	// Validate format
	if formatFlag != "yaml" && formatFlag != "json" && formatFlag != "toml" {
		return fmt.Errorf("unsupported format: %s (use yaml, json, or toml)", formatFlag)
	}

	// Determine output file
	outputFile := outputFlag
	if outputFile == "" {
		outputFile = "tobrew." + formatFlag
	}

	// Check if file already exists
	if _, err := os.Stat(outputFile); err == nil {
		return fmt.Errorf("file already exists: %s (remove it first or use -o to specify different path)", outputFile)
	}

	// Try to detect project name from directory or go.mod
	projectName := detectProjectName()

	// Create default config
	cfg := &config.Config{
		Name:        projectName,
		Description: "Description of your project",
		Homepage:    fmt.Sprintf("https://github.com/USERNAME/%s", projectName),
		License:     "MIT",
		GitHub: config.GitHubConfig{
			User:    "USERNAME",
			Repo:    projectName,
			TapRepo: "homebrew-tap",
		},
		Build: config.BuildConfig{
			Command: "go build -o build/{{.Name}} .",
		},
		Formula: config.FormulaConfig{
			Install: fmt.Sprintf(`system "go", "build", "."
    bin.install "%s"`, projectName),
			Test: fmt.Sprintf(`assert_match "%s", shell_output("#{bin}/%s --version")`, projectName, projectName),
			Caveats: fmt.Sprintf(`%s has been installed!

    Run '%s --help' to get started.`, projectName, projectName),
		},
	}

	// Write config in requested format
	var data []byte
	var err error

	switch formatFlag {
	case "yaml":
		data, err = yaml.Marshal(cfg)
	case "json":
		data, err = json.MarshalIndent(cfg, "", "  ")
	case "toml":
		data, err = toml.Marshal(cfg)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("✓ Created %s\n", outputFile)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Edit the config file and update USERNAME, description, etc.")
	fmt.Println("  2. Create GitHub repository named 'homebrew-tap'")
	fmt.Printf("  3. Create a release: tobrew release\n")

	return nil
}

func detectProjectName() string {
	// Try current directory name
	dir, err := os.Getwd()
	if err == nil {
		return filepath.Base(dir)
	}

	// Try go.mod
	if data, err := os.ReadFile("go.mod"); err == nil {
		// Simple parsing - just get first line "module github.com/user/project"
		lines := string(data)
		var modulePath string
		fmt.Sscanf(lines, "module %s", &modulePath)
		if modulePath != "" {
			parts := filepath.SplitList(modulePath)
			if len(parts) > 0 {
				return filepath.Base(parts[len(parts)-1])
			}
		}
	}

	return "myapp"
}
