package tools

import (
	"context"

	"mcp-middleware/middleware"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ListDashboardsTool is the MCP tool definition for listing dashboards
var ListDashboardsTool = &mcp.Tool{
	Name: "list_dashboards",
	Description: `Get a list of dashboards (i.e. reports) with advanced filtering and pagination support.
	
This tool retrieves dashboards from Middleware.io with support for searching, filtering by various criteria, and pagination. Use this to discover available dashboards, find specific dashboards by name, or filter by ownership and usage patterns.`,
}

type ListDashboardsInput struct {
	Limit        int    `json:"limit,omitempty" jsonschema:"Number of items per page for pagination"`
	Offset       int    `json:"offset,omitempty" jsonschema:"Number of items to skip for pagination (page offset)"`
	Search       string `json:"search,omitempty" jsonschema:"Search query to find dashboards by name or description"`
	FilterBy     string `json:"filter_by,omitempty" jsonschema:"Comma-separated list of filter values. Valid values: custom, created_by_you, favorite, frequently_viewed, or data source names like aws, mysql, postgresql, etc."`
	DisplayScope string `json:"display_scope,omitempty" jsonschema:"Filter dashboards by comma-separated list of display scopes"`
}

// HandleListDashboards handles the list_dashboards tool invocation
func HandleListDashboards(s ServerInterface, ctx context.Context, req *mcp.CallToolRequest, input ListDashboardsInput) (*mcp.CallToolResult, map[string]any, error) {
	params := &middleware.GetDashboardsParams{
		Limit:        input.Limit,
		Offset:       input.Offset,
		Search:       input.Search,
		FilterBy:     input.FilterBy,
		DisplayScope: input.DisplayScope,
	}

	result, err := s.Client().GetDashboards(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	data, err := ToMap(result)
	if err != nil {
		return nil, nil, err
	}

	return nil, data, nil
}

// GetDashboardTool is the MCP tool definition for getting a dashboard
var GetDashboardTool = &mcp.Tool{
	Name: "get_dashboard",
	Description: `Get detailed information about a specific dashboard by its unique key.
	
This tool retrieves complete dashboard configuration including widgets, layout, metadata, and settings. Use this when you need to inspect or work with a specific dashboard's structure and content.`,
}

type GetDashboardInput struct {
	ReportKey string `json:"report_key" jsonschema:"The unique key identifier of the dashboard to retrieve,required"`
}

// HandleGetDashboard handles the get_dashboard tool invocation
func HandleGetDashboard(s ServerInterface, ctx context.Context, req *mcp.CallToolRequest, input GetDashboardInput) (*mcp.CallToolResult, map[string]any, error) {
	result, err := s.Client().GetDashboardByKey(ctx, input.ReportKey)
	if err != nil {
		return nil, nil, err
	}

	data, err := ToMap(result)
	if err != nil {
		return nil, nil, err
	}

	return nil, data, nil
}

// CreateDashboardTool is the MCP tool definition for creating a dashboard
var CreateDashboardTool = &mcp.Tool{
	Name: "create_dashboard",
	Description: `Create a new custom dashboard in Middleware.io.
	
This tool creates a new dashboard with the specified configuration. Dashboards can be public (shared with team) or private (personal). You can organize dashboards using display scopes and provide custom keys for easier identification.`,
}

type CreateDashboardInput struct {
	Label        string `json:"label" jsonschema:"The dashboard name/title. Must be at least 3 characters long,required,minLength=3"`
	Visibility   string `json:"visibility" jsonschema:"Dashboard visibility setting. Must be either 'public' (shared with team) or 'private' (personal only),required,enum=public|private"`
	Description  string `json:"description,omitempty" jsonschema:"Optional detailed description of the dashboard's purpose and contents"`
	DisplayScope string `json:"display_scope,omitempty" jsonschema:"Optional display scope for organizing dashboards into categories or groups"`
	Key          string `json:"key,omitempty" jsonschema:"Optional unique key identifier for the dashboard. If not provided, will be auto-generated"`
}

// HandleCreateDashboard handles the create_dashboard tool invocation
func HandleCreateDashboard(s ServerInterface, ctx context.Context, req *mcp.CallToolRequest, input CreateDashboardInput) (*mcp.CallToolResult, map[string]any, error) {
	dashboardReq := &middleware.UpsertReportRequest{
		Label:        input.Label,
		Visibility:   input.Visibility,
		Description:  input.Description,
		DisplayScope: input.DisplayScope,
		Key:          input.Key,
	}

	result, err := s.Client().CreateDashboard(ctx, dashboardReq)
	if err != nil {
		return nil, nil, err
	}

	data, err := ToMap(result)
	if err != nil {
		return nil, nil, err
	}

	return nil, data, nil
}

// UpdateDashboardTool is the MCP tool definition for updating a dashboard
var UpdateDashboardTool = &mcp.Tool{
	Name: "update_dashboard",
	Description: `Update an existing dashboard's configuration and metadata.
	
This tool modifies an existing dashboard identified by its ID. You can update the name, description, visibility settings, and display scope. Use this to rename dashboards, change sharing settings, or reorganize dashboard categories.`,
}

type UpdateDashboardInput struct {
	ID           int    `json:"id" jsonschema:"The numeric ID of the dashboard to update,required"`
	Label        string `json:"label" jsonschema:"The updated dashboard name/title. Must be at least 3 characters long,required,minLength=3"`
	Visibility   string `json:"visibility" jsonschema:"Updated visibility setting. Must be either 'public' or 'private',required,enum=public|private"`
	Description  string `json:"description,omitempty" jsonschema:"Updated description of the dashboard"`
	DisplayScope string `json:"display_scope,omitempty" jsonschema:"Updated display scope for dashboard organization"`
	Key          string `json:"key,omitempty" jsonschema:"Updated unique key identifier. Must be unique across all dashboards"`
}

// HandleUpdateDashboard handles the update_dashboard tool invocation
func HandleUpdateDashboard(s ServerInterface, ctx context.Context, req *mcp.CallToolRequest, input UpdateDashboardInput) (*mcp.CallToolResult, map[string]any, error) {
	dashboardReq := &middleware.UpsertReportRequest{
		ID:           input.ID,
		Label:        input.Label,
		Visibility:   input.Visibility,
		Description:  input.Description,
		DisplayScope: input.DisplayScope,
		Key:          input.Key,
	}

	result, err := s.Client().UpdateDashboard(ctx, input.ID, dashboardReq)
	if err != nil {
		return nil, nil, err
	}

	data, err := ToMap(result)
	if err != nil {
		return nil, nil, err
	}

	return nil, data, nil
}

// DeleteDashboardTool is the MCP tool definition for deleting a dashboard
var DeleteDashboardTool = &mcp.Tool{
	Name: "delete_dashboard",
	Description: `Permanently delete a dashboard and all its widgets.
	
This tool removes a dashboard from Middleware.io. Warning: This action cannot be undone. All widgets and configurations associated with the dashboard will be permanently deleted.`,
}

type DeleteDashboardInput struct {
	ID int `json:"id" jsonschema:"The numeric ID of the dashboard to delete permanently,required"`
}

// HandleDeleteDashboard handles the delete_dashboard tool invocation
func HandleDeleteDashboard(s ServerInterface, ctx context.Context, req *mcp.CallToolRequest, input DeleteDashboardInput) (*mcp.CallToolResult, map[string]any, error) {
	err := s.Client().DeleteDashboard(ctx, input.ID)
	if err != nil {
		return nil, nil, err
	}

	return nil, map[string]any{"success": true, "message": "Dashboard deleted successfully"}, nil
}

// CloneDashboardTool is the MCP tool definition for cloning a dashboard
var CloneDashboardTool = &mcp.Tool{
	Name: "clone_dashboard",
	Description: `Create a copy of an existing dashboard with all its widgets and configuration.
	
This tool duplicates an existing dashboard, creating a new dashboard with the same widgets, layout, and settings. Useful for creating variations of dashboards or starting from a template. The cloned dashboard will have a new ID and can have different visibility settings.`,
}

type CloneDashboardInput struct {
	Label        string `json:"label" jsonschema:"The name for the new cloned dashboard. Must be at least 3 characters,required,minLength=3"`
	Visibility   string `json:"visibility" jsonschema:"Visibility setting for the cloned dashboard: 'public' or 'private',required,enum=public|private"`
	Description  string `json:"description,omitempty" jsonschema:"Optional description for the cloned dashboard"`
	DisplayScope string `json:"display_scope,omitempty" jsonschema:"Optional display scope for organizing the cloned dashboard"`
	SourceKey    string `json:"source_key,omitempty" jsonschema:"The unique key of the source dashboard to clone from"`
}

// HandleCloneDashboard handles the clone_dashboard tool invocation
func HandleCloneDashboard(s ServerInterface, ctx context.Context, req *mcp.CallToolRequest, input CloneDashboardInput) (*mcp.CallToolResult, map[string]any, error) {
	dashboardReq := &middleware.UpsertReportRequest{
		Label:        input.Label,
		Visibility:   input.Visibility,
		Description:  input.Description,
		DisplayScope: input.DisplayScope,
		Key:          input.SourceKey,
	}

	result, err := s.Client().CloneDashboard(ctx, dashboardReq)
	if err != nil {
		return nil, nil, err
	}

	data, err := ToMap(result)
	if err != nil {
		return nil, nil, err
	}

	return nil, data, nil
}

// SetDashboardFavoriteTool is the MCP tool definition for setting dashboard favorite
var SetDashboardFavoriteTool = &mcp.Tool{
	Name: "set_dashboard_favorite",
	Description: `Mark a dashboard as favorite or remove it from favorites.
	
This tool allows you to favorite dashboards for quick access. Favorited dashboards appear at the top of dashboard lists and can be filtered using the 'favorite' filter in list_dashboards. Use this to bookmark frequently accessed dashboards.`,
}

type SetDashboardFavoriteInput struct {
	ReportID int  `json:"report_id" jsonschema:"The numeric ID of the dashboard to mark as favorite or unfavorite,required"`
	Favorite bool `json:"favorite" jsonschema:"Set to true to add dashboard to favorites, false to remove from favorites,required"`
}

// HandleSetDashboardFavorite handles the set_dashboard_favorite tool invocation
func HandleSetDashboardFavorite(s ServerInterface, ctx context.Context, req *mcp.CallToolRequest, input SetDashboardFavoriteInput) (*mcp.CallToolResult, map[string]any, error) {
	err := s.Client().SetDashboardFavorite(ctx, input.ReportID, input.Favorite)
	if err != nil {
		return nil, nil, err
	}

	return nil, map[string]any{"success": true, "message": "Dashboard favorite status updated"}, nil
}
