BINARY_NAME=go-check-http-methods
VERSION=1.0.0
BUILD_DIR=build
LDFLAGS=-ldflags "-X main.appVersion=${VERSION}"

.PHONY: all build clean test deps linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64 windows-arm64

all: clean linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64 windows-arm64

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@mkdir -p $(BUILD_DIR)

deps:
	@echo "Downloading dependencies..."
	@go mod download

test:
	@echo "Running tests..."
	@go test -v ./...

# Linux builds
linux-amd64:
	@echo "Building for Linux (amd64)..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64

linux-arm64:
	@echo "Building for Linux (arm64)..."
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64

# macOS builds
darwin-amd64:
	@echo "Building for macOS (amd64)..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64

darwin-arm64:
	@echo "Building for macOS (arm64)..."
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64

# Windows builds
windows-amd64:
	@echo "Building for Windows (amd64)..."
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe

windows-arm64:
	@echo "Building for Windows (arm64)..."
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe

# Package builds into zip files
package: all
	@echo "Packaging builds..."
	@cd $(BUILD_DIR) && \
	for file in $(BINARY_NAME)-* ; do \
		zip $$file.zip $$file ; \
	done

# Install locally (current platform)
install:
	@echo "Installing $(BINARY_NAME)..."
	@go install $(LDFLAGS)

# Run locally
run:
	@go run . $(ARGS)
