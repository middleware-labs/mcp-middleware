package server

import "strings"

// ToolAccess classifies a tool as read-only or mutating, so a consent screen can show
// "N read / M write" per category (Datadog-style).
type ToolAccess string

const (
	AccessRead  ToolAccess = "read"
	AccessWrite ToolAccess = "write"
)

// CatalogTool describes a single MCP tool for display on an authorization/consent screen.
type CatalogTool struct {
	Name        string     `json:"name"`
	Access      ToolAccess `json:"access"`
	Description string     `json:"description"`
}

// ToolCategory groups related tools under a human-readable heading.
type ToolCategory struct {
	Name  string        `json:"name"`
	Tools []CatalogTool `json:"tools"`
}

// toolCatalog is the source of truth for the tool list shown to users during authorization.
// Keep it in sync with registerTools(): TestToolCatalog_MatchesRegisteredTools guards drift.
var toolCatalog = []ToolCategory{
	{
		Name: "Dashboards",
		Tools: []CatalogTool{
			{"list_dashboards", AccessRead, "List dashboards / reports"},
			{"get_dashboard", AccessRead, "Get a dashboard by key"},
			{"create_dashboard", AccessWrite, "Create a dashboard"},
			{"update_dashboard", AccessWrite, "Update a dashboard"},
			{"delete_dashboard", AccessWrite, "Delete a dashboard"},
			{"clone_dashboard", AccessWrite, "Clone a dashboard"},
			{"set_dashboard_favorite", AccessWrite, "Mark a dashboard as favorite"},
		},
	},
	{
		Name: "Widgets",
		Tools: []CatalogTool{
			{"list_widgets", AccessRead, "List widgets"},
			{"get_widget_data", AccessRead, "Fetch data for a widget"},
			{"get_multi_widget_data", AccessRead, "Fetch data for multiple widgets"},
			{"create_widget", AccessWrite, "Create a widget"},
			{"update_widget", AccessWrite, "Update a widget"},
			{"delete_widget", AccessWrite, "Delete a widget"},
			{"update_widget_layouts", AccessWrite, "Update widget layout positions"},
		},
	},
	{
		Name: "Metrics",
		Tools: []CatalogTool{
			{"get_metrics", AccessRead, "List available metrics"},
			{"get_resources", AccessRead, "List resource types"},
			{"query", AccessRead, "Run a data query"},
		},
	},
	{
		Name: "Alerts",
		Tools: []CatalogTool{
			{"list_alerts", AccessRead, "List alerts for a rule"},
			{"get_alert_stats", AccessRead, "Get alert statistics"},
			{"create_alert", AccessWrite, "Create an alert"},
		},
	},
	{
		Name: "Error Tracking",
		Tools: []CatalogTool{
			{"list_errors", AccessRead, "List error-tracking incidents"},
			{"get_error_details", AccessRead, "Get details for an error incident"},
		},
	},
}

// ToolCatalog returns the categorized tool catalog (source of truth).
func ToolCatalog() []ToolCategory { return toolCatalog }

// categorySlug turns a display name into a scope-friendly token, e.g. "Error Tracking" → "error_tracking".
func categorySlug(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "_"))
}

// CatalogScopes derives OAuth scopes from the tool catalog: a "<category>:read" and/or
// "<category>:write" scope per category, depending on the access levels of its tools.
// Order is deterministic (category order; read before write). This keeps the advertised
// scopes in lockstep with the tools the server actually exposes.
func CatalogScopes() []string {
	var scopes []string
	for _, cat := range toolCatalog {
		var hasRead, hasWrite bool
		for _, t := range cat.Tools {
			if t.Access == AccessWrite {
				hasWrite = true
			} else {
				hasRead = true
			}
		}
		slug := categorySlug(cat.Name)
		if hasRead {
			scopes = append(scopes, slug+":read")
		}
		if hasWrite {
			scopes = append(scopes, slug+":write")
		}
	}
	return scopes
}

// CatalogScopeString returns the catalog scopes as a single space-separated string,
// for the WWW-Authenticate "scope" parameter.
func CatalogScopeString() string {
	return strings.Join(CatalogScopes(), " ")
}

// ScopeTool pairs a tool with the OAuth scope that governs it, plus a description, for
// rendering on a consent screen.
type ScopeTool struct {
	Name        string     `json:"name"`
	Scope       string     `json:"scope"`
	Access      ToolAccess `json:"access"`
	Description string     `json:"description"`
}

// ScopeCategory groups a category's scopes and the tools each scope covers.
type ScopeCategory struct {
	Category string      `json:"category"`
	Scopes   []string    `json:"scopes"`
	Tools    []ScopeTool `json:"tools"`
}

// CatalogScopeCategories returns the catalog grouped for a consent screen: each category,
// the scopes it introduces, and every tool with the scope that governs it and a description.
// This is the structured companion to the flat `scopes_supported` string list.
func CatalogScopeCategories() []ScopeCategory {
	out := make([]ScopeCategory, 0, len(toolCatalog))
	for _, cat := range toolCatalog {
		slug := categorySlug(cat.Name)
		sc := ScopeCategory{Category: cat.Name}
		var hasRead, hasWrite bool
		for _, t := range cat.Tools {
			sc.Tools = append(sc.Tools, ScopeTool{
				Name:        t.Name,
				Scope:       slug + ":" + string(t.Access),
				Access:      t.Access,
				Description: t.Description,
			})
			if t.Access == AccessWrite {
				hasWrite = true
			} else {
				hasRead = true
			}
		}
		if hasRead {
			sc.Scopes = append(sc.Scopes, slug+":read")
		}
		if hasWrite {
			sc.Scopes = append(sc.Scopes, slug+":write")
		}
		out = append(out, sc)
	}
	return out
}
