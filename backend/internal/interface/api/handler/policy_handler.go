package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/response"
)

// PolicyHandler interface
type PolicyHandler interface {
	// Policies
	ListPolicies(c *gin.Context)
	AddPolicy(c *gin.Context)
	RemovePolicy(c *gin.Context)
	GetPoliciesForRole(c *gin.Context)

	// Role hierarchy
	ListRoles(c *gin.Context)
	AddRole(c *gin.Context)
	RemoveRole(c *gin.Context)
}

type policyHandlerImpl struct {
	authzService service.AuthorizationService
}

// NewPolicyHandler creates a new policy handler
func NewPolicyHandler(authzService service.AuthorizationService) PolicyHandler {
	return &policyHandlerImpl{authzService: authzService}
}

// ListPolicies godoc
// @Summary List all policies
// @Description Get a list of all Casbin policies (role, path, method)
// @Tags Authorization
// @Produce json
// @Success 200 {object} response.APIResponse[[]service.PolicyRule]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /api/v1/policies [get]
func (h *policyHandlerImpl) ListPolicies(c *gin.Context) {
	policies, err := h.authzService.GetAllPolicies()
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithMessage("Không thể lấy danh sách policy").WithError(err))
		return
	}
	response.OK(c, policies, "")
}

// AddPolicy godoc
// @Summary Add a policy
// @Description Add a new Casbin policy
// @Tags Authorization
// @Accept json
// @Produce json
// @Param policy body service.PolicyRule true "Policy rule details"
// @Success 201 {object} response.APIResponse[service.PolicyRule]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 409 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /api/v1/policies [post]
func (h *policyHandlerImpl) AddPolicy(c *gin.Context) {
	var req service.PolicyRule
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Dữ liệu không hợp lệ"))
		return
	}

	if req.Role == "" || req.Path == "" || req.Method == "" {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Role, path và method không được để trống"))
		return
	}

	added, err := h.authzService.AddPolicy(req)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithMessage("Không thể thêm policy").WithError(err))
		return
	}

	if !added {
		response.WriteErrorResponse(c, apperror.ErrConflict.WithMessage("Policy đã tồn tại"))
		return
	}

	response.Created(c, req, "Thêm policy thành công")
}

// RemovePolicy godoc
// @Summary Remove a policy
// @Description Remove an existing Casbin policy
// @Tags Authorization
// @Accept json
// @Produce json
// @Param policy body service.PolicyRule true "Policy rule details to remove"
// @Success 200 {object} response.APIResponse[any]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 404 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /api/v1/policies [delete]
func (h *policyHandlerImpl) RemovePolicy(c *gin.Context) {
	var req service.PolicyRule
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Dữ liệu không hợp lệ"))
		return
	}

	removed, err := h.authzService.RemovePolicy(req)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithMessage("Không thể xóa policy").WithError(err))
		return
	}

	if !removed {
		response.WriteErrorResponse(c, apperror.ErrNotFound.WithMessage("Không tìm thấy policy"))
		return
	}

	response.OK[any](c, nil, "Xóa policy thành công")
}

// GetPoliciesForRole godoc
// @Summary Get policies for a role
// @Description Get all Casbin policies associated with a specific role
// @Tags Authorization
// @Produce json
// @Param role path string true "Role name"
// @Success 200 {object} response.APIResponse[[]service.PolicyRule]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /api/v1/policies/role/{role} [get]
func (h *policyHandlerImpl) GetPoliciesForRole(c *gin.Context) {
	role := c.Param("role")
	if role == "" {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Role không được để trống"))
		return
	}

	policies, err := h.authzService.GetPermissionsForRole(role)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithMessage("Không thể lấy permissions").WithError(err))
		return
	}

	response.OK(c, policies, "")
}

// ListRoles godoc
// @Summary List all role links
// @Description Get a list of all role hierarchy links (child, parent)
// @Tags Authorization
// @Produce json
// @Success 200 {object} response.APIResponse[[]service.RoleLink]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /api/v1/roles [get]
func (h *policyHandlerImpl) ListRoles(c *gin.Context) {
	roles, err := h.authzService.GetAllRoles()
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithMessage("Không thể lấy danh sách roles").WithError(err))
		return
	}
	response.OK(c, roles, "")
}

// AddRole godoc
// @Summary Add a role link
// @Description Add a new role hierarchy link
// @Tags Authorization
// @Accept json
// @Produce json
// @Param role body service.RoleLink true "Role link details"
// @Success 201 {object} response.APIResponse[service.RoleLink]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 409 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /api/v1/roles [post]
func (h *policyHandlerImpl) AddRole(c *gin.Context) {
	var req service.RoleLink
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Dữ liệu không hợp lệ"))
		return
	}

	if req.Child == "" || req.Parent == "" {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Child và parent role không được để trống"))
		return
	}

	added, err := h.authzService.AddRoleLink(req.Child, req.Parent)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithMessage("Không thể thêm role").WithError(err))
		return
	}

	if !added {
		response.WriteErrorResponse(c, apperror.ErrConflict.WithMessage("Role link đã tồn tại"))
		return
	}

	response.Created(c, req, "Thêm role link thành công")
}

// RemoveRole godoc
// @Summary Remove a role link
// @Description Remove an existing role hierarchy link
// @Tags Authorization
// @Accept json
// @Produce json
// @Param role body service.RoleLink true "Role link details to remove"
// @Success 200 {object} response.APIResponse[any]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 404 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /api/v1/roles [delete]
func (h *policyHandlerImpl) RemoveRole(c *gin.Context) {
	var req service.RoleLink
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Dữ liệu không hợp lệ"))
		return
	}

	removed, err := h.authzService.RemoveRoleLink(req.Child, req.Parent)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithMessage("Không thể xóa role").WithError(err))
		return
	}

	if !removed {
		response.WriteErrorResponse(c, apperror.ErrNotFound.WithMessage("Không tìm thấy role link"))
		return
	}

	response.OK[any](c, nil, "Xóa role link thành công")
}
