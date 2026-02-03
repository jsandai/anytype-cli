# Self-Hosted Network Setup

Guide for connecting `anytype-cli` to a self-hosted Anytype network ([any-sync](https://github.com/anyproto/any-sync)).

## Prerequisites

- Self-hosted Anytype network running
- Network config YAML file with your node addresses
- Go 1.24+ and Make (for building from source)

## Quick Start

### 1. Install

**Option A: Install script**
```bash
/usr/bin/env bash -c "$(curl -fsSL https://raw.githubusercontent.com/anyproto/anytype-cli/HEAD/install.sh)"
```

**Option B: Build from source**
```bash
git clone https://github.com/anyproto/anytype-cli.git
cd anytype-cli
make build
# Binary: anytype
```

### 2. Prepare Network Config

Save your network configuration (from your any-sync deployment):

```yaml
# ~/.config/anytype/network.yml
networkId: YOUR_NETWORK_ID
nodes:
  - peerId: 12D3KooW...
    addresses:
      - your.server.com:33021
    types:
      - coordinator
  - peerId: 12D3KooW...
    addresses:
      - your.server.com:33022
    types:
      - consensus
  - peerId: 12D3KooW...
    addresses:
      - your.server.com:33020
    types:
      - tree
  - peerId: 12D3KooW...
    addresses:
      - your.server.com:33023
    types:
      - file
```

### 3. Create Account

```bash
anytype auth create my-bot --network-config ~/.config/anytype/network.yml
```

**⚠️ Save the account key!** It's your only authentication credential.

This persists the network config path to `~/.anytype/config.json` for future commands.

### 4. Start Server

```bash
anytype serve
```

The server must be running for most operations. Use `--quiet` for less output or `--verbose` for debugging.

### 5. Join a Space

Generate an invite link from the Anytype app (connected to your self-hosted network), then:

```bash
anytype space join "<invite-link>"
```

Supported invite formats:
- `https://<host>/<cid>#<key>` (web invite)
- `anytype://invite/?cid=<cid>&key=<key>` (app deep link)

### 6. Use the CLI

```bash
# List spaces
anytype space list

# Create API key for REST access
anytype auth apikey create "my-key"

# Use REST API
curl -H "Authorization: Bearer <api-key>" http://127.0.0.1:31012/v1/spaces
```

## Config Files

| File | Purpose |
|------|---------|
| `~/.anytype/config.json` | Account ID, network config path, cached network ID |
| `~/.config/anytype/network.yml` | Your self-hosted network nodes |
| System keyring | Account key, session tokens |

## Ports

| Port | Service |
|------|---------|
| 31010 | gRPC server |
| 31011 | gRPC Web proxy |
| 31012 | REST API |

## Troubleshooting

| Error | Solution |
|-------|----------|
| `network id mismatch` | Re-create account with `--network-config` pointing to correct YAML |
| `DeadlineExceeded` | Check network connectivity to your self-hosted nodes |
| `no ns peers configured` | Normal for self-hosted — naming service is cloud-only |
| `membership status` errors | Normal for self-hosted — membership is cloud-only |

## API Documentation

Full REST API reference: https://developers.anytype.io
