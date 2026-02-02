package main

import (
	"os"
	"path/filepath"
	"os/exec"
	"strings"
	"testing"
)

// End-to-end test for the Kotlin example using the local CLI + generated models + REST calls.
func TestExampleKotlinE2E(t *testing.T) {
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

	exampleDir := filepath.Join(repo, "examples", "onyx-kotlin")
	schemaPath := filepath.Join(exampleDir, "onyx.schema.json")
	outDir := filepath.Join(exampleDir, "onyx")

	// Fresh artifacts.
	_ = os.RemoveAll(outDir)
	runCmd(t, repo, "go", "run", "./cmd/onyx", "schema", "get", "--out", schemaPath)
	runCmd(t, repo, "go", "run", "./cmd/onyx", "gen", "--kotlin", "--schema", schemaPath, "--out", outDir, "--package", "onyx")

	if _, err := exec.LookPath("gradle"); err != nil {
		t.Skip("gradle not installed; skip Kotlin example")
	}

	// Ensure JAVA_HOME set (Homebrew OpenJDK).
	if os.Getenv("JAVA_HOME") == "" {
		_ = os.Setenv("JAVA_HOME", "/opt/homebrew/opt/openjdk@21")
	}

	// Run Kotlin example (shim gradlew uses system gradle if installed).
	runCmd(t, exampleDir, "sh", "gradlew", "-q", "run")
}
