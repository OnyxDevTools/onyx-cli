package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// End-to-end test for the TypeScript Node example using the local CLI.
func TestExampleNodeE2E(t *testing.T) {
	if os.Getenv("ONYX_E2E") == "" {
		t.Skip("skipped: ONYX_E2E not set (set ONYX_E2E=1 to run integration test)")
	}

	// Populate env vars from root onyx-database.json if present and envs missing.
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

	exampleDir := filepath.Join(repo, "examples", "onyx-node")

	if _, err := os.Stat(filepath.Join(exampleDir, "node_modules")); err != nil {
		runCmd(t, exampleDir, "npm", "install", "--silent")
	}

	runCmd(t, exampleDir, "npm", "run", "test:e2e")
}

func fillEnvFromConfig(t *testing.T, path string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Logf("fillEnvFromConfig: read error: %v", err)
		return
	}
	var cfg struct {
		DatabaseID string `json:"databaseId"`
		BaseURL    string `json:"baseUrl"`
		APIKey     string `json:"apiKey"`
		APISecret  string `json:"apiSecret"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Logf("fillEnvFromConfig: parse error: %v", err)
		return
	}
	setIfEmpty := func(key, val string) {
		if os.Getenv(key) == "" && val != "" {
			_ = os.Setenv(key, val)
		}
	}
	setIfEmpty("ONYX_DATABASE_ID", cfg.DatabaseID)
	setIfEmpty("ONYX_DATABASE_BASE_URL", cfg.BaseURL)
	setIfEmpty("ONYX_DATABASE_API_KEY", cfg.APIKey)
	setIfEmpty("ONYX_DATABASE_API_SECRET", cfg.APISecret)
}
