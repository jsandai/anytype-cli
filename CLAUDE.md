# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the Anytype CLI, a Go-based command-line interface for interacting with Anytype. It includes an embedded gRPC server (from anytype-heart) making it a complete, self-contained binary that provides both server and CLI functionality. The CLI embeds the anytype-heart middleware server directly, eliminating the need for separate server installation or daemon processes.

## Build Commands

```bash
# Build the CLI (includes embedded server, downloads tantivy library automatically)
make build

# Install system-wide
make install

# Install user-local (~/.local/bin)
make install-local

# Clean tantivy libraries
make clean-tantivy

# Run linting
make lint
```

### Build Requirements

- **CGO**: The build requires CGO_ENABLED=1 due to tantivy (full-text search library) dependencies
- **Tantivy Library**: Automatically downloaded for your platform during `make build`
- **C++ Compiler**: Required for linking tantivy library (clang on macOS, gcc on Linux)

## Development Workflow

1. **Initial Setup**:
   ```bash
   make build  # Build the CLI with embedded server
   ```

2. **Running the Application**:
   ```bash
   # Run server interactively (for development)
   ./dist/anytype serve
   
   # Or install as system service
   ./dist/anytype service install
   ./dist/anytype service start
   ```

3. **Code Formatting and Linting**:
   ```bash
   go fmt ./...
   go vet ./...
   make lint  # Uses golangci-lint
   ```

## Architecture Overview

### Command Structure (`/cmd/`)
- Uses Cobra framework for CLI commands
- Each command group has its own directory:
  - `auth/`: Authentication commands (login, logout, status)
    - `apikey/`: API key management (create, list, revoke)
  - `serve/`: Run the embedded Anytype server in foreground
  - `service/`: System service management (install, uninstall, start, stop, restart, status)
  - `space/`: Space management operations
  - `shell/`: Interactive shell mode
  - `update/`: Self-update functionality
  - `version/`: Version information
- `root.go` registers all commands

### Core Logic (`/core/`)
- `client.go`: gRPC client singleton for server communication
- `auth.go`: Authentication logic with keyring integration
- `space.go`: Space management operations
- `stream.go`: Event streaming functionality with EventReceiver
- `keyring.go`: Secure credential storage (tokens and API keys)
- `apikey.go`: API key generation and management
- `config/constants.go`: Centralized configuration constants
- `serviceprogram/`: Service implementation using kardianos/service
- `grpcserver/`: Embedded gRPC server implementation

## Key Dependencies

- `github.com/anyproto/anytype-heart v0.42.0`: The middleware server (embedded)
- `github.com/spf13/cobra v1.8.1`: CLI framework
- `google.golang.org/grpc v1.73.0`: gRPC communication
- `github.com/zalando/go-keyring`: Secure credential storage
- `github.com/cheggaaa/mb/v3 v3.2.0`: Message batching queue for event handling
- `github.com/kardianos/service v1.2.4`: Cross-platform system service management

## Important Notes

1. **Service Architecture**: The CLI includes an embedded gRPC server that runs as a system service or interactively
2. **Cross-Platform Service**: Works on Windows (Service), macOS (launchd), Linux (systemd/upstart/sysv)
3. **Keyring Integration**: Authentication tokens are stored securely in the system keyring
4. **gRPC Communication**: All server interaction happens via gRPC on localhost:31007
5. **Event Streaming**: Uses server-sent events for real-time updates
6. **Version Management**: Version info is injected at build time via ldflags
7. **Self-Updating**: The CLI can update itself using the `anytype update` command
8. **API Keys**: Support for generating API keys for programmatic access

## Common Development Tasks

### Adding a New Command
1. Create a new directory under `/cmd/` for your command group
2. Create a file named after the command (e.g., `config.go` for config command) with a `NewXxxCmd()` function that returns `*cobra.Command`
3. Create subdirectories for each subcommand with their own files
4. Import subcommands with aliases matching the directory name
5. Register the command in `/cmd/root.go` using `NewXxxCmd()`
6. Implement core logic in `/core/` if needed

Directory structure follows the subcommand pattern:
```
cmd/
├── config/
│   ├── config.go      # Main command file (not cmd.go)
│   ├── get/
│   │   └── get.go     # Subcommand with NewGetCmd()
│   ├── set/
│   │   └── set.go     # Subcommand with NewSetCmd()
│   └── reset/
│       └── reset.go   # Subcommand with NewResetCmd()
```

Example main command file:
```go
// cmd/config/config.go
package config

import (
    "github.com/spf13/cobra"
    
    configGetCmd "github.com/anyproto/anytype-cli/cmd/config/get"
    configSetCmd "github.com/anyproto/anytype-cli/cmd/config/set"
    configResetCmd "github.com/anyproto/anytype-cli/cmd/config/reset"
)

func NewConfigCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "config <command>",
        Short: "Manage configuration",
    }
    
    cmd.AddCommand(configGetCmd.NewGetCmd())
    cmd.AddCommand(configSetCmd.NewSetCmd())
    cmd.AddCommand(configResetCmd.NewResetCmd())
    
    return cmd
}
```

### Working with the Service
- Service is managed via the `anytype serve` command
- Service program implementation is in `core/serviceprogram/`
- Supports both interactive mode and system service installation

### Error Handling
- Client connection errors are handled in `core/client.go`
- Server startup errors are managed in `core/serviceprogram/serviceprogram.go`
- Use standard Go error wrapping with context

### API Key Management
- API keys are created and managed by the server via gRPC APIs
- The CLI provides commands to create, list, and revoke API keys
- Keys are generated server-side and can be used for programmatic access
