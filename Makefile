.PHONY: build install clean test cross-compile

BINARY_NAME=tobrew
BUILD_DIR=build
VERSION?=dev

# Build for current platform
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -ldflags="-X main.version=$(VERSION)" -o $(BINARY_NAME) .
	@echo "✓ Build complete: ./$(BINARY_NAME)"

# Install to /usr/local/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	@sudo cp $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "✓ Installed to /usr/local/bin/$(BINARY_NAME)"

# Uninstall from /usr/local/bin
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✓ Uninstalled"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR) $(BINARY_NAME)
	@echo "✓ Clean complete"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Cross-compile for all platforms (used by GitHub Actions)
cross-compile: clean
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)

	@echo "  - darwin/amd64"
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .

	@echo "  - darwin/arm64"
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .

	@echo "  - linux/amd64"
	@GOOS=linux GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .

	@echo "  - linux/arm64"
	@GOOS=linux GOARCH=arm64 go build -ldflags="-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .

	@echo "✓ Release builds complete in $(BUILD_DIR)/"

# Run tobrew on itself (for testing)
self-release:
	@./$(BINARY_NAME) release

# Show help
help:
	@echo "Available targets:"
	@echo "  make build           - Build for current platform"
	@echo "  make install         - Build and install to /usr/local/bin"
	@echo "  make uninstall       - Remove from /usr/local/bin"
	@echo "  make clean           - Remove build artifacts"
	@echo "  make test            - Run tests"
	@echo "  make cross-compile   - Cross-compile for all platforms (used by CI)"
	@echo "  make self-release    - Run tobrew release on itself"
	@echo "  make help            - Show this help"
