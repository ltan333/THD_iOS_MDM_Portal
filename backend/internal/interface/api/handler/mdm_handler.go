package handler

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/apnsconfig"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/response"
)

type MDMHandler interface {
	PushCert(c *gin.Context)
	GetCert(c *gin.Context)
}

type mdmHandler struct {
	client *ent.Client
}

func NewMDMHandler(client *ent.Client) MDMHandler {
	return &mdmHandler{client: client}
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

	dst := filepath.Join("storage", "certs", file.Filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithError(err))
		return
	}

	// Extract topic and expiry from certificate
	certData, err := os.ReadFile(dst)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithError(err))
		return
	}

	topic, expiry, err := parseCert(certData)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid certificate format").WithError(err))
		return
	}

	// Check if exists
	existing, _ := h.client.APNSConfig.Query().Where(apnsconfig.TopicEQ(topic)).First(context.Background())

	var config *ent.APNSConfig
	if existing != nil {
		config, err = h.client.APNSConfig.
			UpdateOne(existing).
			SetCertFilePath(dst).
			SetExpiry(expiry).
			Save(context.Background())
	} else {
		config, err = h.client.APNSConfig.
			Create().
			SetTopic(topic).
			SetCertFilePath(dst).
			SetKeyFilePath(dst).
			SetExpiry(expiry).
			Save(context.Background())
	}

	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithError(err))
		return
	}

	response.OK(c, dto.APNSConfigResponse{
		ID:           config.ID,
		Topic:        config.Topic,
		CertFilePath: config.CertFilePath,
		Expiry:       config.Expiry,
		CreatedAt:    config.CreatedAt,
		UpdatedAt:    config.UpdatedAt,
	}, "Certificate uploaded successfully")
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
	config, err := h.client.APNSConfig.Query().First(context.Background())
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrNotFound.WithMessage("APNs config not found"))
		return
	}

	response.OK(c, dto.APNSConfigResponse{
		ID:           config.ID,
		Topic:        config.Topic,
		CertFilePath: config.CertFilePath,
		Expiry:       config.Expiry,
		CreatedAt:    config.CreatedAt,
		UpdatedAt:    config.UpdatedAt,
	}, "Certificate retrieved successfully")
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
