# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=dwui
OUTPUT_DIR=dist

.PHONY: all build-all clean mac-amd64 mac-arm64 linux-amd64 linux-arm64 windows-amd64

all: build-all

# Build for all target platforms
build-all: mac-amd64 mac-arm64 linux-amd64 linux-arm64 windows-amd64
	@echo "All builds finished."

# Build for specific os/arch
mac-amd64:
	@echo "Building for Mac (Intel amd64)..."
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) -o "$(OUTPUT_DIR)/$(BINARY_NAME)-mac-amd64" main.go

mac-arm64:
	@echo "Building for Mac (Apple arm64)..."
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) -o "$(OUTPUT_DIR)/$(BINARY_NAME)-mac-arm64" main.go

linux-amd64:
	@echo "Building for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 $(GOBUILD) -o "$(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64" main.go

linux-arm64:
	@echo "Building for Linux (arm64)..."
	@GOOS=linux GOARCH=arm64 $(GOBUILD) -o "$(OUTPUT_DIR)/$(BINARY_NAME)-linux-arm64" main.go

windows-amd64:
	@echo "Building for Windows (amd64)..."
	@GOOS=windows GOARCH=amd64 $(GOBUILD) -o "$(OUTPUT_DIR)/$(BINARY_NAME)-windows-amd64.exe" main.go

# Remove the build artifacts from the output directory
clean:
	@echo "Cleaning up..."
	@rm -f \
		"$(OUTPUT_DIR)/$(BINARY_NAME)-mac-amd64" \
		"$(OUTPUT_DIR)/$(BINARY_NAME)-mac-arm64" \
		"$(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64" \
		"$(OUTPUT_DIR)/$(BINARY_NAME)-linux-arm64" \
		"$(OUTPUT_DIR)/$(BINARY_NAME)-windows-amd64.exe"