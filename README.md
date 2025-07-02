# Anytype CLI

## Setup

### Download Pre-built Binaries

The easiest way to get started is to download the pre-built binaries:

```bash
# Interactive mode - select your platform with arrow keys
./setup.sh

# Direct download - specify platform and architecture
./setup.sh linux amd64
./setup.sh darwin arm64
./setup.sh windows amd64
```

Available platforms:
- Linux AMD64 / ARM64
- macOS Apple Silicon (ARM64) / Intel (AMD64)
- Windows AMD64

The setup script will download and extract the binaries to the `dist/` directory.

### Build from Source

If you prefer to build from source:

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
