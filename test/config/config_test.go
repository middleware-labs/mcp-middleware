package config_test

import (
	"os"
	"testing"

	"mcp-middleware/config"
)

func TestLoad(t *testing.T) {
	// Set required environment variables
	os.Setenv("MIDDLEWARE_API_KEY", "test-api-key")
	os.Setenv("MIDDLEWARE_BASE_URL", "https://test.middleware.io")
	defer func() {
		os.Unsetenv("MIDDLEWARE_API_KEY")
		os.Unsetenv("MIDDLEWARE_BASE_URL")
		os.Unsetenv("APP_MODE")
		os.Unsetenv("APP_HOST")
		os.Unsetenv("APP_PORT")
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.MiddlewareAPIKey != "test-api-key" {
		t.Errorf("Expected API key 'test-api-key', got '%s'", cfg.MiddlewareAPIKey)
	}

	if cfg.MiddlewareBaseURL != "https://test.middleware.io" {
		t.Errorf("Expected base URL 'https://test.middleware.io', got '%s'", cfg.MiddlewareBaseURL)
	}

	if cfg.AppMode != "stdio" {
		t.Errorf("Expected default AppMode 'stdio', got '%s'", cfg.AppMode)
	}

	if cfg.AppHost != "localhost" {
		t.Errorf("Expected default AppHost 'localhost', got '%s'", cfg.AppHost)
	}

	if cfg.AppPort != "8080" {
		t.Errorf("Expected default AppPort '8080', got '%s'", cfg.AppPort)
	}
}

func TestLoadMissingAPIKey(t *testing.T) {
	os.Unsetenv("MIDDLEWARE_API_KEY")
	os.Setenv("MIDDLEWARE_BASE_URL", "https://test.middleware.io")
	defer os.Unsetenv("MIDDLEWARE_BASE_URL")

	_, err := config.Load()
	if err == nil {
		t.Error("Expected error when API key is missing, got nil")
	}

	expectedMsg := "MIDDLEWARE_API_KEY is required"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestLoadMissingBaseURL(t *testing.T) {
	os.Setenv("MIDDLEWARE_API_KEY", "test-api-key")
	os.Unsetenv("MIDDLEWARE_BASE_URL")
	defer os.Unsetenv("MIDDLEWARE_API_KEY")

	_, err := config.Load()
	if err == nil {
		t.Error("Expected error when base URL is missing, got nil")
	}

	expectedMsg := "MIDDLEWARE_BASE_URL is required"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestIsToolExcluded(t *testing.T) {
	os.Setenv("MIDDLEWARE_API_KEY", "test-api-key")
	os.Setenv("MIDDLEWARE_BASE_URL", "https://test.middleware.io")
	os.Setenv("EXCLUDED_TOOLS", "delete_dashboard,delete_widget,create_alert")
	defer func() {
		os.Unsetenv("MIDDLEWARE_API_KEY")
		os.Unsetenv("MIDDLEWARE_BASE_URL")
		os.Unsetenv("EXCLUDED_TOOLS")
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	tests := []struct {
		name     string
		toolName string
		want     bool
	}{
		{"delete_dashboard excluded", "delete_dashboard", true},
		{"delete_widget excluded", "delete_widget", true},
		{"create_alert excluded", "create_alert", true},
		{"list_dashboards not excluded", "list_dashboards", false},
		{"get_metrics not excluded", "get_metrics", false},
		{"unknown tool not excluded", "unknown_tool", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cfg.IsToolExcluded(tt.toolName)
			if got != tt.want {
				t.Errorf("IsToolExcluded(%s) = %v, want %v", tt.toolName, got, tt.want)
			}
		})
	}
}

func TestInvalidAppMode(t *testing.T) {
	os.Setenv("MIDDLEWARE_API_KEY", "test-api-key")
	os.Setenv("MIDDLEWARE_BASE_URL", "https://test.middleware.io")
	os.Setenv("APP_MODE", "invalid")
	defer func() {
		os.Unsetenv("MIDDLEWARE_API_KEY")
		os.Unsetenv("MIDDLEWARE_BASE_URL")
		os.Unsetenv("APP_MODE")
	}()

	_, err := config.Load()
	if err == nil {
		t.Error("Expected error for invalid APP_MODE, got nil")
	}
}

func TestCustomAppModes(t *testing.T) {
	tests := []struct {
		name     string
		mode     string
		wantErr  bool
	}{
		{"stdio mode", "stdio", false},
		{"http mode", "http", false},
		{"sse mode", "sse", false},
		{"invalid mode", "grpc", true},
		{"empty mode defaults to stdio", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("MIDDLEWARE_API_KEY", "test-api-key")
			os.Setenv("MIDDLEWARE_BASE_URL", "https://test.middleware.io")
			if tt.mode != "" {
				os.Setenv("APP_MODE", tt.mode)
			} else {
				os.Unsetenv("APP_MODE")
			}
			defer func() {
				os.Unsetenv("MIDDLEWARE_API_KEY")
				os.Unsetenv("MIDDLEWARE_BASE_URL")
				os.Unsetenv("APP_MODE")
			}()

			cfg, err := config.Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				expectedMode := tt.mode
				if expectedMode == "" {
					expectedMode = "stdio"
				}
				if cfg.AppMode != expectedMode {
					t.Errorf("Expected AppMode %s, got %s", expectedMode, cfg.AppMode)
				}
			}
		})
	}
}

func TestCustomHostAndPort(t *testing.T) {
	os.Setenv("MIDDLEWARE_API_KEY", "test-api-key")
	os.Setenv("MIDDLEWARE_BASE_URL", "https://test.middleware.io")
	os.Setenv("APP_HOST", "0.0.0.0")
	os.Setenv("APP_PORT", "9090")
	defer func() {
		os.Unsetenv("MIDDLEWARE_API_KEY")
		os.Unsetenv("MIDDLEWARE_BASE_URL")
		os.Unsetenv("APP_HOST")
		os.Unsetenv("APP_PORT")
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.AppHost != "0.0.0.0" {
		t.Errorf("Expected AppHost '0.0.0.0', got '%s'", cfg.AppHost)
	}

	if cfg.AppPort != "9090" {
		t.Errorf("Expected AppPort '9090', got '%s'", cfg.AppPort)
	}
}

func TestExcludedToolsWithSpaces(t *testing.T) {
	os.Setenv("MIDDLEWARE_API_KEY", "test-api-key")
	os.Setenv("MIDDLEWARE_BASE_URL", "https://test.middleware.io")
	os.Setenv("EXCLUDED_TOOLS", " delete_dashboard , create_widget , delete_widget ")
	defer func() {
		os.Unsetenv("MIDDLEWARE_API_KEY")
		os.Unsetenv("MIDDLEWARE_BASE_URL")
		os.Unsetenv("EXCLUDED_TOOLS")
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if !cfg.IsToolExcluded("delete_dashboard") {
		t.Error("Expected delete_dashboard to be excluded (with space trimming)")
	}

	if !cfg.IsToolExcluded("create_widget") {
		t.Error("Expected create_widget to be excluded (with space trimming)")
	}

	if !cfg.IsToolExcluded("delete_widget") {
		t.Error("Expected delete_widget to be excluded (with space trimming)")
	}
}

func TestEmptyExcludedTools(t *testing.T) {
	os.Setenv("MIDDLEWARE_API_KEY", "test-api-key")
	os.Setenv("MIDDLEWARE_BASE_URL", "https://test.middleware.io")
	os.Setenv("EXCLUDED_TOOLS", "")
	defer func() {
		os.Unsetenv("MIDDLEWARE_API_KEY")
		os.Unsetenv("MIDDLEWARE_BASE_URL")
		os.Unsetenv("EXCLUDED_TOOLS")
	}()

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if len(cfg.ExcludedTools) != 0 {
		t.Errorf("Expected no excluded tools, got %d", len(cfg.ExcludedTools))
	}
}
