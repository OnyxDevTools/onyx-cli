# onyx-node (TypeScript example)

Small Node.js example showing how to:
- fetch a schema with the `onyx` CLI
- generate TypeScript types from that schema
- perform simple CRUD-style calls with type safety

## Prereqs
- Node.js 18+ (for built-in `fetch`)
- `onyx` CLI on your PATH

## Setup
```bash
cd examples/onyx-node
npm install
```

Export your connection info (or put them in `.env`):
```
ONYX_DATABASE_ID=...
ONYX_DATABASE_BASE_URL=https://api.onyx.dev
ONYX_DATABASE_API_KEY=...
ONYX_DATABASE_API_SECRET=...
```

Generate schema + types with the local CLI source (no need to install a global binary):
```bash
npm run gen
```

Run the sample script:
```bash
npm start
```

The script logs basic create/read/update/delete calls against the `User` table using the generated types for payload shape checking.
