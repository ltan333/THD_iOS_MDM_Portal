package dto

import "time"

// AlertResponse represents the response for an alert
type AlertResponse struct {
	ID             uint                   `json:"id"`
	Severity       string                 `json:"severity"` // critical, high, medium, low
	Title          string                 `json:"title"`
	Type           string                 `json:"type"` // security, compliance, connectivity, application, device_health
	Status         string                 `json:"status"` // open, acknowledged, resolved
	DeviceID       string                 `json:"device_id,omitempty"`
	UserID         *uint                  `json:"user_id,omitempty"`
	Details        map[string]interface{} `json:"details,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	AcknowledgedAt *time.Time             `json:"acknowledged_at,omitempty"`
	ResolvedAt     *time.Time             `json:"resolved_at,omitempty"`
}

// CreateAlertRequest represents the request to create an alert
type CreateAlertRequest struct {
	Severity string                 `json:"severity" binding:"required,oneof=critical high medium low"`
	Title    string                 `json:"title" binding:"required,max=255"`
	Type     string                 `json:"type" binding:"required,oneof=security compliance connectivity application device_health"`
	DeviceID string                 `json:"device_id,omitempty"`
	UserID   *uint                  `json:"user_id,omitempty"`
	Details  map[string]interface{} `json:"details,omitempty"`
}

// BulkResolveAlertsRequest represents the request to bulk resolve alerts
type BulkResolveAlertsRequest struct {
	AlertIDs []uint `json:"alert_ids" binding:"required,min=1"`
}

// AlertActionRequest represents a quick action request
type AlertActionRequest struct {
	PolicyID *uint   `json:"policy_id,omitempty"` // For push policy action
	Message  string  `json:"message,omitempty"`   // For send message action
}

// AlertRuleResponse represents the response for an alert rule
type AlertRuleResponse struct {
	ID          uint                   `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Condition   map[string]interface{} `json:"condition,omitempty"`
	Actions     map[string]interface{} `json:"actions,omitempty"`
	Enabled     bool                   `json:"enabled"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// CreateAlertRuleRequest represents the request to create an alert rule
type CreateAlertRuleRequest struct {
	Name        string                 `json:"name" binding:"required,max=255"`
	Description string                 `json:"description,omitempty" binding:"max=500"`
	Condition   map[string]interface{} `json:"condition,omitempty"`
	Actions     map[string]interface{} `json:"actions,omitempty"`
	Enabled     *bool                  `json:"enabled,omitempty"`
}

// UpdateAlertRuleRequest represents the request to update an alert rule
type UpdateAlertRuleRequest struct {
	Name        *string                `json:"name,omitempty" binding:"omitempty,max=255"`
	Description *string                `json:"description,omitempty" binding:"max=500"`
	Condition   map[string]interface{} `json:"condition,omitempty"`
	Actions     map[string]interface{} `json:"actions,omitempty"`
	Enabled     *bool                  `json:"enabled,omitempty"`
}

// AlertsSummaryResponse is defined in dashboard_dtos.go
