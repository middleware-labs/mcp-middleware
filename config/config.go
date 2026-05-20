package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// Application Mode: stdio, http, sse
	AppMode string

	// Server Configuration (for http/sse modes)
	AppHost string
	AppPort string

	// Public URI of the MCP endpoint, advertised in the RFC 9728 Protected Resource
	// Metadata document. From MCP_SERVER_URL if set (use behind a proxy/ingress),
	// otherwise derived from the listen address as http://<AppHost>:<AppPort>/mcp.
	MCPServerURL string

	// OAuth Authorization Server issuer (e.g. https://app.middleware.io).
	// Advertised in the Protected Resource Metadata; Claude builds the AS metadata URL from it.
	AuthServerURL string

	// Space-separated scope list (e.g. "mcp:read mcp:tools").
	MCPScopes string

	// Template for the tenant data-plane base URL, with "{alias}" substituted by the
	// access token's `alias` claim. Default "https://{alias}.middleware.io".
	// For local testing, set to a fixed URL with no placeholder (e.g. "http://localhost:9100")
	// to route every tenant at a mock backend.
	TenantBaseURLTemplate string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppMode:               getEnvOrDefault("APP_MODE", "http"),
		AppHost:               getEnvOrDefault("APP_HOST", "localhost"),
		AppPort:               getEnvOrDefault("APP_PORT", "8080"),
		MCPServerURL:          strings.TrimRight(os.Getenv("MCP_SERVER_URL"), "/"),
		AuthServerURL:         strings.TrimRight(os.Getenv("MW_AUTH_SERVER_URL"), "/"),
		MCPScopes:             getEnvOrDefault("MCP_SCOPES", "mcp:read mcp:tools"),
		TenantBaseURLTemplate: getEnvOrDefault("MW_TENANT_BASE_URL_TEMPLATE", "https://{alias}.middleware.io"),
	}

	// MCP endpoint URL advertised in the metadata. Set MCP_SERVER_URL to the public URL
	// when the server is behind a proxy/ingress; otherwise derive it from the listen address.
	if cfg.MCPServerURL == "" {
		cfg.MCPServerURL = fmt.Sprintf("http://%s:%s/mcp", cfg.AppHost, cfg.AppPort)
	}

	validModes := map[string]bool{"stdio": true, "http": true, "sse": true}
	if !validModes[cfg.AppMode] {
		return nil, fmt.Errorf("invalid APP_MODE: %s (must be stdio, http, or sse)", cfg.AppMode)
	}

	if cfg.AppMode != "stdio" {
		if cfg.AuthServerURL == "" {
			return nil, fmt.Errorf("MW_AUTH_SERVER_URL is required when APP_MODE=%s", cfg.AppMode)
		}
		if _, err := url.Parse(cfg.AuthServerURL); err != nil {
			return nil, fmt.Errorf("MW_AUTH_SERVER_URL is not a valid URL: %w", err)
		}
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
