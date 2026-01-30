# PR: Add `--network-config` flag for self-hosted networks

## Title
feat(auth): add `--network-config` flag for self-hosted network support

## Description

This PR adds support for creating accounts on self-hosted Anytype (any-sync) networks by allowing users to specify a network configuration file during account creation.

### Problem

When running a self-hosted Anytype network, users cannot easily create CLI accounts that connect to their own infrastructure. The CLI defaults to the official Anytype network, and there's no way to specify a custom network during account creation.

This addresses part of [issue #6](https://github.com/anyproto/anytype-cli/issues/6) and complements [PR #8](https://github.com/anyproto/anytype-cli/pull/8) which handles custom invite links.

### Solution

1. **`--network-config` flag for `auth create`**
   ```bash
   anytype auth create my-bot --network-config ~/network.yaml
   ```

2. **Persist network config path** — The path is saved to `~/.anytype/config.json` so subsequent commands (login, serve) automatically use the same network.

3. **Auto-load network ID for `space join`** — When joining a space, the network ID is read from the saved config YAML, eliminating manual extraction.

### Changes

- `cmd/auth/create/create.go` — Add `--network-config` flag
- `cmd/auth/login/login.go` — Load network config on login
- `core/auth.go` — Pass network config to heart initialization
- `core/config/config.go` — Add `NetworkConfigPath` field
- `core/config/config_helper.go` — Add helpers for network config and network ID extraction
- `core/serviceprogram/serviceprogram.go` — Use persisted network config
- `cmd/space/join/join.go` — Auto-load network ID from config YAML

### Usage

```bash
# Create account with custom network
./anytype auth create my-bot --network-config ~/network.yaml

# Login (auto-uses saved network config)
./anytype auth login

# Start server (auto-uses saved network config)  
./anytype serve

# Join space (auto-extracts network ID from config)
./anytype space join <invite-link>
```

### Testing

- Built and tested against a self-hosted any-sync deployment
- Account creation, login, serve, and space join all work with custom network

### Notes

- This is complementary to PR #8 — that handles invite links, this handles account creation
- Network config format is the standard any-sync `network.yml`
