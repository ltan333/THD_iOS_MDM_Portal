package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/response"
)

type DepProfileHandler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	SetAsAssigner(c *gin.Context)
	GetAssigner(c *gin.Context)
}

type depProfileHandlerImpl struct {
	service service.DepProfileService
	depName string // default DEP server name from config
}

func NewDepProfileHandler(service service.DepProfileService, depName string) DepProfileHandler {
	return &depProfileHandlerImpl{
		service: service,
		depName: depName,
	}
}

// Create godoc
// @Summary Create DEP profile
// @Description Create a new DEP enrollment profile and register it with Apple DEP
// @Tags DEP Profile
// @Accept json
// @Produce json
// @Param request body dto.DEPProfileRequest true "Profile configuration"
// @Success 201 {object} response.APIResponse[dto.DEPProfileResponse] "Profile created successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid request body"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 409 {object} response.APIResponse[any] "Profile name already exists"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/dep/profiles [post]
func (h *depProfileHandlerImpl) Create(c *gin.Context) {
	var req dto.DEPProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid request body").WithError(err))
		return
	}

	if req.ProfileName == "" {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("profile_name is required"))
		return
	}

	depName := c.Query("dep_name")
	if depName == "" {
		depName = h.depName
	}

	profile, err := h.service.Create(c.Request.Context(), depName, &req)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.Created(c, mapEntToDEPProfileResponse(profile), "Profile created successfully")
}

// GetByID godoc
// @Summary Get DEP profile by ID
// @Description Retrieve a DEP profile by its ID
// @Tags DEP Profile
// @Produce json
// @Param id path int true "Profile ID"
// @Success 200 {object} response.APIResponse[dto.DEPProfileResponse] "Profile retrieved successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid profile ID"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/dep/profiles/{id} [get]
func (h *depProfileHandlerImpl) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid profile ID"))
		return
	}

	profile, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, mapEntToDEPProfileResponse(profile), "Profile retrieved successfully")
}

// List godoc
// @Summary List DEP profiles
// @Description Get a paginated list of DEP profiles
// @Tags DEP Profile
// @Produce json
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(20)
// @Success 200 {object} response.APIResponse[dto.ListResponse[dto.DEPProfileResponse]] "Profiles retrieved successfully"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/dep/profiles [get]
func (h *depProfileHandlerImpl) List(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Validate pagination parameters
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 20
	}

	profiles, total, err := h.service.List(c.Request.Context(), offset, limit)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	items := make([]dto.DEPProfileResponse, len(profiles))
	for i, p := range profiles {
		items[i] = *mapEntToDEPProfileResponse(p)
	}

	page := (offset / limit) + 1
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response.OK(c, dto.ListResponse[dto.DEPProfileResponse]{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, "Profiles retrieved successfully")
}

// Update godoc
// @Summary Update DEP profile
// @Description Update an existing DEP profile and sync with Apple DEP
// @Tags DEP Profile
// @Accept json
// @Produce json
// @Param id path int true "Profile ID"
// @Param request body dto.DEPProfileRequest true "Updated profile configuration"
// @Success 200 {object} response.APIResponse[dto.DEPProfileResponse] "Profile updated successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid request"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/dep/profiles/{id} [put]
func (h *depProfileHandlerImpl) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid profile ID"))
		return
	}

	var req dto.DEPProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid request body").WithError(err))
		return
	}

	depName := c.Query("dep_name")
	if depName == "" {
		depName = h.depName
	}

	profile, err := h.service.Update(c.Request.Context(), depName, uint(id), &req)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, mapEntToDEPProfileResponse(profile), "Profile updated successfully")
}

// Delete godoc
// @Summary Delete DEP profile
// @Description Delete a DEP profile from local database and Apple DEP
// @Tags DEP Profile
// @Param id path int true "Profile ID"
// @Success 204 "Profile deleted successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid profile ID"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/dep/profiles/{id} [delete]
func (h *depProfileHandlerImpl) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid profile ID"))
		return
	}

	depName := c.Query("dep_name")
	if depName == "" {
		depName = h.depName
	}

	if err := h.service.Delete(c.Request.Context(), depName, uint(id)); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.NoContent(c)
}

// SetAsAssigner godoc
// @Summary Set profile as default assigner
// @Description Set a DEP profile as the default profile for new device enrollments
// @Tags DEP Profile
// @Param id path int true "Profile ID"
// @Success 200 {object} response.APIResponse[any] "Assigner profile set successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid profile ID or profile not registered"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "Profile not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/dep/profiles/{id}/set-assigner [post]
func (h *depProfileHandlerImpl) SetAsAssigner(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid profile ID"))
		return
	}

	depName := c.Query("dep_name")
	if depName == "" {
		depName = h.depName
	}

	if err := h.service.SetAsAssigner(c.Request.Context(), depName, uint(id)); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK[any](c, nil, "Assigner profile set successfully")
}

// GetAssigner godoc
// @Summary Get current assigner profile
// @Description Get the DEP profile currently set as the default for new enrollments
// @Tags DEP Profile
// @Produce json
// @Success 200 {object} response.APIResponse[dto.DEPProfileResponse] "Assigner profile retrieved successfully"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 404 {object} response.APIResponse[any] "No assigner profile configured"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/dep/profiles/assigner [get]
func (h *depProfileHandlerImpl) GetAssigner(c *gin.Context) {
	depName := c.Query("dep_name")
	if depName == "" {
		depName = h.depName
	}

	profile, err := h.service.GetAssigner(c.Request.Context(), depName)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, mapEntToDEPProfileResponse(profile), "Assigner profile retrieved successfully")
}

// mapEntToDEPProfileResponse converts an ent.DepProfile to dto.DEPProfileResponse
func mapEntToDEPProfileResponse(p *ent.DepProfile) *dto.DEPProfileResponse {
	return &dto.DEPProfileResponse{
		ProfileUUID:               p.ProfileUUID,
		Name:                      p.ProfileName,
		AllowPairing:              p.AllowPairing,
		AnchorCerts:               p.AnchorCerts,
		AutoAdvanceSetup:          p.AutoAdvanceSetup,
		AwaitDeviceConfigured:     p.AwaitDeviceConfigured,
		ConfigurationWebURL:       p.ConfigurationWebURL,
		Department:                p.Department,
		Devices:                   p.Devices,
		DoNotUseProfileFromBackup: p.DoNotUseProfileFromBackup,
		IsReturnToService:         p.IsReturnToService,
		IsMandatory:               p.IsMandatory,
		IsMDMRemovable:            p.IsMdmRemovable,
		IsMultiUser:               p.IsMultiUser,
		IsSupervised:              p.IsSupervised,
		Language:                  p.Language,
		OrgMagic:                  p.OrgMagic,
		Region:                    p.Region,
		SkipSetupItems:            p.SkipSetupItems,
		SupervisingHostCerts:      p.SupervisingHostCerts,
		SupportEmailAddress:       p.SupportEmailAddress,
		SupportPhoneNumber:        p.SupportPhoneNumber,
		URL:                       p.URL,
		ProfileData:               p.ProfileData,
	}
}
