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

// @Summary List profiles
// @Description Fetch profiles with pagination and filtering
// @Tags Profiles
// @Produce json
// @Security BearerAuth
// @Router /v1/profiles [get]
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

// @Summary Get profile by ID
// @Description Fetch single profile
// @Tags Profiles
// @Produce json
// @Param id path int true "Profile ID"
// @Security BearerAuth
// @Router /v1/profiles/{id} [get]
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

// @Summary Create profile
// @Description Create a new configuration profile
// @Tags Profiles
// @Produce json
// @Accept json
// @Param request body dto.CreateProfileRequest true "Profile info"
// @Security BearerAuth
// @Router /v1/profiles [post]
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

// @Summary Update profile
// @Description Update an existing profile
// @Tags Profiles
// @Produce json
// @Accept json
// @Param id path int true "Profile ID"
// @Param request body dto.UpdateProfileRequest true "Profile info"
// @Security BearerAuth
// @Router /v1/profiles/{id} [put]
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

// @Summary Delete profile
// @Description Delete an existing profile
// @Tags Profiles
// @Produce json
// @Param id path int true "Profile ID"
// @Security BearerAuth
// @Router /v1/profiles/{id} [delete]
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

func (h *profileHandlerImpl) UpdateSecuritySettings(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}
	var req map[string]interface{}
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

func (h *profileHandlerImpl) UpdateNetworkConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}
	var req map[string]interface{}
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

func (h *profileHandlerImpl) UpdateRestrictions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}
	var req map[string]interface{}
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

func (h *profileHandlerImpl) UpdateContentFilter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}
	var req map[string]interface{}
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

func (h *profileHandlerImpl) UpdateComplianceRules(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}
	var req map[string]interface{}
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
	
	err = h.profileService.Assign(c.Request.Context(), service.AssignProfileCommand{
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
	response.OK[any](c, nil, "Profile assigned successfully")
}

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
