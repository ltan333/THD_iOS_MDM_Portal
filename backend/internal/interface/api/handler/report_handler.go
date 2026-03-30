package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/thienel/go-backend-template/internal/usecase/service"
	"github.com/thienel/go-backend-template/pkg/query"
	"github.com/thienel/go-backend-template/pkg/response"
)

var reportAllowedFields = map[string]bool{
	"search": true,
}

type ReportHandler interface {
	ExportDevices(c *gin.Context)
	ExportAlerts(c *gin.Context)
	ExportApplications(c *gin.Context)
}

type reportHandlerImpl struct {
	reportService service.ReportService
}

func NewReportHandler(reportService service.ReportService) ReportHandler {
	return &reportHandlerImpl{reportService: reportService}
}

// ExportDevices godoc
// @Summary Export devices to CSV
// @Description Generate and download a CSV file containing the current inventory of enrolled devices, with optional search filtering.
// @Tags Reports
// @Produce text/csv
// @Param search query string false "Filter by device name, serial number, or model"
// @Success 200 {string} string "CSV file content"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/reports/devices/export [get]
func (h *reportHandlerImpl) ExportDevices(c *gin.Context) {
	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	opts := query.ParseQueryParams(params, reportAllowedFields)

	csvData, err := h.reportService.ExportDevicesCSV(c.Request.Context(), opts)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	filename := fmt.Sprintf("devices_export_%s.csv", time.Now().Format("20060102150405"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(200, "text/csv", csvData)
}

// ExportAlerts godoc
// @Summary Export alerts to CSV
// @Description Generate and download a CSV file containing all generated security and system alerts.
// @Tags Reports
// @Produce text/csv
// @Param search query string false "Filter alerts by description or severity"
// @Success 200 {string} string "CSV file content"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/reports/alerts/export [get]
func (h *reportHandlerImpl) ExportAlerts(c *gin.Context) {
	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	opts := query.ParseQueryParams(params, reportAllowedFields)

	csvData, err := h.reportService.ExportAlertsCSV(c.Request.Context(), opts)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	filename := fmt.Sprintf("alerts_export_%s.csv", time.Now().Format("20060102150405"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(200, "text/csv", csvData)
}

// ExportApplications godoc
// @Summary Export applications to CSV
// @Description Generate and download a CSV file containing all tracked applications and their basic metadata.
// @Tags Reports
// @Produce text/csv
// @Param search query string false "Filter by bundle ID or application name"
// @Success 200 {string} string "CSV file content"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/reports/applications/export [get]
func (h *reportHandlerImpl) ExportApplications(c *gin.Context) {
	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	opts := query.ParseQueryParams(params, reportAllowedFields)

	csvData, err := h.reportService.ExportApplicationsCSV(c.Request.Context(), opts)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	filename := fmt.Sprintf("applications_export_%s.csv", time.Now().Format("20060102150405"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(200, "text/csv", csvData)
}
