<!-- CONTRIBUTING.md -->

# Contributing to onyx (CLI)

Thanks for helping shape the cross-platform `onyx` CLI. This repo is early and TypeScript-parity is the guiding constraint.

## Quickstart
```bash
go mod tidy
go build -o ./bin/onyx ./cmd/onyx
./bin/onyx --help
./bin/onyx help version
man -M docs/man onyx
./bin/onyx version
```

## Release flow
```
   scripts/bump-version.sh   # will prompt for patch/minor/major
   brew tap OnyxDevTools/onyx-cli
   brew install onyx
   onyx version
```
