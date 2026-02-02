package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// End-to-end test for the Python example using the local CLI + generated SDK.
func TestExamplePythonE2E(t *testing.T) {
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

	exampleDir := filepath.Join(repo, "examples", "onyx-python")
	schemaPath := filepath.Join(exampleDir, "onyx.schema.json")
	genDir := filepath.Join(exampleDir, "onyx")

	// Fresh artifacts for the example (generation happens outside the Python script).
	_ = os.RemoveAll(genDir)
	runCmd(t, repo, "go", "run", "./cmd/onyx", "schema", "get", "--out", schemaPath)
	runCmd(t, repo, "go", "run", "./cmd/onyx", "gen", "--py", "--schema", schemaPath, "--out", genDir)

	venvDir := filepath.Join(exampleDir, ".venv")
	python := filepath.Join(venvDir, "bin", "python")
	pip := filepath.Join(venvDir, "bin", "pip")

	// Create venv if missing to avoid PEP 668 system restrictions.
	if _, err := os.Stat(python); err != nil {
		runCmd(t, exampleDir, "python3", "-m", "venv", ".venv")
	}

	// Install deps into the venv (idempotent).
	runCmd(t, exampleDir, python, "-m", "pip", "install", "--upgrade", "pip")
	runCmd(t, exampleDir, pip, "install", "-r", "requirements.txt")

	// Execute the example script inside the venv so it can import the SDK and generated stubs.
	// Ensure PYTHONPATH points at the example root so imports of the generated `onyx` package succeed when run from src/.
	prev := os.Getenv("PYTHONPATH")
	_ = os.Setenv("PYTHONPATH", exampleDir+string(os.PathListSeparator)+prev)
	runCmd(t, exampleDir, python, "src/e2e.py")
}
