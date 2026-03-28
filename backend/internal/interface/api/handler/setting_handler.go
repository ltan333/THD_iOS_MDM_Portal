package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/response"
)

type SettingHandler interface {
	List(c *gin.Context)
	GetByKey(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type settingHandlerImpl struct {
	settingService service.SettingService
}

func NewSettingHandler(settingService service.SettingService) SettingHandler {
	return &settingHandlerImpl{settingService: settingService}
}

// List godoc
// @Summary List system settings
// @Description Retrieve a complete list of all system-wide configuration settings (key-value pairs).
// @Tags Settings
// @Produce json
// @Success 200 {object} response.APIResponse[[]dto.SettingResponse] "List of system settings"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/settings [get]
func (h *settingHandlerImpl) List(c *gin.Context) {
	settings, err := h.settingService.List(c.Request.Context())
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	res := make([]dto.SettingResponse, 0, len(settings))
	for _, s := range settings {
		res = append(res, mapSettingToResponse(s))
	}

	response.OK(c, res, "")
}

// GetByKey godoc
// @Summary Get system setting
// @Description Retrieve the value and metadata for a specific configuration setting using its unique key.
// @Tags Settings
// @Produce json
// @Param key path string true "Unique setting key"
// @Success 200 {object} response.APIResponse[dto.SettingResponse] "Setting details"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Setting not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/settings/{key} [get]
func (h *settingHandlerImpl) GetByKey(c *gin.Context) {
	key := c.Param("key")

	st, err := h.settingService.GetByKey(c.Request.Context(), key)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, mapSettingToResponse(st), "")
}

// Create godoc
// @Summary Create system setting
// @Description Register a new system-wide configuration setting.
// @Tags Settings
// @Accept json
// @Produce json
// @Param request body dto.CreateSettingRequest true "Setting details"
// @Success 201 {object} response.APIResponse[dto.SettingResponse] "Setting created successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 409 {object} response.APIResponse[any] "Setting key already exists"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/settings [post]
func (h *settingHandlerImpl) Create(c *gin.Context) {
	var req dto.CreateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid data"))
		return
	}

	st, err := h.settingService.Create(c.Request.Context(), service.CreateSettingCommand{
		Key:         req.Key,
		Value:       req.Value,
		Description: req.Description,
	})
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.Created(c, mapSettingToResponse(st), "Setting created")
}

// Update godoc
// @Summary Update system setting
// @Description Modify the value or description of an existing configuration setting.
// @Tags Settings
// @Accept json
// @Produce json
// @Param key path string true "Unique setting key"
// @Param request body dto.UpdateSettingRequest true "Updated setting details"
// @Success 200 {object} response.APIResponse[dto.SettingResponse] "Setting updated successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Setting not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/settings/{key} [put]
func (h *settingHandlerImpl) Update(c *gin.Context) {
	key := c.Param("key")

	var req dto.UpdateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid data"))
		return
	}

	st, err := h.settingService.Update(c.Request.Context(), service.UpdateSettingCommand{
		Key:         key,
		Value:       req.Value,
		Description: req.Description,
	})
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, mapSettingToResponse(st), "Setting updated")
}

// Delete godoc
// @Summary Delete system setting
// @Description Permanently remove a configuration setting from the system.
// @Tags Settings
// @Produce json
// @Param key path string true "Unique setting key"
// @Success 200 {object} response.APIResponse[any] "Setting deleted successfully"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Setting not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/settings/{key} [delete]
func (h *settingHandlerImpl) Delete(c *gin.Context) {
	key := c.Param("key")

	if err := h.settingService.Delete(c.Request.Context(), key); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK[any](c, nil, "Setting deleted")
}

func mapSettingToResponse(s *ent.Setting) dto.SettingResponse {
	return dto.SettingResponse{
		ID:          s.ID,
		Key:         s.Key,
		Value:       s.Value,
		Description: s.Description,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
