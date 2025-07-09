# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the Anytype CLI, a Go-based command-line interface for interacting with Anytype. It uses a client-server architecture where the CLI communicates with a separate middleware server (anytype-heart) via gRPC.

## Build Commands

```bash
# Download the middleware server (required before first use)
make download-server

# Build the CLI
make build

# Install system-wide
make install

# Install user-local
make install-local

# Manual build
go build -o dist/anytype
```

## Development Workflow

1. **Initial Setup**:
   ```bash
   make download-server  # Downloads anytype-heart middleware
   make build           # Builds the CLI
   ```

2. **Running the Application**:
   ```bash
   # Start the daemon (required)
   ./dist/anytype daemon
   
   # In another terminal, start the server
   ./dist/anytype server start
   ```

3. **Code Formatting**:
   ```bash
   go fmt ./...
   go vet ./...
   ```

## Architecture Overview

### Command Structure (`/cmd/`)
- Uses Cobra framework for CLI commands
- Each command group has its own directory (auth/, daemon/, server/, space/, token/, shell/)
- `root.go` registers all commands

### Core Logic (`/internal/`)
- `client.go`: gRPC client singleton for server communication
- `auth.go`: Authentication logic with keyring integration
- `space.go`: Space management operations
- `stream.go`: Event streaming functionality with EventReceiver
- `token.go`: Token management

### Daemon (`/daemon/`)
- `daemon.go`: Main daemon process that manages server lifecycle
- `taskmanager.go`: Schedules and manages background tasks
- `daemon_client.go`: Client for daemon communication

### Tasks (`/tasks/`)
- Background tasks executed by the daemon
- `autoapprove.go`: Auto-approves space join requests
- `server.go`: Manages server startup/shutdown

## Key Dependencies

- `github.com/anyproto/anytype-heart v0.39.5`: The middleware server
- `github.com/spf13/cobra v1.8.1`: CLI framework
- `google.golang.org/grpc v1.70.0`: gRPC communication
- `github.com/99designs/go-keyring v0.2.6`: Secure credential storage

## Important Notes

1. **Two-Process Architecture**: The CLI requires both a daemon process and the middleware server to be running
2. **Keyring Integration**: Credentials are stored securely in the system keyring
3. **gRPC Communication**: All server interaction happens via gRPC on localhost:31007
4. **Event Streaming**: Uses server-sent events for real-time updates in auto-approval
5. **No Tests**: The project currently has no test files

## Common Development Tasks

### Adding a New Command
1. Create a new directory under `/cmd/` for your command group
2. Create a `cmd.go` file with the Cobra command definition
3. Register the command in `/cmd/root.go`
4. Implement core logic in `/internal/` if needed

### Working with the Daemon
- Daemon tasks go in `/tasks/`
- Register new tasks in `daemon/taskmanager.go`
- Use the `Task` interface for new task implementations

### Error Handling
- Client connection errors are handled in `internal/client.go`
- Server startup errors are managed in `daemon/daemon.go:connectToServer`
- Use standard Go error wrapping with context
