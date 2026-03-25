package handler

import (
	"github.com/gin-gonic/gin"

	_ "github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/response"
)

// DashboardHandler interface
type DashboardHandler interface {
	GetStats(c *gin.Context)
	GetDeviceStats(c *gin.Context)
	GetAlertsSummary(c *gin.Context)
	GetChartData(c *gin.Context)
}

type dashboardHandlerImpl struct {
	dashboardService service.DashboardService
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(dashboardService service.DashboardService) DashboardHandler {
	return &dashboardHandlerImpl{
		dashboardService: dashboardService,
	}
}

// GetStats godoc
// @Summary Get dashboard statistics
// @Description Get overall dashboard statistics including devices, users, alerts, and apps
// @Tags Dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse[dto.DashboardStatsResponse]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
// @Router /api/v1/dashboard/stats [get]
func (h *dashboardHandlerImpl) GetStats(c *gin.Context) {
	stats, err := h.dashboardService.GetStats(c.Request.Context())
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, stats, "")
}

// GetDeviceStats godoc
// @Summary Get device statistics
// @Description Get detailed device statistics including counts by platform and status
// @Tags Dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse[dto.DeviceStatsResponse]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
// @Router /api/v1/dashboard/device-stats [get]
func (h *dashboardHandlerImpl) GetDeviceStats(c *gin.Context) {
	stats, err := h.dashboardService.GetDeviceStats(c.Request.Context())
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, stats, "")
}

// GetAlertsSummary godoc
// @Summary Get alerts summary
// @Description Get alerts summary including counts by severity and type
// @Tags Dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse[dto.AlertsSummaryResponse]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
// @Router /api/v1/dashboard/alerts-summary [get]
func (h *dashboardHandlerImpl) GetAlertsSummary(c *gin.Context) {
	summary, err := h.dashboardService.GetAlertsSummary(c.Request.Context())
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, summary, "")
}

// GetChartData godoc
// @Summary Get chart data
// @Description Get chart data for dashboard visualization
// @Tags Dashboard
// @Produce json
// @Security BearerAuth
// @Param type path string true "Chart type (devices, compliance, alerts)"
// @Success 200 {object} response.APIResponse[dto.ChartDataResponse]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
// @Router /api/v1/dashboard/charts/{type} [get]
func (h *dashboardHandlerImpl) GetChartData(c *gin.Context) {
	chartType := c.Param("type")
	if chartType == "" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Loại biểu đồ không được để trống"))
		return
	}

	data, err := h.dashboardService.GetChartData(c.Request.Context(), chartType)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, data, "")
}
