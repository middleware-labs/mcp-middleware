package tools

import (
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

func ToMap(v any) (map[string]any, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return result, nil
}

func ToTextResult(v any) (*mcp.CallToolResult, error) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

// ParseInput parses the arguments from a CallToolRequest into the target struct
func ParseInput[T any](req mcp.CallToolRequest) (T, error) {
	var input T

	// Marshal arguments to JSON and unmarshal into target type
	argsJSON, err := json.Marshal(req.Params.Arguments)
	if err != nil {
		return input, fmt.Errorf("failed to marshal arguments: %w", err)
	}

	if err := json.Unmarshal(argsJSON, &input); err != nil {
		return input, fmt.Errorf("failed to unmarshal arguments: %w", err)
	}

	return input, nil
}
