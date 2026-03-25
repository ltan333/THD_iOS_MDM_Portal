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

func (h *settingHandlerImpl) GetByKey(c *gin.Context) {
	key := c.Param("key")

	st, err := h.settingService.GetByKey(c.Request.Context(), key)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, mapSettingToResponse(st), "")
}

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
