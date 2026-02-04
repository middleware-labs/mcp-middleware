package tools

import (
	"context"
	"fmt"

	"mcp-middleware/middleware"

	"github.com/mark3labs/mcp-go/mcp"
)

func NewListDashboardsTool() mcp.Tool {
	return mcp.NewTool(
		"list_dashboards",
		mcp.WithDescription(`Get a list of dashboards (i.e. reports) with advanced filtering and pagination support.
	
This tool retrieves dashboards from Middleware.io with support for searching, filtering by various criteria, and pagination. Use this to discover available dashboards, find specific dashboards by name, or filter by ownership and usage patterns.`),
		mcp.WithInputSchema[ListDashboardsInput](),
	)
}

type ListDashboardsInput struct {
	Limit        int    `json:"limit,omitempty" jsonschema:"Number of items per page for pagination"`
	Offset       int    `json:"offset,omitempty" jsonschema:"Number of items to skip for pagination (page offset)"`
	Search       string `json:"search,omitempty" jsonschema:"Search query to find dashboards by name or description"`
	FilterBy     string `json:"filter_by,omitempty" jsonschema:"Comma-separated list of filter values. Valid values: custom, created_by_you, favorite, frequently_viewed, or data source names like aws, mysql, postgresql, etc."`
	DisplayScope string `json:"display_scope,omitempty" jsonschema:"Filter dashboards by comma-separated list of display scopes"`
}

func HandleListDashboards(s ServerInterface, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := ParseInput[ListDashboardsInput](req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	params := &middleware.GetDashboardsParams{
		Limit:        input.Limit,
		Offset:       input.Offset,
		Search:       input.Search,
		FilterBy:     input.FilterBy,
		DisplayScope: input.DisplayScope,
	}

	result, err := s.Client().GetDashboards(ctx, params)
	if err != nil {
		return nil, err
	}

	return ToTextResult(result)
}

func NewGetDashboardTool() mcp.Tool {
	return mcp.NewTool(
		"get_dashboard",
		mcp.WithDescription(`Get detailed information about a specific dashboard by its unique key.
	
This tool retrieves complete dashboard configuration including widgets, layout, metadata, and settings. Use this when you need to inspect or work with a specific dashboard's structure and content.`),
		mcp.WithInputSchema[GetDashboardInput](),
	)
}

type GetDashboardInput struct {
	ReportKey string `json:"report_key" jsonschema:"The unique key identifier of the dashboard to retrieve,required"`
}

func HandleGetDashboard(s ServerInterface, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := ParseInput[GetDashboardInput](req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	result, err := s.Client().GetDashboardByKey(ctx, input.ReportKey)
	if err != nil {
		return nil, err
	}

	return ToTextResult(result)
}

func NewCreateDashboardTool() mcp.Tool {
	return mcp.NewTool(
		"create_dashboard",
		mcp.WithDescription(`Create a new custom dashboard in Middleware.io.
	
This tool creates a new dashboard with the specified configuration. Dashboards can be public (shared with team) or private (personal). You can organize dashboards using display scopes and provide custom keys for easier identification.`),
		mcp.WithInputSchema[CreateDashboardInput](),
	)
}

type CreateDashboardInput struct {
	Label       string `json:"label" jsonschema:"The dashboard name/title. Must be at least 3 characters long,required,minLength=3"`
	Visibility  string `json:"visibility" jsonschema:"Dashboard visibility setting. Must be either 'public' (shared with team) or 'private' (personal only),required,enum=public,enum=private"`
	Description string `json:"description,omitempty" jsonschema:"Optional detailed description of the dashboard's purpose and contents"`
	Key         string `json:"key,omitempty" jsonschema:"Optional unique key identifier for the dashboard. If not provided, will be auto-generated"`
}

func HandleCreateDashboard(s ServerInterface, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := ParseInput[CreateDashboardInput](req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	dashboardReq := &middleware.UpsertReportRequest{
		Label:        input.Label,
		Visibility:   input.Visibility,
		Description:  input.Description,
		DisplayScope: "", // Always empty string
		Key:          input.Key,
	}

	result, err := s.Client().CreateDashboard(ctx, dashboardReq)
	if err != nil {
		return nil, err
	}

	return ToTextResult(result)
}

func NewUpdateDashboardTool() mcp.Tool {
	return mcp.NewTool(
		"update_dashboard",
		mcp.WithDescription(`Update an existing dashboard's configuration and metadata.
	
This tool modifies an existing dashboard identified by its ID. You can update the name, description, visibility settings, and display scope. Use this to rename dashboards, change sharing settings, or reorganize dashboard categories.`),
		mcp.WithInputSchema[UpdateDashboardInput](),
	)
}

type UpdateDashboardInput struct {
	ID          int    `json:"id" jsonschema:"The numeric ID of the dashboard to update,required"`
	Label       string `json:"label" jsonschema:"The updated dashboard name/title. Must be at least 3 characters long,required,minLength=3"`
	Visibility  string `json:"visibility" jsonschema:"Updated visibility setting. Must be either 'public' or 'private',required,enum=public,enum=private"`
	Description string `json:"description,omitempty" jsonschema:"Updated description of the dashboard"`
	Key         string `json:"key,omitempty" jsonschema:"Updated unique key identifier. Must be unique across all dashboards"`
}

func HandleUpdateDashboard(s ServerInterface, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := ParseInput[UpdateDashboardInput](req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	dashboardReq := &middleware.UpsertReportRequest{
		ID:           input.ID,
		Label:        input.Label,
		Visibility:   input.Visibility,
		Description:  input.Description,
		DisplayScope: "", // Always empty string
		Key:          input.Key,
	}

	result, err := s.Client().UpdateDashboard(ctx, input.ID, dashboardReq)
	if err != nil {
		return nil, err
	}

	return ToTextResult(result)
}

func NewDeleteDashboardTool() mcp.Tool {
	return mcp.NewTool(
		"delete_dashboard",
		mcp.WithDescription(`Permanently delete a dashboard and all its widgets.
	
This tool removes a dashboard from Middleware.io. Warning: This action cannot be undone. All widgets and configurations associated with the dashboard will be permanently deleted.`),
		mcp.WithInputSchema[DeleteDashboardInput](),
	)
}

type DeleteDashboardInput struct {
	ID int `json:"id" jsonschema:"The numeric ID of the dashboard to delete permanently,required"`
}

func HandleDeleteDashboard(s ServerInterface, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := ParseInput[DeleteDashboardInput](req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	err = s.Client().DeleteDashboard(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return ToTextResult(map[string]any{"success": true, "message": "Dashboard deleted successfully"})
}

func NewCloneDashboardTool() mcp.Tool {
	return mcp.NewTool(
		"clone_dashboard",
		mcp.WithDescription(`Create a copy of an existing dashboard with all its widgets and configuration.
	
This tool duplicates an existing dashboard, creating a new dashboard with the same widgets, layout, and settings. Useful for creating variations of dashboards or starting from a template. The cloned dashboard will have a new ID and can have different visibility settings.`),
		mcp.WithInputSchema[CloneDashboardInput](),
	)
}

type CloneDashboardInput struct {
	Label       string `json:"label" jsonschema:"The name for the new cloned dashboard. Must be at least 3 characters,required,minLength=3"`
	Visibility  string `json:"visibility" jsonschema:"Visibility setting for the cloned dashboard: 'public' or 'private',required,enum=public,enum=private"`
	Description string `json:"description,omitempty" jsonschema:"Optional description for the cloned dashboard"`
	SourceKey   string `json:"source_key,omitempty" jsonschema:"The unique key of the source dashboard to clone from"`
}

func HandleCloneDashboard(s ServerInterface, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := ParseInput[CloneDashboardInput](req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	dashboardReq := &middleware.UpsertReportRequest{
		Label:        input.Label,
		Visibility:   input.Visibility,
		Description:  input.Description,
		DisplayScope: "", // Always empty string
		Key:          input.SourceKey,
	}

	result, err := s.Client().CloneDashboard(ctx, dashboardReq)
	if err != nil {
		return nil, err
	}

	return ToTextResult(result)
}

func NewSetDashboardFavoriteTool() mcp.Tool {
	return mcp.NewTool(
		"set_dashboard_favorite",
		mcp.WithDescription(`Mark a dashboard as favorite or remove it from favorites.
	
This tool allows you to favorite dashboards for quick access. Favorited dashboards appear at the top of dashboard lists and can be filtered using the 'favorite' filter in list_dashboards. Use this to bookmark frequently accessed dashboards.`),
		mcp.WithInputSchema[SetDashboardFavoriteInput](),
	)
}

type SetDashboardFavoriteInput struct {
	ReportID int  `json:"report_id" jsonschema:"The numeric ID of the dashboard to mark as favorite or unfavorite,required"`
	Favorite bool `json:"favorite" jsonschema:"Set to true to add dashboard to favorites, false to remove from favorites,required"`
}

func HandleSetDashboardFavorite(s ServerInterface, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := ParseInput[SetDashboardFavoriteInput](req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	err = s.Client().SetDashboardFavorite(ctx, input.ReportID, input.Favorite)
	if err != nil {
		return nil, err
	}

	return ToTextResult(map[string]any{"success": true, "message": "Dashboard favorite status updated"})
}
