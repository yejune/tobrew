#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🚀 Installing tobrew..."
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    echo ""
    echo "Install Go from: https://golang.org/dl/"
    exit 1
fi

echo "✓ Go found: $(go version)"

# Determine OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Error: Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo "✓ Platform: $OS-$ARCH"
echo ""

# Build tobrew
echo "📦 Building tobrew..."
go build -o tobrew .

if [ ! -f tobrew ]; then
    echo -e "${RED}Error: Build failed${NC}"
    exit 1
fi

echo "✓ Build successful"
echo ""

# Install to /usr/local/bin
echo "📋 Installing to /usr/local/bin/tobrew..."
echo "   (sudo password may be required)"

if sudo rm -f /usr/local/bin/tobrew && sudo cp tobrew /usr/local/bin/tobrew && sudo chmod +x /usr/local/bin/tobrew; then
    echo "✓ tobrew installed to /usr/local/bin/tobrew"
else
    echo -e "${YELLOW}⚠️  Warning: Could not install to /usr/local/bin${NC}"
    echo "   You can still use ./tobrew from this directory"
    echo "   Or manually install: sudo cp tobrew /usr/local/bin/"
    exit 1
fi

# Verify installation
echo ""
echo "🔍 Verifying installation..."
if command -v tobrew &> /dev/null; then
    echo "✓ tobrew is available in PATH"
    tobrew --version
else
    echo -e "${YELLOW}⚠️  Warning: tobrew not found in PATH${NC}"
    echo "   You may need to restart your shell"
fi

echo ""
echo -e "${GREEN}✅ Installation complete!${NC}"
echo ""
echo "Quick start:"
echo "  1. cd your-go-project"
echo "  2. tobrew init"
echo "  3. Edit tobrew.yaml"
echo "  4. tobrew release"
echo ""
echo "Documentation: https://github.com/yejune/tobrew"
