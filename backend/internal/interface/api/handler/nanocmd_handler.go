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
	service       service.NanoCMDService
	deviceService service.DeviceService
}

func NewNanoCMDHandler(svc service.NanoCMDService, deviceService service.DeviceService) NanoCMDHandler {
	return &nanocmdHandler{
		service:       svc,
		deviceService: deviceService,
	}
}

// GetVersion godoc
// @Summary Returns the running NanoCMD server version
// @Description Get the version of the NanoCMD server
// @Tags NanoCMD
// @Produce json
// @Success 200 {object} response.APIResponse[dto.NanoCMDVersionResponse]
// @Failure 401 {object} response.APIResponse[any]
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
// @Summary Start a workflow
// @Description Start a workflow.
// @Tags NanoCMD
// @Produce json
// @Param name path string true "Name of NanoCMD workflow."
// @Param id query []string true "Enrollment ID. Unique identifier of MDM enrollment. Often a device UDID or a user channel UUID."
// @Param context query string false "Workflow-dependent context."
// @Success 200 {object} response.APIResponse[dto.NanoCMDWorkflowStartResponse]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
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
// @Summary Retrieve the event subscription
// @Description Retrieve the event subscription.
// @Tags NanoCMD
// @Produce json
// @Param name path string true "User-defined name of Event Subscription."
// @Success 200 {object} response.APIResponse[dto.EventSubscription]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 404 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
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
// @Summary Store the event subscription
// @Description Store the event subscription provided in the request body.
// @Tags NanoCMD
// @Accept json
// @Produce json
// @Param name path string true "User-defined name of Event Subscription."
// @Param subscription body dto.EventSubscription true "Event Subscription."
// @Success 204
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
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
// @Summary Returns the FileVault enable Configuration Profile template
// @Description Returns the FileVault enable Configuration Profile template.
// @Tags NanoCMD
// @Produce application/x-apple-aspen-config
// @Success 200 {string} string "Apple Configuration Profile"
// @Failure 401 {object} response.APIResponse[any]
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
// @Summary Fetches the named raw profile
// @Description Fetches the named raw profile.
// @Tags NanoCMD
// @Produce application/x-apple-aspen-config
// @Param name path string true "User-defined name of Profile."
// @Success 200 {string} string "Apple Configuration Profile"
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 404 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
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
// @Summary Uploads a raw profile
// @Description Uploads a raw profile. Signed profiles also supported.
// @Tags NanoCMD
// @Accept application/x-apple-aspen-config
// @Param name path string true "User-defined name of Profile."
// @Param data body string true "Raw profile mobileconfig."
// @Success 204
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
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
// @Summary Deletes the named profile
// @Description Deletes the named profile.
// @Tags NanoCMD
// @Param name path string true "User-defined name of Profile."
// @Success 204
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 404 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
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
// @Summary Retrieve profile metadata
// @Description Retrieve profile metadata.
// @Tags NanoCMD
// @Produce json
// @Param name query []string false "User-defined name of profile."
// @Success 200 {object} response.APIResponse[map[string]dto.NanoCMDProfile]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
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
// @Summary Retrieve and return a named command plan
// @Description Retrieve and return a named command plan as JSON.
// @Tags NanoCMD
// @Produce json
// @Param name path string true "User-defined name of Command Plan."
// @Success 200 {object} response.APIResponse[dto.CMDPlan]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 404 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
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
// @Summary Upload a named command plan
// @Description Upload a named JSON command plan.
// @Tags NanoCMD
// @Accept json
// @Produce json
// @Param name path string true "User-defined name of Command Plan."
// @Param plan body dto.CMDPlan true "Command plan."
// @Success 204
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
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
// @Summary Retrieve inventory data for enrollment IDs
// @Description Retrieve inventory data for enrollment IDs.
// @Tags NanoCMD
// @Produce json
// @Param id query []string true "Enrollment ID. Unique identifier of MDM enrollment. Often a device UDID or a user channel UUID."
// @Success 200 {object} response.APIResponse[dto.NanoCMDInventoryResponse]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
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
// @Summary NanoCMD Webhook
// @Description Handler for MicroMDM-compatible webhook callback.
// @Tags NanoCMD
// @Accept json
// @Produce json
// @Description Endpoint to receive MDM/CMD check-in events and state updates
// @Tags Infrastructure
// @Accept json
// @Produce json
// @Param request body dto.NanoCMDWebhook true "Webhook payload"
// @Success 200 {object} response.APIResponse[any]
// @Failure 400 {object} response.APIResponse[any]
// @Router /api/v1/nanocmd/webhook [post]
func (h *nanocmdHandler) Webhook(c *gin.Context) {
	var webhook dto.NanoCMDWebhook
	if err := c.ShouldBindJSON(&webhook); err != nil {
		tlog.Error("Failed to bind webhook", zap.Error(err))
		response.WriteErrorResponse(c, err)
		return
	}

	// Process webhook logic here (e.g., update device status)
	tlog.Info("Received NanoCMD webhook", zap.String("topic", webhook.Topic))

	if err := h.deviceService.HandleWebhook(c.Request.Context(), &webhook); err != nil {
		tlog.Error("Failed to handle device webhook", zap.Error(err))
		// We don't necessarily want to return 500 to the webhook provider if our local DB update fails
		// but for debugging purposes it might be better to know.
	}

	response.OK[any](c, nil, "Webhook processed successfully")
}
