<!-- README.md -->

# onyx (CLI)

A cross-platform CLI for **Onyx Cloud Database** (onyx.dev) focused on:

- **Schema management** (get / validate / diff / publish)
- **Code generation** (TypeScript-first)

> Status: **WIP / not released yet**. This repository is being built incrementally. Packaging (Homebrew, etc) will come after the core CLI behavior is stable.

## Why this exists
The TypeScript SDK currently provides helper CLIs (`onyx-gen` and `onyx-schema`). This repo consolidates that developer workflow into a single, globally-installable binary:

- One executable: `onyx`
- Consistent subcommands across platforms and languages (eventually)
- Minimal friction for new projects

## Current scope (TypeScript-first)
The initial scope is intentionally narrow (TypeScript is the canonical reference; Python and Go generation are implemented to mirror their SDK helpers):

1. Implement `onyx schema ...` commands with behavior matching the current TypeScript CLI semantics.
2. Implement `onyx gen` code generation with TypeScript parity and matching outputs for the published Python (`onyx-database-python`) and Go (`onyx-gen-go`) generators.
3. Use the **same credential/config resolution chain** as the TypeScript SDK.

Schema commands and the credential/config chain must exactly match the current TypeScript SDK behavior; the TS SDK is the canonical reference.

## CLI usage

### Install (preview)

Until official packages are published, you can install a tagged release tarball:

```
curl -fsSL https://raw.githubusercontent.com/OnyxDevTools/onyx-cli/main/scripts/install.sh | bash

# or specific version
curl -fsSL https://raw.githubusercontent.com/OnyxDevTools/onyx-cli/main/scripts/install.sh | bash -s -- v0.1.0
```

Flags/env:
- `BINDIR=/usr/local/bin` to change install location (falls back to `~/.local/bin`).
- `ONYX_TAG=v0.0.0` to pin a specific release (default: latest).
Manual direct install (if you already downloaded a tarball):
```
tar -xzvf onyx_<os>_<arch>.tar.gz
install -m 0755 onyx /usr/local/bin/onyx   # or another dir on PATH
```

### Homebrew (macOS)

Standard (public tap, no SSH needed):
```
brew tap OnyxDevTools/onyx-cli
brew install onyx-cli
onyx version
```
Note: The token `onyx` is already used by the macOS OnyX cask. Install the CLI with `onyx-cli` as shown above.

If you ever get prompted for credentials, the tap likely isn’t public yet—publish the tap repo (`OnyxDevTools/homebrew-onyx-cli`) and retry.

### Releasing (maintainers)

1) Bump version, build artifacts, update tap formula, tag, push, and create GitHub release:
```
scripts/bump-version.sh patch   # or minor/major
```
The script will:
- Ensure main and tap repos are clean and in sync
- Bump `cmd/onyx/cmd/version.go`
- Build darwin/linux (amd64/arm64) tarballs into `dist/`
- Update `homebrew-onyx-cli/Formula/onyx-cli.rb` with version + shas
- Commit/tag/push main and tap repos
- Create the GitHub release and attach the tarballs

2) Verify release assets:
   - GitHub release `vX.Y.Z` exists with four tarballs.
   - Tap repo `OnyxDevTools/homebrew-onyx-cli` has the updated formula.

3) Smoke test:
```
brew uninstall onyx-cli
brew untap OnyxDevTools/onyx-cli
brew tap OnyxDevTools/onyx-cli
brew install onyx-cli
onyx version
```

3) Upgrade flow
```
brew update
brew upgrade onyx-cli
onyx version
```

### Command tables (compact)

**Info**

| Command | Flags (core) | Behavior / defaults |
|---------|--------------|---------------------|
| `onyx info` | *(schema/info flags below)* | Shows resolved config sources, config path, connection check (Schema API ping). |

Shared credential flags (all schema/info commands): `--database-id`, `--base-url`, `--api-key`, `--api-secret`, `--ai-base-url`, `--default-model`, `--config` (overrides `ONYX_CONFIG_PATH` and search chain).

**Init (Go helper)**

| Command | Flags (core) | Behavior / defaults |
|---------|--------------|---------------------|
| `onyx init` | `--schema <path>`, `--out <dir>`, `--package <name>`, `--force` | Writes `generate.go` with `//go:generate onyx gen --go ...`. Defaults: schema `./api/onyx.schema.json`; out `./gen/onyx`; package `onyx`. |

**Codegen**

| Command | Flags (core) | Defaults / notes |
|---------|--------------|------------------|
| `onyx gen` | Language flags: `--typescript`/`--ts`, `--java`, `--kotlin`/`--kt`, `--python`/`--py`, `--go`/`--golang` (TypeScript + Python + Go implemented).<br/>Core flags: `--source auto\|api\|file`, `--schema <path>`, `--out <file\|dir>[,more]`, `--tables a,b` (api), `--name <type>` (TS), `--base <name>` (TS), `--package <name>` (Go), `--overwrite`, `-q/--quiet` | Defaults: source `file`; schema `./onyx.schema.json`; out `./onyx/types.ts` (TS), `./onyx` (Python), `./gen/onyx` (Go); type name `OnyxSchema`; Go package `onyx`; overwrite on.<br/>If no language flag is given, `codegenLanguage` in config or `ONYX_CODEGEN_LANGUAGE` (`typescript`/`ts`/`java`/`kotlin`/`kt`/`python`/`py`/`go`/`golang`) is used. TS output mirrors `onyx-gen`: interfaces per entity, schema mapping type + const, `tables` enum. Python output mirrors onyx-database-python (models.py/tables.py/schema.py). Go output mirrors onyx-gen-go: generates `common.go` plus per-table typed clients (query helpers, updates structs, paging iterators, cascades). |

**Schema**

| Command | Flags (core) | Behavior / defaults |
|---------|--------------|---------------------|
| `onyx schema get [file]` | `--tables a,b` (stdout), `--print` (stdout) | Default file `./onyx.schema.json`. Writes file unless tables/print is used. |
| `onyx schema publish [file]` | *(none)* | Default file `./onyx.schema.json`. Validates first; publishes only if valid. |
| `onyx schema validate [file]` | *(none)* | Default file `./onyx.schema.json`. Exits non-zero on validation errors. |
| `onyx schema diff [file]` | *(none)* | Default file `./onyx.schema.json`. Prints YAML diff vs API schema. |
| `onyx schema info` | *(none)* | Shows resolved config sources, config path, connection check (Schema API ping). |


## Credential & config resolution (must match TypeScript)
This CLI must match the TypeScript SDK’s resolution chain:

**explicit config ➜ environment variables ➜ `ONYX_CONFIG_PATH` file ➜ project config file ➜ home profile**

### Environment variables
- `ONYX_DATABASE_ID`
- `ONYX_DATABASE_BASE_URL` (optional; defaults to `https://api.onyx.dev`)
- `ONYX_DATABASE_API_KEY`
- `ONYX_DATABASE_API_SECRET`
- optional: `ONYX_AI_BASE_URL`
- optional: `ONYX_DEFAULT_MODEL`
- optional: `ONYX_CODEGEN_LANGUAGE` (defaults to `typescript`; aliases: `ts`, `java`, `kotlin`, `kt`, `python`, `py`, `go`, `golang`) — used when no language flag is provided to `onyx gen`

Config JSON keys of interest:
- `codegenLanguage`: optional (defaults to `typescript`; same aliases as above). If no language flag is given, `onyx gen` uses this value.

Defaults (when unspecified):
- Base URL: `https://api.onyx.dev`
- AI Base URL: `https://ai.onyx.dev`
- Default model: `onyx`

### Config file search order
1) If `--config` is provided, that path is used (errors if unreadable).  
2) Else if `ONYX_CONFIG_PATH` is set, that path is used (errors if unreadable).  
3) Otherwise search JSON configs in this order:
   1. `./onyx-database-<databaseId>.json`
   2. `./onyx-database.json`
   3. `./config/onyx-database-<databaseId>.json`
   4. `./config/onyx-database.json`
   5. `~/.onyx/onyx-database-<databaseId>.json`
   6. `~/.onyx/onyx-database.json`
   7. `~/onyx-database.json`

## Schema file conventions
- Default schema file path: `./onyx.schema.json`
- `onyx schema get` overwrites the default file unless output is redirected via `--print`/`--tables` printing behavior.

## Canonical references
- TypeScript SDK (canonical CLI + credential semantics): https://github.com/OnyxDevTools/onyx-database
- Onyx Cloud Console: https://cloud.onyx.dev
- Onyx website: https://onyx.dev
