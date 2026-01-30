# PR: Add `--quiet` and `--verbose` flags to serve command

## Title
feat(serve): add `--quiet` and `--verbose` flags for log control

## Description

Adds flags to control log verbosity when running `anytype serve`.

### Problem

The serve command outputs a lot of debug information by default, which can be noisy in production or when running as a background service.

### Solution

Add two mutually exclusive flags:

- `--quiet` / `-q` — Suppress most output, only show errors
- `--verbose` / `-v` — Show detailed debug output

Default behavior (no flags) remains unchanged.

### Usage

```bash
# Quiet mode - minimal output
anytype serve --quiet

# Verbose mode - debug output  
anytype serve --verbose

# Short flags
anytype serve -q
anytype serve -v
```

### Changes

- `cmd/serve/serve.go` — Add `--quiet` and `--verbose` flags, configure logger based on selection
