package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/mdmcmd"
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
	Restart(c *gin.Context)
	Shutdown(c *gin.Context)
	InstallProfile(c *gin.Context)
	RemoveProfile(c *gin.Context)
	RequestInfo(c *gin.Context)
}

type deviceHandlerImpl struct {
	deviceService  service.DeviceService
	mdmService     service.NanoMDMService
	profileService service.ProfileService
	cmdBuilder     *mdmcmd.CommandBuilder
}

func NewDeviceHandler(
	deviceService service.DeviceService,
	mdmService service.NanoMDMService,
	profileService service.ProfileService,
	cmdBuilder *mdmcmd.CommandBuilder,
) DeviceHandler {
	return &deviceHandlerImpl{
		deviceService:  deviceService,
		mdmService:     mdmService,
		profileService: profileService,
		cmdBuilder:     cmdBuilder,
	}
}

// List godoc
// @Summary List devices
// @Description Retrieve a paginated list of managed devices with support for sorting and filtering by serial number, model, platform, and status.
// @Tags Devices
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 20)"
// @Param serial_number query string false "Filter by serial number"
// @Param model query string false "Filter by device model"
// @Param platform query string false "Filter by platform (iOS, macOS, etc.)Status"
// @Param status query string false "Filter by status"
// @Param search query string false "Search in serial number, name, and model"
// @Success 200 {object} response.APIResponse[dto.ListResponse[dto.DeviceResponse]] "List of devices"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/devices [get]
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

// GetByID godoc
// @Summary Get device by ID
// @Description Fetch detailed information for a single managed device using its unique system ID.
// @Tags Devices
// @Produce json
// @Param id path string true "Device ID"
// @Success 200 {object} response.APIResponse[dto.DeviceResponse] "Device details"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID format"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Device not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/devices/{id} [get]
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

// Export godoc
// @Summary Export devices
// @Description Export all managed devices to a downloadable CSV or JSON file.
// @Tags Devices
// @Produce text/csv
// @Produce application/json
// @Param format query string false "Export format (csv or json, default: csv)"
// @Success 200 {file} file "Exported device data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/devices/export [get]
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

// Lock godoc
// @Summary Lock device
// @Description Enqueue a remote lock command for an MDM-enrolled device. Optionally specify a PIN, message, and phone number to display on the lock screen.
// @Tags Device Actions
// @Accept json
// @Produce json
// @Param id path string true "Device ID (UDID)"
// @Param request body dto.DeviceLockRequest false "Lock configuration options"
// @Success 200 {object} response.APIResponse[dto.APIResult] "Lock command successfully enqueued"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data or device ID"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Device not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/devices/{id}/lock [post]
func (h *deviceHandlerImpl) Lock(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Device ID is required"))
		return
	}

	udid, err := h.deviceService.GetUDID(c.Request.Context(), id)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	var req dto.DeviceLockRequest
	if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid JSON payload").WithError(err))
		return
	}

	opts := &mdmcmd.DeviceLockOptions{
		PIN:         req.PIN,
		Message:     req.Message,
		PhoneNumber: req.PhoneNumber,
	}

	cmdData, _, err := h.cmdBuilder.DeviceLock(opts)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithError(err))
		return
	}

	result, err := h.mdmService.EnqueueCommand(c.Request.Context(), udid, cmdData)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, result, "Lock command queued successfully")
}

// Wipe godoc
// @Summary Wipe device
// @Description Enqueue a remote wipe (factory reset) command for an MDM-enrolled device. This action is irreversible.
// @Tags Device Actions
// @Accept json
// @Produce json
// @Param id path string true "Device ID (UDID)"
// @Param request body dto.DeviceWipeRequest false "Wipe configuration options"
// @Success 200 {object} response.APIResponse[dto.APIResult] "Wipe command successfully enqueued"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Device not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/devices/{id}/wipe [post]
func (h *deviceHandlerImpl) Wipe(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Device ID is required"))
		return
	}

	udid, err := h.deviceService.GetUDID(c.Request.Context(), id)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	var req dto.DeviceWipeRequest
	if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid JSON payload").WithError(err))
		return
	}

	opts := &mdmcmd.EraseDeviceOptions{
		PIN:                    req.PIN,
		PreserveDataPlan:       req.PreserveDataPlan,
		DisallowProximitySetup: req.DisallowProximitySetup,
		ObliterationBehavior:   req.ObliterationBehavior,
	}

	cmdData, _, err := h.cmdBuilder.EraseDevice(opts)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithError(err))
		return
	}

	result, err := h.mdmService.EnqueueCommand(c.Request.Context(), udid, cmdData)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, result, "Wipe command queued successfully")
}

// Restart godoc
// @Summary Restart device
// @Description Enqueue a remote restart command for an MDM-enrolled device.
// @Tags Device Actions
// @Accept json
// @Produce json
// @Param id path string true "Device ID (UDID)"
// @Param request body dto.DeviceRestartRequest false "Restart options"
// @Success 200 {object} response.APIResponse[dto.APIResult] "Restart command successfully enqueued"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Device not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/devices/{id}/restart [post]
func (h *deviceHandlerImpl) Restart(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Device ID is required"))
		return
	}

	udid, err := h.deviceService.GetUDID(c.Request.Context(), id)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	var req dto.DeviceRestartRequest
	if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid JSON payload").WithError(err))
		return
	}

	opts := &mdmcmd.RestartDeviceOptions{
		NotifyUser: req.NotifyUser,
	}

	cmdData, _, err := h.cmdBuilder.RestartDevice(opts)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithError(err))
		return
	}

	result, err := h.mdmService.EnqueueCommand(c.Request.Context(), udid, cmdData)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, result, "Restart command queued successfully")
}

// Shutdown godoc
// @Summary Shutdown device
// @Description Enqueue a remote shutdown command for an MDM-enrolled device.
// @Tags Device Actions
// @Produce json
// @Param id path string true "Device ID (UDID)"
// @Success 200 {object} response.APIResponse[dto.APIResult] "Shutdown command successfully enqueued"
// @Failure 400 {object} response.APIResponse[any] "Invalid device ID"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Device not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/devices/{id}/shutdown [post]
func (h *deviceHandlerImpl) Shutdown(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Device ID is required"))
		return
	}

	udid, err := h.deviceService.GetUDID(c.Request.Context(), id)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	cmdData, _, err := h.cmdBuilder.ShutDownDevice()
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithError(err))
		return
	}

	result, err := h.mdmService.EnqueueCommand(c.Request.Context(), udid, cmdData)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, result, "Shutdown command queued successfully")
}

// InstallProfile godoc
// @Summary Install profile on device
// @Description Enqueue an command to install a specific configuration profile onto an MDM-enrolled device.
// @Tags Device Actions
// @Accept json
// @Produce json
// @Param id path string true "Device ID (UDID)"
// @Param request body dto.DeviceInstallProfileRequest true "Profile selection"
// @Success 200 {object} response.APIResponse[dto.DeviceActionResponse] "Installation command successfully enqueued"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Device or profile not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/devices/{id}/install-profile [post]
func (h *deviceHandlerImpl) InstallProfile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Device ID is required"))
		return
	}

	var req dto.DeviceInstallProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithError(err))
		return
	}

	// Pass portal device ID; profileService resolves UDID internally.
	err := h.profileService.InstallOnDevice(c.Request.Context(), req.ProfileID, id)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, dto.DeviceActionResponse{
		RequestType: "InstallProfile",
		Status:      "queued",
		Message:     "Profile installation queued successfully",
	}, "Profile installation queued successfully")
}

// RemoveProfile godoc
// @Summary Remove profile from device
// @Description Enqueue a command to remove a configuration profile from an MDM-enrolled device using its identifier.
// @Tags Device Actions
// @Accept json
// @Produce json
// @Param id path string true "Device ID (UDID)"
// @Param request body dto.DeviceRemoveProfileRequest true "Profile identifier"
// @Success 200 {object} response.APIResponse[dto.APIResult] "Removal command successfully enqueued"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Device not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/devices/{id}/remove-profile [post]
func (h *deviceHandlerImpl) RemoveProfile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Device ID is required"))
		return
	}

	udid, err := h.deviceService.GetUDID(c.Request.Context(), id)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	var req dto.DeviceRemoveProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithError(err))
		return
	}

	cmdData, _, err := h.cmdBuilder.RemoveProfile(req.ProfileIdentifier)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithError(err))
		return
	}

	result, err := h.mdmService.EnqueueCommand(c.Request.Context(), udid, cmdData)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, result, "Remove profile command queued successfully")
}

// RequestInfo godoc
// @Summary Request device information
// @Description Enqueue an MDM command to query specific hardware and software attributes of an enrolled device.
// @Tags Device Actions
// @Accept json
// @Produce json
// @Param id path string true "Device ID (UDID)"
// @Param request body dto.DeviceInfoRequest false "Optional specific information queries"
// @Success 200 {object} response.APIResponse[dto.APIResult] "Information request successfully enqueued"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Device not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/devices/{id}/request-info [post]
func (h *deviceHandlerImpl) RequestInfo(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Device ID is required"))
		return
	}

	udid, err := h.deviceService.GetUDID(c.Request.Context(), id)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	var req dto.DeviceInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid JSON payload").WithError(err))
		return
	}

	queries := req.Queries
	if len(queries) == 0 {
		queries = mdmcmd.CommonDeviceQueries()
	}

	cmdData, _, err := h.cmdBuilder.DeviceInformation(queries)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithError(err))
		return
	}

	result, err := h.mdmService.EnqueueCommand(c.Request.Context(), udid, cmdData)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, result, "Device information request queued successfully")
}

func mapDeviceToResponse(d *ent.Device) dto.DeviceResponse {
	var ownerID *uint
	if d.OwnerID != 0 {
		o := d.OwnerID
		ownerID = &o
	}

	// Expose the MDM UDID separately from the portal's internal ID.
	udidVal := ""
	if d.Udid != nil {
		udidVal = *d.Udid
	}

	resp := dto.DeviceResponse{
		ID:               d.ID,
		UDID:             udidVal,
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
