# Anytype CLI - Self-Hosted Setup Guide

## Prerequisites

- Self-hosted Anytype network running (any-sync)
- Network config YAML file with your nodes
- Go 1.24+ and Make (for building)

## 1. Build the Fork

```bash
git clone https://github.com/jsandai/anytype-cli.git
cd anytype-cli
make build
# Binary: ./dist/anytype
```

## 2. Create Network Config

Save your network config as `~/network.yaml`:

```yaml
networkId: YOUR_NETWORK_ID
nodes:
  - peerId: 12D3KooW...
    addresses:
      - your.server.com:33021
      - quic://your.server.com:33021
    types:
      - coordinator
  - peerId: 12D3KooW...
    addresses:
      - your.server.com:33022
      - quic://your.server.com:33022
    types:
      - consensus
  - peerId: 12D3KooW...
    addresses:
      - your.server.com:33020
      - quic://your.server.com:33020
    types:
      - tree
  - peerId: 12D3KooW...
    addresses:
      - your.server.com:33023
      - quic://your.server.com:33023
    types:
      - file
```

## 3. Create Bot Account

```bash
./dist/anytype auth create my-bot --network-config ~/network.yaml
```

**Save the account key!** It's your only way to authenticate.

## 4. Run the Server

```bash
./dist/anytype serve
# Runs on localhost:31010-31012
# Keep this running in a terminal or use `service install`
```

## 5. Join a Space

Generate an invite link from Anytype desktop, then:

```bash
# With serve running in another terminal:
./dist/anytype space join "anytype://invite/?cid=...&key=..."
```

**Note:** The network ID is auto-loaded from your saved config.

## 6. Create API Key

```bash
./dist/anytype auth apikey create "my-api-key"
```

## 7. Use the API

```bash
# List spaces
curl -H "Authorization: Bearer YOUR_API_KEY" \
  http://127.0.0.1:31012/v1/spaces

# List objects in a space
curl -H "Authorization: Bearer YOUR_API_KEY" \
  "http://127.0.0.1:31012/v1/spaces/SPACE_ID/objects"

# Create an object
curl -X POST \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"name": "My Note", "type_key": "page", "body": "Hello world"}' \
  "http://127.0.0.1:31012/v1/spaces/SPACE_ID/objects"
```

## Config Files

- `~/.anytype/config.json` — Account ID, network config path
- `~/.config/anytype/` — Account data
- Keyring or config — Account key, session token

## Common Issues

| Error | Solution |
|-------|----------|
| "network id mismatch" | Re-create account with `--network-config` |
| "invite not exists" | Generate fresh invite, ensure same network |
| No sync activity | Check connectivity to nodes (ports 33020-33023) |
| "membership status" errors | Ignore — only for official Anytype cloud |

## API Docs

Full API reference: https://developers.anytype.io
