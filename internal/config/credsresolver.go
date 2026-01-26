package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DefaultSchemaPath is the canonical schema file path for parity with the TypeScript CLI.
const DefaultSchemaPath = "onyx.schema.json"
const (
	DefaultBaseURL   = "https://api.onyx.dev"
	DefaultAIBaseURL = "https://ai.onyx.dev"
	DefaultModel     = "onyx"
	DefaultCodegen   = "typescript"
)

// ResolvedValue tracks a value and where it came from.
type ResolvedValue struct {
	Value  string
	Source string
}

// ResolvedConfig holds the merged configuration along with provenance.
type ResolvedConfig struct {
	DatabaseID   ResolvedValue
	BaseURL      ResolvedValue
	APIKey       ResolvedValue
	APISecret    ResolvedValue
	AIBaseURL    ResolvedValue
	DefaultModel ResolvedValue
	ConfigFile   string // path actually used (if any)
	CodegenLang  ResolvedValue
}

type envValues struct {
	ResolvedConfig
	ConfigPath  ResolvedValue
	CodegenLang ResolvedValue
}

// Options allows callers to provide explicit values/paths that outrank env/config files.
type Options struct {
	DatabaseID   string
	BaseURL      string
	APIKey       string
	APISecret    string
	AIBaseURL    string
	DefaultModel string
	ConfigPath   string // optional explicit config file path
	WorkDir      string // defaults to current working directory if empty
}

type fileConfig struct {
	DatabaseID   string `json:"databaseId"`
	BaseURL      string `json:"baseUrl"`
	APIKey       string `json:"apiKey"`
	APISecret    string `json:"apiSecret"`
	AIBaseURL    string `json:"aiBaseUrl"`
	DefaultModel string `json:"defaultModel"`
	CodegenLang  string `json:"codegenLanguage"`
}

// Resolve merges configuration following the canonical chain:
// explicit options -> environment -> ONYX_CONFIG_PATH file -> project config files -> home profile.
func Resolve(opts Options) (ResolvedConfig, error) {
	wd := opts.WorkDir
	if wd == "" {
		var err error
		wd, err = os.Getwd()
		if err != nil {
			return ResolvedConfig{}, fmt.Errorf("determine working directory: %w", err)
		}
	}

	env := loadEnv()

	// 1) Explicit options (highest)
	cfg := ResolvedConfig{
		DatabaseID:   resolveOne(opts.DatabaseID, "flag"),
		BaseURL:      resolveOne(opts.BaseURL, "flag"),
		APIKey:       resolveOne(opts.APIKey, "flag"),
		APISecret:    resolveOne(opts.APISecret, "flag"),
		AIBaseURL:    resolveOne(opts.AIBaseURL, "flag"),
		DefaultModel: resolveOne(opts.DefaultModel, "flag"),
		CodegenLang:  ResolvedValue{},
	}

	// 2) Environment variables (fill gaps)
	fillIfEmpty(&cfg.DatabaseID, env.DatabaseID.Value, env.DatabaseID.Source)
	fillIfEmpty(&cfg.BaseURL, env.BaseURL.Value, env.BaseURL.Source)
	fillIfEmpty(&cfg.APIKey, env.APIKey.Value, env.APIKey.Source)
	fillIfEmpty(&cfg.APISecret, env.APISecret.Value, env.APISecret.Source)
	fillIfEmpty(&cfg.AIBaseURL, env.AIBaseURL.Value, env.AIBaseURL.Source)
	fillIfEmpty(&cfg.DefaultModel, env.DefaultModel.Value, env.DefaultModel.Source)
	fillIfEmpty(&cfg.CodegenLang, env.CodegenLang.Value, env.CodegenLang.Source)

	// 3) Config files (explicit path -> ONYX_CONFIG_PATH -> search chain)
	configPath := firstNonEmpty(opts.ConfigPath, env.ConfigPath.Value)
	if configPath != "" {
		configPath = expandPath(configPath)
		fc, err := loadFileConfig(configPath)
		if err != nil {
			return cfg, fmt.Errorf("read config %s: %w", configPath, err)
		}
		cfg.ConfigFile = configPath
		mergeFileConfig(&cfg, fc, configPath)
		return cfg, nil
	}

	// search standard locations
	candidates := configSearchOrder(wd, cfg.DatabaseID.Value)
	for _, path := range candidates {
		fc, err := loadFileConfig(path)
		if err != nil {
			continue
		}
		cfg.ConfigFile = path
		mergeFileConfig(&cfg, fc, path)
		break
	}

	// 4) Defaults (lowest precedence)
	if cfg.BaseURL.Value == "" {
		cfg.BaseURL = ResolvedValue{Value: DefaultBaseURL, Source: "default"}
	}
	if cfg.AIBaseURL.Value == "" {
		cfg.AIBaseURL = ResolvedValue{Value: DefaultAIBaseURL, Source: "default"}
	}
	if cfg.DefaultModel.Value == "" {
		cfg.DefaultModel = ResolvedValue{Value: DefaultModel, Source: "default"}
	}
	if cfg.CodegenLang.Value == "" {
		cfg.CodegenLang = ResolvedValue{Value: DefaultCodegen, Source: "default"}
	}

	return cfg, nil
}

// loadEnv gathers env vars into ResolvedConfig form.
func loadEnv() envValues {
	return envValues{
		ResolvedConfig: ResolvedConfig{
			DatabaseID:   resolveOne(os.Getenv("ONYX_DATABASE_ID"), "env"),
			BaseURL:      resolveOne(os.Getenv("ONYX_DATABASE_BASE_URL"), "env"),
			APIKey:       resolveOne(os.Getenv("ONYX_DATABASE_API_KEY"), "env"),
			APISecret:    resolveOne(os.Getenv("ONYX_DATABASE_API_SECRET"), "env"),
			AIBaseURL:    resolveOne(os.Getenv("ONYX_AI_BASE_URL"), "env"),
			DefaultModel: resolveOne(os.Getenv("ONYX_DEFAULT_MODEL"), "env"),
			ConfigFile:   "",
		},
		ConfigPath:  resolveOne(os.Getenv("ONYX_CONFIG_PATH"), "env"),
		CodegenLang: resolveOne(os.Getenv("ONYX_CODEGEN_LANGUAGE"), "env"),
	}
}

func resolveOne(val string, source string) ResolvedValue {
	if val == "" {
		return ResolvedValue{}
	}
	return ResolvedValue{Value: val, Source: source}
}

func fillIfEmpty(target *ResolvedValue, val, source string) {
	if target.Value == "" && val != "" {
		target.Value = val
		target.Source = source
	}
}

func mergeFileConfig(cfg *ResolvedConfig, fc fileConfig, path string) {
	fillIfEmpty(&cfg.DatabaseID, fc.DatabaseID, "config:"+path)
	fillIfEmpty(&cfg.BaseURL, fc.BaseURL, "config:"+path)
	fillIfEmpty(&cfg.APIKey, fc.APIKey, "config:"+path)
	fillIfEmpty(&cfg.APISecret, fc.APISecret, "config:"+path)
	fillIfEmpty(&cfg.AIBaseURL, fc.AIBaseURL, "config:"+path)
	fillIfEmpty(&cfg.DefaultModel, fc.DefaultModel, "config:"+path)
	fillIfEmpty(&cfg.CodegenLang, fc.CodegenLang, "config:"+path)
}

func loadFileConfig(path string) (fileConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return fileConfig{}, err
	}
	var fc fileConfig
	if err := json.Unmarshal(data, &fc); err != nil {
		return fileConfig{}, fmt.Errorf("parse config %s: %w", path, err)
	}
	return fc, nil
}

// configSearchOrder builds the canonical search list for config files.
func configSearchOrder(wd, databaseID string) []string {
	homeDir, _ := os.UserHomeDir()
	var paths []string

	// helpers to add file paths relative to a directory
	add := func(dir, name string) {
		paths = append(paths, filepath.Join(dir, name))
	}

	// working directory paths
	if databaseID != "" {
		add(wd, fmt.Sprintf("onyx-database-%s.json", databaseID))
	}
	add(wd, "onyx-database.json")

	// ./config paths (after working directory files)
	configDir := filepath.Join(wd, "config")
	if databaseID != "" {
		add(configDir, fmt.Sprintf("onyx-database-%s.json", databaseID))
	}
	add(configDir, "onyx-database.json")

	// ~/.onyx paths
	if homeDir != "" {
		onyxDir := filepath.Join(homeDir, ".onyx")
		if databaseID != "" {
			add(onyxDir, fmt.Sprintf("onyx-database-%s.json", databaseID))
		}
		add(onyxDir, "onyx-database.json")

		// legacy path
		add(homeDir, "onyx-database.json")
	}

	// expand ~ in any path we constructed
	for i, p := range paths {
		paths[i] = expandPath(p)
	}
	return paths
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, strings.TrimPrefix(path, "~"))
		}
	}
	return path
}

// Validate ensures required fields are present.
func (rc ResolvedConfig) Validate() error {
	var missing []string
	if rc.DatabaseID.Value == "" {
		missing = append(missing, "databaseId")
	}
	if rc.BaseURL.Value == "" {
		missing = append(missing, "baseUrl")
	}
	if rc.APIKey.Value == "" {
		missing = append(missing, "apiKey")
	}
	if rc.APISecret.Value == "" {
		missing = append(missing, "apiSecret")
	}
	if len(missing) > 0 {
		return errors.New("missing required config: " + strings.Join(missing, ", "))
	}
	return nil
}

// MaskSecret partially masks secrets for display.
func MaskSecret(s string) string {
	if s == "" {
		return ""
	}
	if len(s) <= 4 {
		return "***"
	}
	return s[:2] + strings.Repeat("*", len(s)-4) + s[len(s)-2:]
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
