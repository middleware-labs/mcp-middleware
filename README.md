# Middleware MCP Server

A robust and modular [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server for [Middleware.io](https://middleware.io) built with the [official Go SDK](https://github.com/modelcontextprotocol/go-sdk). This server enables AI assistants like Claude to interact with Middleware's observability platform for monitoring, dashboards, widgets, metrics, and alerts.

## Features

- 🚀 **Complete API Coverage**: Implements all major Middleware API endpoints
- 🔧 **Modular Architecture**: Clean separation of concerns for easy extension
- 🛡️ **Type-Safe**: Fully typed with Go structs matching Middleware API schemas
- ⚙️ **Configurable**: Environment-based configuration with tool exclusion support
- 📊 **Rich Toolset**: 21 tools covering dashboards, widgets, metrics, and alerts
- 🔌 **Official SDK**: Built using the official MCP Go SDK

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

### Metrics & Resources (2 tools)
- `get_metrics` - Get metrics, filters, or groupby tags
- `get_resources` - Get available resources for queries

### Alerts (3 tools)
- `list_alerts` - List alerts for a specific rule
- `create_alert` - Create a new alert
- `get_alert_stats` - Get alert statistics

## Prerequisites

- **Go 1.23 or later** (the project uses Go 1.23.0 with toolchain 1.24.10)
- **Middleware.io account** with API access
- **API Key** from [Middleware API Keys settings](https://app.middleware.io/settings/api-keys)

## Installation

### 1. Clone the repository

```bash
# Clone or download the project
cd mcp-middleware
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Configure environment variables

```bash
cp .env.example .env
```

Edit `.env` with your Middleware credentials:

```env
MIDDLEWARE_API_KEY=your_api_key_here
MIDDLEWARE_BASE_URL=https://your-project.middleware.io
APP_MODE=stdio
```

### 4. Build the server

```bash
go build -o mcp-middleware
```

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

### Running the Server (stdio mode)

```bash
./mcp-middleware
```

### Testing with MCP Inspector

The project includes support for the [MCP Inspector](https://modelcontextprotocol.io/docs/tools/inspector), an interactive tool for testing and debugging your MCP server.

```bash
# Quick start (requires Node.js and npx)
make inspect

# Or with .env configuration
make inspect-env
```

The inspector will open in your browser and connect to your server, allowing you to:
- Test all 21 tools interactively
- View real-time server logs and notifications
- Debug tool inputs and outputs
- Verify tool schemas and descriptions

### Claude Desktop Integration

Add to your Claude Desktop configuration (`~/.config/Claude/claude_desktop_config.json` on Linux/macOS or `%APPDATA%\Claude\claude_desktop_config.json` on Windows):

```json
{
  "mcpServers": {
    "middleware": {
      "command": "/full/path/to/mcp-middleware",
      "env": {
        "MIDDLEWARE_API_KEY": "your_api_key",
        "MIDDLEWARE_BASE_URL": "https://your-project.middleware.io"
      }
    }
  }
}
```

**Important**: Use the full absolute path to your built binary.

### Example Usage with Claude

Once connected, you can ask Claude to:

- "List all my dashboards in Middleware"
- "Create a new dashboard called 'Production Metrics'"
- "Get the data for widget ID 123"
- "Show me all alerts for rule 456"
- "What resources are available in Middleware?"
- "Clone the 'API Performance' dashboard"

## Project Structure

```
mcp-middleware/
├── config/          # Configuration management
│   └── config.go    # Environment variable loading
├── middleware/      # Middleware API client
│   ├── client.go    # HTTP client with authentication
│   ├── types.go     # API type definitions
│   ├── dashboards.go # Dashboard API methods
│   ├── widgets.go   # Widget API methods
│   ├── metrics.go   # Metrics API methods
│   └── alerts.go    # Alerts API methods
├── server/                  # MCP server implementation
│   ├── server.go           # Server initialization and lifecycle
│   ├── register_tools.go   # Tool registration (21 tools)
│   ├── register_resources.go # Resource registration (future)
│   ├── register_prompts.go # Prompt registration (future)
│   └── tools/              # MCP tool definitions
│       ├── dashboards_tools.go # Dashboard MCP tools (7 tools)
│       ├── widgets_tools.go    # Widget MCP tools (6 tools)
│       ├── metrics_tools.go    # Metrics MCP tools (2 tools)
│       └── alerts_tools.go     # Alert MCP tools (3 tools)
├── test/            # Test suite (28 tests)
│   ├── config/      # Config tests (9 tests)
│   ├── middleware/  # Client tests (11 tests)
│   ├── server/      # Server tests (3 tests)
│   ├── integration/ # Integration tests (5 tests)
│   └── README.md    # Testing documentation
├── main.go          # Application entry point
├── go.mod           # Go module definition
├── .env.example     # Example environment configuration
├── .gitignore       # Git ignore rules
└── README.md        # This file
```

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

### Adding New Tools

1. Add the API method to the appropriate file in `middleware/`
2. Define the tool and handler in the corresponding file in `server/`
3. Register the tool in `server/server.go`'s `registerTools()` method

### Code Style

This project follows Go best practices:
- Use `any` instead of `interface{}`
- Proper error handling with wrapped errors
- Context propagation for cancellation
- Clear package separation

## API Documentation

For detailed Middleware API documentation, visit:
- [Middleware API Swagger](https://app.middleware.io/swagger/index.html)
- [Middleware Documentation](https://docs.middleware.io/)

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Troubleshooting

### "MIDDLEWARE_API_KEY is required" error
Make sure you've set the `MIDDLEWARE_API_KEY` environment variable or created a `.env` file with the correct API key.

### "API error (401)" or "API error (403)"
Your API key may be invalid or expired. Generate a new one from [Middleware Settings](https://app.middleware.io/settings/api-keys).

### Claude Desktop doesn't see the server
1. Ensure the path in `claude_desktop_config.json` is absolute
2. Check that the binary has execute permissions (`chmod +x mcp-middleware`)
3. Restart Claude Desktop after configuration changes

### Connection timeouts
Check that your `MIDDLEWARE_BASE_URL` is correct and accessible from your network.

## References

- [Model Context Protocol Documentation](https://modelcontextprotocol.io/docs/getting-started/intro)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [MCP Inspector](https://modelcontextprotocol.io/docs/tools/inspector) - Testing and debugging tool
- [Middleware.io](https://middleware.io)
- [Zerodha Kite MCP Server](https://github.com/zerodha/kite-mcp-server) (inspiration)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with the [official MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- Inspired by [Zerodha's Kite MCP Server](https://github.com/zerodha/kite-mcp-server)
- Thanks to the Middleware.io team for their comprehensive API

## Support

For issues and questions:
- **Middleware Support**: [support@middleware.io](mailto:support@middleware.io)
- **Documentation**: See README.md for full documentation

---

Made with ❤️ for the MCP and Middleware communities

