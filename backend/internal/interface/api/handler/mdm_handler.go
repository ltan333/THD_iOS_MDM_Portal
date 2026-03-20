package handler

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/response"
)

type MDMHandler interface {
	PushCert(c *gin.Context)
	GetCert(c *gin.Context)
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
// @Summary Upload APNs certificate
// @Description Upload or update an APNs certificate for a specific topic
// @Tags MDM
// @Accept multipart/form-data
// @Produce json
// @Param cert formData file true "APNs Certificate file"
// @Success 200 {object} response.APIResponse[dto.APNSConfigResponse]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /mdm/apns/cert [post]
func (h *mdmHandler) PushCert(c *gin.Context) {
	file, err := c.FormFile("cert")
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Certificate file is required"))
		return
	}

	fileReader, err := file.Open()
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithError(err))
		return
	}
	defer fileReader.Close()

	fileBytes, err := io.ReadAll(fileReader)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithError(err))
		return
	}

	if err := h.mdmService.UploadPushCert(c.Request.Context(), fileBytes); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, gin.H{}, "Certificate pushed successfully")
}

// GetCert godoc
// @Summary Get APNs configuration
// @Description Get the current APNs certificate configuration and topic
// @Tags MDM
// @Produce json
// @Success 200 {object} response.APIResponse[dto.APNSConfigResponse]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 404 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /mdm/apns/cert [get]
func (h *mdmHandler) GetCert(c *gin.Context) {
	cert, err := h.mdmService.GetPushCert(c.Request.Context())
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, cert, "Certificate retrieved successfully")
}

func parseCert(data []byte) (string, time.Time, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return "", time.Time{}, fmt.Errorf("failed to decode PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", time.Time{}, err
	}

	// Topic is usually in the Subject Common Name or an extension
	topic := cert.Subject.CommonName
	expiry := cert.NotAfter

	return topic, expiry, nil
}
