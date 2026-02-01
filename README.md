<!-- README.md -->

# onyx (CLI)

A cross-platform CLI for **Onyx Cloud Database** (onyx.dev) focused on:

- **Schema management** (get / validate / diff / publish)
- **Code generation** (typescript, python, go, kotlin)

## Why this exists
The TypeScript SDK currently provides helper CLIs (`onyx gen` and `onyx schema`). This repo consolidates that developer workflow into a single, globally-installable binary:

- One executable: `onyx`
- Consistent subcommands across platforms and languages (eventually)
- Minimal friction for new projects

## CLI usage

### Install

Until official packages are published, you can install a tagged release tarball:

```
curl -fsSL https://raw.githubusercontent.com/OnyxDevTools/onyx-cli/main/scripts/install.sh | bash

# or specific version
curl -fsSL https://raw.githubusercontent.com/OnyxDevTools/onyx-cli/main/scripts/install.sh | bash -s -- v0.1.0
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

### Upgrade flow
```
brew update
brew upgrade onyx-cli
onyx version
```

### Commands

**Version**

| Command | Flags (core) | Behavior / defaults |
|---------|--------------|---------------------|
| `onyx version` | *(none)* | Prints the CLI version string. |

**Info**

| Command | Flags (core) | Behavior / defaults |
|---------|--------------|---------------------|
| `onyx info` | *(none)* | Shows resolved config sources, config path, connection check (Schema API ping). |

Shared credential flags (all schema/info commands): `--database-id`, `--base-url`, `--api-key`, `--api-secret`, `--ai-base-url`, `--default-model`, `--config` (overrides `ONYX_CONFIG_PATH` and search chain).

**Init (Only helpful for GO SDK)**

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
3) Otherwise, the CLI searches for the first readable JSON config in this order (paths with `<databaseId>` are only used if the database ID was already provided via flags or env vars):
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
