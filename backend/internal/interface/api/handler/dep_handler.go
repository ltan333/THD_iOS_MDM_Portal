package handler

import (
	"context"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/deptoken"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/interface/api/middleware"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/response"
)

type DEPHandler interface {
	PutToken(c *gin.Context)
	GetToken(c *gin.Context)
	SyncDevices(c *gin.Context)
	DefineProfile(c *gin.Context)
	GetProfile(c *gin.Context)
	DisownDevice(c *gin.Context)
}

type depHandler struct {
	client       *ent.Client
	authzService service.AuthorizationService
}

func NewDEPHandler(client *ent.Client, authzService service.AuthorizationService) DEPHandler {
	return &depHandler{
		client:       client,
		authzService: authzService,
	}
}

// PutToken godoc
// @Summary Upload DEP token
// @Description Upload or update a DEP token (.p7m file) for a specific name
// @Tags DEP
// @Accept multipart/form-data
// @Produce json
// @Param name path string true "Token name"
// @Param token formData file true "DEP Token file (.p7m)"
// @Success 200 {object} response.APIResponse[dto.DEPTokenResponse]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /dep/token/{name} [post]
func (h *depHandler) PutToken(c *gin.Context) {
	name := c.Param("name")
	file, err := c.FormFile("token")
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Token file is required"))
		return
	}

	dst := filepath.Join("storage", "certs", name+"_token.p7m")
	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithError(err))
		return
	}

	// Check if exists
	existing, _ := h.client.DEPToken.Query().Where(deptoken.NameEQ(name)).Only(context.Background())

	var token *ent.DEPToken
	if existing != nil {
		token, err = h.client.DEPToken.
			UpdateOne(existing).
			SetP7mFilePath(dst).
			Save(context.Background())
	} else {
		token, err = h.client.DEPToken.
			Create().
			SetName(name).
			SetP7mFilePath(dst).
			Save(context.Background())
	}

	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithError(err))
		return
	}

	response.OK(c, dto.DEPTokenResponse{
		ID:          token.ID,
		Name:        token.Name,
		P7mFilePath: token.P7mFilePath,
		Expiry:      token.Expiry,
		LastUsed:    token.LastUsed,
		CreatedAt:   token.CreatedAt,
		UpdatedAt:   token.UpdatedAt,
	}, "Token saved successfully")
}

// GetToken godoc
// @Summary Get DEP token
// @Description Get details of a DEP token by name
// @Tags DEP
// @Produce json
// @Param name path string true "Token name"
// @Success 200 {object} response.APIResponse[dto.DEPTokenResponse]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 404 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /dep/token/{name} [get]
func (h *depHandler) GetToken(c *gin.Context) {
	name := c.Param("name")
	token, err := h.client.DEPToken.
		Query().
		Where(deptoken.NameEQ(name)).
		Only(context.Background())

	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrNotFound.WithMessage("Token not found"))
		return
	}

	response.OK(c, dto.DEPTokenResponse{
		ID:          token.ID,
		Name:        token.Name,
		P7mFilePath: token.P7mFilePath,
		Expiry:      token.Expiry,
		LastUsed:    token.LastUsed,
		CreatedAt:   token.CreatedAt,
		UpdatedAt:   token.UpdatedAt,
	}, "Token retrieved successfully")
}

// SyncDevices godoc
// @Summary Sync DEP devices
// @Description Initiate a sync with Apple DEP servers to fetch new devices
// @Tags DEP
// @Produce json
// @Success 200 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /dep/sync [post]
func (h *depHandler) SyncDevices(c *gin.Context) {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		response.WriteErrorResponse(c, apperror.ErrUnauthorized)
		return
	}

	// Mock syncing a device: Serial SN-MOCK-123
	sn := "SN-MOCK-123"
	_, err := h.authzService.AddResourcePolicy(claims.UserID, "device:"+sn, "read")
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithError(err))
		return
	}
	h.authzService.AddResourcePolicy(claims.UserID, "device:"+sn, "write")

	response.OK(c, gin.H{"status": "sync_initiated", "synced_device": sn}, "Sync initiated and permissions granted")
}

// DefineProfile godoc
// @Summary Define DEP profile
// @Description Create or update a DEP assignment profile
// @Tags DEP
// @Accept json
// @Produce json
// @Success 200 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /dep/profile [post]
func (h *depHandler) DefineProfile(c *gin.Context) {
	// Logic to call nanoMDM /proxy/mdm-dep-server/profile
	response.OK(c, gin.H{"profile_uuid": "MOCK-UUID-12345"}, "Profile defined")
}

// GetProfile godoc
// @Summary Get DEP profile
// @Description Fetch details of a defined DEP profile
// @Tags DEP
// @Produce json
// @Success 200 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /dep/profile [get]
func (h *depHandler) GetProfile(c *gin.Context) {
	// Logic to call nanoMDM /proxy/mdm-dep-server/profile?profile_uuid=xxx
	response.OK(c, gin.H{"profile_name": "Default Profile"}, "Profile retrieved")
}

// DisownDevice godoc
// @Summary Disown DEP device
// @Description Remove a device from DEP management
// @Tags DEP
// @Produce json
// @Success 200 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 403 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /dep/disown [post]
func (h *depHandler) DisownDevice(c *gin.Context) {
	claims := middleware.GetUserClaims(c)
	if claims == nil {
		response.WriteErrorResponse(c, apperror.ErrUnauthorized)
		return
	}

	// Serial number from request (mocked for now)
	sn := "SN-MOCK-123"

	// Check permission
	allowed, err := h.authzService.AuthorizeResource(claims.UserID, "device:"+sn, "write")
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithError(err))
		return
	}

	if !allowed {
		response.WriteErrorResponse(c, apperror.ErrForbidden.WithMessage("Bạn không có quyền thao tác trên thiết bị này"))
		return
	}

	// Logic to call nanoMDM /proxy/mdm-dep-server/devices/disown
	response.OK(c, gin.H{"status": "disowned", "serial_number": sn}, "Device disowned successfully")
}
