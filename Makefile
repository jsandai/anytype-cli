.PHONY: download-server build install

download-server:
	@echo "Downloading Anytype Middleware Server..."
	@./setup.sh

build:
	@echo "Building Anytype CLI..."
	@go build -o dist/anytype
	@echo "Built successfully: dist/anytype"

install: build
	@echo "Installing Anytype CLI..."
	@cp dist/anytype /usr/local/bin/anytype 2>/dev/null || sudo cp dist/anytype /usr/local/bin/anytype
	@echo "Installed to /usr/local/bin/anytype"

install-local: build
	@mkdir -p $$HOME/.local/bin
	@cp dist/anytype $$HOME/.local/bin/anytype
	@echo "Installed to $$HOME/.local/bin/anytype"
	@echo "Make sure $$HOME/.local/bin is in your PATH"
