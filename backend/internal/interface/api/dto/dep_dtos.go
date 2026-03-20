package dto

import "time"

type DEPTokenResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	P7mFilePath string    `json:"p7m_file_path"`
	Expiry      time.Time `json:"expiry"`
	LastUsed    time.Time `json:"last_used"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DEPProfileRequest struct {
	ProfileName string                 `json:"profile_name"`
	ProfileData map[string]interface{} `json:"profile_data"`
}

type DEPProfileResponse struct {
	ProfileUUID string                 `json:"profile_uuid"`
	Name        string                 `json:"name"`
	ProfileData map[string]interface{} `json:"profile_data"`
}

type DeviceResponse struct {
	ID           uint      `json:"id"`
	SerialNumber string    `json:"serial_number"`
	Model        string    `json:"model"`
	OwnerID      uint      `json:"owner_id"`
	IsEnrolled   bool      `json:"is_enrolled"`
	Name         string    `json:"name"`
	LastSync     time.Time `json:"last_sync"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
