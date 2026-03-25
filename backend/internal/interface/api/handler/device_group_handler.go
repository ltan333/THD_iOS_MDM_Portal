package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
	"github.com/thienel/go-backend-template/pkg/response"
)

var deviceGroupAllowedFields = map[string]bool{
	"name":   true,
	"search": true,
}

type DeviceGroupHandler interface {
	List(c *gin.Context)
	GetByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	AddDevices(c *gin.Context)
	RemoveDevice(c *gin.Context)
}

type deviceGroupHandlerImpl struct {
	groupService service.DeviceGroupService
}

func NewDeviceGroupHandler(groupService service.DeviceGroupService) DeviceGroupHandler {
	return &deviceGroupHandlerImpl{groupService: groupService}
}

// @Summary List device groups
// @Description Fetch device groups with pagination and filtering
// @Tags DeviceGroups
// @Produce json
// @Security BearerAuth
// @Router /v1/device-groups [get]
func (h *deviceGroupHandlerImpl) List(c *gin.Context) {
	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	offset, limit := query.GetPagination(params, 20)
	opts := query.ParseQueryParams(params, deviceGroupAllowedFields)

	groups, total, err := h.groupService.List(c.Request.Context(), offset, limit, opts)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	res := make([]dto.DeviceGroupResponse, 0, len(groups))
	for _, g := range groups {
		res = append(res, mapGroupToResponse(g))
	}

	page := (offset / limit) + 1
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response.OK(c, dto.ListResponse[dto.DeviceGroupResponse]{
		Items:      res,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, "")
}

// @Summary Get device group by ID
// @Description Fetch single device group
// @Tags DeviceGroups
// @Produce json
// @Param id path int true "Group ID"
// @Security BearerAuth
// @Router /v1/device-groups/{id} [get]
func (h *deviceGroupHandlerImpl) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("ID không hợp lệ"))
		return
	}

	g, err := h.groupService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, mapGroupToResponse(g), "")
}

// @Summary Create device group
// @Description Create a new device group
// @Tags DeviceGroups
// @Produce json
// @Accept json
// @Param request body dto.CreateDeviceGroupRequest true "Group info"
// @Security BearerAuth
// @Router /v1/device-groups [post]
func (h *deviceGroupHandlerImpl) Create(c *gin.Context) {
	var req dto.CreateDeviceGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Dữ liệu không hợp lệ"))
		return
	}

	g, err := h.groupService.Create(c.Request.Context(), service.CreateDeviceGroupCommand{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.Created(c, mapGroupToResponse(g), "Tạo nhóm thành công")
}

// @Summary Update device group
// @Description Update an existing device group
// @Tags DeviceGroups
// @Produce json
// @Accept json
// @Param id path int true "Group ID"
// @Param request body dto.UpdateDeviceGroupRequest true "Group info"
// @Security BearerAuth
// @Router /v1/device-groups/{id} [put]
func (h *deviceGroupHandlerImpl) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("ID không hợp lệ"))
		return
	}

	var req dto.UpdateDeviceGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Dữ liệu không hợp lệ"))
		return
	}

	g, err := h.groupService.Update(c.Request.Context(), service.UpdateDeviceGroupCommand{
		ID:          uint(id),
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, mapGroupToResponse(g), "Cập nhật nhóm thành công")
}

// @Summary Delete device group
// @Description Delete an existing device group
// @Tags DeviceGroups
// @Produce json
// @Param id path int true "Group ID"
// @Security BearerAuth
// @Router /v1/device-groups/{id} [delete]
func (h *deviceGroupHandlerImpl) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("ID không hợp lệ"))
		return
	}

	if err := h.groupService.Delete(c.Request.Context(), uint(id)); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK[any](c, nil, "Xóa nhóm thành công")
}

// @Summary Add devices to group
// @Description Assign one or more devices to a group
// @Tags DeviceGroups
// @Produce json
// @Accept json
// @Param id path int true "Group ID"
// @Param request body dto.ManageGroupDevicesRequest true "Device IDs"
// @Security BearerAuth
// @Router /v1/device-groups/{id}/devices [post]
func (h *deviceGroupHandlerImpl) AddDevices(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("ID không hợp lệ"))
		return
	}

	var req dto.ManageGroupDevicesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Dữ liệu không hợp lệ"))
		return
	}

	if err := h.groupService.AddDevices(c.Request.Context(), uint(id), req.DeviceIDs); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK[any](c, nil, "Thêm thiết bị vào nhóm thành công")
}

// @Summary Remove device from group
// @Description Unassign a single device from a group
// @Tags DeviceGroups
// @Produce json
// @Param id path int true "Group ID"
// @Param deviceId path string true "Device ID"
// @Security BearerAuth
// @Router /v1/device-groups/{id}/devices/{deviceId} [delete]
func (h *deviceGroupHandlerImpl) RemoveDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("ID không hợp lệ"))
		return
	}

	deviceId := c.Param("deviceId")
	if deviceId == "" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Thiếu tham số deviceId"))
		return
	}

	if err := h.groupService.RemoveDevice(c.Request.Context(), uint(id), deviceId); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK[any](c, nil, "Loại bỏ thiết bị khỏi nhóm thành công")
}


func mapGroupToResponse(g *ent.DeviceGroup) dto.DeviceGroupResponse {
	resp := dto.DeviceGroupResponse{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}

	if g.Edges.Devices != nil {
		resp.DeviceCount = len(g.Edges.Devices)
		for _, d := range g.Edges.Devices {
			resp.Devices = append(resp.Devices, mapDeviceToResponse(d))
		}
	}

	return resp
}
