package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
	"github.com/thienel/go-backend-template/pkg/response"
)

var alertAllowedFields = map[string]bool{
	"severity": true,
	"type":     true,
	"status":   true,
	"search":   true,
}

var alertRuleAllowedFields = map[string]bool{
	"enabled": true,
	"search":  true,
}

type AlertHandler interface {
	List(c *gin.Context)
	GetByID(c *gin.Context)
	Create(c *gin.Context)
	Acknowledge(c *gin.Context)
	Resolve(c *gin.Context)
	BulkResolve(c *gin.Context)
	GetStats(c *gin.Context)

	LockDevice(c *gin.Context)
	WipeDevice(c *gin.Context)
	PushPolicy(c *gin.Context)
	SendMessage(c *gin.Context)

	ListRules(c *gin.Context)
	GetRuleByID(c *gin.Context)
	CreateRule(c *gin.Context)
	UpdateRule(c *gin.Context)
	DeleteRule(c *gin.Context)
	ToggleRule(c *gin.Context)
}

type alertHandlerImpl struct {
	alertService     service.AlertService
	alertRuleService service.AlertRuleService
}

func NewAlertHandler(alertService service.AlertService, alertRuleService service.AlertRuleService) AlertHandler {
	return &alertHandlerImpl{
		alertService:     alertService,
		alertRuleService: alertRuleService,
	}
}

func (h *alertHandlerImpl) List(c *gin.Context) {
	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	offset, limit := query.GetPagination(params, 20)
	opts := query.ParseQueryParams(params, alertAllowedFields)

	alerts, total, err := h.alertService.List(c.Request.Context(), offset, limit, opts)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	res := make([]dto.AlertResponse, 0, len(alerts))
	for _, a := range alerts {
		res = append(res, mapAlertToResponse(a))
	}

	page := (offset / limit) + 1
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response.OK(c, dto.ListResponse[dto.AlertResponse]{
		Items:      res,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, "")
}

func (h *alertHandlerImpl) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	alert, err := h.alertService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, mapAlertToResponse(alert), "")
}

func (h *alertHandlerImpl) Create(c *gin.Context) {
	var req dto.CreateAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid data"))
		return
	}

	alert, err := h.alertService.Create(c.Request.Context(), service.CreateAlertCommand{
		Severity: req.Severity,
		Title:    req.Title,
		Type:     req.Type,
		DeviceID: req.DeviceID,
		UserID:   req.UserID,
		Details:  req.Details,
	})
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.Created(c, mapAlertToResponse(alert), "Alert created successfully")
}

func (h *alertHandlerImpl) Acknowledge(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	if err := h.alertService.Acknowledge(c.Request.Context(), uint(id)); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK[any](c, nil, "Alert acknowledged")
}

func (h *alertHandlerImpl) Resolve(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	if err := h.alertService.Resolve(c.Request.Context(), uint(id)); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK[any](c, nil, "Alert resolved")
}

func (h *alertHandlerImpl) BulkResolve(c *gin.Context) {
	var req dto.BulkResolveAlertsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid data"))
		return
	}

	if err := h.alertService.BulkResolve(c.Request.Context(), req.AlertIDs); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK[any](c, nil, "Alerts bulk resolved")
}

func (h *alertHandlerImpl) GetStats(c *gin.Context) {
	stats, err := h.alertService.GetStats(c.Request.Context())
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, stats, "")
}

func (h *alertHandlerImpl) LockDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	if err := h.alertService.LockDevice(c.Request.Context(), uint(id)); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK[any](c, nil, "Device lock initiated")
}

func (h *alertHandlerImpl) WipeDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	if err := h.alertService.WipeDevice(c.Request.Context(), uint(id)); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK[any](c, nil, "Device wipe initiated")
}

func (h *alertHandlerImpl) PushPolicy(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	var req dto.AlertActionRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.PolicyID == nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid policy ID"))
		return
	}

	if err := h.alertService.PushPolicy(c.Request.Context(), uint(id), *req.PolicyID); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK[any](c, nil, "Policy pushed")
}

func (h *alertHandlerImpl) SendMessage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	var req dto.AlertActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid message data"))
		return
	}

	if err := h.alertService.SendMessage(c.Request.Context(), uint(id), req.Message); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK[any](c, nil, "Message sent")
}

// ---- Rule Methods ----

func (h *alertHandlerImpl) ListRules(c *gin.Context) {
	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	offset, limit := query.GetPagination(params, 20)
	opts := query.ParseQueryParams(params, alertRuleAllowedFields)

	rules, total, err := h.alertRuleService.List(c.Request.Context(), offset, limit, opts)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	res := make([]dto.AlertRuleResponse, 0, len(rules))
	for _, r := range rules {
		res = append(res, mapAlertRuleToResponse(r))
	}

	page := (offset / limit) + 1
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response.OK(c, dto.ListResponse[dto.AlertRuleResponse]{
		Items:      res,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, "")
}

func (h *alertHandlerImpl) GetRuleByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	r, err := h.alertRuleService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, mapAlertRuleToResponse(r), "")
}

func (h *alertHandlerImpl) CreateRule(c *gin.Context) {
	var req dto.CreateAlertRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid data"))
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	r, err := h.alertRuleService.Create(c.Request.Context(), service.CreateAlertRuleCommand{
		Name:        req.Name,
		Description: req.Description,
		Condition:   req.Condition,
		Actions:     req.Actions,
		Enabled:     enabled,
	})
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.Created(c, mapAlertRuleToResponse(r), "Rule created")
}

func (h *alertHandlerImpl) UpdateRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	var req dto.UpdateAlertRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid data"))
		return
	}

	r, err := h.alertRuleService.Update(c.Request.Context(), service.UpdateAlertRuleCommand{
		ID:          uint(id),
		Name:        req.Name,
		Description: req.Description,
		Condition:   req.Condition,
		Actions:     req.Actions,
		Enabled:     req.Enabled,
	})
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, mapAlertRuleToResponse(r), "Rule updated")
}

func (h *alertHandlerImpl) DeleteRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	if err := h.alertRuleService.Delete(c.Request.Context(), uint(id)); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK[any](c, nil, "Rule deleted")
}

func (h *alertHandlerImpl) ToggleRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	if err := h.alertRuleService.Toggle(c.Request.Context(), uint(id)); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK[any](c, nil, "Rule toggled")
}

// Helpers

func mapAlertToResponse(a *ent.Alert) dto.AlertResponse {
	return dto.AlertResponse{
		ID:             a.ID,
		Severity:       string(a.Severity),
		Title:          a.Title,
		Type:           string(a.Type),
		Status:         string(a.Status),
		DeviceID:       a.DeviceID,
		UserID:         a.UserID,
		Details:        a.Details,
		CreatedAt:      a.CreatedAt,
		AcknowledgedAt: a.AcknowledgedAt,
		ResolvedAt:     a.ResolvedAt,
	}
}

func mapAlertRuleToResponse(r *ent.AlertRule) dto.AlertRuleResponse {
	return dto.AlertRuleResponse{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Condition:   r.Condition,
		Actions:     r.Actions,
		Enabled:     r.Enabled,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}
