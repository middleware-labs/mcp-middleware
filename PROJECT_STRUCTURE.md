# Project Structure

This document explains the organization and architecture of the Middleware MCP Server project.

## Directory Layout

```
mcp-middleware/
├── config/                     # Configuration Management
│   └── config.go              # Environment variable loading and validation
│
├── middleware/                 # Middleware.io API Client
│   ├── client.go              # HTTP client with authentication
│   ├── types.go               # API data structures (Dashboard, Widget, Alert, etc.)
│   ├── dashboards.go          # Dashboard API endpoints
│   ├── widgets.go             # Widget API endpoints
│   ├── metrics.go             # Metrics API endpoints
│   └── alerts.go              # Alert API endpoints
│
├── server/                     # MCP Server Implementation
│   ├── server.go              # Server initialization and lifecycle
│   ├── register_tools.go      # Tool registration (21 tools)
│   ├── register_resources.go  # Resource registration (future)
│   ├── register_prompts.go    # Prompt registration (future)
│   └── tools/                 # MCP Tool Definitions
│       ├── server_interface.go # Server interface for tool handlers
│       ├── helpers.go         # Shared utility functions
│       ├── dashboards_tools.go # Dashboard MCP tools (7 tools)
│       ├── widgets_tools.go    # Widget MCP tools (6 tools)
│       ├── metrics_tools.go    # Metrics MCP tools (2 tools)
│       └── alerts_tools.go     # Alert MCP tools (3 tools)
│
├── test/                       # Test Suite
│   ├── config/                # Configuration tests
│   │   └── config_test.go
│   ├── middleware/            # API client tests
│   │   └── client_test.go
│   ├── server/                # Server tests
│   │   └── server_test.go
│   ├── integration/           # Integration tests
│   │   └── integration_test.go
│   └── README.md               # Test documentation
│
├── main.go                     # Application entry point
├── go.mod                      # Go module definition
├── go.sum                      # Go module checksums
├── .env.example                # Example environment configuration
├── .gitignore                  # Git ignore rules
├── Makefile                    # Build and development automation
├── README.md                   # Main project documentation
├── QUICKSTART.md               # Quick start guide
├── TOOLS_DOCUMENTATION.md      # Comprehensive tool reference
└── PROJECT_STRUCTURE.md        # This file
```

## Module Responsibilities

### 1. Configuration (`config/`)

**Purpose:** Centralized configuration management
- Load environment variables from `.env` file
- Validate required configuration
- Provide configuration to other modules

**Key Features:**
- Support for multiple transport modes (stdio, http, sse)
- Tool exclusion for customization
- Default value handling

### 2. Middleware API Client (`middleware/`)

**Purpose:** Abstraction layer for Middleware.io REST API
- HTTP client with authentication
- Type-safe API methods
- Error handling and context support

**Components:**
- **`client.go`**: Base HTTP client, authentication, common request handling
- **`types.go`**: Go structs matching Middleware API data models
- **`dashboards.go`**: CRUD operations for dashboards
- **`widgets.go`**: Widget management and data fetching
- **`metrics.go`**: Metrics metadata and resource discovery
- **`alerts.go`**: Alert instance management

### 3. MCP Server (`server/`)

**Purpose:** Model Context Protocol server implementation
- Register MCP tools
- Handle tool invocations
- Map tool calls to Middleware API

**Structure:**
- **`server.go`**: Core server setup, initialization, and lifecycle management
- **`register_tools.go`**: Registration of all MCP tools (21 tools)
- **`register_resources.go`**: Registration of MCP resources (prepared for future)
- **`register_prompts.go`**: Registration of MCP prompts (prepared for future)
- **`tools/`**: Directory containing all MCP tool definitions
  - **`server_interface.go`**: Interface for tool handlers to access server
  - **`helpers.go`**: Shared utility functions (e.g., ToMap)
  - **`*_tools.go`**: Tool definitions grouped by functionality

**MCP Features (per [MCP Documentation](https://modelcontextprotocol.io/docs/learn/server-concepts)):**
- **Tools** ✅: Functions that AI models can actively call (21 tools implemented)
- **Resources** 🔜: Passive data sources for context (structure prepared)
- **Prompts** 🔜: Pre-built instruction templates (structure prepared)

**Tool Organization:**

#### `dashboards_tools.go` (7 tools)
1. `list_dashboards` - List/search dashboards
2. `get_dashboard` - Get dashboard details
3. `create_dashboard` - Create new dashboard
4. `update_dashboard` - Update dashboard
5. `delete_dashboard` - Delete dashboard
6. `clone_dashboard` - Clone dashboard
7. `set_dashboard_favorite` - Favorite management

#### `widgets_tools.go` (6 tools)
1. `list_widgets` - List widgets
2. `create_widget` - Create/update widget
3. `delete_widget` - Delete widget
4. `get_widget_data` - Fetch widget data
5. `get_multi_widget_data` - Batch widget data
6. `update_widget_layouts` - Update layouts

#### `metrics_tools.go` (2 tools)
1. `get_metrics` - Query metrics/filters/groupby tags
2. `get_resources` - List available resources

#### `alerts_tools.go` (3 tools)
1. `list_alerts` - List alert instances
2. `create_alert` - Create alert instance
3. `get_alert_stats` - Get alert statistics

### 4. Testing (`test/`)

**Purpose:** Comprehensive test coverage
- Unit tests for each module
- Integration tests for full workflows
- HTTP mocking for isolated testing

**Test Organization:**
- `config/`: Configuration loading tests
- `middleware/`: API client tests with httptest
- `server/`: Server initialization tests
- `integration/`: End-to-end workflow tests

## Data Flow

### Tool Invocation Flow

```
1. AI Assistant calls MCP tool
   ↓
2. MCP Go SDK validates input schema
   ↓
3. Server handler receives validated input
   ↓
4. Handler calls middleware client method
   ↓
5. Client makes authenticated HTTP request to Middleware.io API
   ↓
6. API response returned to handler
   ↓
7. Handler formats response as map[string]any
   ↓
8. MCP SDK returns response to AI Assistant
```

### Example: Listing Dashboards

```go
// 1. Tool definition in server/dashboards_tools.go
var listDashboardsTool = &mcp.Tool{
    Name: "list_dashboards",
    Description: "Get a list of dashboards...",
}

// 2. Handler implementation
func (s *Server) handleListDashboards(ctx, req, input) {
    // 3. Call middleware client
    result, err := s.client.GetDashboards(ctx, params)
    
    // 4. Format and return
    return nil, toMap(result), nil
}

// 5. Middleware client in middleware/dashboards.go
func (c *Client) GetDashboards(ctx, params) {
    // 6. Make HTTP request
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    req.Header.Set("ApiKey", c.apiKey)
    
    // 7. Return typed response
    return &DashboardListResponse{}, nil
}
```

## Key Design Principles

### 1. Separation of Concerns
- **Config**: Environment and settings
- **Middleware**: API communication
- **Server**: MCP protocol handling
- **Main**: Application lifecycle

### 2. Type Safety
- Strongly typed structs for all API data
- JSON schema validation for tool inputs
- Compile-time checks where possible

### 3. Testability
- Interfaces for dependency injection
- HTTP mocking with `httptest`
- Context-based cancellation

### 4. Extensibility
- Easy to add new tools (add to `*_tools.go`)
- Tool exclusion for customization
- Modular architecture for new features

### 5. Robustness
- Comprehensive error handling
- Context propagation for cancellation
- Graceful shutdown support
- Input validation at multiple layers

## Adding New Tools

To add a new MCP tool:

1. **Choose appropriate file** based on functionality:
   - Dashboard-related: `server/dashboards_tools.go`
   - Widget-related: `server/widgets_tools.go`
   - Metrics-related: `server/metrics_tools.go`
   - Alert-related: `server/alerts_tools.go`

2. **Define tool and input struct:**
```go
var myNewTool = &mcp.Tool{
    Name: "my_new_tool",
    Description: "Clear description of what this tool does...",
}

type MyNewToolInput struct {
    Param1 string `json:"param1" jsonschema:"Description of param1,required"`
    Param2 int    `json:"param2,omitempty" jsonschema:"Description of param2"`
}
```

3. **Implement handler:**
```go
func (s *Server) handleMyNewTool(ctx context.Context, req *mcp.CallToolRequest, input MyNewToolInput) (*mcp.CallToolResult, map[string]any, error) {
    // Call middleware client
    result, err := s.client.SomeMethod(ctx, params)
    if err != nil {
        return nil, nil, fmt.Errorf("failed: %w", err)
    }
    
    // Format response
    data, err := toMap(result)
    return nil, data, err
}
```

4. **Register in server.go:**
```go
if !s.config.IsToolExcluded("my_new_tool") {
    mcp.AddTool(s.mcpServer, myNewTool, s.handleMyNewTool)
}
```

5. **Add tests in `test/server/`**

6. **Update `TOOLS_DOCUMENTATION.md`**

## File Naming Conventions

### Server Tools Files
- Pattern: `*_tools.go`
- Purpose: Makes it immediately clear the file contains MCP tool definitions
- Examples: `dashboards_tools.go`, `widgets_tools.go`

### Test Files
- Pattern: `*_test.go`
- Location: Parallel structure in `test/` directory
- Examples: `test/server/server_test.go`

### API Client Files
- Pattern: Named after resource (singular)
- Examples: `client.go`, `dashboards.go`, `widgets.go`

## Dependencies

### External Dependencies
- **MCP Go SDK**: `github.com/modelcontextprotocol/go-sdk`
- **godotenv**: `github.com/joho/godotenv`

### Standard Library
- `net/http`: HTTP client
- `context`: Request cancellation
- `encoding/json`: JSON marshaling
- `testing`: Test framework

## Best Practices

### Code Organization
- ✅ One tool per handler function
- ✅ Group related tools in same file
- ✅ Keep files focused and reasonably sized
- ✅ Use clear, descriptive names

### Documentation
- ✅ Comprehensive tool descriptions
- ✅ Detailed parameter documentation
- ✅ Examples and use cases
- ✅ Keep docs in sync with code

### Testing
- ✅ Test each component in isolation
- ✅ Integration tests for workflows
- ✅ Use HTTP mocking for API tests
- ✅ Maintain high test coverage

### Error Handling
- ✅ Wrap errors with context
- ✅ Return meaningful error messages
- ✅ Validate inputs early
- ✅ Handle context cancellation

## Maintenance

### Regular Tasks
1. Keep dependencies updated (`make install`)
2. Run tests frequently (`make test`)
3. Check linting (`make lint`)
4. Update documentation when adding features
5. Review and refactor as needed

### When Adding Features
1. Follow existing patterns
2. Add comprehensive tests
3. Update all documentation
4. Run full test suite
5. Update TOOLS_DOCUMENTATION.md

---

*This structure supports maintainability, extensibility, and clarity for future development.*

