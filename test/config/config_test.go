package config_test

import (
	"os"
	"strings"
	"testing"

	"mcp-middleware/config"
)

// clearEnv unsets every env var the config package reads, so each test starts from a clean slate.
func clearEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{
		"APP_MODE", "APP_HOST", "APP_PORT",
		"MCP_SERVER_URL", "MW_AUTH_SERVER_URL", "MCP_SCOPES",
		"MW_TENANT_BASE_URL_TEMPLATE",
	} {
		os.Unsetenv(k)
	}
}

func TestLoad_Defaults(t *testing.T) {
	clearEnv(t)
	defer clearEnv(t)
	// http is the default mode, which requires the auth server URL.
	os.Setenv("MW_AUTH_SERVER_URL", "https://app.middleware.io")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if cfg.AppMode != "http" {
		t.Errorf("AppMode = %q, want http", cfg.AppMode)
	}
	if cfg.AppHost != "localhost" {
		t.Errorf("AppHost = %q, want localhost", cfg.AppHost)
	}
	if cfg.AppPort != "8080" {
		t.Errorf("AppPort = %q, want 8080", cfg.AppPort)
	}
	if cfg.MCPScopes != "mcp:read mcp:tools" {
		t.Errorf("MCPScopes = %q, want default", cfg.MCPScopes)
	}
}

func TestLoad_HTTPRequiresAuthURL(t *testing.T) {
	clearEnv(t)
	defer clearEnv(t)
	os.Setenv("APP_MODE", "http")

	// The MCP endpoint URL is derived (not an env var), so the only required URL is MW_AUTH_SERVER_URL.
	if _, err := config.Load(); err == nil {
		t.Fatal("expected error when MW_AUTH_SERVER_URL missing in http mode")
	} else if !strings.Contains(err.Error(), "MW_AUTH_SERVER_URL") {
		t.Errorf("error %q does not mention MW_AUTH_SERVER_URL", err)
	}

	os.Setenv("MW_AUTH_SERVER_URL", "https://app.middleware.io")
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if cfg.AuthServerURL != "https://app.middleware.io" {
		t.Errorf("AuthServerURL = %q", cfg.AuthServerURL)
	}
}

func TestLoad_DerivesMCPServerURL(t *testing.T) {
	clearEnv(t)
	defer clearEnv(t)
	os.Setenv("APP_MODE", "http")
	os.Setenv("MW_AUTH_SERVER_URL", "https://app.middleware.io")

	// Defaults: localhost:8080 → derived endpoint URL.
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if cfg.MCPServerURL != "http://localhost:8080/mcp" {
		t.Errorf("derived MCPServerURL = %q, want http://localhost:8080/mcp", cfg.MCPServerURL)
	}

	// Derivation follows the listen address.
	os.Setenv("APP_HOST", "0.0.0.0")
	os.Setenv("APP_PORT", "3000")
	cfg, err = config.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if cfg.MCPServerURL != "http://0.0.0.0:3000/mcp" {
		t.Errorf("derived MCPServerURL = %q, want http://0.0.0.0:3000/mcp", cfg.MCPServerURL)
	}

	// Explicit MCP_SERVER_URL (e.g. public ingress URL) overrides the derived value.
	os.Setenv("MCP_SERVER_URL", "https://mcp.middleware.io/mcp")
	cfg, err = config.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if cfg.MCPServerURL != "https://mcp.middleware.io/mcp" {
		t.Errorf("MCPServerURL = %q, want explicit https://mcp.middleware.io/mcp", cfg.MCPServerURL)
	}
}

func TestLoad_TrimTrailingSlash(t *testing.T) {
	clearEnv(t)
	defer clearEnv(t)
	os.Setenv("APP_MODE", "http")
	os.Setenv("MW_AUTH_SERVER_URL", "https://app.middleware.io/")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if cfg.AuthServerURL != "https://app.middleware.io" {
		t.Errorf("trailing slash not trimmed: %q", cfg.AuthServerURL)
	}
}

func TestLoad_TenantBaseURLTemplate(t *testing.T) {
	clearEnv(t)
	defer clearEnv(t)
	// stdio mode skips the http-only validation; we only care about the template here.
	os.Setenv("APP_MODE", "stdio")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if cfg.TenantBaseURLTemplate != "https://{alias}.middleware.io" {
		t.Errorf("default template = %q, want https://{alias}.middleware.io", cfg.TenantBaseURLTemplate)
	}

	os.Setenv("MW_TENANT_BASE_URL_TEMPLATE", "http://localhost:9100")
	cfg, err = config.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if cfg.TenantBaseURLTemplate != "http://localhost:9100" {
		t.Errorf("template override = %q, want http://localhost:9100", cfg.TenantBaseURLTemplate)
	}
}

func TestLoad_InvalidAppMode(t *testing.T) {
	clearEnv(t)
	defer clearEnv(t)
	os.Setenv("APP_MODE", "grpc")

	if _, err := config.Load(); err == nil {
		t.Fatal("expected error for invalid APP_MODE")
	}
}

func TestLoad_OverrideHostPortAndScopes(t *testing.T) {
	clearEnv(t)
	defer clearEnv(t)
	os.Setenv("APP_MODE", "http")
	os.Setenv("APP_HOST", "0.0.0.0")
	os.Setenv("APP_PORT", "9090")
	os.Setenv("MCP_SCOPES", "mcp:admin")
	os.Setenv("MW_AUTH_SERVER_URL", "https://app.middleware.io")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if cfg.AppHost != "0.0.0.0" || cfg.AppPort != "9090" {
		t.Errorf("host/port not picked up: %q %q", cfg.AppHost, cfg.AppPort)
	}
	if cfg.MCPScopes != "mcp:admin" {
		t.Errorf("scopes override not honored: %q", cfg.MCPScopes)
	}
}
