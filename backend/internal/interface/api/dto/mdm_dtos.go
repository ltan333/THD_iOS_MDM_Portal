package dto

import "time"

type APNSConfigResponse struct {
	ID           uint      `json:"id"`
	Topic        string    `json:"topic"`
	CertFilePath string    `json:"cert_file_path"`
	Expiry       time.Time `json:"expiry"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
