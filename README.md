# Middleware MCP Server

A robust and modular [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server for [Middleware.io](https://middleware.io). This server enables AI assistants like Claude to interact with Middleware's observability platform for monitoring, dashboards, widgets, metrics, and alerts.

## Available Tools

### Dashboard Management (7 tools)
- `list_dashboards` - List all dashboards with filtering and pagination
- `get_dashboard` - Get a specific dashboard by key
- `create_dashboard` - Create a new dashboard
- `update_dashboard` - Update an existing dashboard
- `delete_dashboard` - Delete a dashboard
- `clone_dashboard` - Clone an existing dashboard
- `set_dashboard_favorite` - Mark dashboard as favorite/unfavorite

### Widget Management (6 tools)
- `list_widgets` - List widgets for a report or display scope
- `create_widget` - Create or update a widget
- `delete_widget` - Delete a widget
- `get_widget_data` - Get data for a specific widget
- `get_multi_widget_data` - Get data for multiple widgets at once
- `update_widget_layouts` - Update widget layout positions

### Metrics & Resources (3 tools)
- `get_metrics` - Get metrics, filters, or groupby tags
- `get_resources` - Get available resources for queries
- `query` - Execute flexible queries to retrieve logs, metrics, traces, and other data

### Alerts (3 tools)
- `list_alerts` - List alerts for a specific rule
- `create_alert` - Create a new alert
- `get_alert_stats` - Get alert statistics

### Error/Incident Management (2 tools)
- `list_errors` - List all errors/incidents with filtering and pagination (includes clickable `issue_url` for each incident)
- `get_error_details` - Get detailed information about a specific error/incident by fingerprint

## Quick Start

Get up and running in 5 minutes!

### Step 1: Get Your API Key

1. Log in to your [Middleware.io account](https://app.middleware.io)
2. Navigate to **Settings** → **API Keys**
3. Click **Generate New API Key**
4. Copy your API key and project URL

### Step 2: Install

```bash
# Navigate to the project directory
cd mcp-middleware

# Install dependencies
go mod download

# Build the server
go build -o mcp-middleware .
```

Or using Make:
```bash
make install
make build
```

### Step 3: Configure

Create a `.env` file:
```bash
cp .env.example .env
```

Edit `.env` with your credentials:
```env
MIDDLEWARE_API_KEY=your_api_key_here
MIDDLEWARE_BASE_URL=https://your-project.middleware.io
APP_MODE=stdio
```

### Step 4: Test the Server

Run the server directly:
```bash
./mcp-middleware
```

The server will start in stdio mode. You should see:
```
Middleware MCP Server v1.0.0
Connected to: https://your-project.middleware.io
Starting MCP server in stdio mode...
```

Press `Ctrl+C` to stop.

### Step 5: Connect to Claude Desktop

#### For macOS:

1. Open `~/.config/Claude/claude_desktop_config.json`
2. Add the server configuration:

```json
{
  "mcpServers": {
    "middleware": {
      "command": "/full/path/to/mcp-middleware/mcp-middleware",
      "env": {
        "MIDDLEWARE_API_KEY": "your_api_key",
        "MIDDLEWARE_BASE_URL": "https://your-project.middleware.io"
      }
    }
  }
}
```

#### For Windows:

1. Open `%APPDATA%\Claude\claude_desktop_config.json`
2. Add the server configuration with Windows path:

```json
{
  "mcpServers": {
    "middleware": {
      "command": "C:\\path\\to\\mcp-middleware\\mcp-middleware.exe",
      "env": {
        "MIDDLEWARE_API_KEY": "your_api_key",
        "MIDDLEWARE_BASE_URL": "https://your-project.middleware.io"
      }
    }
  }
}
```

3. Restart Claude Desktop

### Step 6: Test with MCP Inspector (Optional)

Before connecting to Claude, you can test your server with the MCP Inspector:

```bash
# Requires Node.js and npx
make inspect
```

This opens an interactive web interface where you can:
- Test all 21 tools
- View server logs in real-time
- Debug inputs and outputs
- Verify everything works

### Step 7: Try It Out with Claude!

Open Claude Desktop and try these commands:

**List Dashboards:**
```
Can you list all my dashboards in Middleware?
```

**Get Resources:**
```
What resources are available in my Middleware account?
```

**Create a Dashboard:**
```
Create a new dashboard called "Production Metrics" with public visibility
```

**Get Widget Data:**
```
Get the data for widget with builder ID 123
```

**List Errors:**
```
List all errors in the system from the last hour
```

**Get Error Details:**
```
Get detailed information about error with fingerprint 7693967476886782339
```

## Prerequisites

- **Go 1.23 or later** (the project uses Go 1.23.0 with toolchain 1.24.10)
- **Node.js and npx** (for MCP Inspector testing)
- **Middleware.io account** with API access
- **API Key** from [Middleware API Keys settings](https://app.middleware.io/settings/api-keys)

## Transport Modes

The server supports three transport modes:

- **stdio** (default): Standard input/output transport for command-line usage
- **http**: Streamable HTTP transport for web-based clients (uses `NewStreamableHTTPServer`)
- **sse**: Server-Sent Events transport for real-time streaming (uses `NewSSEServer`)

## Configuration

| Environment Variable | Required | Default | Description |
|---------------------|----------|---------|-------------|
| `MIDDLEWARE_API_KEY` | ✅ Yes | - | Your Middleware API key from settings |
| `MIDDLEWARE_BASE_URL` | ✅ Yes | - | Your Middleware project URL (e.g., `https://your-project.middleware.io`) |
| `APP_MODE` | No | `stdio` | Server mode: `stdio`, `http`, or `sse` |
| `APP_HOST` | No | `localhost` | Server host (for http/sse modes) |
| `APP_PORT` | No | `8080` | Server port (for http/sse modes) |
| `EXCLUDED_TOOLS` | No | - | Comma-separated list of tools to exclude |

### Tool Exclusion

You can exclude specific tools for security or functionality reasons:

```env
EXCLUDED_TOOLS=delete_dashboard,delete_widget,create_alert
```

This is useful for creating read-only instances or restricting destructive operations.

## Usage

### Running the Server

#### Stdio Mode (Default)

```bash
./mcp-middleware
# Or set explicitly
APP_MODE=stdio ./mcp-middleware
```

#### HTTP Mode

Start the server in HTTP mode for web-based clients:

```bash
APP_MODE=http APP_HOST=localhost APP_PORT=8080 ./mcp-middleware
```

The server will start on `http://localhost:8080`. Clients can connect using the streamable HTTP transport.

#### SSE Mode

Start the server in SSE (Server-Sent Events) mode:

```bash
APP_MODE=sse APP_HOST=localhost APP_PORT=8080 ./mcp-middleware
```

The server will start on `http://localhost:8080` with SSE support for real-time streaming.

## Project Structure

### Directory Layout

```
mcp-middleware/
├── config/                     # Configuration Management
│   └── config.go              # Environment variable loading and validation
│
├── middleware/                 # Middleware.io API Client
│   ├── client.go              # HTTP client with authentication
│   ├── types.go               # API data structures (Dashboard, Widget, Alert, Incident, etc.)
│   ├── dashboards.go          # Dashboard API endpoints
│   ├── widgets.go             # Widget API endpoints
│   ├── metrics.go             # Metrics API endpoints
│   ├── alerts.go              # Alert API endpoints
│   └── issues.go              # Error/Incident API endpoints
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
│       ├── metrics_tools.go    # Metrics MCP tools (3 tools)
│       ├── alerts_tools.go     # Alert MCP tools (3 tools)
│       ├── errors_tools.go      # Error/Incident MCP tools (2 tools)
│       └── TOOLS_DOCUMENTATION.md # Comprehensive tool reference
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
└── README.md                   # This file
```

### Module Responsibilities

#### 1. Configuration (`config/`)

**Purpose:** Centralized configuration management
- Load environment variables from `.env` file
- Validate required configuration
- Provide configuration to other modules

**Key Features:**
- Support for multiple transport modes (stdio, http, sse)
- Tool exclusion for customization
- Default value handling

#### 2. Middleware API Client (`middleware/`)

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
- **`issues.go`**: Error/incident listing and detail retrieval

#### 3. MCP Server (`server/`)

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
  - **`helpers.go`**: Shared utility functions (e.g., ToMap, ToTextResult)
  - **`*_tools.go`**: Tool definitions grouped by functionality
  - **`TOOLS_DOCUMENTATION.md`**: Comprehensive documentation for all tools

**MCP Features:**
- **Tools** ✅: Functions that AI models can actively call (21 tools implemented)
- **Resources** 🔜: Passive data sources for context (structure prepared)
- **Prompts** 🔜: Pre-built instruction templates (structure prepared)

**Tool Organization:**

- **`dashboards_tools.go`** (7 tools): List, get, create, update, delete, clone dashboards, set favorites
- **`widgets_tools.go`** (6 tools): List, create, delete widgets, get widget data, batch data, update layouts
- **`metrics_tools.go`** (3 tools): Get metrics/filters/groupby tags, list available resources, execute flexible queries
- **`alerts_tools.go`** (3 tools): List alerts, create alerts, get alert statistics
- **`errors_tools.go`** (2 tools): List errors/incidents with clickable URLs, get error details by fingerprint

#### 4. Testing (`test/`)

**Purpose:** Comprehensive test coverage
- Unit tests for each module
- Integration tests for full workflows
- HTTP mocking for isolated testing

**Test Organization:**
- `config/`: Configuration loading tests
- `middleware/`: API client tests with httptest
- `server/`: Server initialization tests
- `integration/`: End-to-end workflow tests

### Key Design Principles

1. **Separation of Concerns**: Config, Middleware, Server, and Main are clearly separated
2. **Type Safety**: Strongly typed structs for all API data with JSON schema validation
3. **Testability**: Interfaces for dependency injection, HTTP mocking, context-based cancellation
4. **Extensibility**: Easy to add new tools, tool exclusion for customization, modular architecture
5. **Robustness**: Comprehensive error handling, context propagation, graceful shutdown, input validation

### Adding New Tools

To add a new MCP tool:

1. **Choose appropriate file** based on functionality (dashboard/widget/metrics/alert)
2. **Define tool using `mcp.NewTool()`** with `mcp.WithDescription()` and `mcp.WithInputSchema[T]()`
3. **Define input struct** with proper JSON schema tags
4. **Implement handler** that calls the middleware client (signature: `func HandleTool(s ServerInterface, ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)`)
5. **Register in `server/register_tools.go`** using `s.mcpServer.AddTool(tools.NewTool(), handler)`
6. **Add tests** in `test/server/`
7. **Update** `server/tools/TOOLS_DOCUMENTATION.md`

## Development

### Running Tests

The project includes **28 comprehensive tests** organized in the `test/` directory:

```bash
# Run all tests (28 tests)
make test

# Run with coverage report
make test-coverage

# Run with race detection
make test-race

# Run specific test suites
make test-config        # Config tests (9 tests)
make test-middleware    # Middleware client tests (11 tests)
make test-server        # Server tests (3 tests)
make test-integration   # Integration tests (5 tests)
```

**Test Coverage:**
- ✅ Config: 9 tests
- ✅ Middleware Client: 11 tests (with HTTP mocking)
- ✅ Server: 3 tests
- ✅ Integration: 5 tests

For detailed testing information, see [test/README.md](test/README.md).

### Testing with MCP Inspector

Use the MCP Inspector for interactive testing during development:

```bash
# Start inspector
make inspect

# In another terminal, make changes and rebuild
make build

# Reconnect in inspector to test changes
```


### Code Style

This project follows Go best practices:
- Use `any` instead of `interface{}`
- Proper error handling with wrapped errors
- Context propagation for cancellation
- Clear package separation

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues and questions:
- **Middleware Support**: [support@middleware.io](mailto:support@middleware.io)
- **Documentation**: See README.md for full documentation

---

Made with ❤️ for the MCP and Middleware communities

