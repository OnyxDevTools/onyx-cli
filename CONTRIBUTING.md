<!-- CONTRIBUTING.md -->

# Contributing to onyx (CLI)

Thanks for helping shape the cross-platform `onyx` CLI. This repo is early and TypeScript-parity is the guiding constraint.

## Prerequisites
- Go 1.24.x (ensure `go version` and `go env GOROOT` point to the same 1.24 toolchain).
- Git.

## Quickstart
```bash
go mod tidy
go build -o ./bin/onyx ./cmd/onyx
./bin/onyx --help
./bin/onyx help version
man -M docs/man onyx

# show version (ldflags can override Version/Commit/Date)
./bin/onyx version
```

The checked-in man page lives at `docs/man/man1/onyx.1` (mdoc) and should render cleanly with `mandoc`. On macOS/BSD `man` does not support `-l`; instead use `man -M docs/man onyx` or `MANPATH=$(pwd)/docs/man man onyx`.
## Contribution flow
1. Open/assign an issue for non-trivial work.
2. Keep changes minimal and scoped to a single concern.
3. Add or update docs/tests alongside code where applicable.
4. Run `go test ./...` (once tests exist) and `go build ./cmd/onyx` before opening a PR.
5. Use clear commit messages; avoid force-pushing shared branches.

## Coding notes
- Cobra is used for the command tree; new commands live under `cmd/onyx/cmd`.
- Default CLI behavior must mirror the TypeScript SDK for schema/codegen until expanded otherwise.
- Avoid breaking changes to the documented command surface without discussion.
