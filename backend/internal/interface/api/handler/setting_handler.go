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
// @Summary List settings
// @Description Fetch all backend configuration kv string settings
// @Tags settings
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse[[]dto.SettingResponse]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
// @Router /settings [get]
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
// @Summary Get setting by key
// @Description Get a specifically named environment configuration value
// @Tags settings
// @Produce json
// @Param key path string true "Setting Key Identifier"
// @Security BearerAuth
// @Success 200 {object} response.APIResponse[dto.SettingResponse]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 404 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
// @Router /settings/{key} [get]
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
// @Summary Create setting
// @Description Create a new key-value system config record
// @Tags settings
// @Accept json
// @Produce json
// @Param request body dto.CreateSettingRequest true "Create Setting Body"
// @Security BearerAuth
// @Success 201 {object} response.APIResponse[dto.SettingResponse]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
// @Router /settings [post]
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
// @Summary Update setting
// @Description Overwrite value and description of a system setting
// @Tags settings
// @Accept json
// @Produce json
// @Param key path string true "Setting Key Identifier"
// @Param request body dto.UpdateSettingRequest true "Update Setting Body"
// @Security BearerAuth
// @Success 200 {object} response.APIResponse[dto.SettingResponse]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 404 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
// @Router /settings/{key} [put]
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
// @Summary Delete setting
// @Description Delete system setting forever
// @Tags settings
// @Produce json
// @Param key path string true "Setting Key String"
// @Security BearerAuth
// @Success 200 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
// @Router /settings/{key} [delete]
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
