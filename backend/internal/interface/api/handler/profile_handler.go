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

var profileAllowedFields = map[string]bool{
	"name":     true,
	"platform": true,
	"scope":    true,
	"status":   true,
	"search":   true,
}

type ProfileHandler interface {
	List(c *gin.Context)
	GetByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	UpdateStatus(c *gin.Context)

	UpdateSecuritySettings(c *gin.Context)
	UpdateNetworkConfig(c *gin.Context)
	UpdateRestrictions(c *gin.Context)
	UpdateContentFilter(c *gin.Context)
	UpdateComplianceRules(c *gin.Context)

	Assign(c *gin.Context)
	Unassign(c *gin.Context)
	ListAssignments(c *gin.Context)

	ListVersions(c *gin.Context)
	Rollback(c *gin.Context)
	GetDeploymentStatus(c *gin.Context)
	Repush(c *gin.Context)
	Duplicate(c *gin.Context)
}

type profileHandlerImpl struct {
	profileService service.ProfileService
}

func NewProfileHandler(profileService service.ProfileService) ProfileHandler {
	return &profileHandlerImpl{profileService: profileService}
}

// List godoc
// @Summary List profiles
// @Description Retrieve a paginated list of configuration profiles with support for filtering by name, platform, scope, and status.
// @Tags Profiles
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 20)"
// @Param name query string false "Filter by name"
// @Param platform query string false "Filter by platform (iOS, macOS, etc.)"
// @Param scope query string false "Filter by scope (system, user)"
// @Param status query string false "Filter by status (active, draft, archived)"
// @Param search query string false "Search in name and description"
// @Success 200 {object} response.APIResponse[dto.ListResponse[dto.ProfileResponse]] "List of profiles"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/profiles [get]
func (h *profileHandlerImpl) List(c *gin.Context) {
	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	offset, limit := query.GetPagination(params, 20)
	opts := query.ParseQueryParams(params, profileAllowedFields)

	profiles, total, err := h.profileService.List(c.Request.Context(), offset, limit, opts)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	res := make([]dto.ProfileResponse, 0, len(profiles))
	for _, p := range profiles {
		res = append(res, mapProfileToResponse(p))
	}

	page := (offset / limit) + 1
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response.OK(c, dto.ListResponse[dto.ProfileResponse]{
		Items:      res,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, "")
}

// GetByID godoc
// @Summary Get profile by ID
// @Description Fetch detailed information for a single configuration profile including its settings and payloads.
// @Tags Profiles
// @Produce json
// @Param id path int true "Profile ID"
// @Success 200 {object} response.APIResponse[dto.ProfileResponse] "Profile details"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID format"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/profiles/{id} [get]
func (h *profileHandlerImpl) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	p, err := h.profileService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, mapProfileToResponse(p), "")
}

// Create godoc
// @Summary Create profile
// @Description Create a new configuration profile with initial settings and assignments.
// @Tags Profiles
// @Accept json
// @Produce json
// @Param request body dto.CreateProfileRequest true "New profile definition"
// @Success 201 {object} response.APIResponse[dto.ProfileResponse] "Profile created successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/profiles [post]
func (h *profileHandlerImpl) Create(c *gin.Context) {
	var req dto.CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid data"))
		return
	}

	p, err := h.profileService.Create(c.Request.Context(), service.CreateProfileCommand{
		Name:             req.Name,
		Platform:         req.Platform,
		Scope:            req.Scope,
		SecuritySettings: req.SecuritySettings,
		NetworkConfig:    req.NetworkConfig,
		Restrictions:     req.Restrictions,
		ContentFilter:    req.ContentFilter,
		ComplianceRules:  req.ComplianceRules,
		Payloads:         req.Payloads,
	})
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.Created(c, mapProfileToResponse(p), "Profile created successfully")
}

// Update godoc
// @Summary Update profile
// @Description Modify an existing configuration profile's settings and metadata.
// @Tags Profiles
// @Accept json
// @Produce json
// @Param id path int true "Profile ID"
// @Param request body dto.UpdateProfileRequest true "Updated profile definition"
// @Success 200 {object} response.APIResponse[dto.ProfileResponse] "Profile updated successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID or request data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/profiles/{id} [put]
func (h *profileHandlerImpl) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid data"))
		return
	}

	p, err := h.profileService.Update(c.Request.Context(), service.UpdateProfileCommand{
		ID:               uint(id),
		Name:             req.Name,
		Platform:         req.Platform,
		Scope:            req.Scope,
		SecuritySettings: req.SecuritySettings,
		NetworkConfig:    req.NetworkConfig,
		Restrictions:     req.Restrictions,
		ContentFilter:    req.ContentFilter,
		ComplianceRules:  req.ComplianceRules,
		Payloads:         req.Payloads,
	})
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, mapProfileToResponse(p), "Profile updated successfully")
}

// Delete godoc
// @Summary Delete profile
// @Description Permanently remove a configuration profile and its associated data.
// @Tags Profiles
// @Produce json
// @Param id path int true "Profile ID"
// @Success 200 {object} response.APIResponse[any] "Profile deleted successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID format"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/profiles/{id} [delete]
func (h *profileHandlerImpl) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	if err := h.profileService.Delete(c.Request.Context(), uint(id)); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK[any](c, nil, "Profile deleted successfully")
}

// UpdateStatus godoc
// @Summary Update profile status
// @Description Update the operational status of a profile (active, draft, archived).
// @Tags Profiles
// @Accept json
// @Produce json
// @Param id path int true "Profile ID"
// @Param request body dto.UpdateProfileStatusRequest true "New status selection"
// @Success 200 {object} response.APIResponse[any] "Profile status updated successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data or status"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/profiles/{id}/status [put]
func (h *profileHandlerImpl) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}
	var req dto.UpdateProfileStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid data"))
		return
	}
	if err := h.profileService.UpdateStatus(c.Request.Context(), uint(id), req.Status); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK[any](c, nil, "Profile status updated successfully")
}

// UpdateSecuritySettings godoc
// @Summary Update security settings
// @Description Modify the security-related platform configurations for a specific profile.
// @Tags Profiles
// @Accept json
// @Produce json
// @Param id path int true "Profile ID"
// @Param request body dto.UpdateSecuritySettingsRequest true "Security settings map"
// @Success 200 {object} response.APIResponse[any] "Security settings updated successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid request dataMap"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/profiles/{id}/settings/security [put]
func (h *profileHandlerImpl) UpdateSecuritySettings(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid data"))
		return
	}
	if err := h.profileService.UpdateSecuritySettings(c.Request.Context(), uint(id), req); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK[any](c, nil, "Security settings updated successfully")
}

// UpdateNetworkConfig godoc
// @Summary Update network configuration
// @Description Modify network connectivity settings (Wi-Fi, VPN, Proxy) for a profile.
// @Tags Profiles
// @Accept json
// @Produce json
// @Param id path int true "Profile ID"
// @Param request body dto.UpdateNetworkConfigRequest true "Network settings map"
// @Success 200 {object} response.APIResponse[any] "Network config updated successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/profiles/{id}/settings/network [put]
func (h *profileHandlerImpl) UpdateNetworkConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid data"))
		return
	}
	if err := h.profileService.UpdateNetworkConfig(c.Request.Context(), uint(id), req); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK[any](c, nil, "Network config updated successfully")
}

// UpdateRestrictions godoc
// @Summary Update device restrictions
// @Description Modify hardware and software usage restrictions applied by this profile.
// @Tags Profiles
// @Accept json
// @Produce json
// @Param id path int true "Profile ID"
// @Param request body dto.UpdateRestrictionsRequest true "Restrictions setting map"
// @Success 200 {object} response.APIResponse[any] "Restrictions updated successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/profiles/{id}/settings/restrictions [put]
func (h *profileHandlerImpl) UpdateRestrictions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid data"))
		return
	}
	if err := h.profileService.UpdateRestrictions(c.Request.Context(), uint(id), req); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK[any](c, nil, "Restrictions updated successfully")
}

// UpdateContentFilter godoc
// @Summary Update content filters
// @Description Modify web content and domain access filters for a profile.
// @Tags Profiles
// @Accept json
// @Produce json
// @Param id path int true "Profile ID"
// @Param request body dto.UpdateContentFilterRequest true "Content filter map"
// @Success 200 {object} response.APIResponse[any] "Content filter updated successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/profiles/{id}/settings/content-filter [put]
func (h *profileHandlerImpl) UpdateContentFilter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid data"))
		return
	}
	if err := h.profileService.UpdateContentFilter(c.Request.Context(), uint(id), req); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK[any](c, nil, "Content filter updated successfully")
}

// UpdateComplianceRules godoc
// @Summary Update compliance rules
// @Description Modify automated compliance verification rules and enforcement actions for a profile.
// @Tags Profiles
// @Accept json
// @Produce json
// @Param id path int true "Profile ID"
// @Param request body dto.UpdateComplianceRulesRequest true "Compliance rules map"
// @Success 200 {object} response.APIResponse[any] "Compliance rules updated successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/profiles/{id}/settings/compliance [put]
func (h *profileHandlerImpl) UpdateComplianceRules(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid data"))
		return
	}
	if err := h.profileService.UpdateComplianceRules(c.Request.Context(), uint(id), req); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK[any](c, nil, "Compliance rules updated successfully")
}

// Assign godoc
// @Summary Assign profile
// @Description Assign a profile to a specific device or group for deployment.
// @Tags Profiles
// @Accept json
// @Produce json
// @Param id path int true "Profile ID"
// @Param request body dto.AssignProfileRequest true "Target and scheduling details"
// @Success 200 {object} response.APIResponse[any] "Profile assignment successful"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data or target"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile or target not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/profiles/{id}/assignments [post]
func (h *profileHandlerImpl) Assign(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}
	var req dto.AssignProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid data"))
		return
	}

	assignment, err := h.profileService.Assign(c.Request.Context(), service.AssignProfileCommand{
		ProfileID:    uint(id),
		TargetType:   req.TargetType,
		DeviceID:     req.DeviceID,
		GroupID:      req.GroupID,
		ScheduleType: req.ScheduleType,
		ScheduledAt:  req.ScheduledAt,
	})
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	res := dto.ProfileAssignmentResponse{
		ID:           assignment.ID,
		ProfileID:    assignment.ProfileID,
		TargetType:   string(assignment.TargetType),
		DeviceID:     nil,
		GroupID:      nil,
		ScheduleType: string(assignment.ScheduleType),
		ScheduledAt:  assignment.ScheduledAt,
		CreatedAt:    assignment.CreatedAt,
	}

	// Set device_id or group_id if present
	if assignment.DeviceID != nil && *assignment.DeviceID != "" {
		res.DeviceID = assignment.DeviceID
	}
	if assignment.GroupID != nil && *assignment.GroupID != 0 {
		res.GroupID = assignment.GroupID
	}

	response.OK(c, res, "Profile assigned successfully")
}

// Unassign godoc
// @Summary Unassign profile
// @Description Remove an existing profile assignment from a device or group.
// @Tags Profiles
// @Produce json
// @Param id path int true "Profile ID"
// @Param assignmentId path int true "Assignment record ID"
// @Success 200 {object} response.APIResponse[any] "Profile unassigned successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID formats"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Assignment not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/profiles/{id}/assignments/{assignmentId} [delete]
func (h *profileHandlerImpl) Unassign(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}
	assignIdStr := c.Param("assignmentId")
	assignId, err := strconv.ParseUint(assignIdStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid Assignment ID"))
		return
	}

	if err := h.profileService.Unassign(c.Request.Context(), uint(id), uint(assignId)); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK[any](c, nil, "Profile unassigned successfully")
}

// ListAssignments godoc
// @Summary List assignments
// @Description Retrieve a list of all current assignments and deployment targets for this profile.
// @Tags Profiles
// @Produce json
// @Param id path int true "Profile ID"
// @Success 200 {object} response.APIResponse[[]dto.ProfileAssignmentResponse] "List of assignments"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID format"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/profiles/{id}/assignments [get]
func (h *profileHandlerImpl) ListAssignments(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	assignments, err := h.profileService.ListAssignments(c.Request.Context(), uint(id))
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	res := make([]dto.ProfileAssignmentResponse, 0, len(assignments))
	for _, a := range assignments {
		res = append(res, dto.ProfileAssignmentResponse{
			ID:           a.ID,
			ProfileID:    a.ProfileID,
			TargetType:   string(a.TargetType),
			DeviceID:     a.DeviceID,
			GroupID:      a.GroupID,
			ScheduleType: string(a.ScheduleType),
			CreatedAt:    a.CreatedAt,
		})
	}

	response.OK(c, res, "")
}

// ListVersions godoc
// @Summary List profile versions
// @Description Retrieve the full version history and change logs for a specific configuration profile.
// @Tags Profiles
// @Produce json
// @Param id path int true "Profile ID"
// @Success 200 {object} response.APIResponse[[]dto.ProfileVersionResponse] "Version history"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID format"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/profiles/{id}/versions [get]
func (h *profileHandlerImpl) ListVersions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	versions, err := h.profileService.ListVersions(c.Request.Context(), uint(id))
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	res := make([]dto.ProfileVersionResponse, 0, len(versions))
	for _, v := range versions {
		res = append(res, dto.ProfileVersionResponse{
			ID:          v.ID,
			ProfileID:   v.ProfileID,
			Version:     v.Version,
			Data:        v.Data,
			ChangeNotes: v.ChangeNotes,
			CreatedAt:   v.CreatedAt,
		})
	}

	response.OK(c, res, "")
}

// Rollback godoc
// @Summary Rollback profile version
// @Description Revert a configuration profile to a previously saved version.
// @Tags Profiles
// @Produce json
// @Param id path int true "Profile ID"
// @Param versionId path int true "Target version ID to rollback to"
// @Success 200 {object} response.APIResponse[any] "Rollback successful"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID formats"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile or version not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/profiles/{id}/versions/{versionId}/rollback [post]
func (h *profileHandlerImpl) Rollback(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}
	versionIdStr := c.Param("versionId")
	versionId, err := strconv.ParseUint(versionIdStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid Version ID"))
		return
	}

	if err := h.profileService.Rollback(c.Request.Context(), uint(id), uint(versionId)); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK[any](c, nil, "Profile rolled back successfully")
}

// GetDeploymentStatus godoc
// @Summary Get profile deployment status
// @Description Fetch real-time deployment and installation status across all assigned devices.
// @Tags Profiles
// @Produce json
// @Param id path int true "Profile ID"
// @Success 200 {object} response.APIResponse[[]dto.ProfileDeploymentStatusResponse] "Deployment statuses"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID format"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/profiles/{id}/deployment-status [get]
func (h *profileHandlerImpl) GetDeploymentStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	statuses, err := h.profileService.GetDeploymentStatus(c.Request.Context(), uint(id))
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	res := make([]dto.ProfileDeploymentStatusResponse, 0, len(statuses))
	for _, s := range statuses {
		res = append(res, dto.ProfileDeploymentStatusResponse{
			ID:           s.ID,
			ProfileID:    s.ProfileID,
			DeviceID:     s.DeviceID,
			Status:       string(s.Status),
			ErrorMessage: s.ErrorMessage,
			CreatedAt:    s.CreatedAt,
			UpdatedAt:    s.UpdatedAt,
		})
	}

	response.OK(c, res, "")
}

// Repush godoc
// @Summary Repush profile
// @Description Force a re-deployment of the profile to all assigned devices that haven't successfully installed it.
// @Tags Profiles
// @Produce json
// @Param id path int true "Profile ID"
// @Success 200 {object} response.APIResponse[any] "Repush command initiated"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID format"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/profiles/{id}/repush [post]
func (h *profileHandlerImpl) Repush(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	if err := h.profileService.Repush(c.Request.Context(), uint(id)); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK[any](c, nil, "Profile repush initiated successfully")
}

// Duplicate godoc
// @Summary Duplicate profile
// @Description Create an exact copy of an existing configuration profile for staging or modification.
// @Tags Profiles
// @Produce json
// @Param id path int true "Profile ID"
// @Success 201 {object} response.APIResponse[dto.ProfileResponse] "Profile duplicated successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID format"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/profiles/{id}/duplicate [post]
func (h *profileHandlerImpl) Duplicate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	p, err := h.profileService.Duplicate(c.Request.Context(), uint(id))
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.Created(c, mapProfileToResponse(p), "Profile duplicated successfully")
}

func mapProfileToResponse(p *ent.Profile) dto.ProfileResponse {
	return dto.ProfileResponse{
		ID:               p.ID,
		Name:             p.Name,
		Platform:         string(p.Platform),
		Scope:            string(p.Scope),
		Status:           string(p.Status),
		SecuritySettings: p.SecuritySettings,
		NetworkConfig:    p.NetworkConfig,
		Restrictions:     p.Restrictions,
		ContentFilter:    p.ContentFilter,
		ComplianceRules:  p.ComplianceRules,
		Payloads:         p.Payloads,
		Version:          p.Version,
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
	}
}
