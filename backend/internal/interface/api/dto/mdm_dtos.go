package dto

import "time"

type APNSConfigResponse struct {
	ID        string    `json:"id"`
	CertPEM   string    `json:"cert_pem"`
	KeyPEM    string    `json:"key_pem"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
