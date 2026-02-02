package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// Integration tests that exercise the CLI against a live Onyx service.
// Opt-in: set ONYX_E2E=1 to run.
func TestVersion(t *testing.T) {
	workdir := prepareE2EWorkspace(t)
	out, _ := runCLI(t, workdir, "version")
	assertOutputContains(t, out, "onyx version")
}

func TestInfo(t *testing.T) {
	workdir := prepareE2EWorkspace(t)
	out, _ := runCLI(t, workdir, "info")
	assertOutputContains(t, out, "databaseId:", "baseUrl:", "apiKey:")
}

func TestInit(t *testing.T) {
	workdir := prepareE2EWorkspace(t)
	// Allow repeat runs by forcing overwrite of generate.go if it exists.
	_ = os.Remove(filepath.Join(workdir, "generate.go"))
	out, _ := runCLI(t, workdir, "init", "--force")
	assertOutputContains(t, out, "Wrote generate.go")
	if _, err := os.Stat(filepath.Join(workdir, "generate.go")); err != nil {
		t.Fatalf("generate.go missing after init: %v", err)
	}
}

func TestSchemaGet(t *testing.T) {
	workdir := prepareE2EWorkspace(t)
	out, _ := runCLI(t, workdir, "schema", "get", "--out", "api/onyx.schema.json")
	assertOutputContains(t, out, "Wrote schema to")
	assertSchemaExists(t, workdir)
}

func TestSchemaDiff(t *testing.T) {
	workdir := prepareE2EWorkspace(t)
	runCLI(t, workdir, "schema", "get", "--out", "api/onyx.schema.json")
	out, _ := runCLI(t, workdir, "schema", "diff", "api/onyx.schema.json")
	assertOutputContainsAny(t, out, "Diff between", "No differences found")
}

func TestSchemaPublish(t *testing.T) {
	workdir := prepareE2EWorkspace(t)
	runCLI(t, workdir, "schema", "get", "--out", "api/onyx.schema.json")
	out, _ := runCLI(t, workdir, "schema", "publish", "api/onyx.schema.json")
	assertOutputContains(t, out, "Published revision")
}

func TestGenGo(t *testing.T) {
	workdir := prepareE2EWorkspace(t)
	schemaPath := ensureSchema(t, workdir)
	tableNames := loadTableNames(t, schemaPath)
	outDir := filepath.Join(workdir, "gen-go")
	out, _ := runCLI(t, workdir, "gen", "--go", "--schema", schemaPath, "--out", outDir, "--package", "onyxgen")
	assertPathNotEmpty(t, outDir)
	logPathContents(t, outDir)
	assertFileContains(t, filepath.Join(outDir, "common.go"), "package onyxgen")
	for _, name := range tableNames {
		goFile := filepath.Join(outDir, strings.ToLower(name)+".go")
		assertFileContains(t, goFile, "type "+exportName(name)+" struct")
	}
	assertOutputContains(t, out, "Go client ->")
}

func TestGenTS(t *testing.T) {
	workdir := prepareE2EWorkspace(t)
	schemaPath := ensureSchema(t, workdir)
	tableNames := loadTableNames(t, schemaPath)
	outPath := filepath.Join(workdir, "gen-ts", "types.ts")
	out, _ := runCLI(t, workdir, "gen", "--ts", "--schema", schemaPath, "--out", outPath, "--name", "OnyxSchema")
	assertFileExists(t, outPath)
	logPathContents(t, outPath)
	assertFileContains(t, outPath, "export type OnyxSchema", "export enum tables")
	for _, name := range tableNames {
		assertFileContains(t, outPath, "export interface "+name)
	}
	assertOutputContains(t, out, "TypeScript types ->")
}

func TestGenPy(t *testing.T) {
	workdir := prepareE2EWorkspace(t)
	schemaPath := ensureSchema(t, workdir)
	tableNames := loadTableNames(t, schemaPath)
	outDir := filepath.Join(workdir, "gen-py")
	out, _ := runCLI(t, workdir, "gen", "--py", "--schema", schemaPath, "--out", outDir)
	assertPathNotEmpty(t, outDir)
	logPathContents(t, outDir)
	assertFileContains(t, filepath.Join(outDir, "models.py"), "class "+tableNames[0])
	assertFileContains(t, filepath.Join(outDir, "tables.py"), tableNames[0]+" = \""+tableNames[0]+"\"")
	assertFileContains(t, filepath.Join(outDir, "schema.py"), "SCHEMA_JSON")
	assertOutputContains(t, out, "Python client ->")
}

func TestGenJava(t *testing.T) {
	workdir := prepareE2EWorkspace(t)
	schemaPath := ensureSchema(t, workdir)
	tableNames := loadTableNames(t, schemaPath)
	outDir := filepath.Join(workdir, "gen-java")
	out, _ := runCLI(t, workdir, "gen", "--java", "--schema", schemaPath, "--out", outDir, "--package", "com.example.onyx")
	assertPathNotEmpty(t, outDir)
	logPathContents(t, outDir)
	for _, name := range tableNames {
		javaFile := filepath.Join(outDir, "com", "example", "onyx", name+".java")
		assertFileContains(t, javaFile, "public class "+name)
	}
	assertOutputContains(t, out, "Java classes ->")
}

func TestGenKotlin(t *testing.T) {
	workdir := prepareE2EWorkspace(t)
	schemaPath := ensureSchema(t, workdir)
	tableNames := loadTableNames(t, schemaPath)
	outDir := filepath.Join(workdir, "gen-kt")
	out, _ := runCLI(t, workdir, "gen", "--kotlin", "--schema", schemaPath, "--out", outDir, "--package", "com.example.onyx")
	assertPathNotEmpty(t, outDir)
	logPathContents(t, outDir)
	assertFileContains(t, filepath.Join(outDir, "Onyx.kt"), "data class "+tableNames[0], "object Tables")
	assertOutputContains(t, out, "Kotlin data classes ->")
}

func prepareE2EWorkspace(t *testing.T) string {
	if os.Getenv("ONYX_E2E") == "" {
		t.Skip("set ONYX_E2E=1 to run integration test")
	}

	repoRoot := moduleRoot(t)

	// Require local config file so we don't bake secrets into the test.
	rootCfg := filepath.Join(repoRoot, "onyx-database.json")
	if _, err := os.Stat(rootCfg); err != nil {
		t.Fatalf("missing %s in repo root: %v", rootCfg, err)
	}

	workdir := filepath.Join(repoRoot, ".e2e-workdir")
	// Keep artifacts by default; set ONYX_E2E_CLEAN=1 to wipe before each run.
	if os.Getenv("ONYX_E2E_CLEAN") == "1" {
		_ = os.RemoveAll(workdir)
	}
	if err := os.MkdirAll(workdir, 0o755); err != nil {
		t.Fatalf("create workdir: %v", err)
	}

	// Copy config into temp workspace and point ONYX_CONFIG_PATH at it.
	cfgBytes, err := os.ReadFile(rootCfg)
	if err != nil {
		t.Fatalf("read %s: %v", rootCfg, err)
	}
	cfgPath := filepath.Join(workdir, "onyx-database.json")
	if err := os.WriteFile(cfgPath, cfgBytes, 0o600); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	t.Setenv("ONYX_CONFIG_PATH", cfgPath)

	return workdir
}

func runCLI(t *testing.T, workdir string, args ...string) (string, string) {
	t.Helper()
	cmd := newRootCmd()
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs(args)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	cwd, _ := os.Getwd()
	if err := os.Chdir(workdir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer os.Chdir(cwd)

	argStr := strings.Join(args, " ")
	t.Logf("RUN %s", argStr)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("%s failed: %v\nstdout:\n%s\nstderr:\n%s", argStr, err, stdout.String(), stderr.String())
	}

	outStr := strings.TrimSpace(stdout.String())
	errStr := strings.TrimSpace(stderr.String())

	if outStr == "" {
		t.Logf("PASS %s\nstdout: <empty>", argStr)
	} else {
		t.Logf("PASS %s\nstdout:\n%s", argStr, stdout.String())
	}
	if errStr != "" {
		t.Logf("stderr:\n%s", stderr.String())
	}
	return stdout.String(), stderr.String()
}

func assertSchemaExists(t *testing.T, workdir string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(workdir, "api", "onyx.schema.json")); err != nil {
		t.Fatalf("schema file missing after get: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(workdir, "api", "onyx.schema.json"))
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	var parsed map[string]any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Fatalf("schema is not valid JSON: %v", err)
	}
	tables, ok := parsed["tables"].([]any)
	if !ok || len(tables) == 0 {
		t.Fatalf("schema missing non-empty \"tables\" array")
	}
	first, _ := tables[0].(map[string]any)
	if first == nil || strings.TrimSpace(fmt.Sprint(first["name"])) == "" {
		t.Fatalf("schema tables[0] missing name")
	}
	if _, hasAttrs := first["attributes"]; !hasAttrs {
		t.Fatalf("schema tables[0] missing attributes")
	}
}

func ensureSchema(t *testing.T, workdir string) string {
	t.Helper()
	schemaPath := filepath.Join(workdir, "api", "onyx.schema.json")
	runCLI(t, workdir, "schema", "get", "--out", schemaPath)
	assertSchemaExists(t, workdir)
	return schemaPath
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if info, err := os.Stat(path); err != nil || info.IsDir() {
		t.Fatalf("expected file %s to exist", path)
	}
}

func assertPathNotEmpty(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("expected path %s to exist: %v", path, err)
	}
	if !info.IsDir() {
		if info.Size() == 0 {
			t.Fatalf("expected file %s to be non-empty", path)
		}
		return
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		t.Fatalf("read dir %s: %v", path, err)
	}
	if len(entries) == 0 {
		t.Fatalf("expected directory %s to contain generated files", path)
	}
}

func logPathContents(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Logf("logPathContents: %s (missing: %v)", path, err)
		return
	}
	if !info.IsDir() {
		t.Logf("ARTIFACT %s (%d bytes)", path, info.Size())
		return
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		t.Logf("logPathContents: %s (read err: %v)", path, err)
		return
	}
	var lines []string
	for _, e := range entries {
		if e.IsDir() {
			lines = append(lines, e.Name()+"/")
		} else {
			if fi, err := e.Info(); err == nil {
				lines = append(lines, fmt.Sprintf("%s (%d bytes)", e.Name(), fi.Size()))
			} else {
				lines = append(lines, e.Name())
			}
		}
	}
	t.Logf("ARTIFACTS in %s:\n  %s", path, strings.Join(lines, "\n  "))
}

func loadTableNames(t *testing.T, schemaPath string) []string {
	t.Helper()
	raw, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read schema for table names: %v", err)
	}
	var parsed map[string]any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Fatalf("parse schema for table names: %v", err)
	}
	tables, ok := parsed["tables"].([]any)
	if !ok || len(tables) == 0 {
		if ents, ok2 := parsed["entities"].([]any); ok2 {
			tables = ents
		}
	}
	if len(tables) == 0 {
		t.Fatalf("no tables/entities in schema for name extraction")
	}
	var names []string
	for _, tval := range tables {
		if m, ok := tval.(map[string]any); ok {
			if n := strings.TrimSpace(fmt.Sprint(m["name"])); n != "" {
				names = append(names, n)
			}
		}
	}
	if len(names) == 0 {
		t.Fatalf("no table names parsed from schema")
	}
	return names
}

func assertFileContains(t *testing.T, path string, substrings ...string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	txt := string(data)
	for _, sub := range substrings {
		if !strings.Contains(txt, sub) {
			t.Fatalf("expected %s to contain %q", path, sub)
		}
	}
}

func assertOutputContains(t *testing.T, output string, substrings ...string) {
	t.Helper()
	for _, sub := range substrings {
		if !strings.Contains(output, sub) {
			t.Fatalf("expected output to contain %q\noutput:\n%s", sub, output)
		}
	}
}

func assertOutputContainsAny(t *testing.T, output string, substrings ...string) {
	t.Helper()
	for _, sub := range substrings {
		if strings.Contains(output, sub) {
			return
		}
	}
	t.Fatalf("expected output to contain any of %v\noutput:\n%s", substrings, output)
}

func runCmd(t *testing.T, dir string, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	t.Logf("RUN %s %s\nstdout/stderr:\n%s", name, strings.Join(args, " "), string(out))
	if err != nil {
		t.Fatalf("command failed: %s %v\noutput:\n%s", name, args, string(out))
	}
	return string(out)
}

// exportName mirrors the generator casing for Go filenames.
func exportName(name string) string {
	if name == "" {
		return ""
	}
	runes := []rune(name)
	runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
	return string(runes)
}

// moduleRoot returns the repo root by walking up from this file location.
func moduleRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("cannot determine caller path")
	}
	// file = .../cmd/onyx/integration_test.go -> repo root is two levels up
	dir := filepath.Dir(file)
	for i := 0; i < 2; i++ {
		dir = filepath.Dir(dir)
	}
	return dir
}
