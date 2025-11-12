package tools

import (
	"context"
	"fmt"

	"mcp-middleware/middleware"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var GetMetricsTool = &mcp.Tool{
	Name: "get_metrics",
	Description: `Get a list of available metrics, filters, or groupby tags for building monitoring queries.
	
This tool retrieves metadata about available metrics, filters, and grouping options in your Middleware.io environment. Use this to discover what metrics are available for a specific resource (e.g., host CPU metrics, database query metrics), what filters can be applied (e.g., filter by service name, environment), or what grouping dimensions are available (e.g., group by region, pod name).

Common use cases:
- Find all metrics for a specific resource type (host, container, database)
- Discover available filters for narrowing down data
- Get groupby tags for aggregating metrics by dimension`,
}

type GetMetricsInput struct {
	DataType                string   `json:"data_type" jsonschema:"Type of data to fetch. Valid values: 'metrics' (metric names), 'filters' (filter dimensions), 'groupby' (grouping tags),required"`
	WidgetType              string   `json:"widget_type" jsonschema:"Widget type for the query. Valid values: 'timeseries', 'list', 'queryValue',required"`
	KpiType                 int      `json:"kpi_type,omitempty" jsonschema:"Single KPI type filter. 1=Metric (infrastructure/APM metrics), 2=Log (log data), 3=Trace (distributed tracing data)"`
	KpiTypes                []int    `json:"kpi_types,omitempty" jsonschema:"Array of KPI types to include. Use this for multi-type queries"`
	Resource                string   `json:"resource,omitempty" jsonschema:"Resource name to filter by (e.g., 'host', 'process', 'container', 'database')"`
	Resources               []string `json:"resources,omitempty" jsonschema:"Array of resource names for multi-resource queries"`
	Metric                  string   `json:"metric,omitempty" jsonschema:"Specific metric name (required when dataType is 'groupby' to get groupby tags for this metric)"`
	Page                    int      `json:"page,omitempty" jsonschema:"Page number for paginated results (default: 1)"`
	Limit                   int      `json:"limit,omitempty" jsonschema:"Number of items per page (default: 100, max: varies by data type)"`
	Search                  string   `json:"search,omitempty" jsonschema:"Search term to filter metrics or resources by name (case-insensitive substring match)"`
	ExcludeMetrics          []string `json:"exclude_metrics,omitempty" jsonschema:"Array of metric names to exclude from results"`
	MandatoryMetrics        []string `json:"mandatory_metrics,omitempty" jsonschema:"Array of metric names to always include at the top of results"`
	ExcludeFilters          []string `json:"exclude_filters,omitempty" jsonschema:"Array of filter names to exclude from results"`
	MandatoryFilters        []string `json:"mandatory_filters,omitempty" jsonschema:"Array of filter names to always include at the top of results"`
	FilterTypes             []int    `json:"filter_types,omitempty" jsonschema:"Array of filter type IDs to include (filters metadata types)"`
	ReturnOnlyMandatoryData bool     `json:"return_only_mandatory_data,omitempty" jsonschema:"Set to true to return only mandatory metrics/filters without additional data"`
}

func HandleGetMetrics(s ServerInterface, ctx context.Context, req *mcp.CallToolRequest, input GetMetricsInput) (*mcp.CallToolResult, map[string]any, error) {
	metricsReq := &middleware.MetricsV2Request{
		DataType:                input.DataType,
		WidgetType:              input.WidgetType,
		KpiType:                 input.KpiType,
		KpiTypes:                input.KpiTypes,
		Resource:                input.Resource,
		Resources:               input.Resources,
		Metric:                  input.Metric,
		Page:                    input.Page,
		Limit:                   input.Limit,
		Search:                  input.Search,
		ExcludeMetrics:          input.ExcludeMetrics,
		MandatoryMetrics:        input.MandatoryMetrics,
		ExcludeFilters:          input.ExcludeFilters,
		MandatoryFilters:        input.MandatoryFilters,
		FilterTypes:             input.FilterTypes,
		ReturnOnlyMandatoryData: input.ReturnOnlyMandatoryData,
	}

	result, err := s.Client().GetMetrics(ctx, metricsReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	return ToTextResult(result)
}

var GetResourcesTool = &mcp.Tool{
	Name: "get_resources",
	Description: `Get a list of all available resource types in your Middleware.io environment.
	
This tool returns all resource types that have data in your monitoring environment. Resources represent the entities you're monitoring (e.g., hosts, containers, databases, services, processes). Use this to discover what resource types are available before querying metrics for specific resources.

Example resources: host, container, pod, service, database, redis, mongodb, postgresql, mysql, nginx, etc.`,
}

type GetResourcesInput struct{}

func HandleGetResources(s ServerInterface, ctx context.Context, req *mcp.CallToolRequest, input GetResourcesInput) (*mcp.CallToolResult, map[string]any, error) {
	result, err := s.Client().GetResources(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get resources: %w", err)
	}

	return ToTextResult(map[string]any{"resources": result})
}

var QueryTool = &mcp.Tool{
	Name: "query",
	Description: `Execute a flexible query to retrieve logs, metrics, traces, and other data from Middleware.io.
	
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
- Query multiple data types in a single request`,
}

type QueryInput struct {
	Queries []QueryInputItem `json:"queries" jsonschema:"Array of query objects to execute. Each query can target different resources and data types,required"`
}

type QueryInputItem struct {
	ChartType string         `json:"chartType" jsonschema:"Type of chart/visualization. Common values: 'data_table', 'timeseries', 'bar_chart', 'pie_chart', 'line_chart',required"`
	Columns   []string       `json:"columns" jsonschema:"Array of column names to retrieve. For logs: ['body', 'timestamp', 'level']. For metrics: metric names. For resources: attribute names,required"`
	Resources []string       `json:"resources" jsonschema:"Array of resource types to query. IMPORTANT: For logs, always use ['log']. For other data types (metrics, traces, etc.), FIRST use get_resources tool to discover available resources, THEN use those resource types here. Examples: ['log'] for logs, ['container'] for container data (discovered via get_resources), ['host'] for host data, ['trace'] for traces, ['k8s.pod'] for Kubernetes pods,required"`
	TimeRange QueryTimeRange `json:"timeRange" jsonschema:"Time range for the query with from and to timestamps in milliseconds,required"`
	Filters   map[string]any `json:"filters,omitempty" jsonschema:"Optional filters to apply. Format: {\"field.name\": {\"=\": \"value\"}} or {\"field.name\": {\"!=\": \"value\"}}"`
	GroupBy   []string       `json:"groupBy,omitempty" jsonschema:"Optional array of field names to group results by (e.g., ['container.id', 'service.name'])"`
}

type QueryTimeRange struct {
	From int64 `json:"from" jsonschema:"Start timestamp in milliseconds (Unix timestamp * 1000),required"`
	To   int64 `json:"to" jsonschema:"End timestamp in milliseconds (Unix timestamp * 1000),required"`
}

func HandleQuery(s ServerInterface, ctx context.Context, req *mcp.CallToolRequest, input QueryInput) (*mcp.CallToolResult, map[string]any, error) {
	queries := make([]middleware.Query, len(input.Queries))
	for i, q := range input.Queries {
		queries[i] = middleware.Query{
			ChartType: q.ChartType,
			Columns:   q.Columns,
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
		return nil, nil, fmt.Errorf("failed to execute query: %w", err)
	}

	if result == nil {
		return nil, nil, fmt.Errorf("query returned nil result")
	}

	return ToTextResult(result)
}
