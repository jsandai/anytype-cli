# Anytype CLI — Self-Hosted Setup Guide

This guide walks through connecting `anytype-cli` to a self-hosted Anytype (any-sync) network.

## Prerequisites

- Self-hosted Anytype network running ([any-sync](https://github.com/anyproto/any-sync))
- Network config YAML file with your node addresses
- Go 1.24+ and Make (for building from source)

## 1. Build the Fork

The upstream `anytype-cli` doesn't fully support self-hosted networks. Use this fork:

```bash
git clone https://github.com/jsandai/anytype-cli.git
cd anytype-cli
git checkout feat/network-config
make build
```

Binary output: `./dist/anytype`

## 2. Prepare Network Config

Save your network configuration as `~/network.yaml`:

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

## 3. Create Bot Account

```bash
./dist/anytype auth create <account-name> --network-config ~/network.yaml
```

Example:
```bash
./dist/anytype auth create Archie --network-config ~/network.yaml
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
   Name:       Archie
   Account Id: A6oGP5omrQMynpvZXcV8rBW6j2vhDNpiTFPpdXUW1MxkiV5s

✓ You are now logged in to your new bot account.
✓ Account key saved to keychain.
```

**⚠️ Save the account key!** It's your only authentication credential.

This creates `~/.anytype/config.json`:
```json
{
  "accountId": "A6oGP5omrQMynpvZXcV8rBW6j2vhDNpiTFPpdXUW1MxkiV5s",
  "techSpaceId": "bafyreif...",
  "networkConfigPath": "/home/dev/network.yaml"
}
```

The `networkConfigPath` is persisted — subsequent commands auto-use your network config.

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

Some warnings are normal for self-hosted:
- `no ns peers configured` — No naming service (expected)
- `can not get membership status` — No membership system (expected)

Keep this running in a terminal or background it.

## 5. Generate Invite Link (from Anytype App)

On your self-hosted network:
1. Open Anytype desktop/mobile app (connected to your self-hosted network)
2. Go to the space you want the bot to join
3. Space Settings → Share → Create invite link
4. Copy the link

Invite formats supported:
- `https://invite.any.coop/{cid}#{key}` (standard)
- `anytype://invite/?cid={cid}&key={key}` (app deep link)

**⚠️ The invite MUST be from your self-hosted network**, not the official Anytype cloud.

## 6. Join a Space

With serve running in another terminal:

```bash
./dist/anytype space join "<invite-link>"
```

Example:
```bash
./dist/anytype space join "anytype://invite/?cid=bafybei...&key=3iunU..."
```

Output:
```
Joining space 'Archie' created by James...
✓ Successfully sent join request to space 'bafyreib...'
```

**Quote the URL** to prevent shell interpretation of `&`.

## 7. Create API Key

```bash
./dist/anytype auth apikey create "my-api-key"
```

## 8. Use the API

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

---

## Config Files Reference

| File | Purpose |
|------|---------|
| `~/.anytype/config.json` | Account ID, tech space ID, network config path |
| `~/network.yaml` | Network nodes and ID (your self-hosted config) |
| System keyring | Account key, session tokens |

## Fork Features (vs Upstream)

| Feature | Upstream | This Fork |
|---------|----------|-----------|
| `--network-config` flag | ❌ | ✅ |
| Persist network config path | ❌ | ✅ |
| Auto-load network ID from YAML | ❌ | ✅ |
| `anytype://invite/` URL format | ❌ | ✅ |
| `--quiet`/`--verbose` for serve | ❌ | ✅ |

## Troubleshooting

| Error | Cause | Solution |
|-------|-------|----------|
| `unknown flag: --network-config` | Using upstream binary | Use `./dist/anytype` from fork |
| `invalid invite link format` | Old binary or wrong format | Rebuild fork; quote URL with `"..."` |
| `No stored account key found` | Serve started before account created | Restart serve after `auth create` |
| `DeadlineExceeded` / `RST_STREAM` | Network mismatch or connectivity | Ensure invite is from YOUR network; check node connectivity |
| `network id mismatch` | Account created on wrong network | Re-create account with `--network-config` |
| `no ns peers configured` | Normal for self-hosted | Ignore — naming service is for official cloud only |
| `membership status` errors | Normal for self-hosted | Ignore — membership is for official cloud only |

## Ports Reference

| Port | Service |
|------|---------|
| 31010 | gRPC server |
| 31011 | gRPC Web proxy |
| 31012 | REST API |
| 33020 | Tree node (sync) |
| 33021 | Coordinator node |
| 33022 | Consensus node |
| 33023 | File node |

## Full API Docs

https://developers.anytype.io
