package dto

import "time"

type DEPTokenResponse struct {
	ID                string     `json:"id"`
	ConsumerKey       string    `json:"consumer_key,omitempty"`
	ConsumerSecret    string    `json:"consumer_secret,omitempty"`
	AccessToken       string    `json:"access_token,omitempty"`
	AccessSecret      string    `json:"access_secret,omitempty"`
	AccessTokenExpiry *time.Time `json:"access_token_expiry,omitempty"`
	ConfigBaseURL     string    `json:"config_base_url,omitempty"`
	TokenpkiCertPem   string    `json:"tokenpki_cert_pem,omitempty"`
	TokenpkiKey_pem   string    `json:"tokenpki_key_pem,omitempty"`
	SyncerCursor      string    `json:"syncer_cursor,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type DEPProfileRequest struct {
	ProfileName                string   `json:"profile_name"`
	AllowPairing               *bool    `json:"allow_pairing,omitempty"`
	AnchorCerts                []string `json:"anchor_certs,omitempty"`
	AutoAdvanceSetup           *bool    `json:"auto_advance_setup,omitempty"`
	AwaitDeviceConfigured      *bool    `json:"await_device_configured,omitempty"`
	ConfigurationWebURL        string   `json:"configuration_web_url,omitempty"`
	Department                 string   `json:"department,omitempty"`
	Devices                    []string `json:"devices,omitempty"`
	DoNotUseProfileFromBackup  *bool    `json:"do_not_use_profile_from_backup,omitempty"`
	IsReturnToService          *bool    `json:"is_return_to_service,omitempty"`
	IsMandatory                *bool    `json:"is_mandatory,omitempty"`
	IsMDMRemovable             *bool    `json:"is_mdm_removable,omitempty"`
	IsMultiUser                *bool    `json:"is_multi_user,omitempty"`
	IsSupervised               *bool    `json:"is_supervised,omitempty"`
	Language                   string   `json:"language,omitempty"`
	OrgMagic                   string   `json:"org_magic,omitempty"`
	Region                     string   `json:"region,omitempty"`
	SkipSetupItems             []string `json:"skip_setup_items,omitempty"`
	SupervisingHostCerts       []string `json:"supervising_host_certs,omitempty"`
	SupportEmailAddress        string   `json:"support_email_address,omitempty"`
	SupportPhoneNumber         string   `json:"support_phone_number,omitempty"`
	URL                        string   `json:"url,omitempty"`
	// ProfileData remains for any additional fields not explicitly mapped
	ProfileData map[string]interface{} `json:"profile_data,omitempty"`
}

type DEPProfileResponse struct {
	ProfileUUID                string                 `json:"profile_uuid"`
	Name                       string                 `json:"name"`
	AllowPairing               bool                   `json:"allow_pairing"`
	AnchorCerts                []string               `json:"anchor_certs"`
	AutoAdvanceSetup           bool                   `json:"auto_advance_setup"`
	AwaitDeviceConfigured      bool                   `json:"await_device_configured"`
	ConfigurationWebURL        string                 `json:"configuration_web_url"`
	Department                 string                 `json:"department"`
	Devices                    []string               `json:"devices"`
	DoNotUseProfileFromBackup  bool                   `json:"do_not_use_profile_from_backup"`
	IsReturnToService          bool                   `json:"is_return_to_service"`
	IsMandatory                bool                   `json:"is_mandatory"`
	IsMDMRemovable             bool                   `json:"is_mdm_removable"`
	IsMultiUser                bool                   `json:"is_multi_user"`
	IsSupervised               bool                   `json:"is_supervised"`
	Language                   string                 `json:"language"`
	OrgMagic                   string                 `json:"org_magic"`
	Region                     string                 `json:"region"`
	SkipSetupItems             []string               `json:"skip_setup_items"`
	SupervisingHostCerts       []string               `json:"supervising_host_certs"`
	SupportEmailAddress        string                 `json:"support_email_address"`
	SupportPhoneNumber         string                 `json:"support_phone_number"`
	URL                        string                 `json:"url"`
	ProfileData                map[string]interface{} `json:"profile_data"`
}

type DeviceResponse struct {
	ID           string    `json:"id"`
	SerialNumber string    `json:"serial_number"`
	Model        string    `json:"model"`
	OwnerID      uint      `json:"owner_id"`
	IsEnrolled   bool      `json:"is_enrolled"`
	Name         string    `json:"name"`
	LastSync     time.Time `json:"last_sync"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
