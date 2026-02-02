# Anytype CLI — Self-Hosted Setup Guide

Complete guide for connecting `anytype-cli` to a self-hosted Anytype network and integrating with OpenClaw.

---

## Why Anytype + OpenClaw?

### For Users

| Benefit | Description |
|---------|-------------|
| **Own your data** | Self-hosted means your conversations and memories never leave your infrastructure |
| **Local-first** | Works offline, syncs when connected — no cloud dependency |
| **End-to-end encrypted** | All data encrypted at rest and in transit |
| **Cross-device access** | Chat with your AI from Anytype desktop, mobile, or web |
| **Unified workspace** | Your AI assistant lives alongside your notes, tasks, and projects |
| **No vendor lock-in** | Open-source stack, export anytime, migrate freely |

### For Your AI Agent

| Benefit | Description |
|---------|-------------|
| **Structured memory** | Store memories as typed objects with relations — not just flat files |
| **Semantic search** | Query memories by type, topic, importance, date — not just text matching |
| **Rich context** | Agent can access your knowledge graph (notes, bookmarks, projects) |
| **Persistent identity** | Same agent, same memory, across all your devices |
| **Native integration** | Feels like chatting in a notes app, not a separate bot interface |

### For Self-Hosters

| Benefit | Description |
|---------|-------------|
| **Full control** | Run on your hardware, your network, your rules |
| **Privacy** | Sensitive conversations stay on-premise |
| **Customization** | Modify the stack, add features, integrate with internal systems |
| **No rate limits** | Your infrastructure, your capacity |
| **Compliance** | Meet data residency and security requirements |

### The Vision

Instead of your AI assistant being a separate chat app, it becomes part of your knowledge workspace:

```
Traditional:                     With Anytype:
┌─────────────┐                  ┌─────────────────────────────┐
│ Chat App    │                  │ Anytype Workspace           │
│ (separate)  │                  │  ├── 📝 Notes               │
└─────────────┘                  │  ├── ✅ Tasks               │
       +                         │  ├── 🔖 Bookmarks           │
┌─────────────┐                  │  ├── 📁 Projects            │
│ Notes App   │       →          │  ├── 💬 AI Chat ← you are here
│ (separate)  │                  │  └── 🧠 AI Memories         │
└─────────────┘                  └─────────────────────────────┘
```

Your AI can read your notes, remember your preferences, and work within the same space where you think and organize.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Build the CLI](#1-build-the-cli)
3. [Prepare Network Config](#2-prepare-network-config)
4. [Create Bot Account](#3-create-bot-account)
5. [Start the Server](#4-start-the-server)
6. [Join a Space](#5-join-a-space)
7. [Find Your Chat ID](#6-find-your-chat-id)
8. [OpenClaw Configuration](#7-openclaw-configuration)
9. [Memory Backend (Optional)](#8-memory-backend-optional)
10. [Testing the Integration](#9-testing-the-integration)
11. [Running as a Service](#10-running-as-a-service)
12. [API Access](#11-api-access)
13. [Troubleshooting](#troubleshooting)
14. [Reference](#reference)

---

## Prerequisites

- **Self-hosted Anytype network** running ([any-sync](https://github.com/anyproto/any-sync))
- **Network config YAML** file with your node addresses
- **Go 1.24+** and **Make** (for building from source)
- **OpenClaw** installed ([github.com/jsandai/openclaw](https://github.com/jsandai/openclaw), `release` branch)

---

## 1. Build the CLI

The upstream `anytype-cli` doesn't fully support self-hosted networks. Use this fork:

```bash
git clone https://github.com/jsandai/anytype-cli.git
cd anytype-cli
git checkout release
make build
```

Binary output: `./dist/anytype`

Verify it works:
```bash
./dist/anytype --version
./dist/anytype --help
```

You should see commands including: `auth`, `serve`, `space`, `chat`, `object`, `type`, `relation`, `memory`, `file`

---

## 2. Prepare Network Config

Create your network configuration file at `~/.config/anytype/network.yml`:

```bash
mkdir -p ~/.config/anytype
```

**Naming convention:** The file should be named `network.yml` (matching anytype-heart's convention). The CLI will reference this path when using `--network-config`.

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

Get this from your any-sync deployment's configuration.

---

## 3. Create Bot Account

```bash
./dist/anytype auth create <account-name> --network-config ~/.config/anytype/network.yml
```

Example:
```bash
./dist/anytype auth create MyBot --network-config ~/.config/anytype/network.yml
```

Output:
```
✓ Bot account created successfully!

⚠️ IMPORTANT: Save your account key in a secure location.
This is the ONLY way to authenticate your bot account.

╔════════════════════════════════════════════════════════════════════╗
║                         BOT ACCOUNT KEY                            ║
╠════════════════════════════════════════════════════════════════════╣
║ nqTDDA7vR7za2XuwXEaU384Bhk45/B36zMVDb/VMpTEEsBnx...               ║
╚════════════════════════════════════════════════════════════════════╝

📋 Bot Account Details:
   Name:       MyBot
   Account Id: YOUR_ACCOUNT_ID_HERE

✓ You are now logged in to your new bot account.
✓ Account key saved to keychain.
```

**⚠️ Save the account key securely!** It's your only authentication credential.

**📝 Note the Account ID** — you'll need this as `botIdentity` in OpenClaw config.

This creates `~/.anytype/config.json` with your network config path persisted.

---

## 4. Start the Server

```bash
./dist/anytype serve
```

**Important:** Start serve AFTER creating the account, or it won't auto-login.

Expected output:
```
Starting anytype-heart...
gRPC Web proxy started at: 127.0.0.1:31011
Starting gRPC server on 127.0.0.1:31010
```

**Normal warnings for self-hosted** (ignore these):
- `no ns peers configured` — No naming service
- `can not get membership status` — No membership system

Keep this running in a terminal, or see [Running as a Service](#10-running-as-a-service).

---

## 5. Join a Space

### Generate Invite Link (from Anytype App)

1. Open Anytype desktop/mobile app (connected to your self-hosted network)
2. Go to the space you want the bot to join
3. Space Settings → Share → Create invite link
4. Copy the link

**⚠️ The invite MUST be from your self-hosted network**, not the official Anytype cloud.

### Join the Space

With `serve` running in another terminal:

```bash
./dist/anytype space join "<invite-link>"
```

Example:
```bash
./dist/anytype space join "anytype://invite/?cid=bafybei...&key=3iunU..."
```

**Quote the URL** to prevent shell interpretation of `&`.

### Verify

```bash
./dist/anytype space list
```

You should see the space you just joined:
```
SPACE ID                                                NAME     STATUS
────────                                                ────     ──────
bafyrei<your-space-id>  MySpace   Active
```

---

## 6. Find Your Chat ID

Every Anytype space has a built-in chat. Find its ID:

```bash
./dist/anytype chat find <space-id>
```

Example:
```bash
./dist/anytype chat find bafyrei<your-space-id>
```

Output:
```
Found 1 chat(s):

  Chat
    ID: bafyrei<your-chat-id>
```

**📝 Note the Chat ID** — you'll need this for OpenClaw config.

### Test the Chat (Optional)

Send a test message:
```bash
./dist/anytype chat send <chat-id> "Hello from CLI!"
```

List recent messages:
```bash
./dist/anytype chat list <chat-id>
```

---

## 7. OpenClaw Configuration

### Install OpenClaw Fork

```bash
git clone https://github.com/jsandai/openclaw.git
cd openclaw
git checkout release
npm install
```

### Configure the Anytype Channel

Add to your OpenClaw config (`~/.openclaw/openclaw.json`):

```json
{
  "channels": {
    "anytype": {
      "enabled": true,
      "cliPath": "/absolute/path/to/anytype-cli/dist/anytype",
      "chatId": "bafyrei<your-chat-id>",
      "spaceId": "bafyrei<your-space-id>",
      "botIdentity": "YOUR_ACCOUNT_ID_HERE",
      "dmPolicy": "open"
    }
  },
  "plugins": {
    "load": {
      "paths": ["/path/to/openclaw/extensions/anytype"]
    },
    "entries": {
      "anytype": { "enabled": true }
    }
  }
}
```

### Configuration Fields

| Field | Required | Description |
|-------|----------|-------------|
| `enabled` | Yes | Enable the Anytype channel |
| `cliPath` | Yes | **Absolute path** to the `anytype` binary |
| `chatId` | Yes | Chat object ID (from step 6) |
| `spaceId` | Yes | Space ID (from step 5) |
| `botIdentity` | Yes | Your Account ID (from step 3) — prevents bot echoing itself |
| `dmPolicy` | No | `"open"`, `"allowlist"`, `"pairing"`, or `"disabled"` |

### Start OpenClaw

```bash
openclaw gateway start
```

Check logs for:
```
[anytype] Config loaded: {"enabled":true,"cliPath":"...","chatId":"..."}
[anytype] Starting subscription to chat...
[anytype] Subscription active
```

---

## 8. Memory Backend (Optional)

Use Anytype as a structured memory store for your AI agent.

### Enable in Config

```json
{
  "channels": {
    "anytype": {
      "enabled": true,
      "cliPath": "...",
      "chatId": "...",
      "spaceId": "...",
      "botIdentity": "...",
      "memoryBackend": {
        "enabled": true
      }
    }
  }
}
```

### Bootstrap Memory Type

On gateway start, you'll see:
```
[anytype] ⚠️  Memory backend enabled but not bootstrapped. Run:
  anytype memory bootstrap --space bafyrei...
Then add the output to your config.
```

Run the bootstrap:
```bash
./dist/anytype memory bootstrap --space <space-id>
```

Output:
```
✓ Memory backend ready!

Type:
  🧠 Memory
     ID:  bafyrei...
     Key: ot-<generated-type-key>

Relations:
  • memoryKind → <relation-key-1>
  • topic → <relation-key-2>
  • source → <relation-key-3>
  • importance → <relation-key-4>

Add to your config:
{
  "memoryBackend": {
    "enabled": true,
    "typeKey": "ot-<generated-type-key>",
    "relations": { ... }
  }
}
```

### Update Config

```json
{
  "channels": {
    "anytype": {
      "memoryBackend": {
        "enabled": true,
        "typeKey": "ot-<generated-type-key>",
        "relations": {
          "memoryKind": "<relation-key-1>",
          "topic": "<relation-key-2>",
          "source": "<relation-key-3>",
          "importance": "<relation-key-4>"
        }
      }
    }
  }
}
```

Restart OpenClaw gateway to apply.

---

## 9. Testing the Integration

### Test 1: Send from Anytype

1. Open Anytype app
2. Go to your space's Chat
3. Send a message: "Hello OpenClaw!"
4. Check OpenClaw logs — you should see the message received
5. OpenClaw should respond in the chat

### Test 2: Send from CLI

```bash
./dist/anytype chat send <chat-id> "Test message from CLI"
```

### Test 3: File Attachments

```bash
./dist/anytype chat send <chat-id> "Check out this image" \
  --file /path/to/image.png --space <space-id>
```

### Test 4: Real-time Subscription

In a separate terminal, watch events:
```bash
./dist/anytype chat subscribe <chat-id> --json
```

Send a message from Anytype app — you should see the event stream.

---

## 10. Running as a Service

### Systemd User Service

Create `~/.config/systemd/user/anytype.service`:

```ini
[Unit]
Description=Anytype CLI Server
After=network.target

[Service]
Type=simple
ExecStart=/path/to/anytype-cli/dist/anytype serve
Restart=on-failure
RestartSec=5

[Install]
WantedBy=default.target
```

Enable and start:
```bash
systemctl --user daemon-reload
systemctl --user enable anytype
systemctl --user start anytype
systemctl --user status anytype
```

### Or Use Built-in Service Command

```bash
./dist/anytype service install
./dist/anytype service start
./dist/anytype service status
```

---

## 11. API Access

### Create API Key

```bash
./dist/anytype auth apikey create "my-api-key"
```

### REST API Examples

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

Full API docs: https://developers.anytype.io

---

## Troubleshooting

### Build Issues

| Error | Solution |
|-------|----------|
| `cannot find -ltantivy_go` | Run `make download-tantivy` first |
| Go version error | Install Go 1.24+ |

### Authentication Issues

| Error | Solution |
|-------|----------|
| `unknown flag: --network-config` | You're using upstream binary, use `./dist/anytype` from fork |
| `No stored account key found` | Restart `serve` after `auth create` |
| `network id mismatch` | Account created on wrong network; re-create with `--network-config` |

### Connection Issues

| Error | Solution |
|-------|----------|
| `DeadlineExceeded` / `RST_STREAM` | Invite is from wrong network; check node connectivity |
| `invalid invite link format` | Quote URL with `"..."` ; rebuild if old binary |
| Server not responding | Ensure `anytype serve` is running |

### OpenClaw Issues

| Error | Solution |
|-------|----------|
| `Channel not enabled` | Check config has `"enabled": true` |
| `cliPath not configured` | Use absolute path to binary |
| Bot echoes its own messages | Set `botIdentity` to your Account ID |
| `Subscription disconnected` | Check `anytype serve` is running; plugin auto-reconnects |

### Normal Warnings (Safe to Ignore)

These appear on self-hosted networks and are expected:
- `no ns peers configured` — No naming service (cloud-only feature)
- `can not get membership status` — No membership system (cloud-only feature)
- `error fetching global names` — No naming service

---

## Reference

### Config Files

| File | Purpose |
|------|---------|
| `~/.anytype/config.json` | Account ID, tech space ID, network config path |
| `~/.config/anytype/network.yml` | Your self-hosted network nodes and ID |
| `~/.openclaw/openclaw.json` | OpenClaw configuration |
| System keyring | Account key, session tokens |

### Ports

| Port | Service |
|------|---------|
| 31010 | gRPC server (internal) |
| 31011 | gRPC Web proxy |
| 31012 | REST API |
| 33020-33023 | Any-sync nodes (your self-hosted network) |

### CLI Commands Quick Reference

```bash
# Account
anytype auth create <name> --network-config ~/.config/anytype/network.yml
anytype auth login --account-key <key>
anytype auth whoami

# Server
anytype serve
anytype serve --quiet
anytype serve --verbose

# Spaces
anytype space list
anytype space join "<invite-url>"

# Chat
anytype chat find <space-id>
anytype chat list <chat-id>
anytype chat send <chat-id> "message"
anytype chat send <chat-id> "message" --file /path/to/file --space <space-id>
anytype chat subscribe <chat-id> --json
anytype chat react <chat-id> <msg-id> "👍"
anytype chat edit <chat-id> <msg-id> "new text"
anytype chat delete <chat-id> <msg-id>

# Objects
anytype object create --space <id> --type <type> --name "Name"
anytype object search --space <id> --query "text"
anytype object get <id> --space <space-id>
anytype object update <id> --space <space-id> --name "New Name"
anytype object delete <id>

# Types & Relations
anytype type list --space <id>
anytype type create --space <id> --name "Name" --icon "🔮"
anytype relation list --space <id>
anytype relation create --space <id> --name "field" --format text

# Memory Backend
anytype memory bootstrap --space <id>

# Files
anytype file upload <space-id> /path/to/file
anytype file download <object-id> /path/to/save
```

### Fork Features (vs Upstream)

| Feature | Upstream | This Fork |
|---------|----------|-----------|
| `--network-config` flag | ❌ | ✅ |
| Persist network config path | ❌ | ✅ |
| Auto-load network ID from YAML | ❌ | ✅ |
| `anytype://invite/` URL format | ❌ | ✅ |
| `--quiet`/`--verbose` for serve | ❌ | ✅ |
| Chat commands (send, list, subscribe) | ❌ | ✅ |
| Object/Type/Relation commands | ❌ | ✅ |
| Memory bootstrap | ❌ | ✅ |
| File upload/download | ❌ | ✅ |

---

## Quick Start Checklist

```
[ ] Build CLI: git clone → checkout release → make build
[ ] Create ~/.config/anytype/network.yml with your self-hosted nodes
[ ] Create bot account: anytype auth create <name> --network-config ~/.config/anytype/network.yml
[ ] Save account key securely
[ ] Note Account ID (for botIdentity)
[ ] Start server: anytype serve
[ ] Generate invite link from Anytype app
[ ] Join space: anytype space join "<invite>"
[ ] Note Space ID: anytype space list
[ ] Find Chat ID: anytype chat find <space-id>
[ ] Configure OpenClaw with all IDs
[ ] Start OpenClaw: openclaw gateway start
[ ] Test: send message in Anytype → check OpenClaw responds
[ ] (Optional) Bootstrap memory: anytype memory bootstrap --space <id>
```

---

*Last updated: 2026-02-02*
