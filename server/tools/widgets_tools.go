package tools

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"mcp-middleware/middleware"

	"github.com/mark3labs/mcp-go/mcp"
)

func NewListWidgetsTool() mcp.Tool {
	return mcp.NewTool(
		"list_widgets",
		mcp.WithDescription(`Get a list of widgets associated with a specific dashboard or display scope.
	
This tool retrieves all widgets (charts, graphs, tables) that belong to a dashboard or scope. Widgets are the building blocks of dashboards - each widget represents a visualization of your monitoring data. Use this to discover what widgets are available in a dashboard or to inspect widget configurations.`),
		mcp.WithInputSchema[ListWidgetsInput](),
	)
}

type ListWidgetsInput struct {
	ReportID     int    `json:"report_id,omitempty" jsonschema:"The numeric ID of the dashboard (report) to filter widgets by"`
	DisplayScope string `json:"display_scope,omitempty" jsonschema:"The display scope to filter widgets by (e.g., 'infrastructure', 'apm', 'logs')"`
	Message      string `json:"message" jsonschema:"Message to know which widgets are being listed. Length should be less than 100 characters."`
}

func HandleListWidgets(s ServerInterface, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := ParseInput[ListWidgetsInput](req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	params := &middleware.GetWidgetsParams{
		ReportID:     input.ReportID,
		DisplayScope: input.DisplayScope,
	}

	result, err := s.Client().GetWidgets(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get widgets: %w", err)
	}

	return ToTextResult(result)
}

func NewCreateWidgetTool() mcp.Tool {
	return mcp.NewTool(
		"create_widget",
		mcp.WithDescription(`Create a new widget or update an existing widget on a dashboard.
	
This tool allows you to add new visualizations (charts, graphs, tables) to dashboards or modify existing ones. 

Creation Workflow:
1. Identify the Resource: Use 'get_resources' to find available resource types (e.g., 'host', 'container').
2. Identify Metrics/Data: Use 'get_metrics' to find available metrics and dimensions for that resource. IMPORTANT: You MUST explore all edge cases! For each metric, explicitly check its supported 'filters' and 'groupby' tags using the get_metrics tool to ensure your widget query is valid and precise.
3. Construct BuilderConfig: Create the widget configuration using the discovered resource and metrics.

BuilderConfig Structure (Critical):
- 'source.name': MUST be an exact resource type returned by 'get_resources' (e.g., 'host', 'container').
- 'metricMetadata': Defines the specific metric to visualize.
- 'columns': Array of column configuration objects. Each object MUST have 'name' (metric name) and can optionally have 'aggregation_method' (e.g., avg, sum, min, max, uniq, count, group) and 'rollup_method' (e.g., avg, sum, min, max, none). This replaces the old string array format.
- 'group_by': (Optional) Dimensions to group data by (discovered via 'get_metrics' with data_type='groupby').
- 'filter_with': (Optional) Conditions to filter data (discovered via 'get_metrics' with data_type='filters').

IMPORTANT - Validation Rules:
- Resource Validation: You CANNOT use arbitrary resource names. You MUST use the exact strings returned by 'get_resources'.
- Dashboard ID: The 'report_id' is REQUIRED to place the widget on a specific dashboard.
- Widget Type: Choose the appropriate visualization type (e.g., 'time_series_chart', 'bar_chart') based on the data.
- Layout Requirements: Based on the widget type, you MUST set proper layout. Width (w) must be minimum 4 (this is a strict minimum requirement) and height (h) must be minimum 6 (this is a strict minimum requirement). The layout dimensions should be appropriate for the widget type to ensure proper visualization.

Use this tool to build rich, data-driven dashboards by combining resources, metrics, and visualizations.`),
		mcp.WithInputSchema[CreateWidgetInput](),
	)
}

func generateWidgetKey(label string) string {
	re := regexp.MustCompile(`[^A-Za-z0-9]`)
	cleaned := re.ReplaceAllString(label, "_")
	cleaned = strings.ToLower(cleaned)
	randomID := fmt.Sprintf("%d", time.Now().UnixNano()%1000000000)
	return fmt.Sprintf("%s_%s", cleaned, randomID)
}

func getWidgetAppID(widgetType string) int {
	widgetTypeMap := map[string]int{
		"time_series_chart": 1,
		"bar_chart":         2,
		"pie_chart":         3,
		"scatter_plot":      4,
		"data_table":        5,
		"count_chart":       7,
		"tree_chart":        8,
		"top_list_chart":    9,
		"heatmap_chart":     10,
		"hexagon_chart":     11,
		"query_value":       12,
	}

	if id, ok := widgetTypeMap[widgetType]; ok {
		return id
	}
	return 1
}

type CreateWidgetInput struct {
	Label             string                   `json:"label" jsonschema:"The display name for the widget (e.g., 'CPU Usage', 'Error Rate'),required"`
	WidgetType        string                   `json:"widget_type" jsonschema:"The type of chart/widget to create,required,enum=time_series_chart|bar_chart|data_table|query_value|pie_chart|scatter_plot|count_chart|tree_chart|top_list_chart|heatmap_chart|hexagon_chart"`
	Key               string                   `json:"key,omitempty" jsonschema:"Optional unique key identifier for the widget"`
	Description       string                   `json:"description,omitempty" jsonschema:"Optional description explaining what the widget displays"`
	BuilderConfig     []BuilderConfigItemInput `json:"builderConfig" jsonschema:"Widget configuration array containing queries, chart type, display settings, and data sources. Each item should have: columns, source, id, meta_data, metricMetadata, key, group_by, and filter_with"`
	ReportID          int                      `json:"report_id" jsonschema:"The numeric ID of the dashboard ID (Report ID) where this widget will be created"`
	ReportKey         string                   `json:"report_key,omitempty" jsonschema:"The unique key identifier of the dashboard (report) where this widget will be created"`
	ReportName        string                   `json:"report_name,omitempty" jsonschema:"The name of the dashboard (report) where this widget will be created"`
	ReportDescription string                   `json:"report_description,omitempty" jsonschema:"Optional description of the dashboard (report)"`
	ReportMetadata    any                      `json:"report_metadata,omitempty" jsonschema:"Optional metadata for the dashboard (report)"`
	DisableUserEdit   bool                     `json:"disable_user_edit,omitempty" jsonschema:"Whether to disable user editing of the widget (default: false)"`
	Layout            *LayoutItemInput         `json:"layout" jsonschema:"Layout for the widget including coordinates and size. Based on the widget type, you MUST set proper layout. Width (w) must be minimum 4 (strict minimum requirement) and height (h) must be minimum 6 (strict minimum requirement),required"`
}

type AggregationMethod string

const (
	AggregationAvg   AggregationMethod = "avg"
	AggregationSum   AggregationMethod = "sum"
	AggregationMin   AggregationMethod = "min"
	AggregationMax   AggregationMethod = "max"
	AggregationAny   AggregationMethod = "any"
	AggregationUniq  AggregationMethod = "uniq"
	AggregationCount AggregationMethod = "count"
	AggregationGroup AggregationMethod = "group"
)

type RollupMethod string

const (
	RollupAvg   RollupMethod = "avg"
	RollupSum   RollupMethod = "sum"
	RollupMin   RollupMethod = "min"
	RollupMax   RollupMethod = "max"
	RollupAny   RollupMethod = "any"
	RollupUniq  RollupMethod = "uniq"
	RollupCount RollupMethod = "count"
	RollupGroup RollupMethod = "group"
	RollupNone  RollupMethod = "none"
)

type ColumnConfig struct {
	Name              string            `json:"name" jsonschema:"The metric or metric attribute name (e.g., 'k8s.node.cpu.utilization', 'host.memory.usage'),required"`
	AggregationMethod AggregationMethod `json:"aggregation_method" jsonschema:"Aggregation method to apply to this column. Supported values: avg, sum, min, max, any (default), uniq, count, group. If empty or 'any', no aggregation is applied.,enum=avg,enum=sum,enum=min,enum=max,enum=any,enum=uniq,enum=count,enum=group"`
	RollupMethod      RollupMethod      `json:"rollup_method,omitempty" jsonschema:"Rollup method to apply to this column. Supported values: avg, sum, min, max, any (default), uniq, count, group, none. If empty, 'none', or 'any', no rollup is applied.,enum=avg,enum=sum,enum=min,enum=max,enum=any,enum=uniq,enum=count,enum=group,enum=none"`
}

func transformColumns(columns []ColumnConfig) []string {
	transformed := make([]string, len(columns))
	for i, col := range columns {
		if col.AggregationMethod == "" || col.AggregationMethod == AggregationAny {
			transformed[i] = col.Name
			continue
		}

		if col.RollupMethod == "" || col.RollupMethod == RollupNone || col.RollupMethod == RollupAny {
			transformed[i] = fmt.Sprintf("%s(%s)", col.AggregationMethod, col.Name)
		} else {
			transformed[i] = fmt.Sprintf("%s(%s, value(%s))", col.AggregationMethod, col.Name, col.RollupMethod)
		}
	}
	return transformed
}

type BuilderConfigItemInput struct {
	Columns        []ColumnConfig                       `json:"columns" jsonschema:"Array of column configurations, each specifying metric/attribute name and its aggregation/rollup methods. Each column can have different aggregation and rollup settings.,required"`
	Source         *middleware.BuilderConfigSource      `json:"source,omitempty" jsonschema:"Data source configuration with name, alias, and dataset_id. IMPORTANT: The source.name field MUST be a resource type that is supported by Middleware and returned by the get_resources tool. You MUST first call the get_resources tool to discover available resource types, then use only those exact resource type names here. Do not use arbitrary or guessed resource names. Examples (if returned by get_resources): 'host', 'container', 'log', 'trace', 'k8s.pod', 'database', 'service', etc. The source.name identifies which resource type the widget will query data from."`
	ID             string                               `json:"id,omitempty" jsonschema:"Unique identifier for this config item (UUID format)"`
	MetaData       *middleware.BuilderConfigMetaData    `json:"meta_data,omitempty" jsonschema:"Metadata containing metricTypes mapping"`
	MetricMetadata map[string]middleware.MetricMetadata `json:"metricMetadata,omitempty" jsonschema:"Map of metric names to their metadata. Each key is a metric name (e.g., \"k8s.node.cpu.utilization_percent\") and value is the metadata object with name, label, resource, type, attributes, and config"`
	Key            string                               `json:"key,omitempty" jsonschema:"Key identifier for this config item"`
	GroupBy        []string                             `json:"group_by,omitempty" jsonschema:"Array of attribute names to group by (e.g., [\"host.cpu.model.id\"]). This will be converted to SELECT_DATA_BY in the 'with' array"`
	FilterWith     any                                  `json:"filter_with,omitempty" jsonschema:"Filter conditions object with 'and' or 'or' arrays (e.g., {\"and\": [{\"host.id\": {\"=\": \"ai-team2\"}}, {\"host.name\": {\"LIKE\": \"%ai%\"}}]}). This will be converted to ATTRIBUTE_FILTER in the 'with' array"`
}

func convertToMiddlewareBuilderConfig(input []BuilderConfigItemInput) []middleware.BuilderConfigItem {
	builderConfig := make([]middleware.BuilderConfigItem, len(input))
	for i, configInput := range input {
		var withItems []middleware.BuilderConfigWith

		if len(configInput.GroupBy) > 0 {
			withItems = append(withItems, middleware.BuilderConfigWith{
				Key:   middleware.BuilderConfigWithKeySelectDataBy,
				Value: configInput.GroupBy,
				IsArg: true,
			})
		}

		if configInput.FilterWith != nil {
			withItems = append(withItems, middleware.BuilderConfigWith{
				Key:   middleware.BuilderConfigWithKeyAttributeFilter,
				Value: configInput.FilterWith,
				IsArg: true,
			})
		}

		var metricMetadata *middleware.MetricMetadata
		if len(configInput.MetricMetadata) > 0 {
			for _, v := range configInput.MetricMetadata {
				metricMetadata = &v
				break
			}
		}

		transformedColumns := transformColumns(configInput.Columns)

		builderConfig[i] = middleware.BuilderConfigItem{
			With:           withItems,
			Columns:        transformedColumns,
			Source:         configInput.Source,
			ID:             configInput.ID,
			MetaData:       configInput.MetaData,
			MetricMetadata: metricMetadata,
			Key:            configInput.Key,
		}
	}
	return builderConfig
}

func HandleCreateWidget(s ServerInterface, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := ParseInput[CreateWidgetInput](req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}
	var builderViewOptions *middleware.BuilderViewOptions
	if input.ReportID > 0 || input.ReportKey != "" || input.ReportName != "" {
		builderViewOptions = &middleware.BuilderViewOptions{
			DisableUserEdit: input.DisableUserEdit,
		}

		if input.ReportID > 0 || input.ReportKey != "" || input.ReportName != "" {
			builderViewOptions.Report = &middleware.ReportView{
				ReportID:          input.ReportID,
				ReportKey:         input.ReportKey,
				ReportName:        input.ReportName,
				ReportDescription: input.ReportDescription,
				Metadata:          input.ReportMetadata,
			}
		}
	}

	widgetKey := input.Key
	if widgetKey == "" {
		widgetKey = generateWidgetKey(input.Label)
	}

	builderConfig := convertToMiddlewareBuilderConfig(input.BuilderConfig)

	widgetAppID := getWidgetAppID(input.WidgetType)

	defaultLayout := middleware.LayoutItem{
		X: 0,
		Y: 0,
		W: 4,
		H: 6,
	}

	var layout *middleware.LayoutItem
	if input.Layout != nil {
		layout = &middleware.LayoutItem{
			X: defaultLayout.X,
			Y: defaultLayout.Y,
			W: input.Layout.W,
			H: input.Layout.H,
		}

		if layout.W < defaultLayout.W {
			layout.W = defaultLayout.W
		}
		if layout.H < defaultLayout.H {
			layout.H = defaultLayout.H
		}
	} else {
		layout = &defaultLayout
	}

	widget := &middleware.CustomWidget{
		Label:              input.Label,
		Key:                widgetKey,
		Description:        input.Description,
		BuilderConfig:      builderConfig,
		BuilderViewOptions: builderViewOptions,
		WidgetAppID:        widgetAppID,

		// Default values (always the same - not exposed in input)
		BuilderID:       -1,
		ScopeID:         -1,
		IsClone:         false,
		Category:        "Metrics",
		Formulas:        []any{},
		DontRefreshData: false,
		Layout:          layout,
	}

	result, err := s.Client().CreateWidget(ctx, widget)
	if err != nil {
		return nil, fmt.Errorf("failed to create widget: %w", err)
	}

	return ToTextResult(result)
}

func NewUpdateWidgetTool() mcp.Tool {
	return mcp.NewTool(
		"update_widget",
		mcp.WithDescription(`Update an existing widget on a dashboard.
	
This tool allows you to modify existing visualizations (charts, graphs, tables) on dashboards. The builderConfig is an array of configuration objects, each containing queries, chart type, and visualization settings. Each builderConfig item should have: with (array), columns (array of column configuration objects, each with 'name', 'aggregation_method', and 'rollup_method'), source (object with name, alias, dataset_id), id (string UUID), meta_data (object with metricTypes), metricMetadata (object with attributes, config, label, name, resource, type), and key (string). You MUST provide the builder_id (widget ID) of the widget you want to update.

IMPORTANT - Source Name (Resource Type):
- The 'source.name' field in builderConfig MUST be a resource type that is supported by Middleware and returned by the get_resources tool
- You MUST first call the get_resources tool to discover available resource types in your environment
- You can ONLY use resource type names that are returned by the get_resources tool
- Do not use arbitrary or guessed resource names - only use the exact resource type names returned by get_resources
- Examples of valid source.name values (if returned by get_resources): 'host', 'container', 'log', 'trace', 'k8s.pod', 'database', 'service', etc.
- The source.name identifies which resource type the widget will query data from, and it must match a resource type that Middleware supports and has data for
IMPORTANT - Builder ID (Widget ID):
- The builder_id field is REQUIRED for updating a widget.
- The builder_id is the widget ID of the widget that needs to be updated.
- This is the unique identifier of the widget you want to update.
- You can get the builder_id (widget ID) from the list_widgets tool or from the widget creation response.
IMPORTANT - Layout Requirements:
- Based on the widget type, you MUST set proper layout. Width (w) must be minimum 4 (this is a strict minimum requirement) and height (h) must be minimum 6 (this is a strict minimum requirement). The layout dimensions should be appropriate for the widget type to ensure proper visualization.
`),
		mcp.WithInputSchema[UpdateWidgetInput](),
	)
}

type UpdateWidgetInput struct {
	BuilderID         int                      `json:"builder_id" jsonschema:"The widget ID (builder ID) of the widget that needs to be updated,required"`
	Label             string                   `json:"label,omitempty" jsonschema:"The display name for the widget (e.g., 'CPU Usage', 'Error Rate')"`
	WidgetType        string                   `json:"widget_type,omitempty" jsonschema:"The type of chart/widget,enum=time_series_chart|bar_chart|data_table|query_value|pie_chart|scatter_plot|count_chart|tree_chart|top_list_chart|heatmap_chart|hexagon_chart"`
	Key               string                   `json:"key,omitempty" jsonschema:"Optional unique key identifier for the widget"`
	Description       string                   `json:"description,omitempty" jsonschema:"Optional description explaining what the widget displays"`
	BuilderConfig     []BuilderConfigItemInput `json:"builderConfig,omitempty" jsonschema:"Widget configuration array containing queries, chart type, display settings, and data sources. Each item should have: columns, source, id, meta_data, metricMetadata, key, group_by, and filter_with"`
	ReportID          int                      `json:"report_id,omitempty" jsonschema:"The numeric ID of the dashboard ID (Report ID) where this widget belongs"`
	ReportKey         string                   `json:"report_key,omitempty" jsonschema:"The unique key identifier of the dashboard (report) where this widget belongs"`
	ReportName        string                   `json:"report_name,omitempty" jsonschema:"The name of the dashboard (report) where this widget belongs"`
	ReportDescription string                   `json:"report_description,omitempty" jsonschema:"Optional description of the dashboard (report)"`
	ReportMetadata    any                      `json:"report_metadata,omitempty" jsonschema:"Optional metadata for the dashboard (report)"`
	DisableUserEdit   bool                     `json:"disable_user_edit,omitempty" jsonschema:"Whether to disable user editing of the widget (default: false)"`
	Layout            *LayoutItemInput         `json:"layout" jsonschema:"Layout for the widget including coordinates and size. Based on the widget type, you MUST set proper layout. Width (w) must be minimum 4 (strict minimum requirement) and height (h) must be minimum 6 (strict minimum requirement),required"`
}

func HandleUpdateWidget(s ServerInterface, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := ParseInput[UpdateWidgetInput](req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	if input.BuilderID <= 0 {
		return nil, fmt.Errorf("builder_id is required for updating a widget")
	}

	var builderViewOptions *middleware.BuilderViewOptions
	if input.ReportID > 0 || input.ReportKey != "" || input.ReportName != "" {
		builderViewOptions = &middleware.BuilderViewOptions{
			DisableUserEdit: input.DisableUserEdit,
		}

		if input.ReportID > 0 || input.ReportKey != "" || input.ReportName != "" {
			builderViewOptions.Report = &middleware.ReportView{
				ReportID:          input.ReportID,
				ReportKey:         input.ReportKey,
				ReportName:        input.ReportName,
				ReportDescription: input.ReportDescription,
				Metadata:          input.ReportMetadata,
			}
		}
	}

	widgetKey := input.Key
	if widgetKey == "" && input.Label != "" {
		widgetKey = generateWidgetKey(input.Label)
	}

	var builderConfig []middleware.BuilderConfigItem
	if len(input.BuilderConfig) > 0 {
		builderConfig = convertToMiddlewareBuilderConfig(input.BuilderConfig)
	}

	widget := &middleware.CustomWidget{
		BuilderID:          input.BuilderID,
		Label:              input.Label,
		Key:                widgetKey,
		Description:        input.Description,
		BuilderConfig:      builderConfig,
		BuilderViewOptions: builderViewOptions,
	}

	// Set widget type if provided
	if input.WidgetType != "" {
		widget.WidgetAppID = getWidgetAppID(input.WidgetType)
	}

	// Set layout if provided
	if input.Layout != nil {
		defaultLayout := middleware.LayoutItem{
			X: 0,
			Y: 0,
			W: 4,
			H: 6,
		}

		layout := &middleware.LayoutItem{
			X: defaultLayout.X,
			Y: defaultLayout.Y,
			W: input.Layout.W,
			H: input.Layout.H,
		}

		if layout.W < defaultLayout.W {
			layout.W = defaultLayout.W
		}
		if layout.H < defaultLayout.H {
			layout.H = defaultLayout.H
		}

		widget.Layout = layout
	}

	result, err := s.Client().UpdateWidget(ctx, widget)
	if err != nil {
		return nil, fmt.Errorf("failed to update widget: %w", err)
	}

	return ToTextResult(result)
}

func NewDeleteWidgetTool() mcp.Tool {
	return mcp.NewTool(
		"delete_widget",
		mcp.WithDescription(`Permanently delete a widget from a dashboard.
	
This tool removes a widget (chart, graph, table) from its dashboard. Warning: This action cannot be undone. The widget configuration and data will be permanently deleted.`),
		mcp.WithInputSchema[DeleteWidgetInput](),
	)
}

type DeleteWidgetInput struct {
	BuilderID   int    `json:"builder_id" jsonschema:"The numeric builder ID of the widget to delete permanently,required"`
	Message     string `json:"message" jsonschema:"Message to know which widget is being deleted."`
	WidgetLabel string `json:"widget_label" jsonschema:"Label of the widget to delete."`
}

func HandleDeleteWidget(s ServerInterface, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := ParseInput[DeleteWidgetInput](req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	err = s.Client().DeleteWidget(ctx, input.BuilderID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete widget: %w", err)
	}

	return ToTextResult(map[string]any{"success": true, "message": "Widget deleted successfully"})
}

func NewGetWidgetDataTool() mcp.Tool {
	return mcp.NewTool(
		"get_widget_data",
		mcp.WithDescription(`Fetch the actual data and metrics displayed by a specific widget.
	
This tool executes the widget's query and returns the visualization data (time series, metrics, logs, traces). Use this to get the current values shown in a widget, analyze trends, or export widget data. The data format depends on the widget type (timeseries, table, single value, etc.).`),
		mcp.WithInputSchema[GetWidgetDataInput](),
	)
}

type GetWidgetDataInput struct {
	BuilderID     int                      `json:"builder_id,omitempty" jsonschema:"The numeric builder ID of the widget to fetch data for"`
	Key           string                   `json:"key,omitempty" jsonschema:"Alternative to builder_id: the unique key identifier of the widget"`
	Label         string                   `json:"label,omitempty" jsonschema:"Alternative to builder_id: the label of the widget"`
	BuilderConfig []BuilderConfigItemInput `json:"builder_config,omitempty" jsonschema:"Widget configuration array containing the query and data source settings. Each item's columns MUST be an object with name, aggregation_method, and rollup_method."`
	UseV2         bool                     `json:"use_v2,omitempty" jsonschema:"Set to true to use the newer v2 data format (default: false)"`
}

func HandleGetWidgetData(s ServerInterface, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := ParseInput[GetWidgetDataInput](req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	var builderConfig []middleware.BuilderConfigItem
	if len(input.BuilderConfig) > 0 {
		builderConfig = convertToMiddlewareBuilderConfig(input.BuilderConfig)
	}

	widget := &middleware.CustomWidget{
		BuilderID:     input.BuilderID,
		Key:           input.Key,
		Label:         input.Label,
		BuilderConfig: builderConfig,
		UseV2:         input.UseV2,
	}

	result, err := s.Client().GetWidgetData(ctx, widget)
	if err != nil {
		return nil, fmt.Errorf("failed to get widget data: %w", err)
	}

	return ToTextResult(result)
}

func NewGetMultiWidgetDataTool() mcp.Tool {
	return mcp.NewTool(
		"get_multi_widget_data",
		mcp.WithDescription(`Fetch data for multiple widgets simultaneously in a single request.
	
This tool is optimized for loading data for multiple widgets at once, such as when refreshing an entire dashboard. It's more efficient than calling get_widget_data multiple times. Returns data for all requested widgets in a single response.`),
		mcp.WithInputSchema[GetMultiWidgetDataInput](),
	)
}

type GetMultiWidgetDataInput struct {
	Widgets []WidgetDataRequest `json:"widgets" jsonschema:"Array of widget specifications to fetch data for. Each widget can be identified by builder_id, key, or label,required"`
}

type WidgetDataRequest struct {
	BuilderID     int                      `json:"builder_id,omitempty" jsonschema:"The numeric builder ID of the widget"`
	Key           string                   `json:"key,omitempty" jsonschema:"The unique key identifier of the widget"`
	Label         string                   `json:"label,omitempty" jsonschema:"The label of the widget"`
	BuilderConfig []BuilderConfigItemInput `json:"builder_config,omitempty" jsonschema:"Widget configuration array containing query and display settings. Each item's columns MUST be an object with name, aggregation_method, and rollup_method."`
	UseV2         bool                     `json:"use_v2,omitempty" jsonschema:"Use v2 data format (default: false)"`
}

func HandleGetMultiWidgetData(s ServerInterface, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := ParseInput[GetMultiWidgetDataInput](req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	widgets := make([]middleware.CustomWidget, len(input.Widgets))
	for i, w := range input.Widgets {
		var builderConfig []middleware.BuilderConfigItem
		if len(w.BuilderConfig) > 0 {
			builderConfig = convertToMiddlewareBuilderConfig(w.BuilderConfig)
		}

		widgets[i] = middleware.CustomWidget{
			BuilderID:     w.BuilderID,
			Key:           w.Key,
			Label:         w.Label,
			BuilderConfig: builderConfig,
			UseV2:         w.UseV2,
		}
	}

	result, err := s.Client().GetMultiWidgetData(ctx, widgets)
	if err != nil {
		return nil, fmt.Errorf("failed to get multi widget data: %w", err)
	}

	// Return as JSON text only (no structuredContent)
	return ToTextResult(map[string]any{"widgets": result})
}

func NewUpdateWidgetLayoutsTool() mcp.Tool {
	return mcp.NewTool(
		"update_widget_layouts",
		mcp.WithDescription(`Update the position and size of widgets on a dashboard.
	
This tool modifies the layout (position, size) of multiple widgets on a dashboard. Use this to rearrange widgets, resize them, or optimize dashboard layout. The dashboard uses a grid system where x,y represent position and w,h represent size in grid units. IMPORTANT: Based on the widget type, you MUST set proper layout. Width (w) must be minimum 4 (this is a strict minimum requirement) and height (h) must be minimum 6 (this is a strict minimum requirement). The layout dimensions should be appropriate for the widget type to ensure proper visualization.`),
		mcp.WithInputSchema[UpdateWidgetLayoutsInput](),
	)
}

type UpdateWidgetLayoutsInput struct {
	Layouts          []LayoutItemInput `json:"layouts" jsonschema:"Array of layout specifications for each widget. Each item defines position and size in the dashboard grid. Based on the widget type, you MUST set proper layout. Width (w) must be minimum 4 (strict minimum requirement) and height (h) must be minimum 6 (strict minimum requirement),required"`
	Message          string            `json:"message" jsonschema:"Message to know which widgets are being updated. Length should be less than 100 characters."`
	OperationMessage string            `json:"operation_message" jsonschema:"Message to know the operation being completed. Example: 'Updating widget CPU Usage layouts successfully' Length should be less than 100 characters."`
}

type LayoutItemInput struct {
	X       int `json:"x" jsonschema:"Horizontal position in the grid (0-based index from left)"`
	Y       int `json:"y" jsonschema:"Vertical position in the grid (0-based index from top)"`
	W       int `json:"w" jsonschema:"Width in grid units between 4 and 12 minimum size is 4"`
	H       int `json:"h" jsonschema:"Height in grid units between 6 and 12 minimum size is 6"`
	ScopeID int `json:"scope_id,omitempty" jsonschema:"The scope ID of the widget to update layout for"`
}

func HandleUpdateWidgetLayouts(s ServerInterface, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := ParseInput[UpdateWidgetLayoutsInput](req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

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

	err = s.Client().UpdateWidgetLayouts(ctx, layoutReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update widget layouts: %w", err)
	}

	return ToTextResult(map[string]any{"success": true, "message": input.OperationMessage})
}
