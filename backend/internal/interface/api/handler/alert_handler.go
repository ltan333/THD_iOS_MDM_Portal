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

// List godoc
// @Summary List alerts
// @Description Fetch tracked system alerts with support for pagination and filtering by severity, type, and status.
// @Tags Alerts
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 20)"
// @Param severity query string false "Filter by severity (info, warning, critical)"
// @Param type query string false "Filter by alert type"
// @Param status query string false "Filter by status (open, acknowledged, resolved)"
// @Param search query string false "Search in title and details"
// @Success 200 {object} response.APIResponse[dto.ListResponse[dto.AlertResponse]] "List of alerts"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/alerts [get]
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

// GetByID godoc
// @Summary Get alert by ID
// @Description Retrieve details of a specific alert by its unique ID.
// @Tags Alerts
// @Produce json
// @Param id path int true "Alert ID"
// @Success 200 {object} response.APIResponse[dto.AlertResponse] "Alert details"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID format"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Alert not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/alerts/{id} [get]
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

// Create godoc
// @Summary Create alert
// @Description Manually create a new system alert.
// @Tags Alerts
// @Accept json
// @Produce json
// @Param alert body dto.CreateAlertRequest true "Alert details"
// @Success 201 {object} response.APIResponse[dto.AlertResponse] "Alert created successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/alerts [post]
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

// Acknowledge godoc
// @Summary Acknowledge alert
// @Description Mark an alert as acknowledged by an administrator.
// @Tags Alerts
// @Produce json
// @Param id path int true "Alert ID"
// @Success 200 {object} response.APIResponse[any] "Alert acknowledged"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID format"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Alert not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/alerts/{id}/acknowledge [put]
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

// Resolve godoc
// @Summary Resolve alert
// @Description Mark an alert as resolved.
// @Tags Alerts
// @Produce json
// @Param id path int true "Alert ID"
// @Success 200 {object} response.APIResponse[any] "Alert resolved"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID format"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Alert not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/alerts/{id}/resolve [put]
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

// BulkResolve godoc
// @Summary Bulk resolve alerts
// @Description Resolve multiple alerts in a single request.
// @Tags Alerts
// @Accept json
// @Produce json
// @Param request body dto.BulkResolveAlertsRequest true "List of alert IDs to resolve"
// @Success 200 {object} response.APIResponse[any] "Alerts bulk resolved"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/alerts/bulk-resolve [post]
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

// GetStats godoc
// @Summary Get alert stats
// @Description Get statistical data about alerts, such as counts by severity and status.
// @Tags Alerts
// @Produce json
// @Success 200 {object} response.APIResponse[dto.AlertsSummaryResponse] "Alert statistics"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/alerts/stats [get]
func (h *alertHandlerImpl) GetStats(c *gin.Context) {
	stats, err := h.alertService.GetStats(c.Request.Context())
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, stats, "")
}

// LockDevice godoc
// @Summary Lock device from alert
// @Description Initiate a remote lock command for the device associated with this alert.
// @Tags Alerts
// @Produce json
// @Param id path int true "Alert ID"
// @Success 200 {object} response.APIResponse[any] "Device lock initiated"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID or device not found"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/alerts/{id}/actions/lock [post]
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

// WipeDevice godoc
// @Summary Wipe device from alert
// @Description Initiate a remote wipe (factory reset) command for the device associated with this alert.
// @Tags Alerts
// @Produce json
// @Param id path int true "Alert ID"
// @Success 200 {object} response.APIResponse[any] "Device wipe initiated"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID or device not found"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/alerts/{id}/actions/wipe [post]
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

// PushPolicy godoc
// @Summary Push policy from alert
// @Description Force push a specific configuration policy to the device associated with this alert.
// @Tags Alerts
// @Accept json
// @Produce json
// @Param id path int true "Alert ID"
// @Param request body dto.AlertActionRequest true "Policy ID"
// @Success 200 {object} response.APIResponse[any] "Policy pushed"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/alerts/{id}/actions/push-policy [post]
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

// SendMessage godoc
// @Summary Send message from alert
// @Description Send a notification message to the device user associated with this alert.
// @Tags Alerts
// @Accept json
// @Produce json
// @Param id path int true "Alert ID"
// @Param request body dto.AlertActionRequest true "Message content"
// @Success 200 {object} response.APIResponse[any] "Message sent"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/alerts/{id}/actions/message [post]
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

// ListRules godoc
// @Summary List alert rules
// @Description Retrieve all alert rules with pagination and filtering by status.
// @Tags Alert-Rules
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 20)"
// @Param enabled query boolean false "Filter by enabled status"
// @Param search query string false "Search in name and description"
// @Success 200 {object} response.APIResponse[dto.ListResponse[dto.AlertRuleResponse]] "List of alert rules"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/alerts/rules [get]
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

// GetRuleByID godoc
// @Summary Get alert rule by ID
// @Description Fetch details of a specific alert rule by its system ID.
// @Tags Alert-Rules
// @Produce json
// @Param id path int true "Rule ID"
// @Success 200 {object} response.APIResponse[dto.AlertRuleResponse] "Alert rule details"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID format"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Rule not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/alerts/rules/{id} [get]
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

// CreateRule godoc
// @Summary Create alert rule
// @Description Define a new alert rule with conditions and automated actions.
// @Tags Alert-Rules
// @Accept json
// @Produce json
// @Param rule body dto.CreateAlertRuleRequest true "Rule details"
// @Success 201 {object} response.APIResponse[dto.AlertRuleResponse] "Rule created successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/alerts/rules [post]
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

// UpdateRule godoc
// @Summary Update alert rule
// @Description Update the configuration of an existing alert rule.
// @Tags Alert-Rules
// @Accept json
// @Produce json
// @Param id path int true "Rule ID"
// @Param rule body dto.UpdateAlertRuleRequest true "Updated rule details"
// @Success 200 {object} response.APIResponse[dto.AlertRuleResponse] "Rule updated successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID or request data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Rule not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/alerts/rules/{id} [put]
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

// DeleteRule godoc
// @Summary Delete alert rule
// @Description Permanently remove an alert rule from the system.
// @Tags Alert-Rules
// @Produce json
// @Param id path int true "Rule ID"
// @Success 200 {object} response.APIResponse[any] "Rule deleted successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID format"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Rule not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/alerts/rules/{id} [delete]
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

// ToggleRule godoc
// @Summary Toggle alert rule active status
// @Description Enable or disable an alert rule.
// @Tags Alert-Rules
// @Produce json
// @Param id path int true "Rule ID"
// @Success 200 {object} response.APIResponse[any] "Rule status toggled"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID format"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Rule not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/alerts/rules/{id}/toggle [put]
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
