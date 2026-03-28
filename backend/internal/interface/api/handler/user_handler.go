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

var userAllowedFields = map[string]bool{
	"id":         true,
	"username":   true,
	"email":      true,
	"role":       true,
	"status":     true,
	"created_at": true,
	"search":     true,
}

// UserHandler interface
type UserHandler interface {
	List(c *gin.Context)
	GetByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type userHandlerImpl struct {
	userService  service.UserService
	authzService service.AuthorizationService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserService, authzService service.AuthorizationService) UserHandler {
	return &userHandlerImpl{
		userService:  userService,
		authzService: authzService,
	}
}

// List godoc
// @Summary List users
// @Description Get a paginated list of users with filtering and sorting capabilities.
// @Tags Users
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 20)"
// @Param search query string false "Search by username or email"
// @Param role query string false "Filter by role"
// @Param status query string false "Filter by status"
// @Param sort query string false "Sort by field (e.g., id, username, created_at)"
// @Param order query string false "Sort order (asc, desc)"
// @Success 200 {object} response.APIResponse[dto.ListResponse[dto.UserResponse]] "List of users"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 403 {object} response.APIResponse[any] "Forbidden - Insufficient permissions"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/users [get]
func (h *userHandlerImpl) List(c *gin.Context) {
	params := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	offset, limit := query.GetPagination(params, 20)
	opts := query.ParseQueryParams(params, userAllowedFields)

	users, total, err := h.userService.List(c.Request.Context(), offset, limit, opts)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	items := make([]dto.UserResponse, len(users))
	for i, u := range users {
		items[i] = h.toUserResponse(u)
	}

	page := (offset / limit) + 1
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response.OK(c, dto.ListResponse[dto.UserResponse]{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, "")
}

// GetByID godoc
// @Summary Get user by ID
// @Description Fetch detailed information for a single user by their system ID.
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.APIResponse[dto.UserResponse] "User details"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID format"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 403 {object} response.APIResponse[any] "Forbidden - Insufficient permissions"
// @Failure 404 {object} response.APIResponse[any] "User not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/users/{id} [get]
func (h *userHandlerImpl) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("ID không hợp lệ"))
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, h.toUserResponse(user), "")
}

// Create godoc
// @Summary Create user
// @Description Register a new user with specified role and credentials.
// @Tags Users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "New user details"
// @Success 201 {object} response.APIResponse[dto.UserResponse] "User created successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data or validation error"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 403 {object} response.APIResponse[any] "Forbidden - Insufficient permissions"
// @Failure 409 {object} response.APIResponse[any] "Username or email already exists"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/users [post]
func (h *userHandlerImpl) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Dữ liệu không hợp lệ"))
		return
	}

	user, err := h.userService.Create(c.Request.Context(), service.CreateUserCommand{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.Created(c, h.toUserResponse(user), "Tạo người dùng thành công")
}

// Update godoc
// @Summary Update user
// @Description Modify an existing user's profile, including email, role, and status.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body dto.UpdateUserRequest true "Updated user details"
// @Success 200 {object} response.APIResponse[dto.UserResponse] "User updated successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID or request data"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 403 {object} response.APIResponse[any] "Forbidden - Insufficient permissions"
// @Failure 404 {object} response.APIResponse[any] "User not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/users/{id} [put]
func (h *userHandlerImpl) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("ID không hợp lệ"))
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Dữ liệu không hợp lệ"))
		return
	}

	user, err := h.userService.Update(c.Request.Context(), service.UpdateUserCommand{
		ID:       uint(id),
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
		Status:   req.Status,
	})
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, h.toUserResponse(user), "Cập nhật thành công")
}

// Delete godoc
// @Summary Delete user
// @Description Remove a user from the system permanently (or mark as deleted).
// @Tags Users
// @Param id path int true "User ID"
// @Success 204 "User deleted successfully"
// @Failure 400 {object} response.APIResponse[any] "Invalid ID format"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 403 {object} response.APIResponse[any] "Forbidden - Insufficient permissions"
// @Failure 404 {object} response.APIResponse[any] "User not found"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/users/{id} [delete]
func (h *userHandlerImpl) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("ID không hợp lệ"))
		return
	}

	if err := h.userService.Delete(c.Request.Context(), uint(id)); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.NoContent(c)
}

func (h *userHandlerImpl) toUserResponse(user *ent.User) dto.UserResponse {
	// Fetch role from Casbin
	role := "USER"
	roles, err := h.authzService.GetRolesForUser(user.ID)
	if err == nil && len(roles) > 0 {
		role = roles[0]
	}

	resp := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	if user.DeletedAt != nil {
		resp.DeletedAt = user.DeletedAt
	}
	return resp
}
