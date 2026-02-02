import "dotenv/config";
import { randomUUID } from "crypto";
import { onyx } from "@onyx.dev/onyx-database";
import { tables, OnyxSchema } from "./types";

type TableName = keyof OnyxSchema;

const db = onyx.init<OnyxSchema>();

async function main() {
  // Demonstrate CRUD on the User table
  const table = tables.User as TableName;

  const created = await db.save(table, {
    id: randomUUID(),
    email: "clitest-ts@example.com",
    name: "clitest-ts",
  });
  console.log("Created:", created);

  const fetched = await db.findById(table, (created as any).id as string);
  console.log("Fetched:", fetched);

  const updated = await db.save(table, {
    id: (created as any).id as string,
    name: "clitest-ts user Updated",
  });
  console.log("Updated:", updated);

  const removed = await db.delete(table, (created as any).id as string);
  console.log("Deleted:", removed);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
