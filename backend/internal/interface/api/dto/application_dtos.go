package dto

import "time"

type CreateApplicationRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	BundleID    string `json:"bundle_id" binding:"required,max=255"`
	Platform    string `json:"platform" binding:"required,oneof=ios android windows macos"`
	Type        string `json:"type" binding:"required,oneof=app_store enterprise web_clip"`
	Description string `json:"description,omitempty"`
	IconURL     string `json:"icon_url,omitempty"`
}

type UpdateApplicationRequest struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,max=255"`
	Platform    *string `json:"platform,omitempty" binding:"omitempty,oneof=ios android windows macos"`
	Type        *string `json:"type,omitempty" binding:"omitempty,oneof=app_store enterprise web_clip"`
	Description *string `json:"description,omitempty"`
	IconURL     *string `json:"icon_url,omitempty"`
}

type CreateAppVersionRequest struct {
	ApplicationID    uint                   `json:"application_id" binding:"required"`
	Version          string                 `json:"version" binding:"required"`
	BuildNumber      string                 `json:"build_number" binding:"required"`
	MinimumOSVersion string                 `json:"minimum_os_version,omitempty"`
	FileURL          string                 `json:"file_url,omitempty"`
	Size             int64                  `json:"size,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

type UpdateAppVersionRequest struct {
	Version          *string                `json:"version,omitempty"`
	BuildNumber      *string                `json:"build_number,omitempty"`
	MinimumOSVersion *string                `json:"minimum_os_version,omitempty"`
	FileURL          *string                `json:"file_url,omitempty"`
	Size             *int64                 `json:"size,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

type CreateAppDeploymentRequest struct {
	AppVersionID uint   `json:"app_version_id" binding:"required"`
	TargetType   string `json:"target_type" binding:"required,oneof=device group user"`
	TargetID     string `json:"target_id" binding:"required"`
}

type ApplicationResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	BundleID    string               `json:"bundle_id"`
	Platform    string               `json:"platform"`
	Type        string               `json:"type"`
	Description string               `json:"description,omitempty"`
	IconURL     string               `json:"icon_url,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	Versions    []AppVersionResponse `json:"versions,omitempty"`
}

type AppVersionResponse struct {
	ID               uint                   `json:"id"`
	ApplicationID    uint                   `json:"application_id"`
	Version          string                 `json:"version"`
	BuildNumber      string                 `json:"build_number"`
	MinimumOSVersion string                 `json:"minimum_os_version,omitempty"`
	FileURL          string                 `json:"file_url,omitempty"`
	Size             int64                  `json:"size,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

type AppDeploymentResponse struct {
	ID           uint       `json:"id"`
	AppVersionID uint       `json:"app_version_id"`
	TargetType   string     `json:"target_type"`
	TargetID     string     `json:"target_id"`
	Status       string     `json:"status"`
	ErrorMessage string     `json:"error_message,omitempty"`
	InstalledAt  *time.Time `json:"installed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
