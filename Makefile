.PHONY: download-server build install

download-server:
	@echo "Downloading anytype middleware server..."
	@./setup.sh

build:
	@echo "Building anytype CLI..."
	@go build -o dist/anytype
	@echo "Built successfully: dist/anytype"

install: build
	@echo "Installing anytype CLI..."
	@cp dist/anytype /usr/local/bin/anytype 2>/dev/null || sudo cp dist/anytype /usr/local/bin/anytype
	@echo "Installed to /usr/local/bin/anytype"

install-local: build
	@mkdir -p $$HOME/.local/bin
	@cp dist/anytype $$HOME/.local/bin/anytype
	@echo "Installed to $$HOME/.local/bin/anytype"
	@echo "Make sure $$HOME/.local/bin is in your PATH"
