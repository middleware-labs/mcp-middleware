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
	client := middleware.NewClient("https://test.middleware.io", "tok")
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
}

func TestRequestSendsBearerHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("Authorization = %q, want %q", got, "Bearer test-token")
		}
		if got := r.Header.Get("ApiKey"); got != "" {
			t.Errorf("ApiKey header should not be set, got %q", got)
		}
		json.NewEncoder(w).Encode(middleware.ReportListResponse{Total: 0})
	}))
	defer server.Close()

	client := middleware.NewClient(server.URL, "test-token")
	if _, err := client.GetDashboards(context.Background(), nil); err != nil {
		t.Fatalf("GetDashboards() error = %v", err)
	}
}

func TestRequestFailsWithoutToken(t *testing.T) {
	client := middleware.NewClient("https://test.middleware.io", "")
	if _, err := client.GetDashboards(context.Background(), nil); err == nil {
		t.Fatal("expected error when token is empty")
	}
}

func TestRequestFailsWithoutBaseURL(t *testing.T) {
	client := middleware.NewClient("", "tok")
	if _, err := client.GetDashboards(context.Background(), nil); err == nil {
		t.Fatal("expected error when baseURL is empty")
	}
}

func TestGetDashboards(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/builder/report" {
			t.Errorf("Expected path /api/v1/builder/report, got %s", r.URL.Path)
		}

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

	client := middleware.NewClient(server.URL, "test-token")
	result, err := client.GetDashboards(context.Background(), nil)
	if err != nil {
		t.Fatalf("GetDashboards() error = %v", err)
	}
	if result.Total != 1 {
		t.Errorf("Expected Total = 1, got %d", result.Total)
	}
	if len(result.Reports) != 1 || result.Reports[0].Label != "Test Dashboard" {
		t.Errorf("unexpected reports: %+v", result.Reports)
	}
}

func TestGetDashboardsWithParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Get("limit") != "5" {
			t.Errorf("Expected limit=5, got %s", query.Get("limit"))
		}
		if query.Get("search") != "production" {
			t.Errorf("Expected search=production, got %s", query.Get("search"))
		}
		json.NewEncoder(w).Encode(middleware.ReportListResponse{Limit: 5})
	}))
	defer server.Close()

	client := middleware.NewClient(server.URL, "test-token")
	_, err := client.GetDashboards(context.Background(), &middleware.GetDashboardsParams{
		Limit:  5,
		Search: "production",
	})
	if err != nil {
		t.Fatalf("GetDashboards() error = %v", err)
	}
}

func TestCreateDashboard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		var req middleware.UpsertReportRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}
		if req.Label != "New Dashboard" {
			t.Errorf("Expected label 'New Dashboard', got '%s'", req.Label)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(middleware.Report{ID: 100, Label: req.Label, Visibility: req.Visibility})
	}))
	defer server.Close()

	client := middleware.NewClient(server.URL, "test-token")
	result, err := client.CreateDashboard(context.Background(), &middleware.UpsertReportRequest{
		Label:      "New Dashboard",
		Visibility: "private",
	})
	if err != nil {
		t.Fatalf("CreateDashboard() error = %v", err)
	}
	if result.ID != 100 || result.Label != "New Dashboard" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestDeleteDashboard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/builder/report/123" {
			t.Errorf("Expected path /api/v1/builder/report/123, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := middleware.NewClient(server.URL, "test-token")
	if err := client.DeleteDashboard(context.Background(), 123); err != nil {
		t.Fatalf("DeleteDashboard() error = %v", err)
	}
}

func TestGetResources(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/builder/resources" {
			t.Errorf("Expected path /api/v1/builder/resources, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]string{"host", "process", "container", "pod"})
	}))
	defer server.Close()

	client := middleware.NewClient(server.URL, "test-token")
	result, err := client.GetResources(context.Background())
	if err != nil {
		t.Fatalf("GetResources() error = %v", err)
	}
	if len(result) != 4 {
		t.Errorf("Expected 4 resources, got %d", len(result))
	}
}

func TestAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(middleware.ErrorResponse{
			Error:   "Invalid token",
			Success: false,
		})
	}))
	defer server.Close()

	client := middleware.NewClient(server.URL, "invalid-token")
	_, err := client.GetDashboards(context.Background(), nil)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	expected := "API error (401): Invalid token"
	if err != nil && err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}
}

func TestContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := middleware.NewClient(server.URL, "tok")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if _, err := client.GetDashboards(ctx, nil); err == nil {
		t.Error("Expected timeout error")
	}
}

func TestGetMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		json.NewEncoder(w).Encode(middleware.MetricsV2Response{
			Items: []map[string]any{{"name": "cpu.usage"}, {"name": "memory.usage"}},
			Page:  1, Limit: 100,
		})
	}))
	defer server.Close()

	client := middleware.NewClient(server.URL, "tok")
	result, err := client.GetMetrics(context.Background(), &middleware.MetricsV2Request{
		DataType:   "metrics",
		WidgetType: "timeseries",
		KpiType:    1,
	})
	if err != nil {
		t.Fatalf("GetMetrics() error = %v", err)
	}
	if len(result.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result.Items))
	}
}

func TestGetAlerts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/rules/456/alerts" {
			t.Errorf("Expected path /api/v1/rules/456/alerts, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(middleware.AlertsResponse{
			Alerts:       []middleware.ViewModelAlert{{ID: 1, Title: "High CPU", Status: 1}},
			LatestStatus: 1,
		})
	}))
	defer server.Close()

	client := middleware.NewClient(server.URL, "tok")
	result, err := client.GetAlerts(context.Background(), 456, nil)
	if err != nil {
		t.Fatalf("GetAlerts() error = %v", err)
	}
	if len(result.Alerts) != 1 || result.Alerts[0].Title != "High CPU" {
		t.Errorf("unexpected alerts: %+v", result.Alerts)
	}
}

func TestCreateAlert(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req middleware.NewAlert
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}
		if req.Title != "Test Alert" {
			t.Errorf("Expected title 'Test Alert', got '%s'", req.Title)
		}
		json.NewEncoder(w).Encode(middleware.Alert{ID: 789, Title: req.Title, Status: req.Status})
	}))
	defer server.Close()

	client := middleware.NewClient(server.URL, "tok")
	result, err := client.CreateAlert(context.Background(), 456, &middleware.NewAlert{
		Title:  "Test Alert",
		Status: 1,
	})
	if err != nil {
		t.Fatalf("CreateAlert() error = %v", err)
	}
	if result.ID != 789 {
		t.Errorf("Expected ID = 789, got %d", result.ID)
	}
}
