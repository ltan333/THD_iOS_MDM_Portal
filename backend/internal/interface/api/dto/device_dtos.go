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
