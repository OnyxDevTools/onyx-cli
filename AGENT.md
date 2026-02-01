<!-- AGENT.md -->

# AI Agent Guide

This repository is the official cross-platform CLI for **Onyx Cloud Database** (onyx.dev), shipped as a single binary named `onyx`.

The goal is to provide a *developer-first*, *scriptable*, *minimal-friction* CLI that works the same on **macOS**, **Linux**, and **Windows**.

## Mission
Unify the functionality that currently exists across SDK helper CLIs into a single, consistent interface:

- TypeScript SDK currently ships helper CLIs: `onyx-gen` + `onyx-schema`.
- This repo will ship: `onyx` (one binary, subcommands).

## Current scope (TypeScript-first)
The initial delivery is **TypeScript-first** and must be compatible with the current TypeScript CLI semantics.

### Command surface (current)
- `onyx gen --typescript`
- `onyx schema get`
- `onyx schema validate`
- `onyx schema diff`
- `onyx schema publish`

### Binary targets (important)
- `cmd/onyx`: production binary **onyx** that gets shipped in releases. Change user-visible behavior here.
- `cmd/localonyx`: dev-only shim for quick local testing. Keep it in sync with `cmd/onyx`, but don’t rely on it being distributed.

## Hard requirements
### 1) Config + credential resolution parity (TypeScript canonical)
This CLI must use the same credential/config chain logic as the TypeScript SDK:

**explicit config ➜ env vars ➜ `ONYX_CONFIG_PATH` file ➜ project config file ➜ home profile**

Environment variables (canonical):
- `ONYX_DATABASE_ID`
- `ONYX_DATABASE_BASE_URL` (optional; defaults to `https://api.onyx.dev`)
- `ONYX_DATABASE_API_KEY`
- `ONYX_DATABASE_API_SECRET`
- optional: `ONYX_AI_BASE_URL`
- optional: `ONYX_DEFAULT_MODEL`
- optional: `ONYX_CODEGEN_LANGUAGE` (e.g., `typescript`) when no language flag is provided to `onyx gen`
  - aliases accepted: `ts`, `java`, `kotlin`, `kt`, `python`, `py`, `go`, `golang`

If not provided anywhere, fall back to:
- Base URL: `https://api.onyx.dev`
- AI Base URL: `https://ai.onyx.dev`
- Default model: `onyx`
- Codegen language: `typescript` (aliases accepted: ts, java, kotlin, kt, python, py, go, golang)

When `ONYX_CONFIG_PATH` is unset, search config files in this order:
1. `./onyx-database-<databaseId>.json`
2. `./onyx-database.json`
3. `./config/onyx-database-<databaseId>.json`
4. `./config/onyx-database.json`
5. `~/.onyx/onyx-database-<databaseId>.json`
6. `~/.onyx/onyx-database.json`
7. `~/onyx-database.json`

### 2) Schema CLI parity (TypeScript canonical)
`onyx schema ...` behavior must match the TypeScript `onyx-schema` behavior:
- default schema file path: `./onyx.schema.json`
- `get` writes the file unless `--tables` or `--print` is provided (then it prints to stdout instead)
- `validate` validates without publishing
- `diff` prints a YAML diff vs the API schema
- `publish` validates before publishing

### 3) Minimal runtime friction
- End users should not need Node/Python/Java to use the CLI.
- Prefer a small dependency footprint (stdlib-first where possible).
- Output must be stable and automation-friendly (CI scripts).

### 4) Cross-platform correctness
- Handle filesystem paths and home directory resolution correctly across macOS/Linux/Windows.
- Never assume POSIX-only behavior (signals, path separators, file permissions).

## Engineering guidelines for agents
- Keep changes tightly scoped to the current task.
- Do not introduce breaking changes to the agreed CLI surface.
- Prefer explicit, testable behavior over “magic”.
- Write error messages that are actionable and include next steps.
- Ensure commands return non-zero on failures (validation errors, API errors, etc).
- Avoid writing secrets to stdout/stderr.

## References (canonical)
- TypeScript SDK (canonical semantics): https://github.com/OnyxDevTools/onyx-database
