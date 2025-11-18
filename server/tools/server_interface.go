package tools

import (
	"context"

	"mcp-middleware/middleware"

	"github.com/mark3labs/mcp-go/mcp"
)

// ServerInterface defines the interface that tool handlers need from the server
type ServerInterface interface {
	Client() *middleware.Client
}

// ToolHandler is a function type for tool handlers
type ToolHandler func(s ServerInterface, ctx context.Context, req mcp.CallToolRequest, input any) (*mcp.CallToolResult, error)
