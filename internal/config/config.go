package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the tobrew configuration file
type Config struct {
	Name        string        `yaml:"name"`
	Language    string        `yaml:"language,omitempty"` // go, rust, python, node, php, binary
	Description string        `yaml:"description"`
	Homepage    string        `yaml:"homepage"`
	License     string        `yaml:"license"`
	GitHub      GitHubConfig  `yaml:"github"`
	Build       BuildConfig   `yaml:"build"`
	Formula     FormulaConfig `yaml:"formula"`
}

type GitHubConfig struct {
	User    string `yaml:"user"`
	Repo    string `yaml:"repo"`
	TapRepo string `yaml:"tap_repo"`
}

type BuildConfig struct {
	Command string `yaml:"command"`
}

type FormulaConfig struct {
	Install string `yaml:"install"`
	Test    string `yaml:"test"`
	Caveats string `yaml:"caveats"`
}

// Load reads and parses the tobrew.yaml config file
func Load(path string) (*Config, error) {
	if path == "" {
		path = "tobrew.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate required fields
	if config.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if config.GitHub.User == "" {
		return nil, fmt.Errorf("github.user is required")
	}
	if config.GitHub.Repo == "" {
		return nil, fmt.Errorf("github.repo is required")
	}
	if config.GitHub.TapRepo == "" {
		return nil, fmt.Errorf("github.tap_repo is required")
	}

	// Default language to "go" if not specified
	if config.Language == "" {
		config.Language = "go"
	}

	return &config, nil
}

// Save writes the config to a file
func (c *Config) Save(path string) error {
	if path == "" {
		path = "tobrew.yaml"
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetTarballURL returns the GitHub tarball URL for a version
func (c *Config) GetTarballURL(version string) string {
	return fmt.Sprintf("https://github.com/%s/%s/archive/refs/tags/%s.tar.gz",
		c.GitHub.User, c.GitHub.Repo, version)
}

// GetTapRepoURL returns the GitHub tap repository URL
func (c *Config) GetTapRepoURL() string {
	return fmt.Sprintf("https://github.com/%s/%s.git",
		c.GitHub.User, c.GitHub.TapRepo)
}

// GetFormulaName returns the Ruby class name for the formula
func (c *Config) GetFormulaName() string {
	return toCamelCase(c.Name)
}

// toCamelCase converts "my-app" to "MyApp"
func toCamelCase(s string) string {
	result := ""
	capitalize := true

	for _, ch := range s {
		if ch == '-' || ch == '_' {
			capitalize = true
			continue
		}

		if capitalize {
			result += string(ch - 32) // Convert to uppercase
			capitalize = false
		} else {
			result += string(ch)
		}
	}

	return result
}

// ProjectRoot finds the project root directory (contains go.mod or .git)
func ProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Check for go.mod or .git
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("project root not found (no go.mod or .git)")
		}
		dir = parent
	}
}
