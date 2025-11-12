package integration_test

import (
	"context"
	"os"
	"testing"
	"time"

	"mcp-middleware/config"
	"mcp-middleware/middleware"
	"mcp-middleware/server"
)

// TestFullServerInitialization tests the complete server initialization flow
func TestFullServerInitialization(t *testing.T) {
	// Set up test environment
	os.Setenv("MIDDLEWARE_API_KEY", "test-api-key")
	os.Setenv("MIDDLEWARE_BASE_URL", "https://test.middleware.io")
	defer func() {
		os.Unsetenv("MIDDLEWARE_API_KEY")
		os.Unsetenv("MIDDLEWARE_BASE_URL")
	}()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Create server
	srv := server.New(cfg)
	if srv == nil {
		t.Fatal("Failed to create server")
	}

	// Server should be ready to run
	// (We don't actually run it to avoid blocking the test)
}

func TestClientServerIntegration(t *testing.T) {
	// This test verifies that the client and server components
	// can be initialized together without conflicts

	os.Setenv("MIDDLEWARE_API_KEY", "test-api-key")
	os.Setenv("MIDDLEWARE_BASE_URL", "https://test.middleware.io")
	defer func() {
		os.Unsetenv("MIDDLEWARE_API_KEY")
		os.Unsetenv("MIDDLEWARE_BASE_URL")
	}()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Config load failed: %v", err)
	}

	// Create client
	client := middleware.NewClient(cfg.MiddlewareBaseURL, cfg.MiddlewareAPIKey)
	if client == nil {
		t.Fatal("Client creation failed")
	}

	// Create server
	srv := server.New(cfg)
	if srv == nil {
		t.Fatal("Server creation failed")
	}

	// Both client and server should coexist
}

func TestConfigToServerFlow(t *testing.T) {
	// Test the complete flow from config to server initialization

	tests := []struct {
		name          string
		envVars       map[string]string
		wantErr       bool
		checkExcluded bool
	}{
		{
			name: "minimal config",
			envVars: map[string]string{
				"MIDDLEWARE_API_KEY":  "key123",
				"MIDDLEWARE_BASE_URL": "https://app.middleware.io",
			},
			wantErr: false,
		},
		{
			name: "with excluded tools",
			envVars: map[string]string{
				"MIDDLEWARE_API_KEY":  "key123",
				"MIDDLEWARE_BASE_URL": "https://app.middleware.io",
				"EXCLUDED_TOOLS":      "delete_dashboard,create_alert",
			},
			wantErr:       false,
			checkExcluded: true,
		},
		{
			name: "missing api key",
			envVars: map[string]string{
				"MIDDLEWARE_BASE_URL": "https://app.middleware.io",
			},
			wantErr: true,
		},
		{
			name: "custom app mode",
			envVars: map[string]string{
				"MIDDLEWARE_API_KEY":  "key123",
				"MIDDLEWARE_BASE_URL": "https://app.middleware.io",
				"APP_MODE":            "http",
				"APP_HOST":            "0.0.0.0",
				"APP_PORT":            "3000",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			}()

			// Load config
			cfg, err := config.Load()
			if (err != nil) != tt.wantErr {
				t.Fatalf("config.Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			// Create server
			srv := server.New(cfg)
			if srv == nil {
				t.Fatal("Failed to create server")
			}

			// Check excluded tools if needed
			if tt.checkExcluded {
				if !cfg.IsToolExcluded("delete_dashboard") {
					t.Error("Expected delete_dashboard to be excluded")
				}
				if !cfg.IsToolExcluded("create_alert") {
					t.Error("Expected create_alert to be excluded")
				}
			}
		})
	}
}

func TestClientContextHandling(t *testing.T) {
	client := middleware.NewClient("https://test.middleware.io", "test-key")

	// Test that client respects context cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// This should fail quickly due to context timeout
	_, err := client.GetDashboards(ctx, nil)
	if err == nil {
		t.Error("Expected error due to context timeout")
	}
}

func TestMultipleServerInstances(t *testing.T) {
	// Verify that multiple server instances can be created
	// (useful for testing isolation)

	os.Setenv("MIDDLEWARE_API_KEY", "test-key")
	os.Setenv("MIDDLEWARE_BASE_URL", "https://test.middleware.io")
	defer func() {
		os.Unsetenv("MIDDLEWARE_API_KEY")
		os.Unsetenv("MIDDLEWARE_BASE_URL")
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Config load failed: %v", err)
	}

	srv1 := server.New(cfg)
	srv2 := server.New(cfg)

	if srv1 == nil || srv2 == nil {
		t.Fatal("Failed to create multiple server instances")
	}

	// Both instances should be independent
}

