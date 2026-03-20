package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/interface/api/middleware"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/response"
)

// AuthHandler interface
type AuthHandler interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	GetMe(c *gin.Context)
}

type authHandlerImpl struct {
	authService  service.AuthService
	userService  service.UserService
	authzService service.AuthorizationService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService, userService service.UserService, authzService service.AuthorizationService) AuthHandler {
	return &authHandlerImpl{
		authService:  authService,
		userService:  userService,
		authzService: authzService,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body dto.LoginRequest true "Login credentials"
// @Success 200 {object} response.APIResponse[dto.LoginResponse]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Router /auth/login [post]
func (h *authHandlerImpl) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Dữ liệu không hợp lệ"))
		return
	}

	loginResp, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, loginResp, "Đăng nhập thành công")
}

// Logout godoc
// @Summary User logout
// @Description Invalidate the current session
// @Tags Authentication
// @Produce json
// @Success 200 {object} response.APIResponse[any]
// @Router /auth/logout [post]
func (h *authHandlerImpl) Logout(c *gin.Context) {
	if err := h.authService.Logout(c.Request.Context()); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK[any](c, nil, "Đăng xuất thành công")
}

// GetMe godoc
// @Summary Get current user info
// @Description Return the profile of the currently authenticated user
// @Tags Authentication
// @Produce json
// @Success 200 {object} response.APIResponse[dto.UserResponse]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /auth/me [get]
func (h *authHandlerImpl) GetMe(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, h.toAuthUserResponse(user), "")
}

func (h *authHandlerImpl) toAuthUserResponse(user *ent.User) dto.UserResponse {
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
