import "dotenv/config";
import assert from "node:assert";
import { randomUUID } from "crypto";
import { onyx } from "@onyx.dev/onyx-database";
import { tables, OnyxSchema } from "./types";

const baseUrl = process.env.ONYX_DATABASE_BASE_URL || "https://api.onyx.dev";
const databaseId = process.env.ONYX_DATABASE_ID;
const apiKey = process.env.ONYX_DATABASE_API_KEY;
const apiSecret = process.env.ONYX_DATABASE_API_SECRET;

if (!databaseId || !apiKey || !apiSecret) {
  console.error("Skipping e2e: missing Onyx env vars");
  process.exit(0);
}

type TableName = keyof OnyxSchema;
type Row<T extends TableName> = OnyxSchema[T];

async function main() {
  const db = onyx.init<OnyxSchema>({
    baseUrl,
    databaseId,
    apiKey,
    apiSecret,
  });
  const id = randomUUID();

  // Create
  const created = await db.save(tables.User, {
    id,
    name: "CLI E2E",
    email: "cli-e2e@example.com",
  });
  assert(created, "create returned falsy");

  // Read
  const fetched = await db.findById(tables.User, id);
  assert.strictEqual((fetched as any).id, id, "fetched id mismatch");

  // Update
  const updated = await db.save(tables.User, { id, name: "CLI E2E Updated" });
  assert.strictEqual((updated as any).name, "CLI E2E Updated", "update name mismatch");

  // Delete
  const deleted = await db.delete(tables.User, id);
  assert(deleted === true, "delete response not true");

  console.log("E2E TypeScript CLI+SDK test passed");
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
