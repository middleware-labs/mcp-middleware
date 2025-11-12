package server

// registerPrompts registers all available MCP prompts with the server.
// Prompts are pre-built instruction templates that tell the model to work with
// specific tools and resources. They are user-controlled and require explicit invocation.
// See: https://modelcontextprotocol.io/docs/learn/server-concepts#prompts
//
// TODO: Implement prompt registration when needed
// Example prompts could include:
// - "Create a monitoring dashboard" - guides through dashboard creation workflow
// - "Analyze alert trends" - provides template for alert analysis
// - "Set up widget monitoring" - step-by-step widget configuration guide
func (s *Server) registerPrompts() {
	// Prompts will be registered here in the future
	// Example:
	// import "github.com/modelcontextprotocol/go-sdk/mcp"
	// mcp.AddPrompt(s.mcpServer, &mcp.Prompt{
	// 	Name:        "create-dashboard",
	// 	Description: "Guide through creating a new monitoring dashboard",
	// 	Arguments: []mcp.PromptArgument{
	// 		{Name: "dashboard_name", Description: "Name for the new dashboard", Required: true},
	// 		{Name: "widgets", Description: "List of widgets to include", Required: false},
	// 	},
	// })
}
