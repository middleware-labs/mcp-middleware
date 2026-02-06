package tools

import (
	"context"
	"fmt"

	"mcp-middleware/middleware"

	"github.com/mark3labs/mcp-go/mcp"
)

type WidgetType string

const (
	WidgetTypeTimeseries WidgetType = "timeseries"
	WidgetTypeList       WidgetType = "list"
	WidgetTypeQueryValue WidgetType = "queryValue"
)

type MetricsDataType string

const (
	MetricsDataTypeMetrics MetricsDataType = "metrics"
	MetricsDataTypeFilters MetricsDataType = "filters"
	MetricsDataTypeGroupby MetricsDataType = "groupby"
)

// ChartType matches widget type keys used by getWidgetAppID in widgets_tools.go.
// Pass this string value in queries; backend expects the same chart type names.
type ChartType string

const (
	ChartTypeTimeSeries  ChartType = "time_series_chart"
	ChartTypeBar         ChartType = "bar_chart"
	ChartTypePie         ChartType = "pie_chart"
	ChartTypeScatter     ChartType = "scatter_plot"
	ChartTypeDataTable   ChartType = "data_table"
	ChartTypeCount       ChartType = "count_chart"
	ChartTypeTree        ChartType = "tree_chart"
	ChartTypeTopList     ChartType = "top_list_chart"
	ChartTypeHeatmap     ChartType = "heatmap_chart"
	ChartTypeHexagon     ChartType = "hexagon_chart"
	ChartTypeQueryValue  ChartType = "query_value"
)

func NewGetMetricsTool() mcp.Tool {
	return mcp.NewTool(
		"get_metrics",
		mcp.WithDescription(`Get a list of available metrics, filters, or groupby tags for building monitoring queries.
	
This tool is essential for discovering the metadata needed to construct accurate queries. Since every metric supports different filters and grouping dimensions (groupby tags), you must use this tool to validate what is available for each specific metric before querying data.

Discovery Workflow:
1. First, identify available resources using the 'get_resources' tool.
2. Find available metrics for a resource: Use data_type='metrics' with the 'resources' parameter.
3. Explore dimensions for a specific metric:
   - To find how you can group data: Use data_type='groupby' AND provide the specific 'metric' name.
   - To find how you can filter data: Use data_type='filters' (optionally with 'resources').

IMPORTANT - Metric-Specific Metadata:
- Filters and GroupBy tags are NOT universal. They vary by metric.
- ALWAYS check 'groupby' options for a specific metric before trying to aggregate by a dimension.
- ALWAYS check 'filters' to see what dimensions are available for narrowing down your search.

Data Type Options:
- 'metrics': List metric names (requires 'resources').
- 'filters': List filter dimensions.
- 'groupby': List grouping tags (requires 'metric').

Resource Selection:
- Use exact resource names returned by 'get_resources'.`),
		mcp.WithInputSchema[GetMetricsInput](),
	)
}

type GetMetricsInput struct {
	DataType   MetricsDataType `json:"data_type" jsonschema:"Type of data to fetch. DataType is the type of data that is being fetched. Must be one of: 'metrics' (metric names), 'filters' (filter dimensions), 'groupby' (grouping tags),required,enum=metrics,enum=filters,enum=groupby"`
	WidgetType WidgetType      `json:"widget_type" jsonschema:"Widget type for the query. Must be one of: 'timeseries' (for timeseries, bar, stackbar, area), 'list' (for table, pie, scatter, tree, toplist, hexagon), or 'queryValue' (for queryvalue),required,enum=timeseries,enum=list,enum=queryValue"`
	// KpiType                 int        `json:"kpi_type,omitempty" jsonschema:"Single KPI type filter. 1=Metric (infrastructure/APM metrics), 2=Log (log data), 3=Trace (distributed tracing data)"`
	// KpiTypes                []int      `json:"kpi_types,omitempty" jsonschema:"Array of KPI types to include. Use this for multi-type queries"`
	// Resource                string   `json:"resource,omitempty" jsonschema:"The resource type name obtained from calling get_resources. This identifies which resource type to filter metrics by and correlates the data source. IMPORTANT: You can ONLY use resource type names that are returned by the get_resources tool. You must first call get_resources to discover available resources, then use only those exact resource type names here. Examples: 'host', 'container', 'log', 'trace', 'k8s.pod', 'database', etc. Always use the exact resource type name returned by get_resources."`
	Resources []string `json:"resources" jsonschema:"Array of resource type names obtained from calling get_resources. Use this for multi-resource queries. IMPORTANT: You can ONLY use resource type names that are returned by the get_resources tool. You must first call get_resources to discover available resources, then use only those exact resource type names here. Each resource name should be the exact resource type name returned by get_resources (e.g., ['host', 'container', 'trace'])."`
	Metric    string   `json:"metric,omitempty" jsonschema:"Specific metric name. IMPORTANT: This field is REQUIRED when data_type is 'groupby'. When fetching groupby tags, you must specify which metric you want to group by to get the available grouping dimensions for that specific metric."`
	Page      int      `json:"page,omitempty" jsonschema:"Page number for paginated results (default: 1)"`
	Limit     int      `json:"limit,omitempty" jsonschema:"Number of items per page (default: 100, max: varies by data type)"`
	Search    string   `json:"search,omitempty" jsonschema:"Search term to filter metrics or resources by name (case-insensitive substring match)"`
	// ExcludeMetrics          []string `json:"exclude_metrics,omitempty" jsonschema:"Array of metric names to exclude from results"`
	// MandatoryMetrics        []string `json:"mandatory_metrics,omitempty" jsonschema:"Array of metric names to always include at the top of results"`
	// ExcludeFilters          []string `json:"exclude_filters,omitempty" jsonschema:"Array of filter names to exclude from results"`
	// MandatoryFilters        []string `json:"mandatory_filters,omitempty" jsonschema:"Array of filter names to always include at the top of results"`
	// FilterTypes             []int    `json:"filter_types,omitempty" jsonschema:"Array of filter type IDs to include (filters metadata types)"`
	// ReturnOnlyMandatoryData bool     `json:"return_only_mandatory_data,omitempty" jsonschema:"Set to true to return only mandatory metrics/filters without additional data"`
}

func HandleGetMetrics(s ServerInterface, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := ParseInput[GetMetricsInput](req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	metricsReq := &middleware.MetricsV2Request{
		DataType:   string(input.DataType),
		WidgetType: string(input.WidgetType),
		// KpiType:                 input.KpiType,
		// KpiTypes:                input.KpiTypes,
		// Resource:                input.Resource,
		Resources: input.Resources,
		Metric:    input.Metric,
		Page:      input.Page,
		Limit:     input.Limit,
		Search:    input.Search,
		// ExcludeMetrics:          input.ExcludeMetrics,
		// MandatoryMetrics:        input.MandatoryMetrics,
		// ExcludeFilters:          input.ExcludeFilters,
		// MandatoryFilters:        input.MandatoryFilters,
		// FilterTypes:             input.FilterTypes,
		// ReturnOnlyMandatoryData: input.ReturnOnlyMandatoryData,
	}

	result, err := s.Client().GetMetrics(ctx, metricsReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	return ToTextResult(result)
}

func NewGetResourcesTool() mcp.Tool {
	return mcp.NewTool(
		"get_resources",
		mcp.WithDescription(`Get a list of all available resource types in your Middleware.io environment.
	
This tool returns all resource types that have data in your monitoring environment. Resources represent the entities you're monitoring (e.g., hosts, containers, databases, services, processes). Use this to discover what resource types are available before querying metrics for specific resources.

Example resources: host, container, pod, service, database, redis, mongodb, postgresql, mysql, nginx, etc.`),
		mcp.WithInputSchema[GetResourcesInput](),
	)
}

type GetResourcesInput struct{}

func HandleGetResources(s ServerInterface, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := ParseInput[GetResourcesInput](req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}
	_ = input // Empty struct, no fields to use

	result, err := s.Client().GetResources(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get resources: %w", err)
	}

	return ToTextResult(map[string]any{"resources": result})
}

func NewQueryTool() mcp.Tool {
	return mcp.NewTool(
		"query",
		mcp.WithDescription(`Execute a flexible query to retrieve logs, metrics, traces, and other data from Middleware.io.
	
This is a powerful tool that allows you to query any type of data from Middleware including logs, metrics, traces, and resource information. You can filter by resource types, time ranges, apply filters, and group results. This tool provides comprehensive access to all your monitoring data.

IMPORTANT - Resource Selection:
- For logs: Always use ["log"] as the resource (no need to check get_resources first)
- For metrics, traces, or other data types: FIRST use the get_resources tool to discover available resource types in your environment, THEN use those specific resource types in this query tool

Workflow for non-log queries:
1. Call get_resources tool to get list of available resources (e.g., ["host", "container", "service", "trace", "k8s.pod", etc.])
2. Use the discovered resource types in this query tool's resources parameter

Use cases:
- Query logs from containers, hosts, or services (use resource: ["log"])
- Retrieve metrics for specific resources (first get resources, then use: ["host"], ["container"], ["service"], etc.)
- Get trace data for distributed systems (first get resources, then use: ["trace"], ["trace.service"], etc.)
- Filter data by any resource attribute
- Group results by dimensions for aggregation
- Query multiple data types in a single request`),
		mcp.WithInputSchema[QueryInput](),
	)
}

type QueryInput struct {
	Queries []QueryInputItem `json:"queries" jsonschema:"Array of query objects to execute. Each query can target different resources and data types,required"`
}

type QueryInputItem struct {
	ChartType ChartType      `json:"chartType" jsonschema:"Type of chart/visualization. Must be one of the supported chart type keys (same as create_widget widget_type),required,enum=time_series_chart,enum=bar_chart,enum=pie_chart,enum=scatter_plot,enum=data_table,enum=count_chart,enum=tree_chart,enum=top_list_chart,enum=heatmap_chart,enum=hexagon_chart,enum=query_value"`
	Columns   []ColumnConfig `json:"columns" jsonschema:"Array of column configs: each has 'name' (metric/attribute name, e.g. 'body', 'timestamp', 'k8s.node.cpu.utilization') and optional 'aggregation_method' (avg, sum, min, max, uniq, count, group) and 'rollup_method' (avg, sum, min, max, none). For logs use name only (e.g. body, timestamp, level). Same format as create_widget columns.,required"`
	Resources []string       `json:"resources" jsonschema:"Array of resource types to query. IMPORTANT: For logs, always use ['log']. For other data types (metrics, traces, etc.), FIRST use get_resources tool to discover available resources, THEN use those resource types here. Examples: ['log'] for logs, ['container'] for container data (discovered via get_resources), ['host'] for host data, ['trace'] for traces, ['k8s.pod'] for Kubernetes pods,required"`
	TimeRange QueryTimeRange `json:"timeRange" jsonschema:"Time range for the query with from and to timestamps in milliseconds,required"`
	Filters   map[string]any `json:"filters,omitempty" jsonschema:"Optional filters to apply. Format: {\"field.name\": {\"=\": \"value\"}} or {\"field.name\": {\"!=\": \"value\"}}"`
	GroupBy   []string       `json:"groupBy,omitempty" jsonschema:"Optional array of field names to group results by (e.g., ['container.id', 'service.name'])"`
}

type QueryTimeRange struct {
	From int64 `json:"from" jsonschema:"Start timestamp in milliseconds (Unix timestamp * 1000),required"`
	To   int64 `json:"to" jsonschema:"End timestamp in milliseconds (Unix timestamp * 1000),required"`
}

func HandleQuery(s ServerInterface, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := ParseInput[QueryInput](req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	queries := make([]middleware.Query, len(input.Queries))
	for i, q := range input.Queries {
		// Convert column configs to request format (same as create_widget: agg(name) or agg(name, value(rollup)))
		columnStrings := transformColumns(q.Columns)
		queries[i] = middleware.Query{
			ChartType: string(q.ChartType),
			Columns:   columnStrings,
			Resources: q.Resources,
			TimeRange: middleware.QueryTimeRange{
				From: q.TimeRange.From,
				To:   q.TimeRange.To,
			},
			Filters: q.Filters,
			GroupBy: q.GroupBy,
		}
	}

	queryReq := &middleware.QueryRequest{
		Queries: queries,
	}

	result, err := s.Client().Query(ctx, queryReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("query returned nil result")
	}

	return ToTextResult(result)
}
