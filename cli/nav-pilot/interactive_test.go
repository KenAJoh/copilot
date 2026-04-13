package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestFindCopilotCLI(t *testing.T) {
	path, name := findCopilotCLI()
	if path != "" {
		if name != "cplt" {
			t.Errorf("expected name 'cplt', got %q", name)
		}
		if _, err := os.Stat(path); err != nil {
			t.Errorf("returned path %q does not exist: %v", path, err)
		}
	} else {
		if name != "" {
			t.Errorf("expected empty name when path is empty, got %q", name)
		}
	}
}

func TestCmdInteractive_NotGitRepo(t *testing.T) {
	origDir, _ := os.Getwd()
	dir := t.TempDir()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	err := cmdInteractive()
	if err == nil {
		t.Fatal("expected error for non-git directory")
	}
	if err.Error() != "not in a git repository" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCmdInteractive_InstalledUpToDate(t *testing.T) {
	origDir, _ := os.Getwd()
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, ".git"), 0o755)
	os.MkdirAll(filepath.Join(dir, ".github"), 0o755)
	os.Chdir(dir)
	defer os.Chdir(origDir)

	state := &StateFile{
		Collection: "test-collection",
		Version:    "2026.04.13.17.00.00-abc1234",
		SourceSHA:  "abc1234",
	}
	if err := writeState(dir, state); err != nil {
		t.Fatal(err)
	}

	// Mock GitHub API returning same version (up-to-date)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]ghRelease{
			{TagName: "nav-pilot/2026.04.13.17.00.00-abc1234"},
		})
	}))
	defer srv.Close()

	origAPI := releasesAPI
	releasesAPI = srv.URL
	defer func() { releasesAPI = origAPI }()

	setupTestCache(t)

	// Should not error — will try to launch cplt (which may not exist, that's ok)
	err := cmdInteractive()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
