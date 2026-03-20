package handler

import (
	"context"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/deptoken"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
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
	ListProfiles(c *gin.Context)
	DisownDevice(c *gin.Context)

	// New methods from apidog
	ListNames(c *gin.Context)
	GetConfig(c *gin.Context)
	GetAssigner(c *gin.Context)
	SetAssigner(c *gin.Context)
	GetAccount(c *gin.Context)
	GetDevices(c *gin.Context)
	GetTokens(c *gin.Context)
}

type depHandler struct {
	client       *ent.Client
	authzService service.AuthorizationService
	mdmService   service.NanoMDMService
}

func NewDEPHandler(client *ent.Client, authzService service.AuthorizationService, mdmService service.NanoMDMService) DEPHandler {
	return &depHandler{
		client:       client,
		authzService: authzService,
		mdmService:   mdmService,
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
// @Router /v1/dep/token/{name} [put]
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

	// Read file for NanoMDM upload
	tokenData, err := os.ReadFile(dst)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithError(err))
		return
	}

	// Upload to NanoMDM
	_, err = h.mdmService.UploadDEPToken(c.Request.Context(), name, tokenData)
	if err != nil {
		response.WriteErrorResponse(c, err)
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
	}, "Token saved and uploaded to MDM successfully")
}

// GetToken godoc
// @Summary Get DEP token pem
// @Description Get PEM certificate for ABM for a specific name
// @Tags DEP
// @Produce json
// @Param name path string true "Token name"
// @Success 200 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 404 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/dep/token/{name} [get]
func (h *depHandler) GetToken(c *gin.Context) {
	name := c.Param("name")
	// For compat with apidog, we might want to return the PEM info from NanoMDM instead of local DB
	// But let's check local DB first
	token, _ := h.client.DEPToken.
		Query().
		Where(deptoken.NameEQ(name)).
		Only(context.Background())

	if token == nil {
		response.WriteErrorResponse(c, apperror.ErrNotFound.WithMessage("Token not found"))
		return
	}

	// Maybe call mdmService.GetDEPTokens if we want exactly what apidog shows
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
// @Description Initiate a sync with Apple DEP servers to fetch new devices via proxy
// @Tags DEP
// @Produce json
// @Param name path string false "DEP name (default: 'default')"
// @Success 200 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/dep/proxy/{name}/devices/sync [post]
func (h *depHandler) SyncDevices(c *gin.Context) {
	depName := c.Param("name")
	if depName == "" {
		depName = "default"
	}

	result, err := h.mdmService.SyncDEPDevices(c.Request.Context(), depName)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, result, "Sync initiated successfully")
}

// DefineProfile godoc
// @Summary Define DEP profile
// @Description Create or update a DEP assignment profile via proxy
// @Tags DEP
// @Accept json
// @Produce json
// @Param name path string false "DEP name (default: 'default')"
// @Success 200 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/dep/proxy/{name}/profile [post]
func (h *depHandler) DefineProfile(c *gin.Context) {
	depName := c.Param("name")
	if depName == "" {
		depName = "default"
	}

	var profile interface{}
	if err := c.ShouldBindJSON(&profile); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	uuid, err := h.mdmService.DefineDEPProfile(c.Request.Context(), depName, profile)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, gin.H{"profile_uuid": uuid}, "Profile defined successfully")
}

// GetProfile godoc
// @Summary Get DEP profile
// @Description Fetch details of a defined DEP profile via proxy
// @Tags DEP
// @Produce json
// @Param name path string false "DEP name (default: 'default')"
// @Param profile_uuid query string true "Profile UUID"
// @Success 200 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/dep/proxy/{name}/profile [get]
func (h *depHandler) GetProfile(c *gin.Context) {
	depName := c.Param("name")
	if depName == "" {
		depName = "default"
	}
	uuid := c.Query("profile_uuid")

	profile, err := h.mdmService.GetDEPProfile(c.Request.Context(), depName, uuid)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, profile, "Profile retrieved successfully")
}

// ListProfiles godoc
// @Summary List Defined DEP profiles
// @Description Fetch all defined DEP profiles from NanoDEP
// @Tags DEP
// @Produce json
// @Success 200 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/dep/profiles [get]
func (h *depHandler) ListProfiles(c *gin.Context) {
	profiles, err := h.mdmService.ListDEPProfiles(c.Request.Context(), "default")
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, profiles, "Profiles retrieved successfully")
}

// DisownDevice godoc
// @Summary Disown DEP device
// @Description Remove a device from DEP management via proxy
// @Tags DEP
// @Produce json
// @Param name path string false "DEP name (default: 'default')"
// @Success 200 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 403 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/dep/proxy/{name}/devices/disown [post]
func (h *depHandler) DisownDevice(c *gin.Context) {
	depName := c.Param("name")
	if depName == "" {
		depName = "default"
	}

	// Serial numbers from request
	var req struct {
		Devices []string `json:"devices"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	result, err := h.mdmService.DisownDEPDevices(c.Request.Context(), depName, req.Devices)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, result, "Devices disowned successfully")
}

// ListNames godoc
// @Summary List DEP names
// @Description List all configured DEP names in NanoDEP
// @Tags DEP
// @Produce json
// @Security BearerAuth
// @Router /v1/dep/names [get]
func (h *depHandler) ListNames(c *gin.Context) {
	result, err := h.mdmService.ListDEPNames(c.Request.Context())
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, result, "DEP names retrieved successfully")
}

// GetConfig godoc
// @Summary Get DEP config
// @Description Get configuration for a specific DEP name
// @Tags DEP
// @Produce json
// @Param name path string true "DEP name"
// @Security BearerAuth
// @Router /v1/dep/config/{name} [get]
func (h *depHandler) GetConfig(c *gin.Context) {
	name := c.Param("name")
	result, err := h.mdmService.GetDEPConfig(c.Request.Context(), name)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, result, "Config retrieved successfully")
}

// GetAssigner godoc
// @Summary Get DEP assigner
// @Description Get assigner for a specific DEP name
// @Tags DEP
// @Produce json
// @Param name path string true "DEP name"
// @Security BearerAuth
// @Router /v1/dep/assigner/{name} [get]
func (h *depHandler) GetAssigner(c *gin.Context) {
	name := c.Param("name")
	result, err := h.mdmService.GetDEPAssigner(c.Request.Context(), name)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, result, "Assigner retrieved successfully")
}

// SetAssigner godoc
// @Summary Set DEP assigner
// @Description Set or update assigner for a DEP name
// @Tags DEP
// @Accept json
// @Produce json
// @Param name path string true "DEP name"
// @Security BearerAuth
// @Router /v1/dep/assigner/{name} [put]
func (h *depHandler) SetAssigner(c *gin.Context) {
	name := c.Param("name")
	var req interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	result, err := h.mdmService.SetDEPAssigner(c.Request.Context(), name, req)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, result, "Assigner updated successfully")
}

// GetAccount godoc
// @Summary Get DEP account
// @Description Fetch Apple DEP account info via proxy
// @Tags DEP
// @Produce json
// @Param name path string true "DEP name"
// @Security BearerAuth
// @Router /v1/dep/proxy/{name}/account [get]
func (h *depHandler) GetAccount(c *gin.Context) {
	name := c.Param("name")
	result, err := h.mdmService.GetDEPAccount(c.Request.Context(), name)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, result, "Account info retrieved successfully")
}

// GetDevices godoc
// @Summary Get DEP devices
// @Description Fetch device list or details from Apple via proxy
// @Tags DEP
// @Produce json
// @Param name path string true "DEP name"
// @Param cursor query string false "Pagination cursor"
// @Security BearerAuth
// @Router /v1/dep/proxy/{name}/devices [post]
func (h *depHandler) GetDevices(c *gin.Context) {
	name := c.Param("name")
	cursor := c.Query("cursor")
	result, err := h.mdmService.GetDEPDevices(c.Request.Context(), name, cursor)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, result, "Devices retrieved successfully")
}

// GetTokens godoc
// @Summary Get DEP tokens info
// @Description Get current token information from NanoDEP for a name
// @Tags DEP
// @Produce json
// @Param name path string true "DEP name"
// @Security BearerAuth
// @Router /v1/dep/tokens/{name} [get]
func (h *depHandler) GetTokens(c *gin.Context) {
	name := c.Param("name")
	result, err := h.mdmService.GetDEPTokens(c.Request.Context(), name)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, result, "Tokens info retrieved successfully")
}
