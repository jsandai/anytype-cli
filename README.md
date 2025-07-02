# Anytype CLI

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

## Setup

### Download Middleware Server

```bash
make download-server
```

This downloads the [Anytype middleware server](https://github.com/anyproto/anytype-heart) (grpc-server) to the `dist/` directory.

### Build from Source

```bash
make build
```

Builds the Anytype CLI binary to `dist/anytype`.

### Install

```bash
# System-wide installation (may require sudo)
make install

# User-local installation (no sudo required)
make install-local
```

The `install` target will install to `/usr/local/bin/anytype`. The `install-local` target installs to `~/.local/bin/anytype`.

### Manual Setup

If you prefer manual setup or need to build from source:

Expected repository structure:

```
parent-directory/
├── anytype-heart/
└── anytype-cli/
```

1. **In `anytype-heart` directory:**

```bash
make install-dev-cli
```

2. **In `anytype-cli` directory:**

```bash
make build
```

## Usage

### Start the daemon

- To run in the foreground:

```bash
./dist/anytype daemon
```

- To run in the background:

```bash
./dist/anytype daemon &
```

### Auto-approve members in a space

```bash
./dist/anytype space autoapprove --role "Editor" --space "<SpaceId>"
```
