# onyx-go example

Simple Go example that uses the local `onyx` CLI to fetch a schema, generate a typed Go client, and perform CRUD with the official Go SDK.

## Setup
```bash
cd examples/onyx-go
# generate schema + client (run from repo root normally)
go run ../../cmd/onyx schema get --out onyx.schema.json
go run ../../cmd/onyx gen --go --schema onyx.schema.json --out onyx --package onyx

# ensure creds via env or onyx-database.json
export ONYX_CONFIG_PATH=../../onyx-database.json  # optional if file exists here
go run .
```

The program will create, fetch, update, and delete a User record using the generated client.
