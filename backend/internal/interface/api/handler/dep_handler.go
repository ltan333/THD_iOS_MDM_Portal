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

func (h *depHandler) DefineProfile(c *gin.Context) {
	// Logic to call nanoMDM /proxy/mdm-dep-server/profile
	response.OK(c, gin.H{"profile_uuid": "MOCK-UUID-12345"}, "Profile defined")
}

func (h *depHandler) GetProfile(c *gin.Context) {
	// Logic to call nanoMDM /proxy/mdm-dep-server/profile?profile_uuid=xxx
	response.OK(c, gin.H{"profile_name": "Default Profile"}, "Profile retrieved")
}

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
