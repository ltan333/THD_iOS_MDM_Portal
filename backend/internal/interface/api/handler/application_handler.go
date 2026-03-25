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

var applicationAllowedFields = map[string]bool{
	"platform": true,
	"type":     true,
	"search":   true,
}

type ApplicationHandler interface {
	List(c *gin.Context)
	GetByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)

	// App Versions
	ListVersions(c *gin.Context)
	CreateVersion(c *gin.Context)
	DeleteVersion(c *gin.Context)

	// App Deployments
	Deploy(c *gin.Context)
	ListDeployments(c *gin.Context)
}

type applicationHandlerImpl struct {
	appService service.ApplicationService
}

func NewApplicationHandler(appService service.ApplicationService) ApplicationHandler {
	return &applicationHandlerImpl{appService: appService}
}

// List godoc
// @Summary List applications
// @Description Fetch tracked applications with pagination
// @Tags applications
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse[any]
// @Router /applications [get]
func (h *applicationHandlerImpl) List(c *gin.Context) {
	params := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	offset, limit := query.GetPagination(params, 20)
	opts := query.ParseQueryParams(params, applicationAllowedFields)

	apps, total, err := h.appService.List(c.Request.Context(), offset, limit, opts)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	res := make([]dto.ApplicationResponse, 0, len(apps))
	for _, app := range apps {
		res = append(res, mapApplicationToResponse(app))
	}

	page := (offset / limit) + 1
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response.OK(c, dto.ListResponse[dto.ApplicationResponse]{
		Items:      res,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, "")
}

// GetByID godoc
// @Summary Get app block by ID
// @Description Fetch details of a specific app
// @Tags applications
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse[any]
// @Router /applications/{id} [get]
func (h *applicationHandlerImpl) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	app, err := h.appService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, mapApplicationToResponse(app), "")
}

// Create godoc
// @Summary Add a new application tracking record
// @Description Create application
// @Tags applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} response.APIResponse[any]
// @Router /applications [post]
func (h *applicationHandlerImpl) Create(c *gin.Context) {
	var req dto.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid data"))
		return
	}

	app, err := h.appService.Create(c.Request.Context(), service.CreateApplicationCommand{
		Name:        req.Name,
		BundleID:    req.BundleID,
		Platform:    req.Platform,
		Type:        req.Type,
		Description: req.Description,
		IconURL:     req.IconURL,
	})
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.Created(c, mapApplicationToResponse(app), "Application created successfully")
}

// Update godoc
// @Summary Update application metadata
// @Description Update application
// @Tags applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse[any]
// @Router /applications/{id} [put]
func (h *applicationHandlerImpl) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	var req dto.UpdateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid data"))
		return
	}

	app, err := h.appService.Update(c.Request.Context(), service.UpdateApplicationCommand{
		ID:          uint(id),
		Name:        req.Name,
		Platform:    req.Platform,
		Type:        req.Type,
		Description: req.Description,
		IconURL:     req.IconURL,
	})
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, mapApplicationToResponse(app), "Application updated successfully")
}

// Delete godoc
// @Summary Delete application
// @Tags applications
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse[any]
// @Router /applications/{id} [delete]
func (h *applicationHandlerImpl) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid ID"))
		return
	}

	if err := h.appService.Delete(c.Request.Context(), uint(id)); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK[any](c, nil, "Application deleted successfully")
}

// ListVersions godoc
// @Summary List application versions
// @Tags applications
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse[any]
// @Router /applications/{id}/versions [get]
func (h *applicationHandlerImpl) ListVersions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid application ID"))
		return
	}

	versions, err := h.appService.ListVersions(c.Request.Context(), uint(id))
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	res := make([]dto.AppVersionResponse, 0, len(versions))
	for _, v := range versions {
		res = append(res, mapAppVersionToResponse(v))
	}

	response.OK(c, res, "")
}

// CreateVersion godoc
// @Summary Upload and create a new version of an app
// @Tags applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} response.APIResponse[any]
// @Router /applications/{id}/versions [post]
func (h *applicationHandlerImpl) CreateVersion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid application ID"))
		return
	}

	var req dto.CreateAppVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid data"))
		return
	}
	req.ApplicationID = uint(id)

	v, err := h.appService.CreateVersion(c.Request.Context(), service.CreateAppVersionCommand{
		ApplicationID:    req.ApplicationID,
		Version:          req.Version,
		BuildNumber:      req.BuildNumber,
		MinimumOSVersion: req.MinimumOSVersion,
		FileURL:          req.FileURL,
		Size:             req.Size,
		Metadata:         req.Metadata,
	})
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.Created(c, mapAppVersionToResponse(v), "Version created successfully")
}

// DeleteVersion godoc
// @Summary Delete app version
// @Tags applications
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse[any]
// @Router /applications/{id}/versions/{versionId} [delete]
func (h *applicationHandlerImpl) DeleteVersion(c *gin.Context) {
	versionIdStr := c.Param("versionId")
	versionId, err := strconv.ParseUint(versionIdStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid version ID"))
		return
	}

	if err := h.appService.DeleteVersion(c.Request.Context(), uint(versionId)); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK[any](c, nil, "Version deleted successfully")
}

// Deploy godoc
// @Summary Push app to devices
// @Description Command the MDM core to push app installations
// @Tags applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse[any]
// @Router /applications/deployments [post]
func (h *applicationHandlerImpl) Deploy(c *gin.Context) {
	var req dto.CreateAppDeploymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Invalid data"))
		return
	}

	deployment, err := h.appService.Deploy(c.Request.Context(), service.CreateAppDeploymentCommand{
		AppVersionID: req.AppVersionID,
		TargetType:   req.TargetType,
		TargetID:     req.TargetID,
	})
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, dto.AppDeploymentResponse{
		ID:           deployment.ID,
		AppVersionID: deployment.AppVersionID,
		TargetType:   string(deployment.TargetType),
		TargetID:     deployment.TargetID,
		Status:       string(deployment.Status),
		ErrorMessage: deployment.ErrorMessage,
		InstalledAt:  deployment.InstalledAt,
		CreatedAt:    deployment.CreatedAt,
		UpdatedAt:    deployment.UpdatedAt,
	}, "Deployment initiated successfully")
}

// ListDeployments godoc
// @Summary List deployment status
// @Description Review installation statuses of pushed tools
// @Tags applications
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse[any]
// @Router /applications/{id}/versions/{versionId}/deployments [get]
func (h *applicationHandlerImpl) ListDeployments(c *gin.Context) {
	versionIdStr := c.Param("versionId")
	versionId, err := strconv.ParseUint(versionIdStr, 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Invalid version ID"))
		return
	}

	deployments, err := h.appService.ListDeployments(c.Request.Context(), uint(versionId))
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	res := make([]dto.AppDeploymentResponse, 0, len(deployments))
	for _, d := range deployments {
		res = append(res, dto.AppDeploymentResponse{
			ID:           d.ID,
			AppVersionID: d.AppVersionID,
			TargetType:   string(d.TargetType),
			TargetID:     d.TargetID,
			Status:       string(d.Status),
			ErrorMessage: d.ErrorMessage,
			InstalledAt:  d.InstalledAt,
			CreatedAt:    d.CreatedAt,
			UpdatedAt:    d.UpdatedAt,
		})
	}

	response.OK(c, res, "")
}


func mapApplicationToResponse(a *ent.Application) dto.ApplicationResponse {
	res := dto.ApplicationResponse{
		ID:          a.ID,
		Name:        a.Name,
		BundleID:    a.BundleID,
		Platform:    string(a.Platform),
		Type:        string(a.Type),
		Description: a.Description,
		IconURL:     a.IconURL,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
	
	if a.Edges.Versions != nil {
		res.Versions = make([]dto.AppVersionResponse, 0, len(a.Edges.Versions))
		for _, v := range a.Edges.Versions {
			res.Versions = append(res.Versions, mapAppVersionToResponse(v))
		}
	}
	return res
}

func mapAppVersionToResponse(v *ent.AppVersion) dto.AppVersionResponse {
	return dto.AppVersionResponse{
		ID:               v.ID,
		ApplicationID:    v.ApplicationID,
		Version:          v.Version,
		BuildNumber:      v.BuildNumber,
		MinimumOSVersion: v.MinimumOsVersion,
		FileURL:          v.FileURL,
		Size:             v.Size,
		Metadata:         v.Metadata,
		CreatedAt:        v.CreatedAt,
		UpdatedAt:        v.UpdatedAt,
	}
}
