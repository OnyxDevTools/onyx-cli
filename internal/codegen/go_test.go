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

func TestParseTables_RespectsNullableFlag(t *testing.T) {
	const schema = `{
  "tables": [
    {
      "name": "Things",
      "fields": [
        {"name": "id", "type": "String", "primaryKey": true},
        {"name": "optional", "type": "String", "nullable": true},
        {"name": "explicitNonNull", "type": "Int", "nullable": false},
        {"name": "legacyNullable", "type": "Timestamp", "isNullable": true}
      ]
    }
  ]
}`
	tables, _, err := parseTables([]byte(schema))
	if err != nil {
		t.Fatalf("parseTables returned error: %v", err)
	}
	fields := tables[0].Fields
	check := func(name string, want bool) {
		for _, f := range fields {
			if f.Name == name {
				if f.IsNullable != want {
					t.Fatalf("field %s nullable=%t, want %t", name, f.IsNullable, want)
				}
				return
			}
		}
		t.Fatalf("field %s not found", name)
	}
	check("id", false)                  // primary key forces false
	check("optional", true)             // nullable:true honored
	check("explicitNonNull", false)     // nullable:false honored
	check("legacyNullable", true)       // legacy isNullable honored
}

func TestMapGoType_UsesPointersForNullable(t *testing.T) {
	cases := []struct {
		schemaType string
		nullable   bool
		want       string
	}{
		{"String", true, "*string"},
		{"Boolean", true, "*bool"},
		{"Int", true, "*int64"},
		{"Float", true, "*float64"},
		{"Timestamp", true, "*time.Time"},
		{"String", false, "string"},
		{"Timestamp", false, "time.Time"},
	}
	for _, tc := range cases {
		if got := mapGoType(tc.schemaType, tc.nullable, false); got != tc.want {
			t.Fatalf("mapGoType(%q, nullable=%t) = %q, want %q", tc.schemaType, tc.nullable, got, tc.want)
		}
	}
}
