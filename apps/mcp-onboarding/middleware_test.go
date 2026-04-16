package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthMiddleware_NoAuthHeader(t *testing.T) {
	store := newTestStore()
	mw := NewAuthMiddleware(store)

	handler := mw.Authenticate(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/mcp", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	if w.Header().Get("WWW-Authenticate") == "" {
		t.Error("expected WWW-Authenticate header")
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	store := newTestStore()
	mw := NewAuthMiddleware(store)

	handler := mw.Authenticate(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/mcp", nil)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNz") // not Bearer
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	store := newTestStore()
	mw := NewAuthMiddleware(store)

	handler := mw.Authenticate(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/mcp", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	store := newTestStore()
	store.SaveToken("expired-token", &TokenData{
		GitHubAccessToken: "gh-token",
		UserLogin:         "testuser",
		UserID:            12345,
		ExpiresAt:         time.Now().Add(-1 * time.Hour),
	})
	mw := NewAuthMiddleware(store)

	handler := mw.Authenticate(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("handler should not be called for expired token")
	}))

	req := httptest.NewRequest("GET", "/mcp", nil)
	req.Header.Set("Authorization", "Bearer expired-token")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	store := newTestStore()
	store.SaveToken("valid-token", &TokenData{
		GitHubAccessToken: "gh-access-token",
		UserLogin:         "octocat",
		UserID:            42,
		ExpiresAt:         time.Now().Add(1 * time.Hour),
	})
	mw := NewAuthMiddleware(store)

	var gotUser *UserContext
	handler := mw.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUser = GetUserFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/mcp", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if gotUser == nil {
		t.Fatal("expected user context to be set")
	}
	if gotUser.Login != "octocat" {
		t.Errorf("expected Login 'octocat', got %q", gotUser.Login)
	}
	if gotUser.ID != 42 {
		t.Errorf("expected ID 42, got %d", gotUser.ID)
	}
	if gotUser.GitHubAccessToken != "gh-access-token" {
		t.Errorf("expected GitHubAccessToken 'gh-access-token', got %q", gotUser.GitHubAccessToken)
	}
}

func TestAuthMiddleware_UnauthorizedResponse(t *testing.T) {
	store := newTestStore()
	mw := NewAuthMiddleware(store)

	handler := mw.Authenticate(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))

	req := httptest.NewRequest("GET", "/mcp", nil)
	req.Host = "example.com"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	var body map[string]string
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["error"] != "unauthorized" {
		t.Errorf("expected error 'unauthorized', got %q", body["error"])
	}
}

func TestGetBaseURL_HTTP(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Host = "example.com"

	url := getBaseURL(req)
	if url != "http://example.com" {
		t.Errorf("expected 'http://example.com', got %q", url)
	}
}

func TestGetBaseURL_XForwardedProto(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Host = "example.com"
	req.Header.Set("X-Forwarded-Proto", "https")

	url := getBaseURL(req)
	if url != "https://example.com" {
		t.Errorf("expected 'https://example.com', got %q", url)
	}
}

// contextKey is needed to test GetUserFromContext
func TestGetUserFromContext_NoUser(t *testing.T) {
	ctx := context.Background()
	user := GetUserFromContext(ctx)
	if user == nil {
		t.Fatal("expected non-nil user from empty context (should return anonymous)")
	}
	if user.Login != "anonymous" {
		t.Errorf("expected Login 'anonymous', got %q", user.Login)
	}
}
