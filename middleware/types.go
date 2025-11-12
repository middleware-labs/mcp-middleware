package middleware

// Dashboard/Report types
type Report struct {
	ID           int         `json:"id,omitempty"`
	Key          string      `json:"key,omitempty"`
	Label        string      `json:"label"`
	Description  string      `json:"description,omitempty"`
	DisplayScope string      `json:"display_scope,omitempty"`
	Visibility   string      `json:"visibility"`
	Favorite     bool        `json:"favorite,omitempty"`
	AccountID    int         `json:"account_id,omitempty"`
	ProjectID    int         `json:"project_id,omitempty"`
	UserID       int         `json:"user_id,omitempty"`
	ViewCount    int         `json:"view_count,omitempty"`
	MetaData     any         `json:"meta_data,omitempty"`
	User         *ReportUser `json:"user,omitempty"`
	CreatedAt    string      `json:"created_at,omitempty"`
	UpdatedAt    string      `json:"updated_at,omitempty"`
}

type ReportUser struct {
	Name string `json:"name"`
}

type ReportListResponse struct {
	Reports []Report `json:"reports"`
	Total   int      `json:"total"`
	Limit   int      `json:"limit"`
	Offset  int      `json:"offset"`
}

type UpsertReportRequest struct {
	ID           int    `json:"id,omitempty"`
	Key          string `json:"key,omitempty"`
	Label        string `json:"label"`
	Description  string `json:"description,omitempty"`
	DisplayScope string `json:"display_scope,omitempty"`
	Visibility   string `json:"visibility"`
	MetaData     any    `json:"metaData,omitempty"`
}

// Widget types
type Widget struct {
	ID          int          `json:"id,omitempty"`
	Key         string       `json:"key,omitempty"`
	Label       string       `json:"label,omitempty"`
	Config      any          `json:"config,omitempty"`
	MetaData    any          `json:"meta_data,omitempty"`
	Scope       *WidgetScope `json:"scope,omitempty"`
	Visibility  string       `json:"visibility,omitempty"`
	Status      string       `json:"status,omitempty"`
	AccountID   int          `json:"account_id,omitempty"`
	ProjectID   int          `json:"project_id,omitempty"`
	UserID      int          `json:"user_id,omitempty"`
	WidgetAppID int          `json:"widget_app_id,omitempty"`
	DatasetID   int          `json:"dataset_id,omitempty"`
	CreatedAt   string       `json:"created_at,omitempty"`
	UpdatedAt   string       `json:"updated_at,omitempty"`
}

type WidgetScope struct {
	ID           int    `json:"id,omitempty"`
	BuilderID    int    `json:"builder_id,omitempty"`
	ReportID     int    `json:"report_id,omitempty"`
	DisplayScope string `json:"display_scope,omitempty"`
	OrderID      int    `json:"order_id,omitempty"`
	ProjectID    int    `json:"project_id,omitempty"`
	MetaData     any    `json:"meta_data,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

type CustomWidget struct {
	BuilderID               int                 `json:"builderId,omitempty"`
	Key                     string              `json:"key,omitempty"`
	Label                   string              `json:"label,omitempty"`
	Description             string              `json:"description,omitempty"`
	BuilderConfig           any                 `json:"builderConfig,omitempty"`
	BuilderMetaData         *WidgetMetaData     `json:"builderMetaData,omitempty"`
	BuilderViewOptions      *BuilderViewOptions `json:"builderViewOptions,omitempty"`
	Layout                  *LayoutItem         `json:"layout,omitempty"`
	Params                  []RequestParam      `json:"params,omitempty"`
	ScopeID                 int                 `json:"scopeId,omitempty"`
	WidgetAppID             int                 `json:"widgetAppId,omitempty"`
	UseV2                   bool                `json:"useV2,omitempty"`
	KeepOldData             bool                `json:"keepOldData,omitempty"`
	ReturnOnlyFormulaResult bool                `json:"returnOnlyFormulaResult,omitempty"`
}

type WidgetMetaData struct {
	ChartType              string   `json:"chartType,omitempty"`
	ColorScheme            string   `json:"colorScheme,omitempty"`
	DefaultKey             string   `json:"default_key,omitempty"`
	Description            string   `json:"description,omitempty"`
	DisplayPreference      string   `json:"display_preference,omitempty"`
	ExpandedLegendColumns  []string `json:"expanedLegendColumns,omitempty"`
	GroupName              string   `json:"group_name,omitempty"`
	GroupOrder             int      `json:"group_order,omitempty"`
	IsDefault              bool     `json:"is_default,omitempty"`
	LegendType             string   `json:"legendType,omitempty"`
	LineStroke             string   `json:"lineStroke,omitempty"`
	LineStyle              string   `json:"lineStyle,omitempty"`
	YAxisAlwaysIncludeZero bool     `json:"yAxisAlwaysIncludeZero,omitempty"`
	YAxisMax               int      `json:"yAxisMax,omitempty"`
	YAxisMin               int      `json:"yAxisMin,omitempty"`
	YAxisType              string   `json:"yAxisType,omitempty"`
}

type BuilderViewOptions struct {
	DisplayScope string      `json:"displayScope,omitempty"`
	Report       *ReportView `json:"report,omitempty"`
	Resource     *Resource   `json:"resource,omitempty"`
}

type ReportView struct {
	ReportID          int    `json:"reportId,omitempty"`
	ReportKey         string `json:"reportKey,omitempty"`
	ReportName        string `json:"reportName,omitempty"`
	ReportDescription string `json:"reportDescription,omitempty"`
	Metadata          any    `json:"metadata,omitempty"`
}

type Resource struct {
	Name string `json:"name"`
}

type LayoutItem struct {
	X             int      `json:"x,omitempty"`
	Y             int      `json:"y,omitempty"`
	W             int      `json:"w,omitempty"`
	H             int      `json:"h,omitempty"`
	ScopeID       int      `json:"_scope_id,omitempty"`
	ResizeHandles []string `json:"resizeHandles,omitempty"`
}

type RequestParam struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type LayoutRequest struct {
	Layouts []LayoutItem `json:"layouts"`
}

type BuilderDataResponse struct {
	Key                 string     `json:"key,omitempty"`
	ChartData           any        `json:"chart_data,omitempty"`
	ChartDataV2         any        `json:"chart_data_v2,omitempty"`
	QueryData           any        `json:"query_data,omitempty"`
	TimeRange           *TimeRange `json:"time_range,omitempty"`
	NotAvailableMetrics []string   `json:"not_available_metrics,omitempty"`
	Error               string     `json:"error,omitempty"`
	ErrorDesc           string     `json:"error_desc,omitempty"`
}

type TimeRange struct {
	FromTs   int64 `json:"fromTs"`
	ToTs     int64 `json:"toTs"`
	Interval int64 `json:"interval"`
}

// Metrics types
type MetricsV2Request struct {
	DataType                string   `json:"dataType"`
	WidgetType              string   `json:"widgetType"`
	KpiType                 int      `json:"kpiType,omitempty"`
	KpiTypes                []int    `json:"kpiTypes,omitempty"`
	Resource                string   `json:"resource,omitempty"`
	Resources               []string `json:"resources,omitempty"`
	Metric                  string   `json:"metric,omitempty"`
	Page                    int      `json:"page,omitempty"`
	Limit                   int      `json:"limit,omitempty"`
	Search                  string   `json:"search,omitempty"`
	ExcludeMetrics          []string `json:"excludeMetrics,omitempty"`
	MandatoryMetrics        []string `json:"mandatoryMetrics,omitempty"`
	ExcludeFilters          []string `json:"excludeFilters,omitempty"`
	MandatoryFilters        []string `json:"mandatoryFilters,omitempty"`
	FilterTypes             []int    `json:"filterTypes,omitempty"`
	ReturnOnlyMandatoryData bool     `json:"returnOnlyMandatoryData,omitempty"`
}

type MetricsV2Response struct {
	Items []map[string]any `json:"items"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}

// Alert types
type Alert struct {
	ID          int               `json:"id,omitempty"`
	RuleID      int               `json:"rule_id"`
	ExecutorID  int               `json:"executor_id,omitempty"`
	ProjectUID  string            `json:"project_uid,omitempty"`
	Title       string            `json:"title"`
	Message     string            `json:"message,omitempty"`
	Status      int               `json:"status"`
	Value       float64           `json:"value,omitempty"`
	Threshold   float64           `json:"threshold,omitempty"`
	Operator    string            `json:"operator,omitempty"`
	Unit        string            `json:"unit,omitempty"`
	Attributes  map[string]string `json:"attributes,omitempty"`
	AttributesB any               `json:"attributesb,omitempty"`
	TotalCount  int               `json:"total_count,omitempty"`
	TriggeredAt string            `json:"triggered_at,omitempty"`
	CreatedAt   string            `json:"created_at,omitempty"`
}

type NewAlert struct {
	RuleID      int               `json:"rule_id"`
	ExecutorID  int               `json:"executor_id,omitempty"`
	ProjectUID  string            `json:"project_uid,omitempty"`
	Title       string            `json:"title"`
	Message     string            `json:"message,omitempty"`
	Status      int               `json:"status"`
	Value       float64           `json:"value,omitempty"`
	Threshold   float64           `json:"threshold,omitempty"`
	Operator    string            `json:"operator,omitempty"`
	Unit        string            `json:"unit,omitempty"`
	Attributes  map[string]string `json:"attributes,omitempty"`
	TriggeredAt string            `json:"triggered_at,omitempty"`
	CreatedAt   string            `json:"created_at,omitempty"`
}

type AlertsResponse struct {
	Alerts            []ViewModelAlert `json:"alerts"`
	Columns           []Column         `json:"columns"`
	LatestStatus      int              `json:"latest_status"`
	LatestTriggeredAt string           `json:"latest_triggered_at"`
}

type ViewModelAlert struct {
	ID          int               `json:"id"`
	ExecutorID  int               `json:"executor_id"`
	Title       string            `json:"title"`
	Message     string            `json:"message"`
	Status      int               `json:"status"`
	Value       float64           `json:"value"`
	Threshold   float64           `json:"threshold"`
	Operator    string            `json:"operator"`
	Unit        string            `json:"unit"`
	Attributes  map[string]string `json:"attributes"`
	TotalCount  int               `json:"total_count"`
	TriggeredAt string            `json:"triggered_at"`
}

type Column struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

type QueryRequest struct {
	Queries []Query `json:"queries"`
}

type Query struct {
	ChartType string         `json:"chartType"`
	Columns   []string       `json:"columns"`
	Resources []string       `json:"resources"`
	TimeRange QueryTimeRange `json:"timeRange"`
	Filters   map[string]any `json:"filters,omitempty"`
	GroupBy   []string       `json:"groupBy,omitempty"`
}

type QueryTimeRange struct {
	From int64 `json:"from"`
	To   int64 `json:"to"`
}

type QueryResponse struct {
	QueryResults []QueryResult `json:"query_results"`
}

type QueryResult struct {
	QueryData QueryData `json:"query_data"`
}

type QueryData struct {
	Columns []QueryColumn    `json:"columns"`
	Data    []map[string]any `json:"data"`
}

type QueryColumn struct {
	Accessor string `json:"accessor"`
	Order    int    `json:"order"`
	Sort     string `json:"sort"`
	IsMetric bool   `json:"isMetric"`
}

type StatsResponse struct {
	CountByStatus     []CountBy `json:"count_by_status"`
	CountByTitle      []CountBy `json:"count_by_title"`
	TimeseriesByTitle []CountBy `json:"timeseries_by_title"`
}

type CountBy struct {
	Name      string  `json:"name"`
	Status    int     `json:"status"`
	Value     float64 `json:"value"`
	Timestamp string  `json:"timestamp"`
}
