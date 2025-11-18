package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"mcp-middleware/config"
	"mcp-middleware/middleware"

	"github.com/mark3labs/mcp-go/server"
)

type Server struct {
	mcpServer *server.MCPServer
	client    *middleware.Client
	config    *config.Config
}

func New(cfg *config.Config) *Server {
	client := middleware.NewClient(cfg.MiddlewareBaseURL, cfg.MiddlewareAPIKey)

	mcpServer := server.NewMCPServer("middleware-mcp-server", "1.0.0")

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

func (s *Server) Client() *middleware.Client {
	return s.client
}

func (s *Server) GetMCPServer() *server.MCPServer {
	return s.mcpServer
}

func (s *Server) RunHTTPMode(ctx context.Context, cfg *config.Config) error {
	// Create streamable HTTP server
	httpServer := server.NewStreamableHTTPServer(s.mcpServer)

	// Create HTTP server
	addr := fmt.Sprintf("%s:%s", cfg.AppHost, cfg.AppPort)
	httpSrv := &http.Server{
		Addr:         addr,
		Handler:      httpServer,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Starting MCP server in HTTP mode on %s", addr)
		log.Printf("Server ready. Connect to: http://%s", addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
		log.Println("Shutting down HTTP server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpSrv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("error shutting down HTTP server: %w", err)
		}
		log.Println("HTTP server stopped")
		return nil
	case err := <-serverErr:
		return fmt.Errorf("HTTP server error: %w", err)
	}
}

func (s *Server) RunSSEMode(ctx context.Context, cfg *config.Config) error {
	// Create SSE server
	sseServer := server.NewSSEServer(s.mcpServer)

	// Create HTTP server
	addr := fmt.Sprintf("%s:%s", cfg.AppHost, cfg.AppPort)
	httpSrv := &http.Server{
		Addr:         addr,
		Handler:      sseServer,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 0, // No timeout for SSE (long-lived connection)
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Starting MCP server in SSE mode on %s", addr)
		log.Printf("Server ready. Connect to: http://%s", addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
		log.Println("Shutting down SSE server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpSrv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("error shutting down SSE server: %w", err)
		}
		log.Println("SSE server stopped")
		return nil
	case err := <-serverErr:
		return fmt.Errorf("SSE server error: %w", err)
	}
}

func (s *Server) RunStdioMode(ctx context.Context) error {
	stdioServer := server.NewStdioServer(s.mcpServer)
	return stdioServer.Listen(ctx, os.Stdin, os.Stdout)
}
