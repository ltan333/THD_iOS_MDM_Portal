package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	"github.com/thienel/go-backend-template/pkg/response"
	"github.com/thienel/tlog"
	"go.uber.org/zap"
)

type NanoCMDHandler interface {
	GetVersion(c *gin.Context)
	StartWorkflow(c *gin.Context)
	GetEvent(c *gin.Context)
	PutEvent(c *gin.Context)
	GetFVEnableProfileTemplate(c *gin.Context)
	GetProfile(c *gin.Context)
	PutProfile(c *gin.Context)
	DeleteProfile(c *gin.Context)
	GetProfiles(c *gin.Context)
	GetCMDPlan(c *gin.Context)
	PutCMDPlan(c *gin.Context)
	GetInventory(c *gin.Context)
	Webhook(c *gin.Context)
}

type nanocmdHandler struct {
	service service.NanoCMDService
}

func NewNanoCMDHandler(svc service.NanoCMDService) NanoCMDHandler {
	return &nanocmdHandler{service: svc}
}

// GetVersion godoc
// @Summary Get NanoCMD version
// @Description Get the version of the NanoCMD server
// @Tags NanoCMD
// @Produce json
// @Success 200 {object} response.APIResponse[dto.NanoCMDVersionResponse]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/nanocmd/version [get]
func (h *nanocmdHandler) GetVersion(c *gin.Context) {
	resp, err := h.service.GetVersion(c.Request.Context())
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, resp, "Version retrieved successfully")
}

// StartWorkflow godoc
// @Summary Start a workflow
// @Description Initiate a NanoCMD workflow on specific devices
// @Tags NanoCMD
// @Produce json
// @Param name path string true "Workflow name"
// @Param id query []string false "Device IDs"
// @Param context query string false "Workflow context"
// @Success 200 {object} response.APIResponse[dto.NanoCMDWorkflowStartResponse]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/nanocmd/workflow/{name}/start [post]
func (h *nanocmdHandler) StartWorkflow(c *gin.Context) {
	name := c.Param("name")
	ids := c.QueryArray("id")
	ctxStr := c.Query("context")

	resp, err := h.service.StartWorkflow(c.Request.Context(), name, ids, ctxStr)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, resp, "Workflow started successfully")
}

// GetEvent godoc
// @Summary Get event subscription
// @Description Get details of an event subscription by name
// @Tags NanoCMD
// @Produce json
// @Param name path string true "Event name"
// @Success 200 {object} response.APIResponse[dto.EventSubscription]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 404 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/nanocmd/event/{name} [get]
func (h *nanocmdHandler) GetEvent(c *gin.Context) {
	name := c.Param("name")
	resp, err := h.service.GetEvent(c.Request.Context(), name)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, resp, "Event retrieved successfully")
}

// PutEvent godoc
// @Summary Update event subscription
// @Description Create or update an event subscription
// @Tags NanoCMD
// @Accept json
// @Produce json
// @Param name path string true "Event name"
// @Param subscription body dto.EventSubscription true "Subscription details"
// @Success 204
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/nanocmd/event/{name} [put]
func (h *nanocmdHandler) PutEvent(c *gin.Context) {
	name := c.Param("name")
	var sub dto.EventSubscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	if err := h.service.PutEvent(c.Request.Context(), name, &sub); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.NoContent(c)
}

// GetFVEnableProfileTemplate godoc
// @Summary Get FileVault enable profile template
// @Description Get the Apple Configuration Profile template for enabling FileVault
// @Tags NanoCMD
// @Produce application/x-apple-aspen-config
// @Success 200 {string} string "Apple Configuration Profile"
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/nanocmd/fvenable/profiletemplate [get]
func (h *nanocmdHandler) GetFVEnableProfileTemplate(c *gin.Context) {
	data, err := h.service.GetFVEnableProfileTemplate(c.Request.Context())
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	c.Data(http.StatusOK, "application/x-apple-aspen-config", data)
}

// GetProfile godoc
// @Summary Get profile by name
// @Description Fetch an Apple Configuration Profile from NanoCMD
// @Tags NanoCMD
// @Produce application/x-apple-aspen-config
// @Param name path string true "Profile name"
// @Success 200 {string} string "Apple Configuration Profile"
// @Failure 401 {object} response.APIResponse[any]
// @Failure 404 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/nanocmd/profile/{name} [get]
func (h *nanocmdHandler) GetProfile(c *gin.Context) {
	name := c.Param("name")
	data, err := h.service.GetProfile(c.Request.Context(), name)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	c.Data(http.StatusOK, "application/x-apple-aspen-config", data)
}

// PutProfile godoc
// @Summary Create or update profile
// @Description Upload an Apple Configuration Profile to NanoCMD
// @Tags NanoCMD
// @Accept application/x-apple-aspen-config
// @Param name path string true "Profile name"
// @Param data body string true "Profile XML data"
// @Success 204
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/nanocmd/profile/{name} [put]
func (h *nanocmdHandler) PutProfile(c *gin.Context) {
	name := c.Param("name")
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	if err := h.service.PutProfile(c.Request.Context(), name, data); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.NoContent(c)
}

// DeleteProfile godoc
// @Summary Delete profile
// @Description Remove a profile from NanoCMD
// @Tags NanoCMD
// @Param name path string true "Profile name"
// @Success 204
// @Failure 401 {object} response.APIResponse[any]
// @Failure 404 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/nanocmd/profile/{name} [delete]
func (h *nanocmdHandler) DeleteProfile(c *gin.Context) {
	name := c.Param("name")
	if err := h.service.DeleteProfile(c.Request.Context(), name); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.NoContent(c)
}

// GetProfiles godoc
// @Summary List profiles
// @Description Get a list of profiles by name
// @Tags NanoCMD
// @Produce json
// @Param name query []string false "Profile names"
// @Success 200 {object} response.APIResponse[[]dto.NanoCMDProfile]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/nanocmd/profiles [get]
func (h *nanocmdHandler) GetProfiles(c *gin.Context) {
	names := c.QueryArray("name")
	resp, err := h.service.GetProfiles(c.Request.Context(), names)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, resp, "Profiles retrieved successfully")
}

// GetCMDPlan godoc
// @Summary Get command plan
// @Description Fetch detailed command plan for a device from NanoCMD
// @Tags NanoCMD
// @Produce json
// @Param name path string true "Device ID"
// @Success 200 {object} response.APIResponse[dto.CMDPlan]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 404 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/nanocmd/cmdplan/{name} [get]
func (h *nanocmdHandler) GetCMDPlan(c *gin.Context) {
	name := c.Param("name")
	resp, err := h.service.GetCMDPlan(c.Request.Context(), name)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, resp, "Command plan retrieved successfully")
}

// PutCMDPlan godoc
// @Summary Update command plan
// @Description Update command plan for a device
// @Tags NanoCMD
// @Accept json
// @Produce json
// @Param name path string true "Device ID"
// @Param plan body dto.CMDPlan true "Command plan details"
// @Success 204
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/nanocmd/cmdplan/{name} [put]
func (h *nanocmdHandler) PutCMDPlan(c *gin.Context) {
	name := c.Param("name")
	var plan dto.CMDPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	if err := h.service.PutCMDPlan(c.Request.Context(), name, &plan); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.NoContent(c)
}

// GetInventory godoc
// @Summary Get device inventory
// @Description Fetch inventory metadata for specific devices
// @Tags NanoCMD
// @Produce json
// @Param id query []string false "Device IDs"
// @Success 200 {object} response.APIResponse[dto.NanoCMDInventoryResponse]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/nanocmd/inventory [get]
func (h *nanocmdHandler) GetInventory(c *gin.Context) {
	ids := c.QueryArray("id")
	resp, err := h.service.GetInventory(c.Request.Context(), ids)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, resp, "Inventory retrieved successfully")
}

// Webhook godoc
// @Summary NanoCMD Webhook
// @Description Endpoint for NanoCMD to send event notifications
// @Tags NanoCMD
// @Accept json
// @Produce json
// @Param webhook body dto.NanoCMDWebhook true "Webhook payload"
// @Success 200 {object} response.APIResponse[any]
// @Failure 400 {object} response.APIResponse[any]
// @Router /v1/nanocmd/webhook [post]
func (h *nanocmdHandler) Webhook(c *gin.Context) {
	var webhook dto.NanoCMDWebhook
	if err := c.ShouldBindJSON(&webhook); err != nil {
		tlog.Error("Failed to bind webhook", zap.Error(err))
		response.WriteErrorResponse(c, err)
		return
	}

	// Process webhook logic here (e.g., update device status)
	tlog.Info("Received NanoCMD webhook", zap.String("topic", webhook.Topic))

	response.OK[any](c, nil, "Webhook processed successfully")
}
