package dto

import "time"

type SettingResponse struct {
	ID          uint      `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateSettingRequest struct {
	Key         string `json:"key" binding:"required,max=255"`
	Value       string `json:"value" binding:"required"`
	Description string `json:"description,omitempty"`
}

type UpdateSettingRequest struct {
	Value       *string `json:"value,omitempty"`
	Description *string `json:"description,omitempty"`
}
