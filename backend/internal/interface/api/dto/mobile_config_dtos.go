package dto

import "time"

type CreateMobileConfigPropertyRequest struct {
	Key       string                 `json:"key" binding:"required"`
	ValueJSON map[string]interface{} `json:"value_json" binding:"required"`
}

type CreateMobileConfigPayloadRequest struct {
	PayloadDescription  string                              `json:"payload_description,omitempty"`
	PayloadDisplayName  string                              `json:"payload_display_name" binding:"required"`
	PayloadIdentifier   string                              `json:"payload_identifier" binding:"required"`
	PayloadOrganization string                              `json:"payload_organization,omitempty"`
	PayloadType         string                              `json:"payload_type" binding:"required"`
	PayloadVersion      int                                 `json:"payload_version,omitempty" binding:"omitempty,min=1"`
	Properties          []CreateMobileConfigPropertyRequest `json:"properties,omitempty"`
}

type CreateMobileConfigRequest struct {
	Name                     string                             `json:"name" binding:"required"`
	PayloadIdentifier        string                             `json:"payload_identifier" binding:"required"`
	PayloadType              string                             `json:"payload_type" binding:"required"`
	PayloadDisplayName       string                             `json:"payload_display_name" binding:"required"`
	PayloadDescription       string                             `json:"payload_description,omitempty"`
	PayloadOrganization      string                             `json:"payload_organization,omitempty"`
	PayloadVersion           int                                `json:"payload_version,omitempty" binding:"omitempty,min=1"`
	PayloadRemovalDisallowed bool                               `json:"payload_removal_disallowed,omitempty"`
	Payloads                 []CreateMobileConfigPayloadRequest `json:"payloads" binding:"required,min=1"`
}

type UpdateMobileConfigPropertyRequest = CreateMobileConfigPropertyRequest

type UpdateMobileConfigPayloadRequest = CreateMobileConfigPayloadRequest

type UpdateMobileConfigRequest struct {
	Name                     string                             `json:"name" binding:"required"`
	PayloadIdentifier        string                             `json:"payload_identifier" binding:"required"`
	PayloadType              string                             `json:"payload_type" binding:"required"`
	PayloadDisplayName       string                             `json:"payload_display_name" binding:"required"`
	PayloadDescription       string                             `json:"payload_description,omitempty"`
	PayloadOrganization      string                             `json:"payload_organization,omitempty"`
	PayloadVersion           int                                `json:"payload_version,omitempty" binding:"omitempty,min=1"`
	PayloadRemovalDisallowed bool                               `json:"payload_removal_disallowed,omitempty"`
	Payloads                 []UpdateMobileConfigPayloadRequest `json:"payloads" binding:"required,min=1"`
}

type MobileConfigPropertyResponse struct {
	ID        uint                   `json:"id"`
	Key       string                 `json:"key"`
	ValueJSON map[string]interface{} `json:"value_json,omitempty"`
}

type MobileConfigPayloadResponse struct {
	ID                  uint                           `json:"id"`
	PayloadDescription  string                         `json:"payload_description,omitempty"`
	PayloadDisplayName  string                         `json:"payload_display_name"`
	PayloadIdentifier   string                         `json:"payload_identifier"`
	PayloadOrganization string                         `json:"payload_organization,omitempty"`
	PayloadType         string                         `json:"payload_type"`
	PayloadUUID         string                         `json:"payload_uuid"`
	PayloadVersion      int                            `json:"payload_version"`
	Properties          []MobileConfigPropertyResponse `json:"properties,omitempty"`
}

type MobileConfigResponse struct {
	ID                       uint                          `json:"id"`
	Name                     string                        `json:"name"`
	PayloadIdentifier        string                        `json:"payload_identifier"`
	PayloadType              string                        `json:"payload_type"`
	PayloadDisplayName       string                        `json:"payload_display_name"`
	PayloadDescription       string                        `json:"payload_description,omitempty"`
	PayloadOrganization      string                        `json:"payload_organization,omitempty"`
	PayloadUUID              string                        `json:"payload_uuid"`
	PayloadVersion           int                           `json:"payload_version"`
	PayloadRemovalDisallowed bool                          `json:"payload_removal_disallowed"`
	Payloads                 []MobileConfigPayloadResponse `json:"payloads,omitempty"`
	CreatedAt                time.Time                     `json:"created_at"`
	UpdatedAt                time.Time                     `json:"updated_at"`
}
