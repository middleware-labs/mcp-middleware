package server

import (
	"context"

	"mcp-middleware/server/tools"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerTools registers all available MCP tools with the server.
// This function is called during server initialization to set up all tool handlers.
// Tools are functions that AI models can actively call to perform actions.
// See: https://modelcontextprotocol.io/docs/learn/server-concepts#tools
func (s *Server) registerTools() {
	// Dashboard tools
	if !s.config.IsToolExcluded("list_dashboards") {
		mcp.AddTool(s.mcpServer, tools.ListDashboardsTool, func(ctx context.Context, req *mcp.CallToolRequest, input tools.ListDashboardsInput) (*mcp.CallToolResult, map[string]any, error) {
			return tools.HandleListDashboards(s, ctx, req, input)
		})
	}
	if !s.config.IsToolExcluded("get_dashboard") {
		mcp.AddTool(s.mcpServer, tools.GetDashboardTool, func(ctx context.Context, req *mcp.CallToolRequest, input tools.GetDashboardInput) (*mcp.CallToolResult, map[string]any, error) {
			return tools.HandleGetDashboard(s, ctx, req, input)
		})
	}
	if !s.config.IsToolExcluded("create_dashboard") {
		mcp.AddTool(s.mcpServer, tools.CreateDashboardTool, func(ctx context.Context, req *mcp.CallToolRequest, input tools.CreateDashboardInput) (*mcp.CallToolResult, map[string]any, error) {
			return tools.HandleCreateDashboard(s, ctx, req, input)
		})
	}
	if !s.config.IsToolExcluded("update_dashboard") {
		mcp.AddTool(s.mcpServer, tools.UpdateDashboardTool, func(ctx context.Context, req *mcp.CallToolRequest, input tools.UpdateDashboardInput) (*mcp.CallToolResult, map[string]any, error) {
			return tools.HandleUpdateDashboard(s, ctx, req, input)
		})
	}
	if !s.config.IsToolExcluded("delete_dashboard") {
		mcp.AddTool(s.mcpServer, tools.DeleteDashboardTool, func(ctx context.Context, req *mcp.CallToolRequest, input tools.DeleteDashboardInput) (*mcp.CallToolResult, map[string]any, error) {
			return tools.HandleDeleteDashboard(s, ctx, req, input)
		})
	}
	if !s.config.IsToolExcluded("clone_dashboard") {
		mcp.AddTool(s.mcpServer, tools.CloneDashboardTool, func(ctx context.Context, req *mcp.CallToolRequest, input tools.CloneDashboardInput) (*mcp.CallToolResult, map[string]any, error) {
			return tools.HandleCloneDashboard(s, ctx, req, input)
		})
	}
	if !s.config.IsToolExcluded("set_dashboard_favorite") {
		mcp.AddTool(s.mcpServer, tools.SetDashboardFavoriteTool, func(ctx context.Context, req *mcp.CallToolRequest, input tools.SetDashboardFavoriteInput) (*mcp.CallToolResult, map[string]any, error) {
			return tools.HandleSetDashboardFavorite(s, ctx, req, input)
		})
	}

	// Widget tools
	if !s.config.IsToolExcluded("list_widgets") {
		mcp.AddTool(s.mcpServer, tools.ListWidgetsTool, func(ctx context.Context, req *mcp.CallToolRequest, input tools.ListWidgetsInput) (*mcp.CallToolResult, map[string]any, error) {
			return tools.HandleListWidgets(s, ctx, req, input)
		})
	}
	if !s.config.IsToolExcluded("create_widget") {
		mcp.AddTool(s.mcpServer, tools.CreateWidgetTool, func(ctx context.Context, req *mcp.CallToolRequest, input tools.CreateWidgetInput) (*mcp.CallToolResult, map[string]any, error) {
			return tools.HandleCreateWidget(s, ctx, req, input)
		})
	}
	if !s.config.IsToolExcluded("delete_widget") {
		mcp.AddTool(s.mcpServer, tools.DeleteWidgetTool, func(ctx context.Context, req *mcp.CallToolRequest, input tools.DeleteWidgetInput) (*mcp.CallToolResult, map[string]any, error) {
			return tools.HandleDeleteWidget(s, ctx, req, input)
		})
	}
	if !s.config.IsToolExcluded("get_widget_data") {
		mcp.AddTool(s.mcpServer, tools.GetWidgetDataTool, func(ctx context.Context, req *mcp.CallToolRequest, input tools.GetWidgetDataInput) (*mcp.CallToolResult, map[string]any, error) {
			return tools.HandleGetWidgetData(s, ctx, req, input)
		})
	}
	if !s.config.IsToolExcluded("get_multi_widget_data") {
		mcp.AddTool(s.mcpServer, tools.GetMultiWidgetDataTool, func(ctx context.Context, req *mcp.CallToolRequest, input tools.GetMultiWidgetDataInput) (*mcp.CallToolResult, map[string]any, error) {
			return tools.HandleGetMultiWidgetData(s, ctx, req, input)
		})
	}
	if !s.config.IsToolExcluded("update_widget_layouts") {
		mcp.AddTool(s.mcpServer, tools.UpdateWidgetLayoutsTool, func(ctx context.Context, req *mcp.CallToolRequest, input tools.UpdateWidgetLayoutsInput) (*mcp.CallToolResult, map[string]any, error) {
			return tools.HandleUpdateWidgetLayouts(s, ctx, req, input)
		})
	}

	// Metrics tools
	if !s.config.IsToolExcluded("get_metrics") {
		mcp.AddTool(s.mcpServer, tools.GetMetricsTool, func(ctx context.Context, req *mcp.CallToolRequest, input tools.GetMetricsInput) (*mcp.CallToolResult, map[string]any, error) {
			return tools.HandleGetMetrics(s, ctx, req, input)
		})
	}
	if !s.config.IsToolExcluded("get_resources") {
		mcp.AddTool(s.mcpServer, tools.GetResourcesTool, func(ctx context.Context, req *mcp.CallToolRequest, input tools.GetResourcesInput) (*mcp.CallToolResult, map[string]any, error) {
			return tools.HandleGetResources(s, ctx, req, input)
		})
	}

	// Alert tools
	if !s.config.IsToolExcluded("list_alerts") {
		mcp.AddTool(s.mcpServer, tools.ListAlertsTool, func(ctx context.Context, req *mcp.CallToolRequest, input tools.ListAlertsInput) (*mcp.CallToolResult, map[string]any, error) {
			return tools.HandleListAlerts(s, ctx, req, input)
		})
	}
	if !s.config.IsToolExcluded("create_alert") {
		mcp.AddTool(s.mcpServer, tools.CreateAlertTool, func(ctx context.Context, req *mcp.CallToolRequest, input tools.CreateAlertInput) (*mcp.CallToolResult, map[string]any, error) {
			return tools.HandleCreateAlert(s, ctx, req, input)
		})
	}
	if !s.config.IsToolExcluded("get_alert_stats") {
		mcp.AddTool(s.mcpServer, tools.GetAlertStatsTool, func(ctx context.Context, req *mcp.CallToolRequest, input tools.GetAlertStatsInput) (*mcp.CallToolResult, map[string]any, error) {
			return tools.HandleGetAlertStats(s, ctx, req, input)
		})
	}
}

