package codegen

import (
	"os"
	"path/filepath"
	"testing"
)

const pythonSchemaJSON = `{
  "entities": [
    {
      "name": "User",
      "identifier": {"name": "id"},
      "attributes": [
        {"name": "id", "type": "String", "isNullable": false},
        {"name": "email", "type": "String", "isNullable": true}
      ]
    },
    {
      "name": "Team",
      "identifier": {"name": "id"},
      "attributes": [
        {"name": "id", "type": "String", "isNullable": false},
        {"name": "name", "type": "String", "isNullable": false}
      ]
    }
  ]
}`

func TestRenderPythonInitExports(t *testing.T) {
	outDir := t.TempDir()

	if err := RenderPython([]byte(pythonSchemaJSON), outDir, true); err != nil {
		t.Fatalf("RenderPython() error = %v", err)
	}

	got, err := os.ReadFile(filepath.Join(outDir, "__init__.py"))
	if err != nil {
		t.Fatalf("read __init__.py: %v", err)
	}

	goldenPath := filepath.Join("testdata", "python", "__init__.py.golden")
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden %s: %v", goldenPath, err)
	}

	if string(got) != string(want) {
		t.Fatalf("__init__.py mismatch\nwant:\n%s\ngot:\n%s", string(want), string(got))
	}
}
