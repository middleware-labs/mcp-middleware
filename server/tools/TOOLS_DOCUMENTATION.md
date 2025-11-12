# MCP Tools Documentation

This document provides comprehensive information about all 21 MCP tools available in the Middleware.io MCP server. Each tool is documented with detailed descriptions and parameter information based on the [official Middleware API](https://app.middleware.io/swagger.json).

## Overview

The Middleware MCP server exposes Middleware.io's monitoring platform capabilities through the Model Context Protocol, allowing AI assistants to interact with dashboards, widgets, metrics, and alerts.

## Tool Categories

### 📊 Dashboard Tools (7 tools)
Manage custom dashboards for organizing monitoring visualizations.

### 📈 Widget Tools (6 tools)
Create and manage visualization widgets (charts, graphs, tables) within dashboards.

### 📉 Metrics Tools (2 tools)
Query available metrics, resources, and metadata for building monitoring queries.

### 🚨 Alert Tools (3 tools)
Manage and analyze alert instances and statistics.

---

## Dashboard Tools

### 1. `list_dashboards`
**Purpose:** Get a list of dashboards with advanced filtering and pagination support.

**Description:** This tool retrieves dashboards from Middleware.io with support for searching, filtering by various criteria, and pagination. Use this to discover available dashboards, find specific dashboards by name, or filter by ownership and usage patterns.

**Parameters:**
- `limit` (integer, optional): Number of items per page for pagination
- `offset` (integer, optional): Number of items to skip for pagination (page offset)
- `search` (string, optional): Search query to find dashboards by name or description
- `filter_by` (string, optional): Comma-separated list of filter values. Valid values: custom, created_by_you, favorite, frequently_viewed, or data source names like aws, mysql, postgresql, etc.
- `display_scope` (string, optional): Filter dashboards by comma-separated list of display scopes

**Example Use Cases:**
- Find all custom dashboards
- Search for dashboards containing "production"
- List your frequently viewed dashboards

---

### 2. `get_dashboard`
**Purpose:** Get detailed information about a specific dashboard by its unique key.

**Description:** This tool retrieves complete dashboard configuration including widgets, layout, metadata, and settings. Use this when you need to inspect or work with a specific dashboard's structure and content.

**Parameters:**
- `report_key` (string, **required**): The unique key identifier of the dashboard to retrieve

**Example Use Cases:**
- Inspect dashboard configuration
- View all widgets in a dashboard
- Export dashboard structure

---

### 3. `create_dashboard`
**Purpose:** Create a new custom dashboard in Middleware.io.

**Description:** This tool creates a new dashboard with the specified configuration. Dashboards can be public (shared with team) or private (personal). You can organize dashboards using display scopes and provide custom keys for easier identification.

**Parameters:**
- `label` (string, **required**): The dashboard name/title. Must be at least 3 characters long
- `visibility` (string, **required**): Dashboard visibility setting. Must be either 'public' (shared with team) or 'private' (personal only)
- `description` (string, optional): Optional detailed description of the dashboard's purpose and contents
- `display_scope` (string, optional): Optional display scope for organizing dashboards into categories or groups
- `key` (string, optional): Optional unique key identifier for the dashboard. If not provided, will be auto-generated

**Example Use Cases:**
- Create a new team dashboard
- Set up personal monitoring dashboards
- Organize dashboards by environment or service

---

### 4. `update_dashboard`
**Purpose:** Update an existing dashboard's configuration and metadata.

**Description:** This tool modifies an existing dashboard identified by its ID. You can update the name, description, visibility settings, and display scope. Use this to rename dashboards, change sharing settings, or reorganize dashboard categories.

**Parameters:**
- `id` (integer, **required**): The numeric ID of the dashboard to update
- `label` (string, **required**): The updated dashboard name/title. Must be at least 3 characters long
- `visibility` (string, **required**): Updated visibility setting. Must be either 'public' or 'private'
- `description` (string, optional): Updated description of the dashboard
- `display_scope` (string, optional): Updated display scope for dashboard organization
- `key` (string, optional): Updated unique key identifier. Must be unique across all dashboards

**Example Use Cases:**
- Rename a dashboard
- Change dashboard from private to public
- Update dashboard description

---

### 5. `delete_dashboard`
**Purpose:** Permanently delete a dashboard and all its widgets.

**Description:** This tool removes a dashboard from Middleware.io. Warning: This action cannot be undone. All widgets and configurations associated with the dashboard will be permanently deleted.

**Parameters:**
- `id` (integer, **required**): The numeric ID of the dashboard to delete permanently

**Example Use Cases:**
- Remove obsolete dashboards
- Clean up test dashboards
- Delete duplicate dashboards

---

### 6. `clone_dashboard`
**Purpose:** Create a copy of an existing dashboard with all its widgets and configuration.

**Description:** This tool duplicates an existing dashboard, creating a new dashboard with the same widgets, layout, and settings. Useful for creating variations of dashboards or starting from a template. The cloned dashboard will have a new ID and can have different visibility settings.

**Parameters:**
- `label` (string, **required**): The name for the new cloned dashboard. Must be at least 3 characters
- `visibility` (string, **required**): Visibility setting for the cloned dashboard: 'public' or 'private'
- `description` (string, optional): Optional description for the cloned dashboard
- `display_scope` (string, optional): Optional display scope for organizing the cloned dashboard
- `source_key` (string, optional): The unique key of the source dashboard to clone from

**Example Use Cases:**
- Create environment-specific dashboards from a template
- Duplicate a dashboard for customization
- Share a personal dashboard with the team

---

### 7. `set_dashboard_favorite`
**Purpose:** Mark a dashboard as favorite or remove it from favorites.

**Description:** This tool allows you to favorite dashboards for quick access. Favorited dashboards appear at the top of dashboard lists and can be filtered using the 'favorite' filter in list_dashboards. Use this to bookmark frequently accessed dashboards.

**Parameters:**
- `report_id` (integer, **required**): The numeric ID of the dashboard to mark as favorite or unfavorite
- `favorite` (boolean, **required**): Set to true to add dashboard to favorites, false to remove from favorites

**Example Use Cases:**
- Bookmark important dashboards
- Organize frequently accessed dashboards
- Quick access to key monitoring views

---

## Widget Tools

### 8. `list_widgets`
**Purpose:** Get a list of widgets associated with a specific dashboard or display scope.

**Description:** This tool retrieves all widgets (charts, graphs, tables) that belong to a dashboard or scope. Widgets are the building blocks of dashboards - each widget represents a visualization of your monitoring data. Use this to discover what widgets are available in a dashboard or to inspect widget configurations.

**Parameters:**
- `report_id` (integer, optional): The numeric ID of the dashboard (report) to filter widgets by
- `display_scope` (string, optional): The display scope to filter widgets by (e.g., 'infrastructure', 'apm', 'logs')

**Example Use Cases:**
- View all widgets in a dashboard
- Discover widget configurations
- Inventory dashboard components

---

### 9. `create_widget`
**Purpose:** Create a new widget or update an existing widget on a dashboard.

**Description:** This tool allows you to add new visualizations (charts, graphs, tables) to dashboards or modify existing ones. The builder_config contains the query, chart type, and visualization settings. If builder_id is provided, it updates the existing widget; otherwise, it creates a new one.

**Parameters:**
- `label` (string, **required**): The display name for the widget (e.g., 'CPU Usage', 'Error Rate')
- `key` (string, optional): Optional unique key identifier for the widget
- `description` (string, optional): Optional description explaining what the widget displays
- `builder_config` (object, optional): Widget configuration object containing queries, chart type, display settings, and data sources. This is a complex object specific to widget type
- `builder_id` (integer, optional): If provided, updates the existing widget with this ID instead of creating a new one

**Example Use Cases:**
- Add new charts to dashboards
- Update widget queries
- Create custom visualizations

---

### 10. `delete_widget`
**Purpose:** Permanently delete a widget from a dashboard.

**Description:** This tool removes a widget (chart, graph, table) from its dashboard. Warning: This action cannot be undone. The widget configuration and data will be permanently deleted.

**Parameters:**
- `builder_id` (integer, **required**): The numeric builder ID of the widget to delete permanently

**Example Use Cases:**
- Remove obsolete widgets
- Clean up dashboard layouts
- Delete duplicate widgets

---

### 11. `get_widget_data`
**Purpose:** Fetch the actual data and metrics displayed by a specific widget.

**Description:** This tool executes the widget's query and returns the visualization data (time series, metrics, logs, traces). Use this to get the current values shown in a widget, analyze trends, or export widget data. The data format depends on the widget type (timeseries, table, single value, etc.).

**Parameters:**
- `builder_id` (integer, optional): The numeric builder ID of the widget to fetch data for
- `key` (string, optional): Alternative to builder_id: the unique key identifier of the widget
- `label` (string, optional): Alternative to builder_id: the label of the widget
- `builder_config` (object, optional): Widget configuration containing the query and data source settings
- `use_v2` (boolean, optional): Set to true to use the newer v2 data format (default: false)

**Example Use Cases:**
- Get current metric values
- Export widget data
- Analyze trends over time

---

### 12. `get_multi_widget_data`
**Purpose:** Fetch data for multiple widgets simultaneously in a single request.

**Description:** This tool is optimized for loading data for multiple widgets at once, such as when refreshing an entire dashboard. It's more efficient than calling get_widget_data multiple times. Returns data for all requested widgets in a single response.

**Parameters:**
- `widgets` (array of objects, **required**): Array of widget specifications to fetch data for. Each widget can be identified by builder_id, key, or label
  - Each widget object contains:
    - `builder_id` (integer, optional): The numeric builder ID of the widget
    - `key` (string, optional): The unique key identifier of the widget
    - `label` (string, optional): The label of the widget
    - `builder_config` (object, optional): Widget configuration containing query and display settings
    - `use_v2` (boolean, optional): Use v2 data format (default: false)

**Example Use Cases:**
- Refresh entire dashboard
- Export multiple widget data sets
- Batch data retrieval

---

### 13. `update_widget_layouts`
**Purpose:** Update the position and size of widgets on a dashboard.

**Description:** This tool modifies the layout (position, size) of multiple widgets on a dashboard. Use this to rearrange widgets, resize them, or optimize dashboard layout. The dashboard uses a grid system where x,y represent position and w,h represent size in grid units.

**Parameters:**
- `layouts` (array of objects, **required**): Array of layout specifications for each widget. Each item defines position and size in the dashboard grid
  - Each layout object contains:
    - `x` (integer): Horizontal position in the grid (0-based index from left)
    - `y` (integer): Vertical position in the grid (0-based index from top)
    - `w` (integer): Width in grid units
    - `h` (integer): Height in grid units
    - `scope_id` (integer, optional): The scope ID of the widget to update layout for

**Example Use Cases:**
- Reorganize dashboard layout
- Resize widgets
- Optimize dashboard appearance

---

## Metrics Tools

### 14. `get_metrics`
**Purpose:** Get a list of available metrics, filters, or groupby tags for building monitoring queries.

**Description:** This tool retrieves metadata about available metrics, filters, and grouping options in your Middleware.io environment. Use this to discover what metrics are available for a specific resource (e.g., host CPU metrics, database query metrics), what filters can be applied (e.g., filter by service name, environment), or what grouping dimensions are available (e.g., group by region, pod name).

**Common Use Cases:**
- Find all metrics for a specific resource type (host, container, database)
- Discover available filters for narrowing down data
- Get groupby tags for aggregating metrics by dimension

**Parameters:**
- `data_type` (string, **required**): Type of data to fetch. Valid values: 'metrics' (metric names), 'filters' (filter dimensions), 'groupby' (grouping tags)
- `widget_type` (string, **required**): Widget type for the query. Valid values: 'timeseries', 'list', 'queryValue'
- `kpi_type` (integer, optional): Single KPI type filter. 1=Metric (infrastructure/APM metrics), 2=Log (log data), 3=Trace (distributed tracing data)
- `kpi_types` (array of integers, optional): Array of KPI types to include. Use this for multi-type queries
- `resource` (string, optional): Resource name to filter by (e.g., 'host', 'process', 'container', 'database')
- `resources` (array of strings, optional): Array of resource names for multi-resource queries
- `metric` (string, optional): Specific metric name (required when dataType is 'groupby' to get groupby tags for this metric)
- `page` (integer, optional): Page number for paginated results (default: 1)
- `limit` (integer, optional): Number of items per page (default: 100)
- `search` (string, optional): Search term to filter metrics or resources by name (case-insensitive substring match)
- `exclude_metrics` (array of strings, optional): Array of metric names to exclude from results
- `mandatory_metrics` (array of strings, optional): Array of metric names to always include at the top of results
- `exclude_filters` (array of strings, optional): Array of filter names to exclude from results
- `mandatory_filters` (array of strings, optional): Array of filter names to always include at the top of results
- `filter_types` (array of integers, optional): Array of filter type IDs to include
- `return_only_mandatory_data` (boolean, optional): Set to true to return only mandatory metrics/filters without additional data

**Example Use Cases:**
- Discover CPU metrics for hosts
- Find available filters for containers
- Get groupby options for database metrics

---

### 15. `get_resources`
**Purpose:** Get a list of all available resource types in your Middleware.io environment.

**Description:** This tool returns all resource types that have data in your monitoring environment. Resources represent the entities you're monitoring (e.g., hosts, containers, databases, services, processes). Use this to discover what resource types are available before querying metrics for specific resources.

**Example resources:** host, container, pod, service, database, redis, mongodb, postgresql, mysql, nginx, etc.

**Parameters:** None

**Example Use Cases:**
- Discover available resource types
- Find monitored services
- Explore data sources

---

## Alert Tools

### 16. `list_alerts`
**Purpose:** Get a list of triggered alerts for a specific alert rule with pagination and sorting.

**Description:** This tool retrieves all alert instances that have been triggered for a specific alert rule. Each alert instance represents a time when the alert condition was met. Use this to review alert history, analyze alert patterns, or investigate recent incidents. Results can be paginated and ordered by various fields.

**Parameters:**
- `rule_id` (integer, **required**): The numeric ID of the alert rule to fetch alerts for
- `page` (integer, optional): Page number for pagination. 0-based index (default: 0 for first page)
- `order_by` (string, optional): Field name to sort results by (e.g., 'created_at', 'triggered_at', 'status'). Default: 'created_at' in descending order

**Example Use Cases:**
- Review recent alerts
- Analyze alert frequency
- Investigate incidents

---

### 17. `create_alert`
**Purpose:** Manually create a new alert instance for a specific alert rule.

**Description:** This tool allows you to programmatically create alert instances, typically used for custom alerting logic or integrations with external monitoring systems. The alert will be associated with an existing alert rule and will appear in the alerts list and trigger configured notification channels.

**Note:** In most cases, alerts are automatically created when rule conditions are met. Use this tool for custom alerting workflows or manual alert creation.

**Parameters:**
- `rule_id` (integer, **required**): The numeric ID of the alert rule this alert instance belongs to
- `title` (string, **required**): Alert title/summary describing what triggered (e.g., 'High CPU Usage on prod-server-01')
- `message` (string, optional): Detailed alert message with additional context and information
- `status` (integer, **required**): Alert status code. Typically: 0=OK/Resolved, 1=Warning, 2=Critical, 3=Unknown
- `value` (number, optional): The actual measured value that triggered the alert (e.g., 95.5 for 95.5% CPU usage)
- `threshold` (number, optional): The threshold value that was exceeded (e.g., 80.0 for 80% threshold)
- `operator` (string, optional): Comparison operator used (e.g., '>', '<', '>=', '<=', '==', '!=')
- `unit` (string, optional): Unit of measurement for the value (e.g., 'percent', 'ms', 'requests/sec', 'GB')
- `attributes` (object, optional): Additional key-value pairs with context (e.g., {'hostname': 'prod-01', 'region': 'us-east-1'})
- `project_uid` (string, optional): Project unique identifier if alert is project-specific
- `executor_id` (integer, optional): ID of the executor/rule evaluator that triggered the alert
- `triggered_at` (string, optional): Timestamp when the alert was triggered (ISO 8601 format)

**Example Use Cases:**
- Custom alerting workflows
- External system integrations
- Manual alert creation for testing

---

### 18. `get_alert_stats`
**Purpose:** Get aggregated statistics and metrics for alerts of a specific rule.

**Description:** This tool provides statistical analysis of alert history including counts by status (OK, Warning, Critical), counts by alert title, and time series data showing alert trends over time. Use this to understand alert patterns, identify frequently triggered alerts, and analyze alert distribution.

**Returns:**
- Count by status: Number of alerts in each status (OK, Warning, Critical)
- Count by title: Distribution of alerts by their titles
- Timeseries by title: Historical alert counts over time grouped by title

**Parameters:**
- `rule_id` (integer, **required**): The numeric ID of the alert rule to fetch statistics for

**Example Use Cases:**
- Analyze alert patterns
- Identify noisy alerts
- Report on alert trends

---

## Best Practices

### For AI Assistants Using These Tools:

1. **Start with Discovery Tools:**
   - Use `get_resources` to understand available data sources
   - Use `get_metrics` to discover available metrics for specific resources
   - Use `list_dashboards` to find existing dashboards

2. **Building Queries:**
   - Always specify required parameters
   - Use pagination for large result sets
   - Leverage search and filter capabilities

3. **Dashboard Management:**
   - Use descriptive labels and descriptions
   - Set appropriate visibility (public/private)
   - Organize with display scopes

4. **Widget Operations:**
   - Fetch widget data to understand current state
   - Use multi-widget data fetching for efficiency
   - Keep widget layouts organized

5. **Alert Analysis:**
   - Check alert stats before listing all alerts
   - Use pagination and ordering for alert lists
   - Understand alert status codes

---

## Support

For questions or issues:
- **Middleware Support:** support@middleware.io
- **Project Documentation:** See [README.md](README.md)

---

*Generated from Middleware API Swagger specification v1.0*
*Last Updated: November 2025*

