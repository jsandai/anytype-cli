.PHONY: all build install uninstall install-local uninstall-local update lint lint-fix install-linter download-tantivy clean-tantivy

all: download-tantivy build

GOLANGCI_LINT_VERSION := v2.2.1

VERSION ?= $(shell git describe --tags 2>/dev/null)
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d %H:%M:%S')
GIT_STATE ?= $(shell git diff --quiet 2>/dev/null && echo "clean" || echo "dirty")
LDFLAGS := -s -w \
           -X 'github.com/anyproto/anytype-cli/core.Version=$(VERSION)' \
           -X 'github.com/anyproto/anytype-cli/core.Commit=$(COMMIT)' \
           -X 'github.com/anyproto/anytype-cli/core.BuildTime=$(BUILD_TIME)' \
           -X 'github.com/anyproto/anytype-cli/core.GitState=$(GIT_STATE)'

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
OUTPUT ?= dist/anytype

# Tantivy library configuration
TANTIVY_VERSION := v1.0.4
TANTIVY_LIB_PATH := deps/libs
HEART_MODULE_PATH := $(shell go list -m -f '{{.Dir}}' github.com/anyproto/anytype-heart)
HEART_TANTIVY_PATH := $(HEART_MODULE_PATH)/deps/libs
CGO_LDFLAGS := -L$(TANTIVY_LIB_PATH)

build: download-tantivy
	@echo "Building Anytype CLI with embedded server..."
	@CGO_ENABLED=1 CGO_LDFLAGS="$(CGO_LDFLAGS)" GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -o $(OUTPUT)
	@echo "Built successfully: $(OUTPUT)"

install: build
	@echo "Installing Anytype CLI..."
	@cp dist/anytype /usr/local/bin/anytype 2>/dev/null || sudo cp dist/anytype /usr/local/bin/anytype
	@echo "Installed to /usr/local/bin/"
	@echo ""
	@echo "Usage:"
	@echo "  anytype serve              # Run server in foreground"
	@echo "  anytype service install    # Install as system service"

uninstall:
	@echo "Uninstalling Anytype CLI..."
	@rm -f /usr/local/bin/anytype 2>/dev/null || sudo rm -f /usr/local/bin/anytype
	@echo "Uninstalled from /usr/local/bin/"

install-local: build
	@mkdir -p $$HOME/.local/bin
	@cp dist/anytype $$HOME/.local/bin/anytype
	@echo "Installed to $$HOME/.local/bin/"
	@echo "Make sure $$HOME/.local/bin is in your PATH"
	@echo ""
	@echo "Usage:"
	@echo "  anytype serve              # Run server in foreground"
	@echo "  anytype service install    # Install as system service"

uninstall-local:
	@echo "Uninstalling Anytype CLI from local..."
	@rm -f $$HOME/.local/bin/anytype
	@echo "Uninstalled from $$HOME/.local/bin/"

install-linter:
	@echo "Installing golangci-lint..."
	@go install github.com/daixiang0/gci@latest
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	@echo "golangci-lint installed successfully"

lint:
	@golangci-lint run ./...

lint-fix:
	@golangci-lint run --fix ./...

download-tantivy:
	@if [ ! -f "$(TANTIVY_LIB_PATH)/libtantivy_go.a" ]; then \
		echo "Downloading tantivy library $(TANTIVY_VERSION) for $(GOOS)/$(GOARCH)..."; \
		mkdir -p $(TANTIVY_LIB_PATH); \
		if [ "$(GOOS)" = "darwin" ]; then \
			if [ "$(GOARCH)" = "amd64" ]; then \
				curl -L "https://github.com/anyproto/tantivy-go/releases/download/$(TANTIVY_VERSION)/darwin-amd64.tar.gz" | tar xz -C $(TANTIVY_LIB_PATH); \
			elif [ "$(GOARCH)" = "arm64" ]; then \
				curl -L "https://github.com/anyproto/tantivy-go/releases/download/$(TANTIVY_VERSION)/darwin-arm64.tar.gz" | tar xz -C $(TANTIVY_LIB_PATH); \
			else \
				echo "Unsupported architecture: $(GOARCH) for macOS"; \
				exit 1; \
			fi; \
		elif [ "$(GOOS)" = "linux" ]; then \
			if [ "$(GOARCH)" = "amd64" ]; then \
				curl -L "https://github.com/anyproto/tantivy-go/releases/download/$(TANTIVY_VERSION)/linux-amd64.tar.gz" | tar xz -C $(TANTIVY_LIB_PATH); \
			elif [ "$(GOARCH)" = "arm64" ]; then \
				curl -L "https://github.com/anyproto/tantivy-go/releases/download/$(TANTIVY_VERSION)/linux-arm64.tar.gz" | tar xz -C $(TANTIVY_LIB_PATH); \
			else \
				echo "Unsupported architecture: $(GOARCH) for Linux"; \
				exit 1; \
			fi; \
		elif [ "$(GOOS)" = "windows" ]; then \
			if [ "$(GOARCH)" = "amd64" ]; then \
				curl -L "https://github.com/anyproto/tantivy-go/releases/download/$(TANTIVY_VERSION)/windows-amd64.tar.gz" | tar xz -C $(TANTIVY_LIB_PATH); \
			else \
				echo "Unsupported architecture: $(GOARCH) for Windows"; \
				exit 1; \
			fi; \
		else \
			echo "Unsupported OS: $(GOOS)"; \
			exit 1; \
		fi; \
		echo "Tantivy library downloaded successfully"; \
	else \
		echo "Tantivy library already exists"; \
	fi

clean-tantivy:
	@echo "Cleaning tantivy libraries..."
	@rm -rf $(TANTIVY_LIB_PATH)
	@echo "Tantivy libraries cleaned"