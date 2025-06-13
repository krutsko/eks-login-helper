# EKS Login Helper Makefile

BINARY_NAME=eks-login
VERSION=1.0.0
BUILD_DIR=build
INSTALL_PATH=/usr/local/bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

.PHONY: all build clean test deps install uninstall help

all: clean deps build

# Build the binary
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
build-all: clean deps
	@echo "🔨 Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# Linux
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	
	# macOS
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	
	# Windows
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	
	@echo "✅ Multi-platform build complete"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	$(GOMOD) tidy
	$(GOMOD) download

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Run tests
test:
	@echo "🧪 Running tests..."
	$(GOTEST) -v ./...

# Install the binary system-wide
install: build
	@echo "📦 Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/
	@sudo chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✅ Installation complete. You can now run: $(BINARY_NAME)"

# Uninstall the binary
uninstall:
	@echo "🗑️  Uninstalling $(BINARY_NAME)..."
	@sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✅ Uninstallation complete"

# Install locally (in current directory)
install-local: build
	@echo "📦 Installing $(BINARY_NAME) locally..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) ./
	@chmod +x ./$(BINARY_NAME)
	@echo "✅ Local installation complete. Run with: ./$(BINARY_NAME)"

# Run the application
run: build
	@echo "🚀 Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

# Run with example parameters
run-example: build
	@echo "🚀 Running $(BINARY_NAME) with example parameters..."
	@$(BUILD_DIR)/$(BINARY_NAME) --help

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@$(GOCMD) fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "🔍 Linting code..."
	@golangci-lint run

# Check for security issues (requires gosec)
security:
	@echo "🔒 Running security checks..."
	@gosec ./...

# Create release archives
release: build-all
	@echo "📦 Creating release archives..."
	@mkdir -p $(BUILD_DIR)/releases
	
	# Create tar.gz for Unix systems
	@cd $(BUILD_DIR) && tar -czf releases/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	@cd $(BUILD_DIR) && tar -czf releases/$(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64
	@cd $(BUILD_DIR) && tar -czf releases/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	@cd $(BUILD_DIR) && tar -czf releases/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	
	# Create zip for Windows
	@cd $(BUILD_DIR) && zip releases/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	
	@echo "✅ Release archives created in $(BUILD_DIR)/releases/"

# Development workflow
dev: deps fmt test build

# Show help
help:
	@echo "EKS Login Helper - Makefile Help"
	@echo ""
	@echo "Available targets:"
	@echo "  build         Build the binary for current platform"
	@echo "  build-all     Build for all supported platforms"
	@echo "  deps          Install Go dependencies"
	@echo "  clean         Clean build artifacts"
	@echo "  test          Run tests"
	@echo "  install       Install binary system-wide (requires sudo)"
	@echo "  install-local Install binary in current directory"
	@echo "  uninstall     Remove installed binary"
	@echo "  run           Build and run the application"
	@echo "  run-example   Build and run with --help"
	@echo "  fmt           Format Go code"
	@echo "  lint          Run golangci-lint (if installed)"
	@echo "  security      Run gosec security checks (if installed)"
	@echo "  release       Create release archives for all platforms"
	@echo "  dev           Run development workflow (deps, fmt, test, build)"
	@echo "  help          Show this help message"
	@echo ""
	@echo "Quick start:"
	@echo "  make build    # Build the binary"
	@echo "  make install  # Install system-wide"
	@echo "  eks-login     # Run the tool"