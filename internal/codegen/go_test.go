package codegen

import (
	"testing"
)

const legacyTableSchema = `{
  "tables": [
    {
      "name": "Users",
      "fields": [
        {"name": "id", "type": "String", "primaryKey": true},
        {"name": "email", "type": "String"}
      ],
      "resolvers": [{"name": "profile"}]
    }
  ]
}`

func TestParseTables_LegacyTablesShape(t *testing.T) {
	tables, resolvers, err := parseTables([]byte(legacyTableSchema))
	if err != nil {
		t.Fatalf("parseTables returned error: %v", err)
	}
	if len(tables) != 1 || tables[0].Name != "Users" {
		t.Fatalf("expected one Users table, got %#v", tables)
	}
	if len(tables[0].Fields) != 2 {
		t.Fatalf("expected 2 fields, got %#v", tables[0].Fields)
	}
	// primaryKey should force non-nullable
	if tables[0].Fields[0].IsNullable {
		t.Fatalf("primary key field should not be nullable")
	}
	if got := resolvers["Users"]; len(got) != 1 || got[0] != "profile" {
		t.Fatalf("expected resolver 'profile', got %#v", got)
	}
}
