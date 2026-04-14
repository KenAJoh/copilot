package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

const syncConfigPath = ".github/copilot-sync.json"

// SyncConfig holds optional per-repo sync configuration.
// Teams create .github/copilot-sync.json to customize sync behavior.
type SyncConfig struct {
	// Overrides lists files that the team maintains locally.
	// These files are skipped during sync — no hash comparison, no PR diff.
	Overrides []string `json:"overrides,omitempty"`
}

// readSyncConfig reads .github/copilot-sync.json from the given directory.
// Returns nil (no error) if the file does not exist.
func readSyncConfig(dir string) (*SyncConfig, error) {
	path := filepath.Join(dir, syncConfigPath)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var cfg SyncConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// overrideSet builds a lookup set from the config's overrides list.
// Paths are canonicalized with filepath.Clean and forward slashes for
// consistent matching against syncFile.localPath.
func overrideSet(cfg *SyncConfig) map[string]bool {
	if cfg == nil || len(cfg.Overrides) == 0 {
		return nil
	}
	m := make(map[string]bool, len(cfg.Overrides))
	for _, p := range cfg.Overrides {
		clean := filepath.ToSlash(filepath.Clean(p))
		m[clean] = true
		// Also match with trailing slash for directory overrides
		if !strings.HasSuffix(clean, "/") {
			m[clean+"/"] = true
		}
	}
	return m
}
