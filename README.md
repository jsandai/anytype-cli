# Anytype CLI

A command-line interface for interacting with [Anytype](https://github.com/anyproto/anytype-ts). This CLI includes an embedded gRPC server from [anytype-heart](https://github.com/anyproto/anytype-heart), making it a complete, self-contained solution for developers to work with Anytype.

## Quick Start

```bash
# Build the CLI (includes embedded server)
make build

# Install the CLI
make install

# Run the Anytype server
anytype serve
```

## Installation

### Quick Install (Recommended)

Install the latest release with a single command:

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/anyproto/anytype-cli/HEAD/install.sh)"
```

### Build from Source

#### Prerequisites

- Go 1.20 or later
- Git
- Make
- C++ compiler (for tantivy library)

#### Build and Install

```bash
# Build only (automatically downloads tantivy library)
make build

# Build and install system-wide (may require sudo)
make install

# Build and install to ~/.local/bin (no sudo required)
make install-local
```

### Uninstall

```bash
# Remove system-wide installation
make uninstall

# Remove user-local installation
make uninstall-local
```

## Usage

```
anytype <command> <subcommand> [flags]

Commands:
  auth        Authenticate with Anytype
  serve       Run the Anytype server
  service     Manage Anytype as a system service
  shell       Start the Anytype interactive shell
  space       Manage spaces
  update      Update anytype CLI to the latest version
  version     Show version information

Examples:
  anytype serve                     # Run server in foreground
  anytype service install           # Install as system service
  anytype service start             # Start the service
  anytype auth login                # Login with mnemonic
  anytype space list                # List available spaces

Use "anytype <command> --help" for more information about a command.
```

### Running the Server

The CLI includes an embedded gRPC server that can be run in two ways:

#### 1. Interactive Mode (for development)
```bash
anytype serve
```
This runs the server in the foreground with logs output to stdout, similar to `ollama serve`.

#### 2. System Service (for production)
```bash
# Install as system service
anytype service install

# Start the service
anytype service start

# Check service status
anytype service status

# Stop the service
anytype service stop

# Uninstall the service
anytype service uninstall
```

The service management works across platforms:
- **macOS**: Uses launchd
- **Linux**: Uses systemd/upstart/sysv
- **Windows**: Uses Windows Service

### Authentication

After starting the server, authenticate with your Anytype account:

```bash
# Login with mnemonic
anytype auth login

# Check authentication status
anytype auth status

# Logout
anytype auth logout
```

### API Keys

Generate API keys for programmatic access:

```bash
# Create a new API key
anytype auth apikey create --name "my-app"

# List API keys
anytype auth apikey list

# Revoke an API key
anytype auth apikey revoke <key-id>
```

## Development

### Project Structure

```
anytype-cli/
├── cmd/              # CLI commands
│   ├── auth/         # Authentication commands
│   ├── serve/        # Server command
│   ├── service/      # Service management
│   ├── space/        # Space management
│   └── ...
├── core/             # Core business logic
│   ├── grpcserver/   # Embedded gRPC server
│   ├── serviceprogram/ # Service implementation
│   └── ...
└── dist/             # Build output
```

### Building from Source

```bash
# Clone the repository
git clone https://github.com/anyproto/anytype-cli.git
cd anytype-cli

# Build (CGO is automatically enabled for tantivy)
make build

# Run tests
go test ./...

# Run linting
make lint
```

## Contribution

Thank you for your desire to develop Anytype together!

❤️ This project and everyone involved in it is governed by the [Code of Conduct](https://github.com/anyproto/.github/blob/main/docs/CODE_OF_CONDUCT.md).

🧑‍💻 Check out our [contributing guide](https://github.com/anyproto/.github/blob/main/docs/CONTRIBUTING.md) to learn about asking questions, creating issues, or submitting pull requests.

🫢 For security findings, please email [security@anytype.io](mailto:security@anytype.io) and refer to our [security guide](https://github.com/anyproto/.github/blob/main/docs/SECURITY.md) for more information.

🤝 Follow us on [Github](https://github.com/anyproto) and join the [Contributors Community](https://github.com/orgs/anyproto/discussions).

---

Made by Any — a Swiss association 🇨🇭

Licensed under [MIT](./LICENSE.md).