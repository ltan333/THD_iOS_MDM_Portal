package handler

import "github.com/gin-gonic/gin"

// Health godoc
// @Summary Health check
// @Description Check service liveness
// @Tags System
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func Health(c *gin.Context) {
	c.JSON(200, HealthResponse{Status: "ok"})
}
