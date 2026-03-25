package dto

import "time"

type APNSConfigResponse struct {
	ID        string    `json:"id"`
	CertPEM   string    `json:"cert_pem"`
	KeyPEM    string    `json:"key_pem"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PushCertResponse is the response for APNs push certificate upload/retrieval
type PushCertResponse struct {
	Topic    string    `json:"topic"`
	NotAfter time.Time `json:"not_after"`
}

// APIResult is the generic result for push and enqueue commands
type APIResult struct {
	NoPush       bool                  `json:"no_push"`
	PushError    string                `json:"push_error,omitempty"`
	CommandError string                `json:"command_error,omitempty"`
	CommandUUID  string                `json:"command_uuid,omitempty"`
	RequestType  string                `json:"request_type,omitempty"`
	Status       map[string]PushStatus `json:"status,omitempty"`
}

// PushStatus represents the status of a push/command for a specific ID
type PushStatus struct {
	PushError    string `json:"push_error,omitempty"`
	PushResult   string `json:"push_result,omitempty"`
	CommandError string `json:"command_error,omitempty"`
}

// NanoMDMErrorResponse matches the error response from NanoMDM
type NanoMDMErrorResponse struct {
	Error string `json:"error"`
}

// EscrowKeyUnlockRequest is the request body for Apple Escrow Key Unlock
type EscrowKeyUnlockRequest struct {
	Topic       string `form:"topic" binding:"required"`
	Serial      string `form:"serial" binding:"required"`
	ProductType string `form:"productType" binding:"required"`
	EscrowKey   string `form:"escrowKey" binding:"required"`
	OrgName     string `form:"orgName" binding:"required"`
	Guid        string `form:"guid" binding:"required"`
	IMEI        string `form:"imei"`
	IMEI2       string `form:"imei2"`
	MEID        string `form:"meid"`
}

// NanoMDMVersionResponse is the response for NanoMDM version
type NanoMDMVersionResponse struct {
	Version string `json:"version"`
}
