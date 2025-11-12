package server

// registerResources registers all available MCP resources with the server.
// Resources provide structured access to information that AI applications can retrieve
// and provide to models as context (read-only data sources).
// See: https://modelcontextprotocol.io/docs/learn/server-concepts#resources
//
// TODO: Implement resource registration when needed
// Example resources could include:
// - Dashboard configurations as resources
// - Widget templates as resources
// - Metric definitions as resources
// - Alert rule templates as resources
func (s *Server) registerResources() {
	// Resources will be registered here in the future
	// Example:
	// import "github.com/modelcontextprotocol/go-sdk/mcp"
	// mcp.AddResource(s.mcpServer, &mcp.Resource{
	// 	URI:         "middleware://dashboards",
	// 	Name:        "dashboards",
	// 	Description: "List of all available dashboards",
	// 	MimeType:    "application/json",
	// })
}
