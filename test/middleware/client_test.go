package middleware_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"mcp-middleware/middleware"
)

func TestNewClient(t *testing.T) {
	baseURL := "https://test.middleware.io"
	apiKey := "test-api-key"

	client := middleware.NewClient(baseURL, apiKey)
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
}

func TestNewClientWithAuthorization(t *testing.T) {
	baseURL := "https://test.middleware.io"
	authHeader := "Bearer test-token"

	client := middleware.NewClientWithAuth(baseURL, "", authHeader)
	if client == nil {
		t.Fatal("NewClientWithAuth() returned nil")
	}
}

func TestGetDashboards(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/builder/report" {
			t.Errorf("Expected path /api/v1/builder/report, got %s", r.URL.Path)
		}

		// Check API key header
		if apiKey := r.Header.Get("ApiKey"); apiKey != "test-key" {
			t.Errorf("Expected ApiKey header 'test-key', got '%s'", apiKey)
		}

		// Return mock response
		response := middleware.ReportListResponse{
			Reports: []middleware.Report{
				{ID: 1, Label: "Test Dashboard", Visibility: "public"},
			},
			Total:  1,
			Limit:  10,
			Offset: 0,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := middleware.NewClient(server.URL, "test-key")
	ctx := context.Background()

	result, err := client.GetDashboards(ctx, nil)
	if err != nil {
		t.Fatalf("GetDashboards() error = %v", err)
	}

	if result.Total != 1 {
		t.Errorf("Expected Total = 1, got %d", result.Total)
	}

	if len(result.Reports) != 1 {
		t.Errorf("Expected 1 report, got %d", len(result.Reports))
	}

	if result.Reports[0].Label != "Test Dashboard" {
		t.Errorf("Expected label 'Test Dashboard', got '%s'", result.Reports[0].Label)
	}
}

func TestGetDashboardsWithAuthorizationHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-auth" {
			t.Errorf("Expected Authorization header 'Bearer test-auth', got '%s'", auth)
		}
		if apiKey := r.Header.Get("ApiKey"); apiKey != "" {
			t.Errorf("Expected ApiKey header to be empty when using Authorization, got '%s'", apiKey)
		}

		response := middleware.ReportListResponse{
			Reports: []middleware.Report{},
			Total:   0,
			Limit:   10,
			Offset:  0,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := middleware.NewClientWithAuth(server.URL, "", "Bearer test-auth")
	ctx := context.Background()

	if _, err := client.GetDashboards(ctx, nil); err != nil {
		t.Fatalf("GetDashboards() error = %v", err)
	}
}

func TestGetDashboardsPrefersAuthorizationOverApiKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth := r.Header.Get("Authorization"); auth != "Bearer preferred-auth" {
			t.Errorf("Expected Authorization header 'Bearer preferred-auth', got '%s'", auth)
		}
		if apiKey := r.Header.Get("ApiKey"); apiKey != "" {
			t.Errorf("Expected ApiKey header to be empty when Authorization is provided, got '%s'", apiKey)
		}

		response := middleware.ReportListResponse{
			Reports: []middleware.Report{},
			Total:   0,
			Limit:   10,
			Offset:  0,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := middleware.NewClientWithAuth(server.URL, "should-not-send", "Bearer preferred-auth")
	ctx := context.Background()

	if _, err := client.GetDashboards(ctx, nil); err != nil {
		t.Fatalf("GetDashboards() error = %v", err)
	}
}

func TestGetDashboardsWithParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		query := r.URL.Query()
		if query.Get("limit") != "5" {
			t.Errorf("Expected limit=5, got %s", query.Get("limit"))
		}

		if query.Get("search") != "production" {
			t.Errorf("Expected search=production, got %s", query.Get("search"))
		}

		response := middleware.ReportListResponse{
			Reports: []middleware.Report{},
			Total:   0,
			Limit:   5,
			Offset:  0,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := middleware.NewClient(server.URL, "test-key")
	ctx := context.Background()

	params := &middleware.GetDashboardsParams{
		Limit:  5,
		Search: "production",
	}

	_, err := client.GetDashboards(ctx, params)
	if err != nil {
		t.Fatalf("GetDashboards() error = %v", err)
	}
}

func TestCreateDashboard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/builder/report" {
			t.Errorf("Expected path /api/v1/builder/report, got %s", r.URL.Path)
		}

		// Decode request body
		var req middleware.UpsertReportRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if req.Label != "New Dashboard" {
			t.Errorf("Expected label 'New Dashboard', got '%s'", req.Label)
		}

		// Return created dashboard
		response := middleware.Report{
			ID:         100,
			Label:      req.Label,
			Visibility: req.Visibility,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := middleware.NewClient(server.URL, "test-key")
	ctx := context.Background()

	req := &middleware.UpsertReportRequest{
		Label:      "New Dashboard",
		Visibility: "private",
	}

	result, err := client.CreateDashboard(ctx, req)
	if err != nil {
		t.Fatalf("CreateDashboard() error = %v", err)
	}

	if result.ID != 100 {
		t.Errorf("Expected ID = 100, got %d", result.ID)
	}

	if result.Label != "New Dashboard" {
		t.Errorf("Expected label 'New Dashboard', got '%s'", result.Label)
	}
}

func TestDeleteDashboard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/builder/report/123" {
			t.Errorf("Expected path /api/v1/builder/report/123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := middleware.NewClient(server.URL, "test-key")
	ctx := context.Background()

	err := client.DeleteDashboard(ctx, 123)
	if err != nil {
		t.Fatalf("DeleteDashboard() error = %v", err)
	}
}

func TestGetResources(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/builder/resources" {
			t.Errorf("Expected path /api/v1/builder/resources, got %s", r.URL.Path)
		}

		resources := []string{"host", "process", "container", "pod"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resources)
	}))
	defer server.Close()

	client := middleware.NewClient(server.URL, "test-key")
	ctx := context.Background()

	result, err := client.GetResources(ctx)
	if err != nil {
		t.Fatalf("GetResources() error = %v", err)
	}

	if len(result) != 4 {
		t.Errorf("Expected 4 resources, got %d", len(result))
	}

	expectedResources := map[string]bool{
		"host": true, "process": true, "container": true, "pod": true,
	}

	for _, r := range result {
		if !expectedResources[r] {
			t.Errorf("Unexpected resource: %s", r)
		}
	}
}

func TestAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(middleware.ErrorResponse{
			Error:   "Invalid API key",
			Success: false,
		})
	}))
	defer server.Close()

	client := middleware.NewClient(server.URL, "invalid-key")
	ctx := context.Background()

	_, err := client.GetDashboards(ctx, nil)
	if err == nil {
		t.Error("Expected error for invalid API key, got nil")
	}

	expectedMsg := "API error (401): Invalid API key"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestContextCancellation(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := middleware.NewClient(server.URL, "test-key")

	// Create context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := client.GetDashboards(ctx, nil)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestGetMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/builder/metrics-v2" {
			t.Errorf("Expected path /api/v1/builder/metrics-v2, got %s", r.URL.Path)
		}

		var req middleware.MetricsV2Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if req.DataType != "metrics" {
			t.Errorf("Expected DataType 'metrics', got '%s'", req.DataType)
		}

		response := middleware.MetricsV2Response{
			Items: []map[string]any{
				{"name": "cpu.usage", "type": "gauge"},
				{"name": "memory.usage", "type": "gauge"},
			},
			Page:  1,
			Limit: 100,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := middleware.NewClient(server.URL, "test-key")
	ctx := context.Background()

	req := &middleware.MetricsV2Request{
		DataType:   "metrics",
		WidgetType: "timeseries",
		KpiType:    1,
	}

	result, err := client.GetMetrics(ctx, req)
	if err != nil {
		t.Fatalf("GetMetrics() error = %v", err)
	}

	if len(result.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result.Items))
	}
}

func TestGetAlerts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/api/v1/rules/456/alerts" {
			t.Errorf("Expected path /api/v1/rules/456/alerts, got %s", r.URL.Path)
		}

		response := middleware.AlertsResponse{
			Alerts: []middleware.ViewModelAlert{
				{ID: 1, Title: "High CPU", Status: 1},
			},
			LatestStatus: 1,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := middleware.NewClient(server.URL, "test-key")
	ctx := context.Background()

	result, err := client.GetAlerts(ctx, 456, nil)
	if err != nil {
		t.Fatalf("GetAlerts() error = %v", err)
	}

	if len(result.Alerts) != 1 {
		t.Errorf("Expected 1 alert, got %d", len(result.Alerts))
	}

	if result.Alerts[0].Title != "High CPU" {
		t.Errorf("Expected title 'High CPU', got '%s'", result.Alerts[0].Title)
	}
}

func TestCreateAlert(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var req middleware.NewAlert
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		if req.Title != "Test Alert" {
			t.Errorf("Expected title 'Test Alert', got '%s'", req.Title)
		}

		response := middleware.Alert{
			ID:     789,
			Title:  req.Title,
			Status: req.Status,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := middleware.NewClient(server.URL, "test-key")
	ctx := context.Background()

	req := &middleware.NewAlert{
		Title:  "Test Alert",
		Status: 1,
	}

	result, err := client.CreateAlert(ctx, 456, req)
	if err != nil {
		t.Fatalf("CreateAlert() error = %v", err)
	}

	if result.ID != 789 {
		t.Errorf("Expected ID = 789, got %d", result.ID)
	}
}
