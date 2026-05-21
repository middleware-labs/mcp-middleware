package integration_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"mcp-middleware/config"
	"mcp-middleware/middleware"
	"mcp-middleware/server"
)

func clearEnv() {
	for _, k := range []string{
		"APP_MODE", "APP_HOST", "APP_PORT",
		"MCP_SERVER_URL", "MW_AUTH_SERVER_URL",
		"MW_TENANT_BASE_URL_TEMPLATE",
	} {
		os.Unsetenv(k)
	}
}

func TestFullServerInitialization(t *testing.T) {
	clearEnv()
	defer clearEnv()
	os.Setenv("APP_MODE", "http")
	os.Setenv("MW_AUTH_SERVER_URL", "https://app.middleware.io")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	srv := server.New(cfg)
	if srv == nil {
		t.Fatal("Failed to create server")
	}
}

func TestConfigToServerFlow(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
	}{
		{
			name: "stdio defaults",
			envVars: map[string]string{
				"APP_MODE": "stdio",
			},
		},
		{
			name: "http minimal",
			envVars: map[string]string{
				"APP_MODE":           "http",
				"MW_AUTH_SERVER_URL": "https://app.middleware.io",
			},
		},
		{
			name: "http missing MW_AUTH_SERVER_URL",
			envVars: map[string]string{
				"APP_MODE": "http",
			},
			wantErr: true,
		},
		{
			name: "custom host/port",
			envVars: map[string]string{
				"APP_MODE":           "http",
				"APP_HOST":           "0.0.0.0",
				"APP_PORT":           "3000",
				"MW_AUTH_SERVER_URL": "https://app.middleware.io",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv()
			defer clearEnv()
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			cfg, err := config.Load()
			if (err != nil) != tt.wantErr {
				t.Fatalf("config.Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			if srv := server.New(cfg); srv == nil {
				t.Fatal("Failed to create server")
			}
		})
	}
}

// TestEndToEndAuthFlow walks the whole HTTP handler chain:
//   - unauth request to /mcp returns 401 + WWW-Authenticate pointing at protected-resource metadata
//   - the metadata URL serves the RFC 9728 doc
//   - a request with a JWT carrying an `alias` claim passes the gate and the per-request
//     middleware.Client built from ctx talks to the tenant base URL derived from it.
func TestEndToEndAuthFlow(t *testing.T) {
	cfg := &config.Config{
		AppMode:               "http",
		AppHost:               "localhost",
		AppPort:               "0",
		MCPServerURL:          "https://mcp.middleware.io/mcp",
		AuthServerURL:         "https://app.middleware.io",
		TenantBaseURLTemplate: "https://{alias}.middleware.io",
	}

	// Spin up a fake tenant data plane.
	var seenAuth string
	tenant := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(middleware.ReportListResponse{Total: 0})
	}))
	defer tenant.Close()

	// Point the template at the fake backend (no {alias} placeholder) so the derived
	// URL routes there.
	cfg.TenantBaseURLTemplate = tenant.URL

	mux := http.NewServeMux()
	mux.Handle("/.well-known/oauth-protected-resource/mcp", server.ProtectedResourceMetadataHandler(cfg))
	mux.Handle("/mcp", server.RequireBearer(cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a tool handler: build the per-request client from ctx and hit the tenant.
		srv := server.New(cfg)
		c := srv.Client(r.Context())
		if _, err := c.GetDashboards(r.Context(), nil); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.WriteHeader(http.StatusOK)
	})))

	httpSrv := httptest.NewServer(mux)
	defer httpSrv.Close()

	// 1. Unauthenticated request → 401 with proper WWW-Authenticate.
	resp, err := http.Post(httpSrv.URL+"/mcp", "application/json", strings.NewReader(`{}`))
	if err != nil {
		t.Fatalf("unauth POST: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", resp.StatusCode)
	}
	wa := resp.Header.Get("WWW-Authenticate")
	if !strings.Contains(wa, "resource_metadata=") {
		t.Errorf("WWW-Authenticate missing resource_metadata: %q", wa)
	}
	resp.Body.Close()

	// 2. Fetch metadata.
	resp, err = http.Get(httpSrv.URL + "/.well-known/oauth-protected-resource/mcp")
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("metadata fetch failed: %v %v", err, resp.StatusCode)
	}
	resp.Body.Close()

	// 3. Authenticated request with JWT carrying the alias → tools/call hits the tenant.
	token := makeJWT(map[string]any{"alias": "acme"})
	req, _ := http.NewRequest(http.MethodPost, httpSrv.URL+"/mcp", strings.NewReader(`{}`))
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("auth POST: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		body := make([]byte, 1024)
		n, _ := resp.Body.Read(body)
		t.Fatalf("status = %d, want 200; body=%s", resp.StatusCode, string(body[:n]))
	}
	resp.Body.Close()
	if seenAuth != "Bearer "+token {
		t.Errorf("backend Authorization = %q, want %q", seenAuth, "Bearer "+token)
	}
}

func TestClientContextHandling(t *testing.T) {
	client := middleware.NewClient("https://test.middleware.io", "tok")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	if _, err := client.GetDashboards(ctx, nil); err == nil {
		t.Error("Expected error due to context timeout")
	}
}

func makeJWT(claims map[string]any) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payloadBytes, _ := json.Marshal(claims)
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	return header + "." + payload + ".sig"
}
