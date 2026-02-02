package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// End-to-end test for the Go example using the local CLI + generated client.
func TestExampleGoE2E(t *testing.T) {
	if os.Getenv("ONYX_E2E") == "" {
		t.Skip("skipped: ONYX_E2E not set (set ONYX_E2E=1 to run integration test)")
	}

	repo := moduleRoot(t)
	cfgPath := filepath.Join(repo, "onyx-database.json")
	if _, err := os.Stat(cfgPath); err == nil {
		fillEnvFromConfig(t, cfgPath)
	}

	required := []string{"ONYX_DATABASE_ID", "ONYX_DATABASE_API_KEY", "ONYX_DATABASE_API_SECRET"}
	var missing []string
	for _, key := range required {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		t.Skipf("skipping: missing required envs: %s", strings.Join(missing, ", "))
	}

	exampleDir := filepath.Join(repo, "examples", "onyx-go")
	schemaPath := filepath.Join(exampleDir, "onyx.schema.json")
	genDir := filepath.Join(exampleDir, "onyx")

	// Fresh artifacts
	_ = os.RemoveAll(genDir)
	runCmd(t, repo, "go", "run", "./cmd/onyx", "schema", "get", "--out", schemaPath)
	runCmd(t, repo, "go", "run", "./cmd/onyx", "gen", "--go", "--schema", schemaPath, "--out", genDir, "--package", "onyx")

	// Run the example program.
	runCmd(t, exampleDir, "go", "run", ".")
}

