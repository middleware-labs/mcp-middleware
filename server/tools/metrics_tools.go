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

	// Return as JSON text only (no structuredContent)
	return ToTextResult(map[string]any{"resources": result})
}
