package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	checkInterval    = 24 * time.Hour
	stalenessTimeout = 2 * time.Second
)

// stalenessCache persists the last update check result outside the repo.
type stalenessCache struct {
	LastChecked   string `json:"last_checked"`
	LatestVersion string `json:"latest_version"`
}

// cacheHome can be overridden in tests.
var cacheHome = ""

func cacheFilePath() string {
	if cacheHome != "" {
		return filepath.Join(cacheHome, "cache.json")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".nav-pilot", "cache.json")
}

func readCache() *stalenessCache {
	path := cacheFilePath()
	if path == "" {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var c stalenessCache
	if err := json.Unmarshal(data, &c); err != nil {
		return nil
	}
	return &c
}

func writeCache(c *stalenessCache) {
	path := cacheFilePath()
	if path == "" {
		return
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	data, _ := json.MarshalIndent(c, "", "  ")
	data = append(data, '\n')
	os.WriteFile(path, data, 0o644)
}

// checkStaleness returns the latest available version if the installed
// collection is outdated. Returns "" if up-to-date, check was skipped
// (within cooldown), or any error occurred (network, API, etc).
// Designed to be fast and never block — uses a 2s HTTP timeout.
func checkStaleness(installedVersion string) string {
	if installedVersion == "" || installedVersion == "dev" {
		return ""
	}

	// Check cooldown
	cache := readCache()
	if cache != nil && cache.LastChecked != "" {
		if t, err := time.Parse(time.RFC3339, cache.LastChecked); err == nil {
			if time.Since(t) < checkInterval {
				// Within cooldown — use cached result
				if cache.LatestVersion != "" && versionNewer(cache.LatestVersion, installedVersion) {
					return cache.LatestVersion
				}
				return ""
			}
		}
	}

	// Use a short timeout client for staleness checks
	client := &http.Client{Timeout: stalenessTimeout}
	origClient := httpClient
	httpClient = client
	defer func() { httpClient = origClient }()

	latest, _, err := fetchLatestVersion()
	if err != nil {
		return ""
	}

	// Only write cache on successful check
	writeCache(&stalenessCache{
		LastChecked:   time.Now().UTC().Format(time.RFC3339),
		LatestVersion: latest,
	})

	if versionNewer(latest, installedVersion) {
		return latest
	}
	return ""
}
