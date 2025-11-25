package version

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const lockFile = "tobrew.lock"

// Lock represents the tobrew.lock file
type Lock struct {
	Version     string    `yaml:"version"`
	LastRelease time.Time `yaml:"last_release"`
	SHA256      string    `yaml:"sha256,omitempty"`
}

// Load reads the lock file
func LoadLock() (*Lock, error) {
	data, err := os.ReadFile(lockFile)
	if err != nil {
		if os.IsNotExist(err) {
			// No lock file yet - start with v0.0.0
			return &Lock{
				Version: "v0.0.0",
			}, nil
		}
		return nil, err
	}

	var lock Lock
	if err := yaml.Unmarshal(data, &lock); err != nil {
		return nil, fmt.Errorf("failed to parse lock file: %w", err)
	}

	return &lock, nil
}

// Save writes the lock file
func (l *Lock) Save() error {
	data, err := yaml.Marshal(l)
	if err != nil {
		return err
	}

	return os.WriteFile(lockFile, data, 0644)
}

// BumpType represents version bump type
type BumpType int

const (
	BumpPatch BumpType = iota // 0.0.+1
	BumpMinor                 // 0.+1.0
	BumpMajor                 // +1.0.0
)

// Bump increments the version according to bump type
func (l *Lock) Bump(bumpType BumpType) (string, error) {
	current := l.Version
	if !strings.HasPrefix(current, "v") {
		return "", fmt.Errorf("invalid version format: %s (must start with 'v')", current)
	}

	// Parse version: v1.2.3 -> [1, 2, 3]
	parts := strings.Split(current[1:], ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid version format: %s (expected v1.2.3)", current)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", fmt.Errorf("invalid major version: %s", parts[0])
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid minor version: %s", parts[1])
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", fmt.Errorf("invalid patch version: %s", parts[2])
	}

	// Bump version
	switch bumpType {
	case BumpMajor:
		major++
		minor = 0
		patch = 0
	case BumpMinor:
		minor++
		patch = 0
	case BumpPatch:
		patch++
	}

	newVersion := fmt.Sprintf("v%d.%d.%d", major, minor, patch)
	l.Version = newVersion
	l.LastRelease = time.Now()

	return newVersion, nil
}

// UpdateSHA256 updates the SHA256 in lock file
func (l *Lock) UpdateSHA256(sha256 string) {
	l.SHA256 = sha256
}
