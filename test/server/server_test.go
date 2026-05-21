package server_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mcp-middleware/config"
	"mcp-middleware/server"
)

func newTestCfg() *config.Config {
	return &config.Config{
		AppMode:               "http",
		AppHost:               "localhost",
		AppPort:               "8080",
		MCPServerURL:          "https://mcp.middleware.io/mcp",
		AuthServerURL:         "https://app.middleware.io",
		TenantBaseURLTemplate: "https://{alias}.middleware.io",
	}
}

func TestNewServer(t *testing.T) {
	srv := server.New(newTestCfg())
	if srv == nil {
		t.Fatal("New() returned nil")
	}
	if srv.GetMCPServer() == nil {
		t.Fatal("GetMCPServer() returned nil")
	}
}

func TestProtectedResourceMetadataHandler(t *testing.T) {
	cfg := newTestCfg()
	h := server.ProtectedResourceMetadataHandler(cfg)

	req := httptest.NewRequest(http.MethodGet, "/.well-known/oauth-protected-resource/mcp", nil)
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}

	var doc map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &doc); err != nil {
		t.Fatalf("body is not JSON: %v", err)
	}
	if doc["resource"] != cfg.MCPServerURL {
		t.Errorf("resource = %v, want %s", doc["resource"], cfg.MCPServerURL)
	}
	authServers, _ := doc["authorization_servers"].([]any)
	if len(authServers) != 1 || authServers[0] != cfg.AuthServerURL {
		t.Errorf("authorization_servers = %v", doc["authorization_servers"])
	}
	bearerMethods, _ := doc["bearer_methods_supported"].([]any)
	if len(bearerMethods) != 1 || bearerMethods[0] != "header" {
		t.Errorf("bearer_methods_supported = %v", doc["bearer_methods_supported"])
	}
}

func TestRequireBearer_MissingHeader_Returns401(t *testing.T) {
	cfg := newTestCfg()
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })
	h := server.RequireBearer(cfg, next)

	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(`{}`))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
	if called {
		t.Fatal("next handler should not have been invoked")
	}
	wa := w.Header().Get("WWW-Authenticate")
	if !strings.HasPrefix(wa, "Bearer ") {
		t.Fatalf("WWW-Authenticate = %q, want Bearer scheme", wa)
	}
	if !strings.Contains(wa, "resource_metadata=") {
		t.Errorf("WWW-Authenticate missing resource_metadata: %q", wa)
	}
	if !strings.Contains(wa, "mcp.middleware.io/.well-known/oauth-protected-resource/mcp") {
		t.Errorf("WWW-Authenticate resource_metadata wrong: %q", wa)
	}
}

func TestRequireBearer_NonJWT_Returns401InvalidToken(t *testing.T) {
	cfg := newTestCfg()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	h := server.RequireBearer(cfg, next)

	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(`{}`))
	req.Header.Set("Authorization", "Bearer not-a-jwt")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
	if !strings.Contains(w.Header().Get("WWW-Authenticate"), `error="invalid_token"`) {
		t.Errorf("expected invalid_token error in WWW-Authenticate: %q", w.Header().Get("WWW-Authenticate"))
	}
}

func TestRequireBearer_ValidJWT_PutsTokenAndTenantInContext(t *testing.T) {
	cfg := newTestCfg()
	token := makeJWT(t, map[string]any{
		"alias":      "github-ai",
		"account_id": 2,
	})

	var gotToken, gotTenant string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotToken, _ = server.TokenFromContext(r.Context())
		gotTenant, _ = server.TenantBaseURLFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})
	h := server.RequireBearer(cfg, next)

	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(`{}`))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		body, _ := io.ReadAll(w.Body)
		t.Fatalf("status = %d, want 200; body=%s", w.Code, string(body))
	}
	if gotToken != token {
		t.Errorf("token in ctx = %q, want %q", gotToken, token)
	}
	// Derived from the alias claim via the default template.
	if gotTenant != "https://github-ai.middleware.io" {
		t.Errorf("tenant in ctx = %q, want https://github-ai.middleware.io", gotTenant)
	}
}

func TestRequireBearer_MissingAliasClaim_Returns401(t *testing.T) {
	cfg := newTestCfg()
	token := makeJWT(t, map[string]any{"sub": "user-1"})
	h := server.RequireBearer(cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(`{}`))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
}

func TestRequireBearer_RejectsAliasInjection(t *testing.T) {
	cfg := newTestCfg()
	// An alias that tries to smuggle a different host/path must be rejected before
	// it can be substituted into the URL template.
	for _, bad := range []string{"evil.com/", "foo/../bar", "a.b", "has space", "-leadinghyphen"} {
		token := makeJWT(t, map[string]any{"alias": bad})
		h := server.RequireBearer(cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(`{}`))
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("alias %q: status = %d, want 401 (injection guard)", bad, w.Code)
		}
	}
}

func TestRequireBearer_ParsesExcludeToolsQueryParam(t *testing.T) {
	cfg := newTestCfg()
	token := makeJWT(t, map[string]any{"alias": "acme"})

	var got map[string]struct{}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = server.ExcludedToolsFromContext(r.Context())
	})
	h := server.RequireBearer(cfg, next)

	req := httptest.NewRequest(http.MethodPost, "/mcp?exclude_tools=create_dashboard,%20delete_widget", strings.NewReader(`{}`))
	req.Header.Set("Authorization", "Bearer "+token)
	h.ServeHTTP(httptest.NewRecorder(), req)

	if _, ok := got["create_dashboard"]; !ok {
		t.Errorf("create_dashboard not in excluded set: %v", got)
	}
	if _, ok := got["delete_widget"]; !ok {
		t.Errorf("delete_widget not in excluded set: %v", got)
	}
}

func TestClient_BuiltFromContext(t *testing.T) {
	srv := server.New(newTestCfg())
	ctx := context.Background()
	ctx = server.WithToken(ctx, "tok")
	ctx = server.WithTenantBaseURL(ctx, "https://acme.middleware.io")
	c := srv.Client(ctx)
	if c == nil {
		t.Fatal("Client(ctx) returned nil")
	}
}

// makeJWT builds an unsigned-style JWT (header.payload.sig) where the payload carries the given
// claims. The signature is bogus — we don't verify it.
func makeJWT(t *testing.T, claims map[string]any) string {
	t.Helper()
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		t.Fatalf("marshal claims: %v", err)
	}
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	return header + "." + payload + ".sig"
}
