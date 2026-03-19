package handler

import (
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	"github.com/thienel/go-backend-template/pkg/response"
)

// APIErrorResponse documents the common error envelope.
type APIErrorResponse struct {
	IsSuccess bool            `json:"is_success"`
	Error     *response.Error `json:"error"`
}

type LoginSuccessResponse struct {
	IsSuccess bool              `json:"is_success"`
	Data      dto.LoginResponse `json:"data"`
	Message   string            `json:"message,omitempty"`
}

type UserSuccessResponse struct {
	IsSuccess bool             `json:"is_success"`
	Data      dto.UserResponse `json:"data"`
	Message   string           `json:"message,omitempty"`
}

type UserListData struct {
	Items      []dto.UserResponse `json:"items"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

type UserListSuccessResponse struct {
	IsSuccess bool         `json:"is_success"`
	Data      UserListData `json:"data"`
	Message   string       `json:"message,omitempty"`
}

type EmptySuccessResponse struct {
	IsSuccess bool   `json:"is_success"`
	Message   string `json:"message,omitempty"`
}

type PoliciesSuccessResponse struct {
	IsSuccess bool                 `json:"is_success"`
	Data      []service.PolicyRule `json:"data"`
	Message   string               `json:"message,omitempty"`
}

type RolesSuccessResponse struct {
	IsSuccess bool               `json:"is_success"`
	Data      []service.RoleLink `json:"data"`
	Message   string             `json:"message,omitempty"`
}

type PolicySuccessResponse struct {
	IsSuccess bool               `json:"is_success"`
	Data      service.PolicyRule `json:"data"`
	Message   string             `json:"message,omitempty"`
}

type RoleLinkSuccessResponse struct {
	IsSuccess bool             `json:"is_success"`
	Data      service.RoleLink `json:"data"`
	Message   string           `json:"message,omitempty"`
}

type HealthResponse struct {
	Status string `json:"status"`
}
