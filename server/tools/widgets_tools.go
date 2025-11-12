package tools

import (
	"context"
	"fmt"

	"mcp-middleware/middleware"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var ListWidgetsTool = &mcp.Tool{
	Name: "list_widgets",
	Description: `Get a list of widgets associated with a specific dashboard or display scope.
	
This tool retrieves all widgets (charts, graphs, tables) that belong to a dashboard or scope. Widgets are the building blocks of dashboards - each widget represents a visualization of your monitoring data. Use this to discover what widgets are available in a dashboard or to inspect widget configurations.`,
}

type ListWidgetsInput struct {
	ReportID     int    `json:"report_id,omitempty" jsonschema:"The numeric ID of the dashboard (report) to filter widgets by"`
	DisplayScope string `json:"display_scope,omitempty" jsonschema:"The display scope to filter widgets by (e.g., 'infrastructure', 'apm', 'logs')"`
}

func HandleListWidgets(s ServerInterface, ctx context.Context, req *mcp.CallToolRequest, input ListWidgetsInput) (*mcp.CallToolResult, map[string]any, error) {
	params := &middleware.GetWidgetsParams{
		ReportID:     input.ReportID,
		DisplayScope: input.DisplayScope,
	}

	result, err := s.Client().GetWidgets(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get widgets: %w", err)
	}

	return ToTextResult(result)
}

var CreateWidgetTool = &mcp.Tool{
	Name: "create_widget",
	Description: `Create a new widget or update an existing widget on a dashboard.
	
This tool allows you to add new visualizations (charts, graphs, tables) to dashboards or modify existing ones. The builder_config contains the query, chart type, and visualization settings. If builder_id is provided, it updates the existing widget; otherwise, it creates a new one.`,
}

type CreateWidgetInput struct {
	Label         string `json:"label" jsonschema:"The display name for the widget (e.g., 'CPU Usage', 'Error Rate'),required"`
	Key           string `json:"key,omitempty" jsonschema:"Optional unique key identifier for the widget"`
	Description   string `json:"description,omitempty" jsonschema:"Optional description explaining what the widget displays"`
	BuilderConfig any    `json:"builder_config,omitempty" jsonschema:"Widget configuration object containing queries, chart type, display settings, and data sources. This is a complex object specific to widget type"`
	BuilderID     int    `json:"builder_id,omitempty" jsonschema:"If provided, updates the existing widget with this ID instead of creating a new one"`
}

func HandleCreateWidget(s ServerInterface, ctx context.Context, req *mcp.CallToolRequest, input CreateWidgetInput) (*mcp.CallToolResult, map[string]any, error) {
	widget := &middleware.CustomWidget{
		Label:         input.Label,
		Key:           input.Key,
		Description:   input.Description,
		BuilderConfig: input.BuilderConfig,
		BuilderID:     input.BuilderID,
	}

	result, err := s.Client().CreateWidget(ctx, widget)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create widget: %w", err)
	}

	return ToTextResult(result)
}

var DeleteWidgetTool = &mcp.Tool{
	Name: "delete_widget",
	Description: `Permanently delete a widget from a dashboard.
	
This tool removes a widget (chart, graph, table) from its dashboard. Warning: This action cannot be undone. The widget configuration and data will be permanently deleted.`,
}

type DeleteWidgetInput struct {
	BuilderID int `json:"builder_id" jsonschema:"The numeric builder ID of the widget to delete permanently,required"`
}

func HandleDeleteWidget(s ServerInterface, ctx context.Context, req *mcp.CallToolRequest, input DeleteWidgetInput) (*mcp.CallToolResult, map[string]any, error) {
	err := s.Client().DeleteWidget(ctx, input.BuilderID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to delete widget: %w", err)
	}

	return ToTextResult(map[string]any{"success": true, "message": "Widget deleted successfully"})
}

var GetWidgetDataTool = &mcp.Tool{
	Name: "get_widget_data",
	Description: `Fetch the actual data and metrics displayed by a specific widget.
	
This tool executes the widget's query and returns the visualization data (time series, metrics, logs, traces). Use this to get the current values shown in a widget, analyze trends, or export widget data. The data format depends on the widget type (timeseries, table, single value, etc.).`,
}

type GetWidgetDataInput struct {
	BuilderID     int    `json:"builder_id,omitempty" jsonschema:"The numeric builder ID of the widget to fetch data for"`
	Key           string `json:"key,omitempty" jsonschema:"Alternative to builder_id: the unique key identifier of the widget"`
	Label         string `json:"label,omitempty" jsonschema:"Alternative to builder_id: the label of the widget"`
	BuilderConfig any    `json:"builder_config,omitempty" jsonschema:"Widget configuration containing the query and data source settings"`
	UseV2         bool   `json:"use_v2,omitempty" jsonschema:"Set to true to use the newer v2 data format (default: false)"`
}

func HandleGetWidgetData(s ServerInterface, ctx context.Context, req *mcp.CallToolRequest, input GetWidgetDataInput) (*mcp.CallToolResult, map[string]any, error) {
	widget := &middleware.CustomWidget{
		BuilderID:     input.BuilderID,
		Key:           input.Key,
		Label:         input.Label,
		BuilderConfig: input.BuilderConfig,
		UseV2:         input.UseV2,
	}

	result, err := s.Client().GetWidgetData(ctx, widget)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get widget data: %w", err)
	}

	return ToTextResult(result)
}

var GetMultiWidgetDataTool = &mcp.Tool{
	Name: "get_multi_widget_data",
	Description: `Fetch data for multiple widgets simultaneously in a single request.
	
This tool is optimized for loading data for multiple widgets at once, such as when refreshing an entire dashboard. It's more efficient than calling get_widget_data multiple times. Returns data for all requested widgets in a single response.`,
}

type GetMultiWidgetDataInput struct {
	Widgets []WidgetDataRequest `json:"widgets" jsonschema:"Array of widget specifications to fetch data for. Each widget can be identified by builder_id, key, or label,required"`
}

type WidgetDataRequest struct {
	BuilderID     int    `json:"builder_id,omitempty" jsonschema:"The numeric builder ID of the widget"`
	Key           string `json:"key,omitempty" jsonschema:"The unique key identifier of the widget"`
	Label         string `json:"label,omitempty" jsonschema:"The label of the widget"`
	BuilderConfig any    `json:"builder_config,omitempty" jsonschema:"Widget configuration containing query and display settings"`
	UseV2         bool   `json:"use_v2,omitempty" jsonschema:"Use v2 data format (default: false)"`
}

func HandleGetMultiWidgetData(s ServerInterface, ctx context.Context, req *mcp.CallToolRequest, input GetMultiWidgetDataInput) (*mcp.CallToolResult, map[string]any, error) {
	widgets := make([]middleware.CustomWidget, len(input.Widgets))
	for i, w := range input.Widgets {
		widgets[i] = middleware.CustomWidget{
			BuilderID:     w.BuilderID,
			Key:           w.Key,
			Label:         w.Label,
			BuilderConfig: w.BuilderConfig,
			UseV2:         w.UseV2,
		}
	}

	result, err := s.Client().GetMultiWidgetData(ctx, widgets)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get multi widget data: %w", err)
	}

	// Return as JSON text only (no structuredContent)
	return ToTextResult(map[string]any{"widgets": result})
}

var UpdateWidgetLayoutsTool = &mcp.Tool{
	Name: "update_widget_layouts",
	Description: `Update the position and size of widgets on a dashboard.
	
This tool modifies the layout (position, size) of multiple widgets on a dashboard. Use this to rearrange widgets, resize them, or optimize dashboard layout. The dashboard uses a grid system where x,y represent position and w,h represent size in grid units.`,
}

type UpdateWidgetLayoutsInput struct {
	Layouts []LayoutItemInput `json:"layouts" jsonschema:"Array of layout specifications for each widget. Each item defines position and size in the dashboard grid,required"`
}

type LayoutItemInput struct {
	X       int `json:"x" jsonschema:"Horizontal position in the grid (0-based index from left)"`
	Y       int `json:"y" jsonschema:"Vertical position in the grid (0-based index from top)"`
	W       int `json:"w" jsonschema:"Width in grid units"`
	H       int `json:"h" jsonschema:"Height in grid units"`
	ScopeID int `json:"scope_id,omitempty" jsonschema:"The scope ID of the widget to update layout for"`
}

func HandleUpdateWidgetLayouts(s ServerInterface, ctx context.Context, req *mcp.CallToolRequest, input UpdateWidgetLayoutsInput) (*mcp.CallToolResult, map[string]any, error) {
	layouts := make([]middleware.LayoutItem, len(input.Layouts))
	for i, l := range input.Layouts {
		layouts[i] = middleware.LayoutItem{
			X:       l.X,
			Y:       l.Y,
			W:       l.W,
			H:       l.H,
			ScopeID: l.ScopeID,
		}
	}

	layoutReq := &middleware.LayoutRequest{
		Layouts: layouts,
	}

	err := s.Client().UpdateWidgetLayouts(ctx, layoutReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update widget layouts: %w", err)
	}

	return ToTextResult(map[string]any{"success": true, "message": "Widget layouts updated successfully"})
}
