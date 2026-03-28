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
	Refresh(c *gin.Context)
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
// @Description Authenticate user with username and password to receive JWT tokens.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body dto.LoginRequest true "Login credentials"
// @Success 200 {object} response.APIResponse[dto.LoginResponse] "Login successful"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data"
// @Failure 401 {object} response.APIResponse[any] "Invalid username or password"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Router /api/v1/auth/login [post]
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

// Refresh godoc
// @Summary Refresh access token
// @Description Refresh existing access token using a valid refresh token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param refresh body dto.TokenRefreshRequest true "Refresh token"
// @Success 200 {object} response.APIResponse[dto.LoginResponse] "Tokens refreshed"
// @Failure 400 {object} response.APIResponse[any] "Invalid request data"
// @Failure 401 {object} response.APIResponse[any] "Invalid or expired refresh token"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Router /api/v1/auth/refresh [post]
func (h *authHandlerImpl) Refresh(c *gin.Context) {
	var req dto.TokenRefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Dữ liệu không hợp lệ"))
		return
	}

	loginResp, err := h.authService.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, loginResp, "Làm mới token thành công")
}

// Logout godoc
// @Summary User logout
// @Description Invalidate the current session and blacklists the token.
// @Tags Authentication
// @Produce json
// @Success 200 {object} response.APIResponse[any] "Successfully logged out"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/auth/logout [post]
func (h *authHandlerImpl) Logout(c *gin.Context) {
	// Extract token from Authorization header
	token := middleware.GetToken(c)
	if token == "" {
		response.WriteErrorResponse(c, apperror.ErrUnauthorized)
		return
	}

	if err := h.authService.Logout(c.Request.Context(), token); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK[any](c, nil, "Đăng xuất thành công")
}

// GetMe godoc
// @Summary Get current user info
// @Description Returns the profile and permissions of the currently authenticated user.
// @Tags Authentication
// @Produce json
// @Success 200 {object} response.APIResponse[dto.UserResponse] "User information"
// @Failure 401 {object} response.APIResponse[any] "Unauthorized"
// @Failure 500 {object} response.APIResponse[any] "Internal server error"
// @Security BearerAuth
// @Router /api/v1/auth/me [get]
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
