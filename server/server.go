package server

import (
	"context"

	"mcp-middleware/config"
	"mcp-middleware/middleware"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server wraps the MCP server with Middleware client
type Server struct {
	mcpServer *mcp.Server
	client    *middleware.Client
	config    *config.Config
}

// New creates a new MCP server with all tools registered
func New(cfg *config.Config) *Server {
	client := middleware.NewClient(cfg.MiddlewareBaseURL, cfg.MiddlewareAPIKey)

	impl := &mcp.Implementation{
		Name:    "middleware-mcp-server",
		Version: "1.0.0",
	}
	mcpServer := mcp.NewServer(impl, nil)

	s := &Server{
		mcpServer: mcpServer,
		client:    client,
		config:    cfg,
	}

	// Register all MCP features
	s.registerTools()
	s.registerResources()
	s.registerPrompts()

	return s
}

// Run starts the MCP server with the given transport
func (s *Server) Run(ctx context.Context, transport mcp.Transport) error {
	return s.mcpServer.Run(ctx, transport)
}

// Client returns the middleware client for use by tool handlers
func (s *Server) Client() *middleware.Client {
	return s.client
}
