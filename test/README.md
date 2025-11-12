# Test Suite

This directory contains comprehensive tests for the Middleware MCP Server organized by component.

## Test Structure

```
test/
├── config/          # Configuration tests (9 tests)
│   └── config_test.go
├── middleware/      # API client tests (11 tests)
│   └── client_test.go
├── server/          # Server initialization tests (3 tests)
│   └── server_test.go
└── integration/     # Integration tests (5 tests)
    └── integration_test.go
```

## Running Tests

### Run All Tests (28 tests)
```bash
make test
# or
go test -v ./test/...
```

### Run Specific Test Suites
```bash
# Config tests only (9 tests)
make test-config

# Middleware client tests only (11 tests)
make test-middleware

# Server tests only (3 tests)
make test-server

# Integration tests only (5 tests)
make test-integration
```

### Run with Coverage
```bash
make test-coverage
```

This generates:
- `coverage.out` - Coverage data
- `coverage.html` - Visual coverage report (open in browser)

### Run with Race Detection
```bash
make test-race
```

Detects potential race conditions in concurrent code.

## Test Coverage

### Config Tests (`test/config/config_test.go`) - 9 tests

Tests for configuration loading and validation:

| Test | Description |
|------|-------------|
| `TestLoad` | Basic configuration loading |
| `TestLoadMissingAPIKey` | Error handling for missing API key |
| `TestLoadMissingBaseURL` | Error handling for missing base URL |
| `TestIsToolExcluded` | Tool exclusion functionality |
| `TestInvalidAppMode` | Invalid app mode validation |
| `TestCustomAppModes` | Custom app mode configurations |
| `TestCustomHostAndPort` | Custom host/port settings |
| `TestExcludedToolsWithSpaces` | Space handling in tool exclusion |
| `TestEmptyExcludedTools` | Empty exclusion list handling |

### Middleware Client Tests (`test/middleware/client_test.go`) - 11 tests

Tests for Middleware API client functionality:

| Test | Description |
|------|-------------|
| `TestNewClient` | Client initialization |
| `TestGetDashboards` | Dashboard listing |
| `TestGetDashboardsWithParams` | Dashboard filtering with params |
| `TestCreateDashboard` | Dashboard creation |
| `TestDeleteDashboard` | Dashboard deletion |
| `TestGetResources` | Resource listing |
| `TestAPIError` | API error handling |
| `TestContextCancellation` | Context cancellation support |
| `TestGetMetrics` | Metrics retrieval |
| `TestGetAlerts` | Alert listing |
| `TestCreateAlert` | Alert creation |

**Note:** These tests use `httptest` to mock API responses, so no actual API calls are made.

### Server Tests (`test/server/server_test.go`) - 3 tests

Tests for MCP server initialization:

| Test | Description |
|------|-------------|
| `TestNewServer` | Basic server creation |
| `TestNewServerWithExcludedTools` | Server with excluded tools |
| `TestServerConfiguration` | Various configuration scenarios |

### Integration Tests (`test/integration/integration_test.go`) - 5 tests

End-to-end integration tests:

| Test | Description |
|------|-------------|
| `TestFullServerInitialization` | Complete server init flow |
| `TestClientServerIntegration` | Client + server together |
| `TestConfigToServerFlow` | Config → server pipeline |
| `TestClientContextHandling` | Context propagation |
| `TestMultipleServerInstances` | Multiple server instances |

## Test Principles

### 1. **Isolation**
Each test is independent and doesn't rely on external services.

### 2. **Mocking**
- HTTP requests are mocked with `httptest`
- No real Middleware API calls
- Fast test execution

### 3. **Table-Driven Tests**
Many tests use table-driven approach for comprehensive coverage:

```go
tests := []struct {
    name    string
    input   Input
    want    Output
    wantErr bool
}{
    // test cases...
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // test logic...
    })
}
```

### 4. **Environment Management**
Tests properly set and unset environment variables:

```go
os.Setenv("KEY", "value")
defer os.Unsetenv("KEY")
```

## Writing New Tests

### Adding Config Tests

Add to `test/config/config_test.go`:

```go
func TestNewConfigFeature(t *testing.T) {
    os.Setenv("MIDDLEWARE_API_KEY", "test-key")
    os.Setenv("MIDDLEWARE_BASE_URL", "https://test.middleware.io")
    defer func() {
        os.Unsetenv("MIDDLEWARE_API_KEY")
        os.Unsetenv("MIDDLEWARE_BASE_URL")
    }()

    cfg, err := config.Load()
    if err != nil {
        t.Fatalf("Load() failed: %v", err)
    }

    // Your assertions here
}
```

### Adding Middleware Client Tests

Add to `test/middleware/client_test.go`:

```go
func TestNewEndpoint(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Mock response
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(YourResponse{})
    }))
    defer server.Close()

    client := middleware.NewClient(server.URL, "test-key")
    ctx := context.Background()

    result, err := client.YourMethod(ctx, params)
    if err != nil {
        t.Fatalf("YourMethod() error = %v", err)
    }

    // Your assertions here
}
```

### Adding Server Tests

Add to `test/server/server_test.go`:

```go
func TestNewServerFeature(t *testing.T) {
    cfg := &config.Config{
        MiddlewareAPIKey:  "test-key",
        MiddlewareBaseURL: "https://test.middleware.io",
        AppMode:           "stdio",
        ExcludedTools:     make(map[string]bool),
    }

    srv := server.New(cfg)
    if srv == nil {
        t.Fatal("Server creation failed")
    }

    // Your assertions here
}
```

### Adding Integration Tests

Add to `test/integration/integration_test.go`:

```go
func TestNewIntegrationScenario(t *testing.T) {
    os.Setenv("MIDDLEWARE_API_KEY", "test-key")
    os.Setenv("MIDDLEWARE_BASE_URL", "https://test.middleware.io")
    defer func() {
        os.Unsetenv("MIDDLEWARE_API_KEY")
        os.Unsetenv("MIDDLEWARE_BASE_URL")
    }()

    // Test your integration scenario
}
```

## Continuous Integration

Tests are designed to run in CI environments:

```yaml
# Example GitHub Actions
- name: Run tests
  run: make test

- name: Run tests with coverage
  run: make test-coverage

- name: Run race detection
  run: make test-race
```

## Test Results

Expected output when all tests pass:

```
=== RUN   TestLoad
--- PASS: TestLoad (0.00s)
...
PASS
ok      mcp-middleware/test/config        0.002s
PASS
ok      mcp-middleware/test/middleware    2.018s
PASS
ok      mcp-middleware/test/server        0.017s
PASS
ok      mcp-middleware/test/integration   0.023s
```

**Total: 28 tests, all passing ✅**

## Troubleshooting

### Tests fail with "context deadline exceeded"
- Check if `TestContextCancellation` timeout is too short
- Increase timeout in test if needed

### Tests fail with "address already in use"
- Tests use `httptest` which picks random ports
- Should not happen, but restart if it does

### Tests fail with "MIDDLEWARE_API_KEY is required"
- Tests properly set/unset env vars
- If this fails, check for env var pollution from previous tests

## Best Practices

1. ✅ **Always clean up** - Use `defer` to unset env vars
2. ✅ **Use table-driven tests** - More comprehensive, easier to maintain
3. ✅ **Mock external calls** - Use `httptest` for HTTP mocking
4. ✅ **Test errors** - Verify error conditions, not just happy paths
5. ✅ **Keep tests fast** - All 28 tests run in < 3 seconds
6. ✅ **Test isolation** - Each test is independent

## Contributing

When adding new functionality:

1. Write tests first (TDD)
2. Ensure tests pass: `make test`
3. Check coverage: `make test-coverage`
4. Run race detector: `make test-race`
5. Update this README if adding new test files

---

**Happy Testing! 🧪**

