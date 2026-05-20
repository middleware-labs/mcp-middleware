package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"mcp-middleware/config"
	"mcp-middleware/middleware"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Server struct {
	mcpServer *server.MCPServer
	config    *config.Config
}

func New(cfg *config.Config) *Server {
	mcpServer := server.NewMCPServer(
		"middleware-mcp-server",
		"1.0.0",
		server.WithToolFilter(excludedToolsFilter),
		server.WithToolHandlerMiddleware(excludedToolsCallGuard),
	)

	s := &Server{
		mcpServer: mcpServer,
		config:    cfg,
	}

	s.registerTools()
	s.registerResources()
	s.registerPrompts()

	return s
}

// Client returns a per-request middleware.Client built from the Bearer token and tenant base URL
// that the auth middleware placed in ctx. In stdio mode (no auth middleware) the resulting client
// will fail with a clear error on first call.
func (s *Server) Client(ctx context.Context) *middleware.Client {
	token, _ := TokenFromContext(ctx)
	baseURL, _ := TenantBaseURLFromContext(ctx)
	return middleware.NewClient(baseURL, token)
}

func (s *Server) GetMCPServer() *server.MCPServer {
	return s.mcpServer
}

// excludedToolsFilter drops tools whose names appear in the per-request exclude_tools query param.
func excludedToolsFilter(ctx context.Context, tools []mcp.Tool) []mcp.Tool {
	excluded := ExcludedToolsFromContext(ctx)
	if len(excluded) == 0 {
		return tools
	}
	out := tools[:0]
	for _, t := range tools {
		if _, skip := excluded[t.Name]; skip {
			continue
		}
		out = append(out, t)
	}
	return out
}

// excludedToolsCallGuard rejects tools/call for any tool name in the per-request exclude_tools set.
func excludedToolsCallGuard(next server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		excluded := ExcludedToolsFromContext(ctx)
		if _, skip := excluded[req.Params.Name]; skip {
			return mcp.NewToolResultError(fmt.Sprintf("tool %q is excluded for this connection", req.Params.Name)), nil
		}
		return next(ctx, req)
	}
}

func (s *Server) RunHTTPMode(ctx context.Context, cfg *config.Config) error {
	streamable := server.NewStreamableHTTPServer(s.mcpServer,
		server.WithHTTPContextFunc(injectRequestContext),
	)

	mux := http.NewServeMux()
	metadataHandler := ProtectedResourceMetadataHandler(cfg)
	mux.Handle("/.well-known/oauth-protected-resource", metadataHandler)
	mux.Handle("/.well-known/oauth-protected-resource/mcp", metadataHandler)
	mux.Handle("/mcp", RequireBearer(cfg, streamable))
	mux.Handle("/mcp/", RequireBearer(cfg, streamable))

	addr := fmt.Sprintf("%s:%s", cfg.AppHost, cfg.AppPort)
	httpSrv := &http.Server{
		Addr:        addr,
		Handler:     logRequests(mux),
		ReadTimeout: 30 * time.Second,
		// Streamable HTTP holds SSE streams open — no write timeout.
		WriteTimeout: 0,
		IdleTimeout:  120 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Starting MCP server in HTTP mode on %s", addr)
		log.Printf("MCP endpoint:    %s", cfg.MCPServerURL)
		log.Printf("Auth server:     %s", cfg.AuthServerURL)
		log.Printf("Resource metadata: http://%s/.well-known/oauth-protected-resource/mcp", addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

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
	sseServer := server.NewSSEServer(s.mcpServer,
		server.WithSSEContextFunc(injectRequestContext),
	)

	mux := http.NewServeMux()
	metadataHandler := ProtectedResourceMetadataHandler(cfg)
	mux.Handle("/.well-known/oauth-protected-resource", metadataHandler)
	mux.Handle("/.well-known/oauth-protected-resource/mcp", metadataHandler)
	mux.Handle("/", RequireBearer(cfg, sseServer))

	addr := fmt.Sprintf("%s:%s", cfg.AppHost, cfg.AppPort)
	httpSrv := &http.Server{
		Addr:         addr,
		Handler:      logRequests(mux),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 0,
		IdleTimeout:  120 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Starting MCP server in SSE mode on %s", addr)
		log.Printf("MCP endpoint:    %s", cfg.MCPServerURL)
		log.Printf("Auth server:     %s", cfg.AuthServerURL)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

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

// logRequests logs every incoming HTTP request: the entry line is emitted immediately
// (so long-lived SSE/streamable connections show up as soon as they arrive), and a
// completion line with the status code and duration is emitted when the handler returns.
func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("→ %s %s from %s", r.Method, r.URL.RequestURI(), clientIP(r))

		lw := &loggingResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(lw, r)

		log.Printf("← %s %s %d (%s)", r.Method, r.URL.Path, lw.status, time.Since(start).Round(time.Millisecond))
	})
}

// clientIP returns the originating client address, honoring X-Forwarded-For when set
// by a proxy/ingress (the first entry is the original client).
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := strings.IndexByte(xff, ','); i >= 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}
	return r.RemoteAddr
}

// loggingResponseWriter captures the response status code while transparently
// supporting the streaming interfaces (Flush) the MCP transports rely on.
type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *loggingResponseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Unwrap lets http.ResponseController reach the underlying writer (e.g. for Hijack).
func (w *loggingResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

// injectRequestContext carries the auth-middleware-derived values across the boundary into
// whatever context the MCP library passes to tool handlers.
func injectRequestContext(ctx context.Context, r *http.Request) context.Context {
	if token, ok := TokenFromContext(r.Context()); ok {
		ctx = WithToken(ctx, token)
	}
	if tenantURL, ok := TenantBaseURLFromContext(r.Context()); ok {
		ctx = WithTenantBaseURL(ctx, tenantURL)
	}
	if excluded := ExcludedToolsFromContext(r.Context()); len(excluded) > 0 {
		ctx = WithExcludedTools(ctx, excluded)
	}
	return ctx
}
