package dto

import "time"

// NanoCMD Version
type NanoCMDVersionResponse struct {
	Version string `json:"version"`
}

// NanoCMD Workflow
type NanoCMDWorkflowStartResponse struct {
	InstanceID string `json:"instance_id"`
}

// NanoCMD Event
type EventSubscription struct {
	Event        string `json:"event"`
	Workflow     string `json:"workflow"`
	Context      string `json:"context,omitempty"`
	EventContext string `json:"event_context,omitempty"`
}

// NanoCMD Profile
type NanoCMDProfile struct {
	Identifier string `json:"identifier"`
	UUID       string `json:"uuid"`
}

// NanoCMD CMDPlan
type CMDPlan struct {
	ProfileNames     []string `json:"profile_names"`
	ManifestURLs     []string `json:"manifest_urls"`
	DeviceConfigured bool     `json:"device_configured"`
}

// NanoCMD Inventory
type NanoCMDInventoryResponse map[string]NanoCMDInventoryDevice

type NanoCMDInventoryDevice struct {
	SerialNumber string `json:"serial_number"`
	Model        string `json:"model"`
}

// NanoCMD Webhook
type NanoCMDWebhook struct {
	Topic            string                 `json:"topic"`
	EventID          string                 `json:"event_id"`
	CreatedAt        time.Time              `json:"created_at"`
	AcknowledgeEvent map[string]interface{} `json:"acknowledge_event"`
	Checkin_event    map[string]interface{} `json:"checkin_event"`
}

// NanoCMD Error
type NanoCMDError struct {
	Error string `json:"error"`
}
