package dto

// DashboardStatsResponse represents the overall dashboard statistics
type DashboardStatsResponse struct {
	TotalDevices     int64 `json:"total_devices"`
	ActiveDevices    int64 `json:"active_devices"`
	TotalUsers       int64 `json:"total_users"`
	ActiveUsers      int64 `json:"active_users"`
	TotalAlerts      int64 `json:"total_alerts"`
	PendingAlerts    int64 `json:"pending_alerts"`
	TotalApps        int64 `json:"total_apps"`
	DeployedApps     int64 `json:"deployed_apps"`
	ComplianceRate   int   `json:"compliance_rate"`    // percentage 0-100
	NonCompliantRate int   `json:"non_compliant_rate"` // percentage 0-100
}

// DeviceStatsResponse represents device-related statistics
type DeviceStatsResponse struct {
	Total        int64            `json:"total"`
	Active       int64            `json:"active"`
	Inactive     int64            `json:"inactive"`
	Enrolled     int64            `json:"enrolled"`
	ByPlatform   map[string]int64 `json:"by_platform"` // ios, android, windows
	ByStatus     map[string]int64 `json:"by_status"`   // active, inactive, pending
	Compliant    int64            `json:"compliant"`
	NonCompliant int64            `json:"non_compliant"`
}

// AlertsSummaryResponse represents alerts summary statistics
type AlertsSummaryResponse struct {
	Total        int64            `json:"total"`
	Open         int64            `json:"open"`
	Acknowledged int64            `json:"acknowledged"`
	Resolved     int64            `json:"resolved"`
	BySeverity   map[string]int64 `json:"by_severity"` // critical, high, medium, low
	ByType       map[string]int64 `json:"by_type"`     // security, compliance, connectivity, etc.
}

// ChartDataPoint represents a single data point for charts
type ChartDataPoint struct {
	Label string `json:"label"`
	Value any    `json:"value"`
}

// ChartDataResponse represents chart data for visualization
type ChartDataResponse struct {
	ChartType string         `json:"chart_type"` // line, bar, pie, doughnut
	Title     string         `json:"title"`
	Labels    []string       `json:"labels"`
	Datasets  []ChartDataset `json:"datasets"`
}

// ChartDataset represents a dataset for charts
type ChartDataset struct {
	Label string `json:"label"`
	Data  []any  `json:"data"`
	Color string `json:"color,omitempty"`
}
