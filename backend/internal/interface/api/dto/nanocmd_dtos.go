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
	Topic               string               `json:"topic"`
	EventID             string               `json:"event_id"`
	CreatedAt           time.Time            `json:"created_at"`
	AcknowledgeEvent    map[string]any       `json:"acknowledge_event"`
	Checkin_event       map[string]any       `json:"checkin_event"`
	DeviceResponseEvent *DeviceResponseEvent `json:"device_response_event,omitempty"`
}

// DeviceResponseEvent represents the device_response_event in DEP webhooks
type DeviceResponseEvent struct {
	DEPName        string                 `json:"dep_name"`
	DeviceResponse *DEPSyncDeviceResponse `json:"device_response,omitempty"`
}

// DEPSyncDeviceResponse represents the device_response in DEP webhooks (FetchDevices/SyncDevices)
type DEPSyncDeviceResponse struct {
	Cursor       string      `json:"cursor"`
	FetchedUntil string      `json:"fetched_until"`
	MoreToFollow bool        `json:"more_to_follow"`
	Devices      []DEPDevice `json:"devices"`
}

// DEPDevice represents a device in DEP fetch/sync response
type DEPDevice struct {
	SerialNumber       string `json:"serial_number"`
	Model              string `json:"model"`
	Description        string `json:"description"`
	Color              string `json:"color"`
	AssetTag           string `json:"asset_tag"`
	ProfileStatus      string `json:"profile_status"`
	ProfileUUID        string `json:"profile_uuid"`
	ProfileAssignTime  string `json:"profile_assign_time"`
	ProfilePushTime    string `json:"profile_push_time"`
	DeviceAssignedBy   string `json:"device_assigned_by"`
	DeviceAssignedDate string `json:"device_assigned_date"`
	OS                 string `json:"os"`
	DeviceFamily       string `json:"device_family"`
	OpType             string `json:"op_type"`
	OpDate             string `json:"op_date"`
}

// DEPProfileAssignRequest is the request body for POST /proxy/{name}/profile/devices
type DEPProfileAssignRequest struct {
	ProfileUUID string   `json:"profile_uuid"`
	Devices     []string `json:"devices"`
}

// DEPProfileAssignResponse is the response from POST /proxy/{name}/profile/devices
type DEPProfileAssignResponse struct {
	Devices map[string]string `json:"devices"`
}

// NanoCMD Error
type NanoCMDError struct {
	Error string `json:"error"`
}
