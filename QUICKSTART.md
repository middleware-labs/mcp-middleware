# Quick Start Guide

Get up and running with Middleware MCP Server in 5 minutes!

## Step 1: Get Your API Key

1. Log in to your [Middleware.io account](https://app.middleware.io)
2. Navigate to **Settings** → **API Keys**
3. Click **Generate New API Key**
4. Copy your API key and project URL

## Step 2: Install

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

## Step 3: Configure

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

## Step 4: Test the Server

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

## Step 5: Connect to Claude Desktop

### For Linux/macOS:

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

### For Windows:

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

## Step 6: Test with MCP Inspector (Optional)

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

## Step 7: Try It Out with Claude!

Open Claude Desktop and try these commands:

### List Dashboards
```
Can you list all my dashboards in Middleware?
```

### Get Resources
```
What resources are available in my Middleware account?
```

### Create a Dashboard
```
Create a new dashboard called "Production Metrics" with public visibility
```

### Get Widget Data
```
Get the data for widget with builder ID 123
```

## Common Issues

### "MIDDLEWARE_API_KEY is required"
- Make sure your `.env` file exists and contains the API key
- Or set the environment variable directly: `export MIDDLEWARE_API_KEY=your_key`

### "Connection refused" or timeout errors
- Verify your `MIDDLEWARE_BASE_URL` is correct
- Check that you can access your Middleware project in a browser
- Ensure your network allows outbound HTTPS connections

### Claude Desktop doesn't see the server
- Use the **full absolute path** to the binary in the config
- Make sure the binary has execute permissions: `chmod +x mcp-middleware`
- Restart Claude Desktop after configuration changes
- Check Claude Desktop logs for errors

### "API error (401)" or "API error (403)"
- Your API key may be invalid or expired
- Generate a new API key from Middleware settings
- Make sure you're using the correct API key (not the agent installation key)

## Next Steps

- Read the full [README.md](README.md) for detailed documentation
- Check out [CONTRIBUTING.md](CONTRIBUTING.md) if you want to contribute
- Explore all 21 available tools in the main documentation
- Configure tool exclusions for read-only access

## Need Help?

- 📖 [Full Documentation](README.md)
- 📧 [Middleware Support](mailto:support@middleware.io)
- 🧪 [Testing Documentation](test/README.md)

## What's Next?

Now that you're set up, you can:

1. **Explore all tools**: Try different commands to see what's available
2. **Create dashboards**: Build custom dashboards through Claude
3. **Query metrics**: Get insights from your monitoring data
4. **Manage alerts**: Create and monitor alerts programmatically
5. **Automate workflows**: Use Claude to automate your observability tasks

Happy monitoring! 🚀

