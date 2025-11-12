package server_test

import (
	"testing"

	"mcp-middleware/config"
	"mcp-middleware/server"
)

func TestNewServer(t *testing.T) {
	cfg := &config.Config{
		MiddlewareAPIKey:  "test-key",
		MiddlewareBaseURL: "https://test.middleware.io",
		AppMode:           "stdio",
		AppHost:           "localhost",
		AppPort:           "8080",
		ExcludedTools:     make(map[string]bool),
	}

	srv := server.New(cfg)
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestNewServerWithExcludedTools(t *testing.T) {
	cfg := &config.Config{
		MiddlewareAPIKey:  "test-key",
		MiddlewareBaseURL: "https://test.middleware.io",
		AppMode:           "stdio",
		AppHost:           "localhost",
		AppPort:           "8080",
		ExcludedTools: map[string]bool{
			"delete_dashboard": true,
			"delete_widget":    true,
			"create_alert":     true,
		},
	}

	srv := server.New(cfg)
	if srv == nil {
		t.Fatal("New() returned nil with excluded tools")
	}

	// Server should initialize successfully with excluded tools
	// The actual tool exclusion is tested at runtime
}

func TestServerConfiguration(t *testing.T) {
	tests := []struct {
		name string
		cfg  *config.Config
	}{
		{
			name: "stdio mode",
			cfg: &config.Config{
				MiddlewareAPIKey:  "test-key",
				MiddlewareBaseURL: "https://test.middleware.io",
				AppMode:           "stdio",
				ExcludedTools:     make(map[string]bool),
			},
		},
		{
			name: "http mode",
			cfg: &config.Config{
				MiddlewareAPIKey:  "test-key",
				MiddlewareBaseURL: "https://test.middleware.io",
				AppMode:           "http",
				AppHost:           "0.0.0.0",
				AppPort:           "9090",
				ExcludedTools:     make(map[string]bool),
			},
		},
		{
			name: "with multiple excluded tools",
			cfg: &config.Config{
				MiddlewareAPIKey:  "test-key",
				MiddlewareBaseURL: "https://test.middleware.io",
				AppMode:           "stdio",
				ExcludedTools: map[string]bool{
					"delete_dashboard":     true,
					"delete_widget":        true,
					"create_alert":         true,
					"update_dashboard":     true,
					"create_widget":        true,
					"update_widget_layouts": true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := server.New(tt.cfg)
			if srv == nil {
				t.Errorf("%s: New() returned nil", tt.name)
			}
		})
	}
}

