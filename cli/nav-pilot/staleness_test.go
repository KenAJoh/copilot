package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"
)

func TestCheckStaleness_SkipsDevVersion(t *testing.T) {
	setupTestCache(t)
	result := checkStaleness("dev")
	if result != "" {
		t.Errorf("expected empty for dev version, got %q", result)
	}
}

func TestCheckStaleness_SkipsEmptyVersion(t *testing.T) {
	setupTestCache(t)
	result := checkStaleness("")
	if result != "" {
		t.Errorf("expected empty for empty version, got %q", result)
	}
}

func TestCheckStaleness_DetectsUpdate(t *testing.T) {
	setupTestCache(t)

	// Mock GitHub API
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]ghRelease{
			{TagName: "nav-pilot/2026.04.13-abc1234"},
		})
	}))
	defer srv.Close()

	origAPI := releasesAPI
	releasesAPI = srv.URL
	defer func() { releasesAPI = origAPI }()

	result := checkStaleness("2026.01.01-old1234")
	if result != "2026.04.13-abc1234" {
		t.Errorf("expected update version, got %q", result)
	}
}

func TestCheckStaleness_UpToDate(t *testing.T) {
	setupTestCache(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]ghRelease{
			{TagName: "nav-pilot/2026.04.13-abc1234"},
		})
	}))
	defer srv.Close()

	origAPI := releasesAPI
	releasesAPI = srv.URL
	defer func() { releasesAPI = origAPI }()

	result := checkStaleness("2026.04.13-abc1234")
	if result != "" {
		t.Errorf("expected empty for up-to-date version, got %q", result)
	}
}

func TestCheckStaleness_UsesCachedResult(t *testing.T) {
	setupTestCache(t)

	// Write a recent cache entry
	writeCache(&stalenessCache{
		LastChecked:   time.Now().UTC().Format(time.RFC3339),
		LatestVersion: "2026.05.01-new1234",
	})

	// Server should NOT be hit (cache is fresh)
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		json.NewEncoder(w).Encode([]ghRelease{
			{TagName: "nav-pilot/2026.05.01-new1234"},
		})
	}))
	defer srv.Close()

	origAPI := releasesAPI
	releasesAPI = srv.URL
	defer func() { releasesAPI = origAPI }()

	result := checkStaleness("2026.01.01-old1234")
	if called {
		t.Error("expected cache hit, but server was called")
	}
	if result != "2026.05.01-new1234" {
		t.Errorf("expected cached version, got %q", result)
	}
}

func TestCheckStaleness_ExpiredCacheRefetches(t *testing.T) {
	setupTestCache(t)

	// Write an expired cache entry
	writeCache(&stalenessCache{
		LastChecked:   time.Now().Add(-25 * time.Hour).UTC().Format(time.RFC3339),
		LatestVersion: "2026.03.01-old",
	})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]ghRelease{
			{TagName: "nav-pilot/2026.05.01-new1234"},
		})
	}))
	defer srv.Close()

	origAPI := releasesAPI
	releasesAPI = srv.URL
	defer func() { releasesAPI = origAPI }()

	result := checkStaleness("2026.01.01-old1234")
	if result != "2026.05.01-new1234" {
		t.Errorf("expected new version from API, got %q", result)
	}
}

func TestCheckStaleness_NetworkErrorSkips(t *testing.T) {
	setupTestCache(t)

	origAPI := releasesAPI
	releasesAPI = "http://127.0.0.1:1" // connection refused
	defer func() { releasesAPI = origAPI }()

	result := checkStaleness("2026.01.01-old1234")
	if result != "" {
		t.Errorf("expected empty on network error, got %q", result)
	}
}

func TestCacheFilePath(t *testing.T) {
	path := cacheFilePath()
	if path == "" {
		t.Skip("no home directory available")
	}
	if !filepath.IsAbs(path) {
		t.Errorf("expected absolute path, got %q", path)
	}
	if filepath.Base(path) != "cache.json" {
		t.Errorf("expected cache.json, got %q", filepath.Base(path))
	}
}

// setupTestCache sets cacheHome to a temp dir and restores it after the test.
func setupTestCache(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	origHome := cacheHome
	cacheHome = dir
	t.Cleanup(func() { cacheHome = origHome })
}
