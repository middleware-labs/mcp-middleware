package tools

import (
	"context"
	"fmt"

	"mcp-middleware/middleware"

	"github.com/mark3labs/mcp-go/mcp"
)

func NewListAlertsTool() mcp.Tool {
	return mcp.NewTool(
		"list_alerts",
		mcp.WithDescription(`Get a list of triggered alerts for a specific alert rule with pagination and sorting.	
This tool retrieves all alert instances that have been triggered for a specific alert rule. Each alert instance represents a time when the alert condition was met. Use this to review alert history, analyze alert patterns, or investigate recent incidents. Results can be paginated and ordered by various fields.`),
		mcp.WithInputSchema[ListAlertsInput](),
	)
}

type ListAlertsInput struct {
	RuleID  int    `json:"rule_id" jsonschema:"The numeric ID of the alert rule to fetch alerts for,required"`
	Page    int    `json:"page,omitempty" jsonschema:"Page number for pagination. 0-based index (default: 0 for first page)"`
	OrderBy string `json:"order_by,omitempty" jsonschema:"Field name to sort results by (e.g., 'created_at', 'triggered_at', 'status'). Default: 'created_at' in descending order"`
}

func HandleListAlerts(s ServerInterface, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := ParseInput[ListAlertsInput](req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	params := &middleware.GetAlertsParams{
		Page:    input.Page,
		OrderBy: input.OrderBy,
	}

	result, err := s.Client().GetAlerts(ctx, input.RuleID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts: %w", err)
	}

	return ToTextResult(result)
}

func NewCreateAlertTool() mcp.Tool {
	return mcp.NewTool(
		"create_alert",
		mcp.WithDescription(`Manually create a new alert instance for a specific alert rule.
	
This tool allows you to programmatically create alert instances, typically used for custom alerting logic or integrations with external monitoring systems. The alert will be associated with an existing alert rule and will appear in the alerts list and trigger configured notification channels.

Note: In most cases, alerts are automatically created when rule conditions are met. Use this tool for custom alerting workflows or manual alert creation.`),
		mcp.WithInputSchema[CreateAlertInput](),
	)
}

type CreateAlertInput struct {
	RuleID      int               `json:"rule_id" jsonschema:"The numeric ID of the alert rule this alert instance belongs to,required"`
	Title       string            `json:"title" jsonschema:"Alert title/summary describing what triggered (e.g., 'High CPU Usage on prod-server-01'),required"`
	Message     string            `json:"message,omitempty" jsonschema:"Detailed alert message with additional context and information"`
	Status      int               `json:"status" jsonschema:"Alert status code. Typically: 0=OK/Resolved, 1=Warning, 2=Critical, 3=Unknown,required"`
	Value       float64           `json:"value,omitempty" jsonschema:"The actual measured value that triggered the alert (e.g., 95.5 for 95.5% CPU usage)"`
	Threshold   float64           `json:"threshold,omitempty" jsonschema:"The threshold value that was exceeded (e.g., 80.0 for 80% threshold)"`
	Operator    string            `json:"operator,omitempty" jsonschema:"Comparison operator used (e.g., '>', '<', '>=', '<=', '==', '!=')"`
	Unit        string            `json:"unit,omitempty" jsonschema:"Unit of measurement for the value (e.g., 'percent', 'ms', 'requests/sec', 'GB')"`
	Attributes  map[string]string `json:"attributes,omitempty" jsonschema:"Additional key-value pairs with context (e.g., {'hostname': 'prod-01', 'region': 'us-east-1'})"`
	ProjectUID  string            `json:"project_uid,omitempty" jsonschema:"Project unique identifier if alert is project-specific"`
	ExecutorID  int               `json:"executor_id,omitempty" jsonschema:"ID of the executor/rule evaluator that triggered the alert"`
	TriggeredAt string            `json:"triggered_at,omitempty" jsonschema:"Timestamp when the alert was triggered (ISO 8601 format, e.g., '2024-01-15T10:30:00Z')"`
}

func HandleCreateAlert(s ServerInterface, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := ParseInput[CreateAlertInput](req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	alert := &middleware.NewAlert{
		RuleID:      input.RuleID,
		Title:       input.Title,
		Message:     input.Message,
		Status:      input.Status,
		Value:       input.Value,
		Threshold:   input.Threshold,
		Operator:    input.Operator,
		Unit:        input.Unit,
		Attributes:  input.Attributes,
		ProjectUID:  input.ProjectUID,
		ExecutorID:  input.ExecutorID,
		TriggeredAt: input.TriggeredAt,
	}

	result, err := s.Client().CreateAlert(ctx, input.RuleID, alert)
	if err != nil {
		return nil, fmt.Errorf("failed to create alert: %w", err)
	}

	return ToTextResult(result)
}

func NewGetAlertStatsTool() mcp.Tool {
	return mcp.NewTool(
		"get_alert_stats",
		mcp.WithDescription(`Get aggregated statistics and metrics for alerts of a specific rule.
	
This tool provides statistical analysis of alert history including counts by status (OK, Warning, Critical), counts by alert title, and time series data showing alert trends over time. Use this to understand alert patterns, identify frequently triggered alerts, and analyze alert distribution.

Returns:
- Count by status: Number of alerts in each status (OK, Warning, Critical)
- Count by title: Distribution of alerts by their titles
- Timeseries by title: Historical alert counts over time grouped by title`),
		mcp.WithInputSchema[GetAlertStatsInput](),
	)
}

type GetAlertStatsInput struct {
	RuleID int `json:"rule_id" jsonschema:"The numeric ID of the alert rule to fetch statistics for,required"`
}

func HandleGetAlertStats(s ServerInterface, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := ParseInput[GetAlertStatsInput](req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	result, err := s.Client().GetAlertStats(ctx, input.RuleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert stats: %w", err)
	}

	return ToTextResult(result)
}
