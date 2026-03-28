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
// @Description Fetch overall system statistics for the dashboard, including total devices, users, active alerts, and managed applications.
// @Tags Dashboard
// @Produce json
// @Success 200 {object} response.APIResponse[dto.DashboardStatsResponse] "Dashboard statistics"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
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
// @Description Retrieve a breakdown of device statistics, including counts by OS platform (iOS, macOS, etc.) and enrollment status.
// @Tags Dashboard
// @Produce json
// @Success 200 {object} response.APIResponse[dto.DeviceStatsResponse] "Device breakdown statistics"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
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
// @Description Get a summary of system alerts, categorized by severity (Critical, Warning, Info) and resolution status.
// @Tags Dashboard
// @Produce json
// @Success 200 {object} response.APIResponse[dto.AlertsSummaryResponse] "Alerts summary data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
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
// @Description Fetch time-series or categorical data for dashboard charts (e.g., device enrollment trends, compliance rates).
// @Tags Dashboard
// @Produce json
// @Param type path string true "Chart type (devices, compliance, alerts)"
// @Success 200 {object} response.APIResponse[dto.ChartDataResponse] "Chart visualization data"
// @Failure 400 {object} response.APIResponse[any] "Invalid chart type requested"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
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
