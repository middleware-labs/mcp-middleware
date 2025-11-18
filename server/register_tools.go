package server

import (
	"context"

	"mcp-middleware/server/tools"

	"github.com/mark3labs/mcp-go/mcp"
)

// See: https://modelcontextprotocol.io/docs/learn/server-concepts#tools
func (s *Server) registerTools() {
	// Dashboard tools
	if !s.config.IsToolExcluded("list_dashboards") {
		s.mcpServer.AddTool(tools.NewListDashboardsTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.HandleListDashboards(s, ctx, req)
		})
	}
	if !s.config.IsToolExcluded("get_dashboard") {
		s.mcpServer.AddTool(tools.NewGetDashboardTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.HandleGetDashboard(s, ctx, req)
		})
	}
	if !s.config.IsToolExcluded("create_dashboard") {
		s.mcpServer.AddTool(tools.NewCreateDashboardTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.HandleCreateDashboard(s, ctx, req)
		})
	}
	if !s.config.IsToolExcluded("update_dashboard") {
		s.mcpServer.AddTool(tools.NewUpdateDashboardTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.HandleUpdateDashboard(s, ctx, req)
		})
	}
	if !s.config.IsToolExcluded("delete_dashboard") {
		s.mcpServer.AddTool(tools.NewDeleteDashboardTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.HandleDeleteDashboard(s, ctx, req)
		})
	}
	if !s.config.IsToolExcluded("clone_dashboard") {
		s.mcpServer.AddTool(tools.NewCloneDashboardTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.HandleCloneDashboard(s, ctx, req)
		})
	}
	if !s.config.IsToolExcluded("set_dashboard_favorite") {
		s.mcpServer.AddTool(tools.NewSetDashboardFavoriteTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.HandleSetDashboardFavorite(s, ctx, req)
		})
	}

	// Widget tools
	if !s.config.IsToolExcluded("list_widgets") {
		s.mcpServer.AddTool(tools.NewListWidgetsTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.HandleListWidgets(s, ctx, req)
		})
	}
	if !s.config.IsToolExcluded("create_widget") {
		s.mcpServer.AddTool(tools.NewCreateWidgetTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.HandleCreateWidget(s, ctx, req)
		})
	}
	if !s.config.IsToolExcluded("delete_widget") {
		s.mcpServer.AddTool(tools.NewDeleteWidgetTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.HandleDeleteWidget(s, ctx, req)
		})
	}
	if !s.config.IsToolExcluded("get_widget_data") {
		s.mcpServer.AddTool(tools.NewGetWidgetDataTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.HandleGetWidgetData(s, ctx, req)
		})
	}
	if !s.config.IsToolExcluded("get_multi_widget_data") {
		s.mcpServer.AddTool(tools.NewGetMultiWidgetDataTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.HandleGetMultiWidgetData(s, ctx, req)
		})
	}
	if !s.config.IsToolExcluded("update_widget_layouts") {
		s.mcpServer.AddTool(tools.NewUpdateWidgetLayoutsTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.HandleUpdateWidgetLayouts(s, ctx, req)
		})
	}

	// Metrics tools
	if !s.config.IsToolExcluded("get_metrics") {
		s.mcpServer.AddTool(tools.NewGetMetricsTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.HandleGetMetrics(s, ctx, req)
		})
	}
	if !s.config.IsToolExcluded("get_resources") {
		s.mcpServer.AddTool(tools.NewGetResourcesTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.HandleGetResources(s, ctx, req)
		})
	}
	if !s.config.IsToolExcluded("query") {
		s.mcpServer.AddTool(tools.NewQueryTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.HandleQuery(s, ctx, req)
		})
	}

	// Alert tools
	if !s.config.IsToolExcluded("list_alerts") {
		s.mcpServer.AddTool(tools.NewListAlertsTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.HandleListAlerts(s, ctx, req)
		})
	}
	if !s.config.IsToolExcluded("create_alert") {
		s.mcpServer.AddTool(tools.NewCreateAlertTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.HandleCreateAlert(s, ctx, req)
		})
	}
	if !s.config.IsToolExcluded("get_alert_stats") {
		s.mcpServer.AddTool(tools.NewGetAlertStatsTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.HandleGetAlertStats(s, ctx, req)
		})
	}

	// Error/Incident tools
	if !s.config.IsToolExcluded("list_errors") {
		s.mcpServer.AddTool(tools.NewListErrorsTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.HandleListErrors(s, ctx, req)
		})
	}
	if !s.config.IsToolExcluded("get_error_details") {
		s.mcpServer.AddTool(tools.NewGetErrorDetailsTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return tools.HandleGetErrorDetails(s, ctx, req)
		})
	}
}
