package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
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

func (h *nanocmdHandler) GetVersion(c *gin.Context) {
	resp, err := h.service.GetVersion(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *nanocmdHandler) StartWorkflow(c *gin.Context) {
	name := c.Param("name")
	ids := c.QueryArray("id")
	ctxStr := c.Query("context")

	resp, err := h.service.StartWorkflow(c.Request.Context(), name, ids, ctxStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *nanocmdHandler) GetEvent(c *gin.Context) {
	name := c.Param("name")
	resp, err := h.service.GetEvent(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *nanocmdHandler) PutEvent(c *gin.Context) {
	name := c.Param("name")
	var sub dto.EventSubscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.PutEvent(c.Request.Context(), name, &sub); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *nanocmdHandler) GetFVEnableProfileTemplate(c *gin.Context) {
	data, err := h.service.GetFVEnableProfileTemplate(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/x-apple-aspen-config", data)
}

func (h *nanocmdHandler) GetProfile(c *gin.Context) {
	name := c.Param("name")
	data, err := h.service.GetProfile(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/x-apple-aspen-config", data)
}

func (h *nanocmdHandler) PutProfile(c *gin.Context) {
	name := c.Param("name")
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.PutProfile(c.Request.Context(), name, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *nanocmdHandler) DeleteProfile(c *gin.Context) {
	name := c.Param("name")
	if err := h.service.DeleteProfile(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *nanocmdHandler) GetProfiles(c *gin.Context) {
	names := c.QueryArray("name")
	resp, err := h.service.GetProfiles(c.Request.Context(), names)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *nanocmdHandler) GetCMDPlan(c *gin.Context) {
	name := c.Param("name")
	resp, err := h.service.GetCMDPlan(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *nanocmdHandler) PutCMDPlan(c *gin.Context) {
	name := c.Param("name")
	var plan dto.CMDPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.PutCMDPlan(c.Request.Context(), name, &plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *nanocmdHandler) GetInventory(c *gin.Context) {
	ids := c.QueryArray("id")
	resp, err := h.service.GetInventory(c.Request.Context(), ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *nanocmdHandler) Webhook(c *gin.Context) {
	var webhook dto.NanoCMDWebhook
	if err := c.ShouldBindJSON(&webhook); err != nil {
		tlog.Error("Failed to bind webhook", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	// Process webhook logic here (e.g., update device status)
	tlog.Info("Received NanoCMD webhook", zap.String("topic", webhook.Topic))

	c.Status(http.StatusOK)
}
