package dto

import "time"

// CreateProfileRequest represents the request to create a profile
type CreateProfileRequest struct {
	Name             string         `json:"name" binding:"required,max=255"`
	Platform         string         `json:"platform,omitempty"` // ios, android, windows, macos, all
	Scope            string         `json:"scope,omitempty"`    // device, user, group
	SecuritySettings map[string]any `json:"security_settings,omitempty"`
	NetworkConfig    map[string]any `json:"network_config,omitempty"`
	Restrictions     map[string]any `json:"restrictions,omitempty"`
	ContentFilter    map[string]any `json:"content_filter,omitempty"`
	ComplianceRules  map[string]any `json:"compliance_rules,omitempty"`
	Payloads         map[string]any `json:"payloads,omitempty"`
}

// UpdateProfileRequest represents the request to update a profile
type UpdateProfileRequest struct {
	Name             *string        `json:"name,omitempty" binding:"omitempty,max=255"`
	Platform         *string        `json:"platform,omitempty"`
	Scope            *string        `json:"scope,omitempty"`
	SecuritySettings map[string]any `json:"security_settings,omitempty"`
	NetworkConfig    map[string]any `json:"network_config,omitempty"`
	Restrictions     map[string]any `json:"restrictions,omitempty"`
	ContentFilter    map[string]any `json:"content_filter,omitempty"`
	ComplianceRules  map[string]any `json:"compliance_rules,omitempty"`
	Payloads         map[string]any `json:"payloads,omitempty"`
}

// UpdateProfileStatusRequest represents the request to update profile status
type UpdateProfileStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active draft archived"`
}

// UpdateSecuritySettingsRequest represents security settings
type UpdateSecuritySettingsRequest struct {
	Passcode   *PasscodeSettings   `json:"passcode,omitempty"`
	Encryption *EncryptionSettings `json:"encryption,omitempty"`
	Biometrics *BiometricsSettings `json:"biometrics,omitempty"`
	ScreenLock *ScreenLockSettings `json:"screen_lock,omitempty"`
}

type PasscodeSettings struct {
	AutoLock            int  `json:"auto_lock,omitempty"` // minutes
	MinLength           int  `json:"min_length,omitempty"`
	RetryLimit          int  `json:"retry_limit,omitempty"`
	RequireAlphanumeric bool `json:"require_alphanumeric,omitempty"`
}

type EncryptionSettings struct {
	Required bool `json:"required,omitempty"`
}

type BiometricsSettings struct {
	FaceIDEnabled      bool `json:"face_id_enabled,omitempty"`
	FingerprintEnabled bool `json:"fingerprint_enabled,omitempty"`
}

type ScreenLockSettings struct {
	Enabled bool `json:"enabled,omitempty"`
	Timeout int  `json:"timeout,omitempty"` // seconds
}

// UpdateNetworkConfigRequest represents network configuration
type UpdateNetworkConfigRequest struct {
	Wifi  *WifiConfig  `json:"wifi,omitempty"`
	VPN   *VPNConfig   `json:"vpn,omitempty"`
	Proxy *ProxyConfig `json:"proxy,omitempty"`
}

type WifiConfig struct {
	SSID     string `json:"ssid,omitempty"`
	Password string `json:"password,omitempty"`
	AutoJoin bool   `json:"auto_join,omitempty"`
}

type VPNConfig struct {
	Type         string         `json:"type,omitempty"` // IKEv2, L2TP, etc.
	ServerConfig map[string]any `json:"server_config,omitempty"`
}

type ProxyConfig struct {
	Enabled  bool   `json:"enabled,omitempty"`
	Server   string `json:"server,omitempty"`
	Port     int    `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// UpdateRestrictionsRequest represents device restrictions
type UpdateRestrictionsRequest struct {
	CameraEnabled             *bool `json:"camera_enabled,omitempty"`
	BluetoothEnabled          *bool `json:"bluetooth_enabled,omitempty"`
	AirdropEnabled            *bool `json:"airdrop_enabled,omitempty"`       // iOS only
	USBDebuggingEnabled       *bool `json:"usb_debugging_enabled,omitempty"` // Android only
	ExternalAppInstallAllowed *bool `json:"external_app_install_allowed,omitempty"`
}

// UpdateContentFilterRequest represents web & content filtering
type UpdateContentFilterRequest struct {
	BlockedWebsites []string `json:"blocked_websites,omitempty"`
	AllowedDomains  []string `json:"allowed_domains,omitempty"`
	SafeBrowsing    *bool    `json:"safe_browsing,omitempty"`
}

// UpdateComplianceRulesRequest represents compliance rules
type UpdateComplianceRulesRequest struct {
	SendAlert   *bool `json:"send_alert,omitempty"`
	LockDevice  *bool `json:"lock_device,omitempty"`
	BlockAccess *bool `json:"block_access,omitempty"`
}

// AssignProfileRequest represents the request to assign a profile
type AssignProfileRequest struct {
	TargetType   string     `json:"target_type" binding:"required,oneof=device group"`
	DeviceID     *string    `json:"device_id,omitempty"`
	GroupID      *uint      `json:"group_id,omitempty"`
	ScheduleType string     `json:"schedule_type,omitempty"` // immediate, scheduled
	ScheduledAt  *time.Time `json:"scheduled_at,omitempty"`
}

// ProfileResponse represents the response for a profile
type ProfileResponse struct {
	ID               uint           `json:"id"`
	Name             string         `json:"name"`
	Platform         string         `json:"platform"`
	Scope            string         `json:"scope"`
	Status           string         `json:"status"`
	SecuritySettings map[string]any `json:"security_settings,omitempty"`
	NetworkConfig    map[string]any `json:"network_config,omitempty"`
	Restrictions     map[string]any `json:"restrictions,omitempty"`
	ContentFilter    map[string]any `json:"content_filter,omitempty"`
	ComplianceRules  map[string]any `json:"compliance_rules,omitempty"`
	Payloads         map[string]any `json:"payloads,omitempty"`
	Version          int            `json:"version"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// ProfileAssignmentResponse represents a profile assignment
type ProfileAssignmentResponse struct {
	ID           uint       `json:"id"`
	ProfileID    uint       `json:"profile_id"`
	TargetType   string     `json:"target_type"`
	DeviceID     *string    `json:"device_id,omitempty"`
	GroupID      *uint      `json:"group_id,omitempty"`
	ScheduleType string     `json:"schedule_type"`
	ScheduledAt  *time.Time `json:"scheduled_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// ProfileVersionResponse represents a profile version
type ProfileVersionResponse struct {
	ID          uint           `json:"id"`
	ProfileID   uint           `json:"profile_id"`
	Version     int            `json:"version"`
	Data        map[string]any `json:"data,omitempty"`
	ChangeNotes string         `json:"change_notes,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
}

// ProfileDeploymentStatusResponse represents deployment status
type ProfileDeploymentStatusResponse struct {
	ID           uint       `json:"id"`
	ProfileID    uint       `json:"profile_id"`
	DeviceID     string     `json:"device_id"`
	Status       string     `json:"status"` // pending, success, failed
	ErrorMessage string     `json:"error_message,omitempty"`
	AppliedAt    *time.Time `json:"applied_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
