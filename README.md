# Anytype CLI

A command-line interface for interacting with [Anytype](https://github.com/anyproto/anytype-ts). Built for developers to enable [headless instances](https://github.com/anyproto/anytype-heart) or server-side exposure of the [Anytype API](https://github.com/anyproto/anytype-api).

## Quick Start

```bash
# Download the middleware server (required before first use)
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

The Anytype CLI requires the [Anytype middleware server](https://github.com/anyproto/anytype-heart) to function. You must download the server before first use:

```bash
make download-server
```

This downloads the appropriate server binary for your platform. You can also specify a different platform:

```bash
# Download for Linux AMD64
make download-server GOOS=linux GOARCH=amd64

# Download for macOS ARM64
make download-server GOOS=darwin GOARCH=arm64
```

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

## Contribution

Thank you for your desire to develop Anytype together!

❤️ This project and everyone involved in it is governed by the [Code of Conduct](https://github.com/anyproto/.github/blob/main/docs/CODE_OF_CONDUCT.md).

🧑‍💻 Check out our [contributing guide](https://github.com/anyproto/.github/blob/main/docs/CONTRIBUTING.md) to learn about asking questions, creating issues, or submitting pull requests.

🫢 For security findings, please email [security@anytype.io](mailto:security@anytype.io) and refer to our [security guide](https://github.com/anyproto/.github/blob/main/docs/SECURITY.md) for more information.

🤝 Follow us on [Github](https://github.com/anyproto) and join the [Contributors Community](https://github.com/orgs/anyproto/discussions).

---

Made by Any — a Swiss association 🇨🇭

Licensed under [MIT](./LICENSE.md).
