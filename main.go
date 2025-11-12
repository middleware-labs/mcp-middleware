package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"mcp-middleware/config"
	"mcp-middleware/server"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	srv := server.New(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down gracefully...")
		cancel()
	}()

	log.Printf("Middleware MCP Server v1.0.0")
	log.Printf("Connected to: %s", cfg.MiddlewareBaseURL)
	if len(cfg.ExcludedTools) > 0 {
		log.Printf("Excluded tools: %v", getExcludedToolsList(cfg))
	}

	switch cfg.AppMode {
	case "stdio":
		transport := &mcp.StdioTransport{}
		log.Println("Starting MCP server in stdio mode...")
		if err := srv.Run(ctx, transport); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	case "http":
		if err := srv.RunHTTPMode(ctx, cfg); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	case "sse":
		if err := srv.RunSSEMode(ctx, cfg); err != nil {
			log.Fatalf("SSE server error: %v", err)
		}
	default:
		log.Fatalf("Invalid APP_MODE: %s", cfg.AppMode)
	}
}

func getExcludedToolsList(cfg *config.Config) []string {
	tools := make([]string, 0, len(cfg.ExcludedTools))
	for tool := range cfg.ExcludedTools {
		tools = append(tools, tool)
	}
	return tools
}
