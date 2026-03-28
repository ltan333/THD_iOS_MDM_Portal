package handler

import "github.com/gin-gonic/gin"

// Health godoc
// @Summary System health check
// @Description Verify that the MDM backend service is operational and reachable.
// @Tags System
// @Produce json
// @Success 200 {object} HealthResponse "Service is healthy"
// @Router /api/v1/health [get]
func Health(c *gin.Context) {
	c.JSON(200, HealthResponse{Status: "ok"})
}
