package handler

import (
	"io"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/response"
)

type MDMHandler interface {
	PushCert(c *gin.Context)
	GetCert(c *gin.Context)
	Push(c *gin.Context)
	EnqueueCommand(c *gin.Context)
	EscrowKeyUnlock(c *gin.Context)
	GetVersion(c *gin.Context)
}

type mdmHandler struct {
	client     *ent.Client
	mdmService service.NanoMDMService // Added field
}

func NewMDMHandler(client *ent.Client, mdmService service.NanoMDMService) MDMHandler { // Updated signature
	return &mdmHandler{
		client:     client,
		mdmService: mdmService, // Initialized field
	}
}

// PushCert godoc
// @Summary Upload APNs certificate and private key
// @Description Upload APNs certificate and private key. Concatenated PEM format.
// @Tags MDM
// @Accept text/plain
// @Produce json
// @Param data body string true "PEM-encoded certificate and private key"
// @Success 200 {object} response.APIResponse[dto.PushCertResponse]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /api/v1/mdm/pushcert [put]
func (h *mdmHandler) PushCert(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	result, err := h.mdmService.UploadPushCert(c.Request.Context(), data)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, result, "Certificate uploaded successfully")
}

// GetCert godoc
// @Summary Retrieve APNs push certificate info
// @Description Retrieve the topic and expiry of the stored APNs push certificate.
// @Tags MDM
// @Produce json
// @Param topic query string true "APNs topic"
// @Success 200 {object} response.APIResponse[dto.PushCertResponse]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 404 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /api/v1/mdm/pushcert [get]
func (h *mdmHandler) GetCert(c *gin.Context) {
	topic := c.Query("topic")
	if topic == "" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Topic query parameter is required"))
		return
	}

	cert, err := h.mdmService.GetPushCert(c.Request.Context(), topic)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, cert, "Certificate info retrieved successfully")
}

// Push godoc
// @Summary Send APNs push notifications
// @Description Send APNs push notifications to MDM enrollments.
// @Tags MDM
// @Produce json
// @Param id path string true "Enrollment ID(s) comma-separated"
// @Success 200 {object} response.APIResponse[dto.APIResult]
// @Success 207 {object} response.APIResponse[dto.APIResult]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[dto.APIResult]
// @Security BearerAuth
// @Router /api/v1/mdm/push/{id} [get]
func (h *mdmHandler) Push(c *gin.Context) {
	ids := strings.Split(c.Param("id"), ",")
	result, err := h.mdmService.Push(c.Request.Context(), ids)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, result, "Push notification sent")
}

// EnqueueCommand godoc
// @Summary Enqueue MDM commands
// @Description Enqueue MDM commands to MDM enrollments and (optionally) send APNs push notifications.
// @Tags MDM
// @Accept text/plain
// @Produce json
// @Param id path string true "Enrollment ID"
// @Param nopush query string false "Do not send push (1 to enable)"
// @Param command body string true "XML-encoded MDM command plist"
// @Success 200 {object} response.APIResponse[dto.APIResult]
// @Success 207 {object} response.APIResponse[dto.APIResult]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[dto.APIResult]
// @Security BearerAuth
// @Router /api/v1/mdm/enqueue/{id} [put]
func (h *mdmHandler) EnqueueCommand(c *gin.Context) {
	udid := c.Param("id")
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	result, err := h.mdmService.EnqueueCommand(c.Request.Context(), udid, data)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, result, "Command enqueued successfully")
}

// EscrowKeyUnlock godoc
// @Summary Perform an Escrow Key Unlock
// @Description Perform an Escrow Key Unlock against Apple's API.
// @Tags MDM
// @Accept x-www-form-urlencoded
// @Produce json
// @Param topic formData string true "APNs Push topic"
// @Param serial formData string true "Device serial number"
// @Param productType formData string true "Apple product type"
// @Param escrowKey formData string true "Bypass Code"
// @Param orgName formData string true "Organization name"
// @Param guid formData string true "Requester identifier"
// @Success 200 {string} string "Success"
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Router /api/v1/mdm/escrowkeyunlock [post]
func (h *mdmHandler) EscrowKeyUnlock(c *gin.Context) {
	var req dto.EscrowKeyUnlockRequest
	if err := c.ShouldBind(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithError(err))
		return
	}

	body, headers, status, err := h.mdmService.EscrowKeyUnlock(c.Request.Context(), &req)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	for k, v := range headers {
		c.Header(k, v[0])
	}
	c.Data(status, "application/json", body)
}

// GetVersion godoc
// @Summary Returns the running NanoMDM version
// @Description Returns the running NanoMDM version.
// @Tags MDM
// @Produce json
// @Success 200 {object} response.APIResponse[dto.NanoMDMVersionResponse]
// @Router /api/v1/mdm/version [get]
func (h *mdmHandler) GetVersion(c *gin.Context) {
	resp, err := h.mdmService.GetVersion(c.Request.Context())
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, resp, "Version retrieved successfully")
}
