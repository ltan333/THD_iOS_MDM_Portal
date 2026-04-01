package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	"github.com/thienel/go-backend-template/pkg/config"
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
	service          service.NanoCMDService
	deviceService    service.DeviceService
	depDeviceService service.DepDeviceService
	nanomdmService   service.NanoMDMService
	cfg              *config.Config
	httpClient       *http.Client
}

func NewNanoCMDHandler(svc service.NanoCMDService, deviceService service.DeviceService, depDeviceService service.DepDeviceService, nanomdmService service.NanoMDMService, cfg *config.Config) NanoCMDHandler {
	return &nanocmdHandler{
		service:          svc,
		deviceService:    deviceService,
		depDeviceService: depDeviceService,
		nanomdmService:   nanomdmService,
		cfg:              cfg,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// GetVersion godoc
// @Summary Get NanoCMD version
// @Description Retrieve the current version of the running NanoCMD service.
// @Tags NanoCMD
// @Produce json
// @Success 200 {object} response.APIResponse[dto.NanoCMDVersionResponse] "NanoCMD version information"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/nanocmd/version [get]
func (h *nanocmdHandler) GetVersion(c *gin.Context) {
	resp, err := h.service.GetVersion(c.Request.Context())
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, resp, "Version retrieved successfully")
}

// StartWorkflow godoc
// @Summary Start NanoCMD workflow
// @Description Initiate a pre-defined command workflow for one or more enrolled devices or user channels.
// @Tags NanoCMD
// @Produce json
// @Param name path string true "Workflow name"
// @Param id query []string true "One or more Enrollment IDs (UDIDs or User Channel UUIDs)"
// @Param context query string false "Optional context data for the workflow"
// @Success 200 {object} response.APIResponse[dto.NanoCMDWorkflowStartResponse] "Workflow started successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid request parameters"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/nanocmd/workflow/{name}/start [post]
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
// @Description Retrieve the details of an existing event subscription by its name.
// @Tags NanoCMD
// @Produce json
// @Param name path string true "User-defined subscription name"
// @Success 200 {object} response.APIResponse[dto.EventSubscription] "Event subscription details"
// @Failure 400 {object} response.APIResponse[any] "Invalid subscription name"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Subscription not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/nanocmd/event/{name} [get]
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
// @Summary Create or update event subscription
// @Description Store an event subscription definition. Subsequent state changes will trigger callbacks based on this subscription.
// @Tags NanoCMD
// @Accept json
// @Produce json
// @Param name path string true "User-defined subscription name"
// @Param subscription body dto.EventSubscription true "Subscription definition"
// @Success 204 "Subscription stored successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid request body"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/nanocmd/event/{name} [put]
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
// @Summary Get FileVault profile template
// @Description Retrieve the Apple Configuration Profile template used for enabling FileVault encryption.
// @Tags NanoCMD
// @Produce application/x-apple-aspen-config
// @Success 200 {file} file "Apple Configuration Profile XML"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/nanocmd/fvenable/profiletemplate [get]
func (h *nanocmdHandler) GetFVEnableProfileTemplate(c *gin.Context) {
	data, err := h.service.GetFVEnableProfileTemplate(c.Request.Context())
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	c.Data(http.StatusOK, "application/x-apple-aspen-config", data)
}

// GetProfile godoc
// @Summary Get raw NanoCMD profile
// @Description Retrieve the raw XML content of a stored Apple Configuration Profile.
// @Tags NanoCMD
// @Produce application/x-apple-aspen-config
// @Param name path string true "Profile name"
// @Success 200 {file} file "Apple Configuration Profile XML"
// @Failure 400 {object} response.APIResponse[any] "Invalid profile name"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/nanocmd/profile/{name} [get]
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
// @Summary Upload raw NanoCMD profile
// @Description Store a raw Apple .mobileconfig XML file. Signed profiles are also supported.
// @Tags NanoCMD
// @Accept application/x-apple-aspen-config
// @Param name path string true "User-defined profile name"
// @Param data body string true "Raw .mobileconfig content"
// @Success 204 "Profile stored successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid request body"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/nanocmd/profile/{name} [put]
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
// @Summary Delete NanoCMD profile
// @Description Permanently remove a stored Apple Configuration Profile definition.
// @Tags NanoCMD
// @Param name path string true "Profile name to delete"
// @Success 204 "Profile deleted successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid profile name"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/nanocmd/profile/{name} [delete]
func (h *nanocmdHandler) DeleteProfile(c *gin.Context) {
	name := c.Param("name")
	if err := h.service.DeleteProfile(c.Request.Context(), name); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.NoContent(c)
}

// GetProfiles godoc
// @Summary List profile metadata
// @Description Retrieve metadata for one or more stored configuration profiles.
// @Tags NanoCMD
// @Produce json
// @Param name query []string false "Optional filter by list of profile names"
// @Success 200 {object} response.APIResponse[map[string]dto.NanoCMDProfile] "Map of profile metadata"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/nanocmd/profiles [get]
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
// @Description Retrieve a specific NanoCMD command plan definition by its name.
// @Tags NanoCMD
// @Produce json
// @Param name path string true "Command plan name"
// @Success 200 {object} response.APIResponse[dto.CMDPlan] "Command plan details"
// @Failure 400 {object} response.APIResponse[any] "Invalid command plan name"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Command plan not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/nanocmd/cmdplan/{name} [get]
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
// @Summary Create or update command plan
// @Description Store a command plan definition for automated command issuance.
// @Tags NanoCMD
// @Accept json
// @Produce json
// @Param name path string true "User-defined command plan name"
// @Param plan body dto.CMDPlan true "Plan definition"
// @Success 204 "Command plan stored successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid request body"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/nanocmd/cmdplan/{name} [put]
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
// @Description Retrieve cached inventory data for one or more MDM-enrolled devices.
// @Tags NanoCMD
// @Produce json
// @Param id query []string true "List of Enrollment IDs (UDIDs)"
// @Success 200 {object} response.APIResponse[dto.NanoCMDInventoryResponse] "Device inventory data"
// @Failure 400 {object} response.APIResponse[any] "Invalid Enrollment IDs"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/nanocmd/inventory [get]
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
// @Summary NanoCMD Webhook callback
// @Description Internal endpoint for NanoCMD to report command status changes and state updates. This consumes MicroMDM-compatible webhook payloads.
// @Tags Infrastructure
// @Accept json
// @Produce json
// @Param request body dto.NanoCMDWebhook true "Webhook payload"
// @Success 200 {object} response.APIResponse[any] "Webhook processed successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid webhook signature or data"
// @Failure 500 {object} response.APIResponse[any] "Internal processing error"
// @Router /api/v1/nanocmd/webhook [post]
func (h *nanocmdHandler) Webhook(c *gin.Context) {
	// 1. Capture the Payload (Body Draining)
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		tlog.Error("Failed to read webhook body", zap.Error(err))
		response.WriteErrorResponse(c, err)
		return
	}

	// Reconstruct Request.Body for subsequent binding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var webhook dto.NanoCMDWebhook
	if err := c.ShouldBindJSON(&webhook); err != nil {
		tlog.Error("Failed to bind webhook", zap.Error(err))
		response.WriteErrorResponse(c, err)
		return
	}

	// 2. Log and route based on topic
	tlog.Info("Received NanoCMD webhook", zap.String("topic", webhook.Topic))

	// 3. Check for DEP device events - these are handled locally, NOT forwarded to NanoCMD
	if webhook.Topic == "dep.FetchDevices" || webhook.Topic == "dep.SyncDevices" {
		// Return 200 OK immediately, process in background
		go h.handleDEPDeviceEvent(&webhook)
		response.OK[any](c, nil, "DEP device event received, processing in background")
		return
	}

	// 4. For other topics, keep existing logic (local processing + forward to NanoCMD)
	if err := h.deviceService.HandleWebhook(c.Request.Context(), &webhook); err != nil {
		tlog.Error("Failed to handle device webhook", zap.Error(err))
	}

	// 5. Asynchronous Forwarding to NanoCMD (only for non-DEP topics)
	go h.forwardToNanoCMD(body, c.Request.Header)

	// 6. Return 200 OK to NanoMDM immediately
	response.OK[any](c, nil, "Webhook processed and fan-out initiated")
}

// handleDEPDeviceEvent processes dep.FetchDevices and dep.SyncDevices webhooks.
// It fetches the assigner profile UUID once, then processes each device.
func (h *nanocmdHandler) handleDEPDeviceEvent(webhook *dto.NanoCMDWebhook) {
	ctx := context.Background()

	// Validate device_response_event
	if webhook.DeviceResponseEvent == nil || webhook.DeviceResponseEvent.DeviceResponse == nil {
		tlog.Warn("DEP webhook missing device_response_event", zap.String("topic", webhook.Topic))
		return
	}

	devices := webhook.DeviceResponseEvent.DeviceResponse.Devices
	if len(devices) == 0 {
		tlog.Info("DEP webhook has no devices to process", zap.String("topic", webhook.Topic))
		return
	}

	depName := webhook.DeviceResponseEvent.DEPName
	if depName == "" {
		depName = h.cfg.NanoMDM.DEPServerName
	}

	tlog.Info("Processing DEP device event",
		zap.String("topic", webhook.Topic),
		zap.String("dep_name", depName),
		zap.Int("device_count", len(devices)))

	// Step 1: Get assigner profile UUID (one call per webhook request)
	assignerResp, err := h.nanomdmService.GetDEPAssigner(ctx, depName)
	if err != nil {
		tlog.Error("Failed to get DEP assigner profile",
			zap.String("dep_name", depName),
			zap.Error(err))
		return
	}

	assignerProfileUUID := assignerResp.ProfileUUID
	if assignerProfileUUID == "" {
		tlog.Warn("No assigner profile configured for DEP server",
			zap.String("dep_name", depName))
	}

	tlog.Info("Retrieved assigner profile UUID",
		zap.String("dep_name", depName),
		zap.String("profile_uuid", assignerProfileUUID))

	// Step 2: Process devices - check and reassign profiles if needed, then upsert to DB
	if err := h.depDeviceService.HandleDEPDeviceEvent(ctx, depName, devices, assignerProfileUUID, h.nanomdmService); err != nil {
		tlog.Error("Failed to handle DEP device event",
			zap.String("topic", webhook.Topic),
			zap.Error(err))
	}

	tlog.Info("Completed DEP device event processing",
		zap.String("topic", webhook.Topic),
		zap.Int("device_count", len(devices)))
}

func (h *nanocmdHandler) forwardToNanoCMD(payload []byte, originalHeaders http.Header) {
	if h.cfg.NanoCMD.BaseURL == "" {
		tlog.Warn("NanoCMD BaseURL is not configured, skipping fan-out")
		return
	}

	url := fmt.Sprintf("%s/webhook", strings.TrimSuffix(h.cfg.NanoCMD.BaseURL, "/"))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		tlog.Error("Failed to create forward request", zap.Error(err))
		return
	}

	// Copy essential headers
	req.Header.Set("Content-Type", "application/json")
	if auth := originalHeaders.Get("Authorization"); auth != "" {
		req.Header.Set("Authorization", auth)
	}

	// Ensure Basic Auth if configured (as requested)
	if h.cfg.NanoCMD.Username != "" && h.cfg.NanoCMD.Password != "" {
		req.SetBasicAuth(h.cfg.NanoCMD.Username, h.cfg.NanoCMD.Password)
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		tlog.Error("Failed to forward webhook to NanoCMD", zap.String("url", url), zap.Error(err))
		return
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		tlog.Warn("NanoCMD returned non-OK status for forwarded webhook",
			zap.String("url", url),
			zap.Int("status", resp.StatusCode))
	} else {
		tlog.Info("Successfully forwarded webhook to NanoCMD", zap.String("url", url))
	}
}
