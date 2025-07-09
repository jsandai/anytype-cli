# Anytype CLI

A command-line interface for interacting with [Anytype](https://github.com/anyproto/anytype-ts). Built for developers to enable [headless instances](https://github.com/anyproto/anytype-heart) or server-side exposure of the [Anytype API](https://github.com/anyproto/anytype-api).

## Quick Start

```bash
# Download middleware server
make download-server

# Build and install CLI
make install

# Start the daemon
anytype daemon

# In another terminal, start the server
anytype server start
```

## Installation

### Prerequisites

- Go 1.20 or later
- Git
- Make

### Download Middleware Server

The Anytype CLI requires the [Anytype middleware server](https://github.com/anyproto/anytype-heart) to function:

```bash
make download-server
```

This downloads the latest binary for your platform to the `dist/` directory.

### Build and Install

```bash
# Build only
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

<details>
<summary>Manual Setup (Advanced)</summary>

If you prefer manual setup or need to build the middleware from source:

#### Expected repository structure:

```
parent-directory/
├── anytype-heart/
└── anytype-cli/
```

#### Steps:

1. **In `anytype-heart` directory:**

```bash
make install-dev-cli
```

2. **In `anytype-cli` directory:**

```bash
go build -o dist/anytype
```

</details>

## Usage

```
anytype <command> <subcommand> [flags]

Commands:
  auth        Authenticate with Anytype
  daemon      Run the Anytype background daemon
  server      Manage the middleware server
  shell       Start the Anytype interactive shell
  space       Manage spaces
  version     Show version information

Examples:
  anytype daemon                    # Run daemon in foreground
  anytype server start              # Start the middleware server
  anytype auth login                # Login with mnemonic
  anytype space autoapprove         # Auto-approve space join requests

Use "anytype <command> --help" for more information about a command.
```
