# tobrew

**tobrew** - Automated Homebrew tap release tool for CLI projects

Automate your Homebrew tap releases with a single command. No more manual version management, SHA256 calculations, or tap repository updates.

## Features

- ✅ **Automatic version management** - `tobrew.lock` tracks your current version
- 🚀 **One-command releases** - `tobrew release` does everything
- 🌍 **Multi-language support** - Go, Rust, Python, Node.js, PHP, and prebuilt binaries
- 📝 **Multiple config formats** - YAML, JSON, or TOML
- 🔐 **Automatic SHA256** calculation from GitHub releases
- 🍺 **Homebrew formula** generation
- 📦 **Auto-update** homebrew-tap repository
- 🎯 **Simple workflow** - Only 2 commands needed

## Installation

### Method 1: Homebrew (Recommended)

```bash
brew install yejune/tap/tobrew
```

### Method 2: Using go install

```bash
go install github.com/yejune/tobrew@latest
```

### Method 3: From source

```bash
git clone https://github.com/yejune/tobrew.git
cd tobrew
go build
./tobrew install
```

### Updating

```bash
brew upgrade tobrew
# or
tobrew self-update
```

## Quick Start

### 1. Initialize configuration (once)

In your project directory:

```bash
# Go project (default)
tobrew init

# Rust project
tobrew init --language rust

# Python project with specific version
tobrew init --language python@3.11

# PHP project with specific version
tobrew init --language php@8.4

# Node.js project
tobrew init --language node

# Prebuilt binary
tobrew init --language binary
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

### `tobrew sync`

Sync lock file with remote git tags.

```bash
tobrew sync
```

Useful when:
- Lock file is out of sync with actual releases
- Working on a different machine
- Recovering from a failed release

## Version Management

tobrew uses a `tobrew.lock` file to track your project version:

```yaml
version: v1.2.3
last_release: 2025-11-25T15:30:00+09:00
sha256: abc123...
fingerprint: XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
```

- **First release**: Starts at `v0.0.1`
- **Automatic bumping**: No need to specify version numbers
- **Semantic versioning**: Follows semver (MAJOR.MINOR.PATCH)
- **Git tracked**: Commit `tobrew.lock` to your repository
- **Auto-sync**: Automatically syncs with remote when fingerprint differs (different machine) or tag conflict occurs

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

### Go Project

```yaml
name: mycli
language: go
description: "My awesome Go CLI tool"
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

### Rust Project

```yaml
name: rustcli
language: rust
description: "Rust CLI application"
homepage: https://github.com/user/rustcli
license: MIT

github:
  user: user
  repo: rustcli
  tap_repo: homebrew-tap

build:
  command: cargo build --release

formula:
  install: |
    system "cargo", "install", *std_cargo_args

  test: |
    assert_match "rustcli", shell_output("#{bin}/rustcli --version")
```

### Python Project

```yaml
name: pycli
language: python@3.11  # Specify Python version
description: "Python CLI tool"
homepage: https://github.com/user/pycli
license: MIT

github:
  user: user
  repo: pycli
  tap_repo: homebrew-tap

build:
  command: python -m build

formula:
  install: |
    virtualenv_install_with_resources

  test: |
    assert_match "pycli", shell_output("#{bin}/pycli --version")
```

### Node.js Project

```yaml
name: nodecli
language: node  # or node@20 for specific version
description: "Node.js CLI application"
homepage: https://github.com/user/nodecli
license: MIT

github:
  user: user
  repo: nodecli
  tap_repo: homebrew-tap

build:
  command: npm run build

formula:
  install: |
    system "npm", "install", *Language::Node.std_npm_install_args(libexec)
    bin.install_symlink Dir["#{libexec}/bin/*"]

  test: |
    assert_match "nodecli", shell_output("#{bin}/nodecli --version")
```

### PHP Project

```yaml
name: phpcli
language: php@8.4  # Specify PHP version
description: "PHP CLI tool"
homepage: https://github.com/user/phpcli
license: MIT

github:
  user: user
  repo: phpcli
  tap_repo: homebrew-tap

build:
  command: composer install --no-dev --optimize-autoloader

formula:
  install: |
    libexec.install Dir["*"]
    bin.install_symlink libexec/"phpcli"

  test: |
    assert_match "phpcli", shell_output("#{bin}/phpcli --version")
```

### Prebuilt Binary

```yaml
name: binarycli
language: binary
description: "Precompiled binary (built by CI/CD)"
homepage: https://github.com/user/binarycli
license: MIT

github:
  user: user
  repo: binarycli
  tap_repo: homebrew-tap

build:
  command: "# Binary built by GitHub Actions"

formula:
  install: |
    bin.install "binarycli"

  test: |
    assert_match "binarycli", shell_output("#{bin}/binarycli --version")
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

tobrew is designed for simplicity:

- ✅ **Simpler** - Only 2 commands (`init` and `release`)
- ✅ **Focused** - Just Homebrew tap management
- ✅ **Automatic** - Version management without manual input
- ✅ **Lightweight** - Minimal configuration required
- ✅ **Multi-language** - Supports Go, Rust, Python, Node.js, PHP, and more

Perfect for CLI tools that just need simple Homebrew distribution.

## License

MIT License - see LICENSE file for details

## Contributing

Contributions welcome! Please open an issue or PR.

---

Made with ❤️ for easier Homebrew releases
