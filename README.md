# tobrew

**tobrew** - Automated Homebrew tap release tool for Go projects

Automate your Homebrew tap releases with a single command. No more manual version management, SHA256 calculations, or tap repository updates.

## Features

- ✅ **Automatic version management** - `tobrew.lock` tracks your current version
- 🚀 **One-command releases** - `tobrew release` does everything
- 📝 **Multiple config formats** - YAML, JSON, or TOML
- 🔐 **Automatic SHA256** calculation from GitHub releases
- 🍺 **Homebrew formula** generation
- 📦 **Auto-update** homebrew-tap repository
- 🎯 **Simple workflow** - Only 2 commands needed

## Installation

### Method 1: Install script (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/yejune/tobrew/main/install.sh | bash
```

Or download and run:

```bash
git clone https://github.com/yejune/tobrew.git
cd tobrew
./install.sh
```

### Method 2: Using go install

```bash
go install github.com/yejune/tobrew@latest
```

### Method 3: From source (for development)

```bash
git clone https://github.com/yejune/tobrew.git
cd tobrew
go build
./tobrew install
```

### Updating

```bash
tobrew self-update
```

## Quick Start

### 1. Initialize configuration (once)

In your Go project directory:

```bash
tobrew init
```

This creates a `tobrew.yaml` file. You can also use JSON or TOML:

```bash
tobrew init --format json
tobrew init --format toml
```

### 2. Edit configuration

Edit `tobrew.yaml` and update the placeholders:

```yaml
name: myapp
description: "My awesome CLI tool"
homepage: https://github.com/username/myapp
license: MIT

github:
  user: username           # Your GitHub username
  repo: myapp             # Your repo name
  tap_repo: homebrew-tap  # Your tap repo name (must start with "homebrew-")

build:
  command: go build -o build/{{.Name}} .

formula:
  install: |
    system "go", "build", "."
    bin.install "myapp"

  test: |
    assert_match "myapp", shell_output("#{bin}/myapp --version")

  caveats: |
    myapp has been installed!
    Run 'myapp --help' to get started.
```

### 3. Create GitHub tap repository

Create a new repository on GitHub named `homebrew-tap` (must start with `homebrew-`).

### 4. Release!

```bash
# First release (creates v0.0.1)
tobrew release

# Patch release (v0.0.1 → v0.0.2)
tobrew release

# Minor release (v0.0.2 → v0.1.0)
tobrew release --minor

# Major release (v0.1.0 → v1.0.0)
tobrew release --major
```

This will:
1. **Load** current version from `tobrew.lock`
2. **Bump** version (patch/minor/major)
3. **Build** your project
4. **Create** and push git tag
5. **Download** release tarball and calculate SHA256
6. **Generate** Homebrew formula
7. **Update** your homebrew-tap repository
8. **Save** new version to `tobrew.lock`

Your users can now install with:

```bash
brew install username/tap/myapp
```

## Commands

### `tobrew init`

Initialize a new configuration file.

```bash
tobrew init                    # Creates tobrew.yaml
tobrew init --format json      # Creates tobrew.json
tobrew init --format toml      # Creates tobrew.toml
tobrew init -o custom.yaml     # Custom output path
```

### `tobrew release`

Create a release with automatic version bumping.

```bash
tobrew release              # Patch: v1.0.0 → v1.0.1 (default)
tobrew release --patch      # Patch: v1.0.0 → v1.0.1 (explicit)
tobrew release --minor      # Minor: v1.0.1 → v1.1.0
tobrew release --major      # Major: v1.1.0 → v2.0.0
```

## Version Management

tobrew uses a `tobrew.lock` file to track your project version:

```yaml
version: v1.2.3
last_release: 2025-11-25T15:30:00+09:00
sha256: abc123...
```

- **First release**: Starts at `v0.0.1`
- **Automatic bumping**: No need to specify version numbers
- **Semantic versioning**: Follows semver (MAJOR.MINOR.PATCH)
- **Git tracked**: Commit `tobrew.lock` to your repository

## Configuration Reference

### Full `tobrew.yaml` example

```yaml
name: docker-bootapp
description: "Docker Compose multi-project manager"
homepage: https://github.com/yejune/docker-bootapp
license: MIT

github:
  user: yejune
  repo: docker-bootapp
  tap_repo: homebrew-tap

build:
  # Go template - {{.Name}} replaced with project name
  command: go build -o build/{{.Name}} .

formula:
  # Ruby code for Homebrew formula install section
  install: |
    system "go", "build", "."
    bin.install "docker-bootapp"

    # Install as Docker CLI plugin
    docker_plugins = "#{ENV["HOME"]}/.docker/cli-plugins"
    mkdir_p docker_plugins
    cp bin/"docker-bootapp", "#{docker_plugins}/docker-bootapp"
    chmod 0755, "#{docker_plugins}/docker-bootapp"

  # Ruby code for formula test section
  test: |
    assert_match "bootapp", shell_output("#{bin}/docker-bootapp help")

  # Message shown after installation
  caveats: |
    docker-bootapp has been installed!

    You can use it in two ways:
      docker bootapp [command]  # As Docker CLI plugin
      bootapp [command]         # As standalone binary
```

## How It Works

1. **Load Version**: Read current version from `tobrew.lock` (or start at v0.0.0)
2. **Bump**: Increment version according to flags (default: patch)
3. **Build**: Run configured build command
4. **Tag**: Create and push git tag to GitHub
5. **Download**: Fetch the release tarball from GitHub
6. **Hash**: Calculate SHA256 checksum
7. **Generate**: Create Homebrew formula from template
8. **Push**: Update your homebrew-tap repository
9. **Save**: Write new version to `tobrew.lock`

## Examples

### Simple Go CLI

```yaml
name: mycli
description: "My awesome CLI tool"
homepage: https://github.com/user/mycli
license: MIT

github:
  user: user
  repo: mycli
  tap_repo: homebrew-tap

build:
  command: go build -o build/{{.Name}} .

formula:
  install: |
    system "go", "build", "."
    bin.install "mycli"

  test: |
    assert_match "mycli version", shell_output("#{bin}/mycli --version")
```

### Multi-binary Project

```yaml
formula:
  install: |
    system "go", "build", "-o", "bin/server", "./cmd/server"
    system "go", "build", "-o", "bin/client", "./cmd/client"
    bin.install "bin/server"
    bin.install "bin/client"
```

## Typical Workflow

```bash
# Initial setup (once)
cd myproject
tobrew init
# Edit tobrew.yaml with your info
git add tobrew.yaml tobrew.lock
git commit -m "Add tobrew config"

# Development cycle
# ... make changes to your code ...
git add -A
git commit -m "Add new feature"
git push

# Release!
tobrew release              # Patch release
# or
tobrew release --minor      # Minor release
# or
tobrew release --major      # Major release
```

## Troubleshooting

### "failed to download tarball"

- Make sure the git tag exists on GitHub
- Wait a few seconds after pushing the tag
- Check that your GitHub repository is public or you have access

### "tap update failed"

- Check that your `homebrew-tap` repository exists
- Ensure you have push access to the tap repository
- Verify the repository name starts with `homebrew-`

### "invalid version format"

- Check `tobrew.lock` has valid version (e.g., `v1.2.3`)
- Delete `tobrew.lock` to start fresh from `v0.0.1`

## Why tobrew?

Other tools like [goreleaser](https://goreleaser.com/) are comprehensive but complex. tobrew is:

- ✅ **Simpler** - Only 2 commands (`init` and `release`)
- ✅ **Focused** - Just Homebrew tap management
- ✅ **Automatic** - Version management without manual input
- ✅ **Lightweight** - Minimal configuration required

Perfect for Go CLI tools that just need simple Homebrew distribution.

## License

MIT License - see LICENSE file for details

## Contributing

Contributions welcome! Please open an issue or PR.

## Related Projects

- [docker-bootapp](https://github.com/yejune/docker-bootapp) - Example project using tobrew
- [goreleaser](https://goreleaser.com/) - Comprehensive release automation

---

Made with ❤️ for easier Homebrew releases
