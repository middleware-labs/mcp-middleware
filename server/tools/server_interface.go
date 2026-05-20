package tools

import (
	"context"

	"mcp-middleware/middleware"

	"github.com/mark3labs/mcp-go/mcp"
)

// ServerInterface defines the interface that tool handlers need from the server.
// Client takes a context so it can build a per-request middleware client from the
// Bearer token and tenant base URL placed in ctx by the auth middleware.
type ServerInterface interface {
	Client(ctx context.Context) *middleware.Client
}

// ToolHandler is a function type for tool handlers.
type ToolHandler func(s ServerInterface, ctx context.Context, req mcp.CallToolRequest, input any) (*mcp.CallToolResult, error)
