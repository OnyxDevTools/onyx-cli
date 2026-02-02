# onyx-python example

This example uses the local `onyx` CLI to fetch a schema and generate Python types, then performs a simple CRUD flow against the Onyx API using the official Python SDK.

## Prereqs
- Python 3.9+
- Go toolchain (used to run the local CLI)

## Setup
```bash
cd examples/onyx-python
python3 -m venv .venv
. .venv/bin/activate
pip install -r requirements.txt
```

Generate the schema + stubs (run from this directory):
```bash
go run ../../cmd/onyx schema get --out onyx.schema.json
go run ../../cmd/onyx gen --py --schema onyx.schema.json --out onyx
```

Set your credentials (or ensure `onyx-database.json` exists in the repo root):
```
export ONYX_DATABASE_ID=...
export ONYX_DATABASE_API_KEY=...
export ONYX_DATABASE_API_SECRET=...
export ONYX_DATABASE_BASE_URL=https://api.onyx.dev   # optional, defaults to https://api.onyx.dev
```

## Run
```bash
PYTHONPATH=. python src/e2e.py
```

The script imports the generated models, initializes `onyx-database` using `MODEL_MAP`, and runs create/read/update/delete against the API.
