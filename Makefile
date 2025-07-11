.PHONY: all download-server build install uninstall install-local uninstall-local update lint lint-fix install-linter

all: build

GOLANGCI_LINT_VERSION := v2.2.1

VERSION ?= $(shell git describe --tags 2>/dev/null)
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d %H:%M:%S')
GIT_STATE ?= $(shell git diff --quiet 2>/dev/null && echo "clean" || echo "dirty")
LDFLAGS := -X 'github.com/anyproto/anytype-cli/core.Version=$(VERSION)' \
           -X 'github.com/anyproto/anytype-cli/core.Commit=$(COMMIT)' \
           -X 'github.com/anyproto/anytype-cli/core.BuildTime=$(BUILD_TIME)' \
           -X 'github.com/anyproto/anytype-cli/core.GitState=$(GIT_STATE)'

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
OUTPUT ?= dist/anytype

download-server:
	@echo "Downloading Anytype Middleware Server..."
	@./setup.sh

build:
	@if [ ! -f dist/anytype-grpc-server ]; then \
		echo "Server binary not found, downloading..."; \
		$(MAKE) download-server; \
	fi
	@echo "Building Anytype CLI..."
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -o $(OUTPUT)
	@echo "Built successfully: $(OUTPUT)"

install: build
	@echo "Installing Anytype CLI..."
	@cp dist/anytype /usr/local/bin/anytype 2>/dev/null || sudo cp dist/anytype /usr/local/bin/anytype
	@cp dist/anytype-grpc-server /usr/local/bin/anytype-grpc-server 2>/dev/null || sudo cp dist/anytype-grpc-server /usr/local/bin/anytype-grpc-server
	@echo "Installed to /usr/local/bin/"

uninstall:
	@echo "Uninstalling Anytype CLI..."
	@rm -f /usr/local/bin/anytype 2>/dev/null || sudo rm -f /usr/local/bin/anytype
	@rm -f /usr/local/bin/anytype-grpc-server 2>/dev/null || sudo rm -f /usr/local/bin/anytype-grpc-server
	@echo "Uninstalled from /usr/local/bin/"

install-local: build
	@mkdir -p $$HOME/.local/bin
	@cp dist/anytype $$HOME/.local/bin/anytype
	@cp dist/anytype-grpc-server $$HOME/.local/bin/anytype-grpc-server
	@echo "Installed to $$HOME/.local/bin/"
	@echo "Make sure $$HOME/.local/bin is in your PATH"

uninstall-local:
	@echo "Uninstalling Anytype CLI from local..."
	@rm -f $$HOME/.local/bin/anytype
	@rm -f $$HOME/.local/bin/anytype-grpc-server
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