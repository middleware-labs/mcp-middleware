package server

import (
	"context"

	"mcp-middleware/server/tools"

	"github.com/mark3labs/mcp-go/mcp"
)

// registerTools adds every MCP tool to the server unconditionally. Per-connection tool
// exclusion is handled at request time via the `?exclude_tools=` query param — see the
// tool filter and call guard wired in server.New().
//
// See: https://modelcontextprotocol.io/docs/learn/server-concepts#tools
func (s *Server) registerTools() {
	// Dashboard tools
	s.mcpServer.AddTool(tools.NewListDashboardsTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleListDashboards(s, ctx, req)
	})
	s.mcpServer.AddTool(tools.NewGetDashboardTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleGetDashboard(s, ctx, req)
	})
	s.mcpServer.AddTool(tools.NewCreateDashboardTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleCreateDashboard(s, ctx, req)
	})
	s.mcpServer.AddTool(tools.NewUpdateDashboardTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleUpdateDashboard(s, ctx, req)
	})
	s.mcpServer.AddTool(tools.NewDeleteDashboardTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleDeleteDashboard(s, ctx, req)
	})
	s.mcpServer.AddTool(tools.NewCloneDashboardTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleCloneDashboard(s, ctx, req)
	})
	s.mcpServer.AddTool(tools.NewSetDashboardFavoriteTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleSetDashboardFavorite(s, ctx, req)
	})

	// Widget tools
	s.mcpServer.AddTool(tools.NewListWidgetsTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleListWidgets(s, ctx, req)
	})
	s.mcpServer.AddTool(tools.NewCreateWidgetTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleCreateWidget(s, ctx, req)
	})
	s.mcpServer.AddTool(tools.NewUpdateWidgetTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleUpdateWidget(s, ctx, req)
	})
	s.mcpServer.AddTool(tools.NewDeleteWidgetTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleDeleteWidget(s, ctx, req)
	})
	s.mcpServer.AddTool(tools.NewGetWidgetDataTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleGetWidgetData(s, ctx, req)
	})
	s.mcpServer.AddTool(tools.NewGetMultiWidgetDataTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleGetMultiWidgetData(s, ctx, req)
	})
	s.mcpServer.AddTool(tools.NewUpdateWidgetLayoutsTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleUpdateWidgetLayouts(s, ctx, req)
	})

	// Metrics tools
	s.mcpServer.AddTool(tools.NewGetMetricsTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleGetMetrics(s, ctx, req)
	})
	s.mcpServer.AddTool(tools.NewGetResourcesTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleGetResources(s, ctx, req)
	})
	s.mcpServer.AddTool(tools.NewQueryTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleQuery(s, ctx, req)
	})

	// Alert tools
	s.mcpServer.AddTool(tools.NewListAlertsTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleListAlerts(s, ctx, req)
	})
	s.mcpServer.AddTool(tools.NewCreateAlertTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleCreateAlert(s, ctx, req)
	})
	s.mcpServer.AddTool(tools.NewGetAlertStatsTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleGetAlertStats(s, ctx, req)
	})

	// Error/Incident tools
	s.mcpServer.AddTool(tools.NewListErrorsTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleListErrors(s, ctx, req)
	})
	s.mcpServer.AddTool(tools.NewGetErrorDetailsTool(), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tools.HandleGetErrorDetails(s, ctx, req)
	})
}
