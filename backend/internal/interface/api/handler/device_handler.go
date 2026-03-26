package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
	"github.com/thienel/go-backend-template/pkg/response"
)

var deviceAllowedFields = map[string]bool{
	"serial_number": true,
	"model":         true,
	"platform":      true,
	"status":        true,
	"is_enrolled":   true,
	"search":        true,
}

type DeviceHandler interface {
	List(c *gin.Context)
	GetByID(c *gin.Context)
	Export(c *gin.Context)
	Lock(c *gin.Context)
	Wipe(c *gin.Context)
}

type deviceHandlerImpl struct {
	deviceService service.DeviceService
	mdmService    service.NanoMDMService
}

func NewDeviceHandler(deviceService service.DeviceService, mdmService service.NanoMDMService) DeviceHandler {
	return &deviceHandlerImpl{
		deviceService: deviceService,
		mdmService:    mdmService,
	}
}

// @Summary List devices
// @Description Fetch devices with pagination, sorting and filtering
// @Tags Devices
// @Produce json
// @Security BearerAuth
// @Router /v1/devices [get]
func (h *deviceHandlerImpl) List(c *gin.Context) {
	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	offset, limit := query.GetPagination(params, 20)
	opts := query.ParseQueryParams(params, deviceAllowedFields)

	devices, total, err := h.deviceService.List(c.Request.Context(), offset, limit, opts)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	res := make([]dto.DeviceResponse, 0, len(devices))
	for _, dev := range devices {
		res = append(res, mapDeviceToResponse(dev))
	}

	page := (offset / limit) + 1
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response.OK(c, dto.ListResponse[dto.DeviceResponse]{
		Items:      res,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, "")
}

// @Summary Get device by ID
// @Description Fetch single device details
// @Tags Devices
// @Produce json
// @Param id path string true "Device ID"
// @Security BearerAuth
// @Router /v1/devices/{id} [get]
func (h *deviceHandlerImpl) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Thiếu tham số ID"))
		return
	}

	dev, err := h.deviceService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, mapDeviceToResponse(dev), "")
}

// @Summary Export devices
// @Description Export devices to CSV or JSON
// @Tags Devices
// @Produce text/csv
// @Param format query string false "Format (csv or json)"
// @Security BearerAuth
// @Router /v1/devices/export [get]
func (h *deviceHandlerImpl) Export(c *gin.Context) {
	format := c.Query("format")
	if format == "" {
		format = "csv"
	}

	data, err := h.deviceService.Export(c.Request.Context(), format)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	if format == "csv" {
		c.Header("Content-Disposition", "attachment; filename=devices.csv")
		c.Data(200, "text/csv", data)
	} else {
		c.Header("Content-Disposition", "attachment; filename=devices.json")
		c.Data(200, "application/json", data)
	}
}

// @Summary Lock device
// @Description Queues MDM device lock command
// @Tags Devices
// @Produce json
// @Param id path string true "Device ID"
// @Security BearerAuth
// @Router /v1/devices/{id}/lock [post]
func (h *deviceHandlerImpl) Lock(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Thiếu tham số ID"))
		return
	}

	// Simple DeviceLock XML command
	cmdXML := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Command</key>
	<dict>
		<key>RequestType</key>
		<string>DeviceLock</string>
	</dict>
	<key>CommandUUID</key>
	<string>Lock-` + id + `</string>
</dict>
</plist>`

	result, err := h.mdmService.EnqueueCommand(c.Request.Context(), id, []byte(cmdXML))
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, result, "Lệnh Lock đã được đẩy vào hàng đợi")
}

// @Summary Wipe device
// @Description Queues MDM erase device command
// @Tags Devices
// @Produce json
// @Param id path string true "Device ID"
// @Security BearerAuth
// @Router /v1/devices/{id}/wipe [post]
func (h *deviceHandlerImpl) Wipe(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Thiếu tham số ID"))
		return
	}

	// Simple EraseDevice XML command
	cmdXML := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Command</key>
	<dict>
		<key>RequestType</key>
		<string>EraseDevice</string>
	</dict>
	<key>CommandUUID</key>
	<string>Wipe-` + id + `</string>
</dict>
</plist>`

	result, err := h.mdmService.EnqueueCommand(c.Request.Context(), id, []byte(cmdXML))
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, result, "Lệnh Wipe đã được đẩy vào hàng đợi")
}

func mapDeviceToResponse(d *ent.Device) dto.DeviceResponse {
	var ownerID *uint
	if d.OwnerID != 0 {
		o := d.OwnerID
		ownerID = &o
	}

	resp := dto.DeviceResponse{
		ID:               d.ID,
		SerialNumber:     d.SerialNumber,
		Model:            d.Model,
		Name:             d.Name,
		Platform:         string(d.Platform),
		Status:           string(d.Status),
		ComplianceStatus: string(d.ComplianceStatus),
		OsVersion:        d.OsVersion,
		DeviceType:       d.DeviceType,
		MacAddress:       d.MACAddress,
		IpAddress:        d.IPAddress,
		BatteryLevel:     d.BatteryLevel,
		StorageCapacity:  d.StorageCapacity,
		StorageUsed:      d.StorageUsed,
		IsJailbroken:     d.IsJailbroken,
		IsEnrolled:       d.IsEnrolled,
		EnrollmentType:   string(d.EnrollmentType),
		OwnerID:          ownerID,
		CreatedAt:        d.CreatedAt,
		UpdatedAt:        d.UpdatedAt,
	}

	if !d.LastSeen.IsZero() {
		ls := d.LastSeen
		resp.LastSeen = &ls
	}
	if !d.EnrolledAt.IsZero() {
		ea := d.EnrolledAt
		resp.EnrolledAt = &ea
	}

	return resp
}
