package server_test

import (
	"sort"
	"strings"
	"testing"

	"mcp-middleware/server"
)

// expectedTools is the full set of tools registered in registerTools(). The catalog must
// cover exactly these — this test fails if a tool is added/removed without updating both.
var expectedTools = []string{
	"list_dashboards", "get_dashboard", "create_dashboard", "update_dashboard",
	"delete_dashboard", "clone_dashboard", "set_dashboard_favorite",
	"list_widgets", "create_widget", "update_widget", "delete_widget",
	"get_widget_data", "get_multi_widget_data", "update_widget_layouts",
	"get_metrics", "get_resources", "query",
	"list_alerts", "create_alert", "get_alert_stats",
	"list_errors", "get_error_details",
}

func TestToolCatalog_MatchesRegisteredTools(t *testing.T) {
	seen := map[string]bool{}
	for _, cat := range server.ToolCatalog() {
		if cat.Name == "" {
			t.Error("category with empty name")
		}
		for _, tool := range cat.Tools {
			if tool.Name == "" || tool.Description == "" {
				t.Errorf("tool with empty name/description in %q: %+v", cat.Name, tool)
			}
			if tool.Access != server.AccessRead && tool.Access != server.AccessWrite {
				t.Errorf("tool %q has invalid access %q", tool.Name, tool.Access)
			}
			if seen[tool.Name] {
				t.Errorf("tool %q listed more than once", tool.Name)
			}
			seen[tool.Name] = true
		}
	}

	got := make([]string, 0, len(seen))
	for name := range seen {
		got = append(got, name)
	}
	want := append([]string(nil), expectedTools...)
	sort.Strings(got)
	sort.Strings(want)

	if len(got) != len(want) {
		t.Fatalf("catalog has %d tools, registered set has %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("catalog tools differ from registered set:\n got=%v\nwant=%v", got, want)
		}
	}
}

func TestCatalogScopes(t *testing.T) {
	if len(server.ToolCatalog()) == 0 {
		t.Fatal("empty catalog")
	}

	got := server.CatalogScopes()
	want := []string{
		"dashboards:read", "dashboards:write",
		"widgets:read", "widgets:write",
		"metrics:read",
		"alerts:read", "alerts:write",
		"error_tracking:read",
	}
	if len(got) != len(want) {
		t.Fatalf("CatalogScopes() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("CatalogScopes()[%d] = %q, want %q (full: %v)", i, got[i], want[i], got)
		}
	}

	if s := server.CatalogScopeString(); s != strings.Join(want, " ") {
		t.Errorf("CatalogScopeString() = %q", s)
	}
}

func TestCatalogScopeCategories(t *testing.T) {
	cats := server.CatalogScopeCategories()
	if len(cats) != 5 {
		t.Fatalf("got %d scope categories, want 5", len(cats))
	}

	total := 0
	for _, c := range cats {
		if c.Category == "" || len(c.Scopes) == 0 || len(c.Tools) == 0 {
			t.Errorf("category %q is missing name/scopes/tools", c.Category)
		}
		slug := strings.ToLower(strings.ReplaceAll(c.Category, " ", "_"))
		for _, tool := range c.Tools {
			total++
			if tool.Name == "" || tool.Description == "" {
				t.Errorf("%s tool missing name/description: %+v", c.Category, tool)
			}
			// Each tool's scope must be "<category>:<access>".
			want := slug + ":" + string(tool.Access)
			if tool.Scope != want {
				t.Errorf("tool %q scope = %q, want %q", tool.Name, tool.Scope, want)
			}
		}
	}
	if total != len(expectedTools) {
		t.Errorf("scope categories cover %d tools, want %d", total, len(expectedTools))
	}
}
