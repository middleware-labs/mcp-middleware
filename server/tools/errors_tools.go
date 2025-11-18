package tools

import (
	"context"
	"fmt"

	"mcp-middleware/middleware"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var ListErrorsTool = &mcp.Tool{
	Name: "list_errors",
	Description: `List all errors/incidents currently happening in the system.
	
This tool retrieves all error incidents from the Middleware.io system. Use this to monitor system health, identify ongoing issues, and track error patterns. Results can be filtered by time range, status, and search terms, and support pagination.

IMPORTANT: Each error/incident in the response includes an 'issue_url' field that contains a direct, clickable URL link to view the issue details in the Middleware.io web interface. This URL can be used to redirect users to the full issue details page where they can see complete context, occurrence history, related information, and all technical details. The URL format is: https://[base-url]/ops-ai?fingerprint=[fingerprint]. Always include this URL when presenting error information to users so they can easily navigate to view more details.`,
}

type ListErrorsInput struct {
	FromTs int64  `json:"from_ts" jsonschema:"Start timestamp in milliseconds (Unix timestamp * 1000),required"`
	ToTs   int64  `json:"to_ts" jsonschema:"End timestamp in milliseconds (Unix timestamp * 1000),required"`
	Page   int    `json:"page" jsonschema:"Page number for pagination (default: 1),required"`
	Filter string `json:"filter,omitempty" jsonschema:"Optional filter string to narrow down results"`
	Status string `json:"status" jsonschema:"Filter by status,required,enum=all|for_review|resolved|reviewed|ignored"`
	Search string `json:"search,omitempty" jsonschema:"Search term to filter incidents by title or description"`
}

func HandleListErrors(s ServerInterface, ctx context.Context, req *mcp.CallToolRequest, input ListErrorsInput) (*mcp.CallToolResult, map[string]any, error) {
	// Default page to 1 if not provided or 0
	page := input.Page
	if page <= 0 {
		page = 1
	}

	params := &middleware.GetIncidentsParams{
		FromTs: input.FromTs,
		ToTs:   input.ToTs,
		Page:   page,
		Filter: input.Filter,
		Status: input.Status,
		Search: input.Search,
	}

	result, err := s.Client().GetIncidents(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get errors/incidents: %w", err)
	}

	return ToTextResult(result)
}

var GetErrorDetailsTool = &mcp.Tool{
	Name: "get_error_details",
	Description: `Get detailed information about a specific error/incident by its fingerprint.
	
This tool retrieves comprehensive details about a specific error incident from the Middleware.io system. Use this to investigate a particular error, view its full context, occurrence history, and related information.`,
}

type GetErrorDetailsInput struct {
	Fingerprint string `json:"fingerprint" jsonschema:"The unique fingerprint identifier of the error/incident,required"`
	FromTs      int64  `json:"from_ts" jsonschema:"Start timestamp in milliseconds (Unix timestamp * 1000),required"`
	ToTs        int64  `json:"to_ts" jsonschema:"End timestamp in milliseconds (Unix timestamp * 1000),required"`
	Filter      string `json:"filter,omitempty" jsonschema:"Optional filter string to narrow down results"`
}

func HandleGetErrorDetails(s ServerInterface, ctx context.Context, req *mcp.CallToolRequest, input GetErrorDetailsInput) (*mcp.CallToolResult, map[string]any, error) {
	params := &middleware.GetIncidentDetailParams{
		Fingerprint: input.Fingerprint,
		FromTs:      input.FromTs,
		ToTs:        input.ToTs,
		Filter:      input.Filter,
	}

	result, err := s.Client().GetIncidentDetail(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get error details: %w", err)
	}

	return ToTextResult(result)
}
