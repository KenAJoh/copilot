package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestDownloadAndParseNDJSON(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		status      int
		wantRecords int
		wantErr     bool
	}{
		{
			name:        "single record",
			body:        `{"day":"2025-10-11","daily_active_users":30}`,
			status:      http.StatusOK,
			wantRecords: 1,
		},
		{
			name:        "multiple records",
			body:        "{\"a\":1}\n{\"b\":2}\n{\"c\":3}\n",
			status:      http.StatusOK,
			wantRecords: 3,
		},
		{
			name:        "blank lines skipped",
			body:        "{\"a\":1}\n\n{\"b\":2}\n\n",
			status:      http.StatusOK,
			wantRecords: 2,
		},
		{
			name:        "invalid JSON lines skipped",
			body:        "{\"a\":1}\nnot json\n{\"b\":2}\n",
			status:      http.StatusOK,
			wantRecords: 2,
		},
		{
			name:        "empty body",
			body:        "",
			status:      http.StatusOK,
			wantRecords: 0,
		},
		{
			name:    "non-200 status",
			body:    "forbidden",
			status:  http.StatusForbidden,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") != "" {
					t.Error("download request should NOT have Authorization header")
				}
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			})

			client := &GitHubClient{
				httpClient:     &http.Client{},
				downloadClient: mockClient(handler),
			}

			records, err := client.downloadAndParseNDJSON(context.Background(), "http://test/download")
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(records) != tt.wantRecords {
				t.Errorf("got %d records, want %d", len(records), tt.wantRecords)
			}
		})
	}
}

func TestDownloadAndParseNDJSON_LargeRecord(t *testing.T) {
	largeValue := strings.Repeat("x", 500_000)
	body := `{"data":"` + largeValue + `"}` + "\n"

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	})

	client := &GitHubClient{
		downloadClient: mockClient(handler),
	}

	records, err := client.downloadAndParseNDJSON(context.Background(), "http://test/download")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 1 {
		t.Errorf("got %d records, want 1", len(records))
	}
}

func TestDownloadUsesDownloadClient(t *testing.T) {
	authHandler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("httpClient (auth) should NOT be used for downloads")
	})

	downloadHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}` + "\n"))
	})

	client := &GitHubClient{
		httpClient:     mockClient(authHandler),
		downloadClient: mockClient(downloadHandler),
	}

	records, err := client.downloadAndParseNDJSON(context.Background(), "http://test/download")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 1 {
		t.Errorf("got %d records, want 1", len(records))
	}
}

func TestFetchMetricsFromURL_StringResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`"No report available for this day"`))
	})

	client := &GitHubClient{
		httpClient:     mockClient(handler),
		downloadClient: &http.Client{},
	}

	_, err := client.fetchMetricsFromURL(context.Background(), "http://test/report")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrReportNotAvailable) {
		t.Errorf("expected ErrReportNotAvailable, got: %v", err)
	}
}

func TestFetchMetricsFromURL_ValidResponse(t *testing.T) {
	downloadHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{\"users\":42}\n"))
	})

	reportResp := MetricsReportResponse{
		DownloadLinks: []string{"http://test/ndjson"},
		ReportDay:     "2026-03-20",
	}
	reportJSON, _ := json.Marshal(reportResp)

	reportHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Accept"); got != "application/vnd.github+json" {
			t.Errorf("Accept header = %q, want %q", got, "application/vnd.github+json")
		}
		if got := r.Header.Get("X-GitHub-Api-Version"); got != "2026-03-10" {
			t.Errorf("X-GitHub-Api-Version = %q, want %q", got, "2026-03-10")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(reportJSON)
	})

	client := &GitHubClient{
		httpClient:     mockClient(reportHandler),
		downloadClient: mockClient(downloadHandler),
	}

	records, err := client.fetchMetricsFromURL(context.Background(), "http://test/report")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 1 {
		t.Errorf("got %d records, want 1", len(records))
	}
}

func TestFetchMetricsFromURLWithRetry_NoRetryOnStringResponse(t *testing.T) {
	attempts := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`"Report is being generated"`))
	})

	client := &GitHubClient{
		httpClient:     mockClient(handler),
		downloadClient: &http.Client{},
	}

	_, err := client.fetchMetricsFromURLWithRetry(context.Background(), "http://test/report")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt (no retries), got %d", attempts)
	}
}
