package dto

import "time"

// Device DTOs
type CreateDeviceRequest struct {
	ID           string `json:"id" binding:"required"`
	SerialNumber string `json:"serial_number"`
	Model        string `json:"model"`
	Name         string `json:"name"`
	Platform     string `json:"platform"`
	OwnerID      *uint  `json:"owner_id"`
}

type UpdateDeviceRequest struct {
	SerialNumber     *string `json:"serial_number"`
	Model            *string `json:"model"`
	Name             *string `json:"name"`
	Platform         *string `json:"platform"`
	Status           *string `json:"status"`
	ComplianceStatus *string `json:"compliance_status"`
	IsEnrolled       *bool   `json:"is_enrolled"`
	OwnerID          *uint   `json:"owner_id"`
	OsVersion        *string `json:"os_version"`
	DeviceType       *string `json:"device_type"`
}

type DeviceResponse struct {
	ID               string     `json:"id"`
	// UDID is the Apple MDM enrollment identifier. Empty for DEP devices not yet enrolled.
	UDID             string     `json:"udid,omitempty"`
	SerialNumber     string     `json:"serial_number"`
	Model            string     `json:"model"`
	Name             string     `json:"name"`
	Platform         string     `json:"platform"`
	Status           string     `json:"status"`
	ComplianceStatus string     `json:"compliance_status"`
	IsEnrolled       bool       `json:"is_enrolled"`
	OwnerID          *uint      `json:"owner_id"`
	OsVersion        string     `json:"os_version"`
	DeviceType       string     `json:"device_type"`
	MacAddress       string     `json:"mac_address,omitempty"`
	IpAddress        string     `json:"ip_address,omitempty"`
	BatteryLevel     float64    `json:"battery_level,omitempty"`
	StorageCapacity  uint64     `json:"storage_capacity,omitempty"`
	StorageUsed      uint64     `json:"storage_used,omitempty"`
	IsJailbroken     bool       `json:"is_jailbroken"`
	EnrollmentType   string     `json:"enrollment_type"`
	LastSeen         *time.Time `json:"last_seen"`
	EnrolledAt       *time.Time `json:"enrolled_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// DeviceGroup DTOs
type CreateDeviceGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateDeviceGroupRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type DeviceGroupResponse struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	DeviceCount int              `json:"device_count"`
	Devices     []DeviceResponse `json:"devices,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type ManageGroupDevicesRequest struct {
	DeviceIDs []string `json:"device_ids" binding:"required,min=1"`
}

// Device Action DTOs

// DeviceLockRequest contains optional parameters for the lock device action (EnableLostMode).
type DeviceLockRequest struct {
	Message     string `json:"message,omitempty" example:"Hãy trả thiết bị này cho tui"`
	PhoneNumber string `json:"phone_number,omitempty" example:"+84123456789"`
	Footnote    string `json:"footnote,omitempty" example:"THD"`
}

// DeviceWipeRequest contains optional parameters for the wipe device action.
type DeviceWipeRequest struct {
	PIN                    string `json:"pin,omitempty" example:"123456"`
	PreserveDataPlan       bool   `json:"preserve_data_plan,omitempty"`
	DisallowProximitySetup bool   `json:"disallow_proximity_setup,omitempty"`
	ObliterationBehavior   string `json:"obliteration_behavior,omitempty" enums:"Default,DoNotObliterate,ObliterateWithWarning,Always"`
}

// DeviceRestartRequest contains optional parameters for the restart device action.
type DeviceRestartRequest struct {
	NotifyUser bool `json:"notify_user,omitempty"`
}

// DeviceShutdownRequest is a placeholder for shutdown action (no parameters needed).
type DeviceShutdownRequest struct{}

// DeviceInstallProfileRequest contains parameters to install a profile on a device.
type DeviceInstallProfileRequest struct {
	ProfileID uint `json:"profile_id" binding:"required"`
}

// DeviceRemoveProfileRequest contains parameters to remove a profile from a device.
type DeviceRemoveProfileRequest struct {
	ProfileIdentifier string `json:"profile_identifier" binding:"required"`
}

// DeviceInfoRequest contains optional parameters for requesting device information.
type DeviceInfoRequest struct {
	Queries []string `json:"queries,omitempty"`
}

// DeviceActionResponse is the standard response for device action commands.
type DeviceActionResponse struct {
	CommandUUID string `json:"command_uuid"`
	RequestType string `json:"request_type"`
	Status      string `json:"status"`
	Message     string `json:"message,omitempty"`
}
