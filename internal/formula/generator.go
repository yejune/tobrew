package formula

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/yejune/tobrew/internal/config"
)

const defaultTemplate = `class {{.ClassName}} < Formula
  desc "{{.Description}}"
  homepage "{{.Homepage}}"
  url "{{.URL}}"
  sha256 "{{.SHA256}}"
  license "{{.License}}"
  head "{{.HeadURL}}", branch: "main"
{{if .DependsOn}}
  {{.DependsOn}}
{{end}}
  def install
    {{.InstallScript}}
  end
{{if .TestScript}}
  def test
    {{.TestScript}}
  end
{{end}}{{if .Caveats}}
  def caveats
    <<~EOS
      {{.Caveats}}
    EOS
  end
{{end}}end
`

type TemplateData struct {
	ClassName     string
	Description   string
	Homepage      string
	URL           string
	SHA256        string
	License       string
	HeadURL       string
	DependsOn     string
	InstallScript string
	TestScript    string
	Caveats       string
}

// Generate creates a Homebrew formula from config
func Generate(cfg *config.Config, version string, sha256sum string) (string, error) {
	data := TemplateData{
		ClassName:     cfg.GetFormulaName(),
		Description:   cfg.Description,
		Homepage:      cfg.Homepage,
		URL:           cfg.GetTarballURL(version),
		SHA256:        sha256sum,
		License:       cfg.License,
		HeadURL:       fmt.Sprintf("https://github.com/%s/%s.git", cfg.GitHub.User, cfg.GitHub.Repo),
		DependsOn:     getDependency(cfg.Language),
		InstallScript: indentScript(cfg.Formula.Install, 4),
		TestScript:    indentScript(cfg.Formula.Test, 4),
		Caveats:       indentLines(cfg.Formula.Caveats, 6),
	}

	tmpl, err := template.New("formula").Parse(defaultTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// indentScript adds proper indentation to Ruby code
func indentScript(script string, spaces int) string {
	if script == "" {
		return ""
	}

	lines := strings.Split(script, "\n")
	indent := strings.Repeat(" ", spaces)

	var result []string
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			result = append(result, "")
		} else {
			result = append(result, indent+line)
		}
	}

	return strings.Join(result, "\n")
}

// indentLines indents each line for caveats
func indentLines(text string, spaces int) string {
	if text == "" {
		return ""
	}

	lines := strings.Split(text, "\n")
	indent := strings.Repeat(" ", spaces)

	var result []string
	for _, line := range lines {
		result = append(result, indent+line)
	}

	return strings.Join(result, "\n")
}

// getDependency returns the appropriate depends_on line for the language
func getDependency(language string) string {
	// Check for version-specific formats (e.g., php@8.4, python@3.11)
	if strings.HasPrefix(language, "php@") {
		return fmt.Sprintf(`depends_on "%s"`, language)
	}
	if strings.HasPrefix(language, "python@") {
		return fmt.Sprintf(`depends_on "%s"`, language)
	}
	if strings.HasPrefix(language, "node@") {
		return fmt.Sprintf(`depends_on "%s"`, language)
	}

	switch language {
	case "go":
		return `depends_on "go" => :build`
	case "rust":
		return `depends_on "rust" => :build`
	case "python":
		return `depends_on "python@3.11"`
	case "node":
		return `depends_on "node"`
	case "php":
		return `depends_on "php"`
	case "binary":
		return "" // No build dependency for prebuilt binaries
	default:
		return `depends_on "go" => :build` // Default to Go
	}
}
