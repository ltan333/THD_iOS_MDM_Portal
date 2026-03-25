package handler

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/deptoken"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
	"github.com/thienel/go-backend-template/pkg/response"
)

type DEPHandler interface {
	PutTokenPKI(c *gin.Context)
	GetTokenPKI(c *gin.Context)
	GetToken(c *gin.Context)
	SyncDevices(c *gin.Context)
	DefineProfile(c *gin.Context)
	GetProfile(c *gin.Context)
	ListProfiles(c *gin.Context)
	DisownDevice(c *gin.Context)

	// New methods from apidog / NanoDEP spec
	ListNames(c *gin.Context)
	GetConfig(c *gin.Context)
	PutConfig(c *gin.Context)
	GetAssigner(c *gin.Context)
	SetAssigner(c *gin.Context)
	GetAccount(c *gin.Context)
	GetDevices(c *gin.Context)
	GetTokens(c *gin.Context)
	UpdateTokens(c *gin.Context)
	GetMAIDJWT(c *gin.Context)
	GetBypassCode(c *gin.Context)
	GetVersion(c *gin.Context)
}

type depHandler struct {
	client            *ent.Client
	authzService      service.AuthorizationService
	mdmService        service.NanoMDMService
	depProfileService service.DepProfileService
}

func NewDEPHandler(
	client *ent.Client,
	authzService service.AuthorizationService,
	mdmService service.NanoMDMService,
	depProfileService service.DepProfileService,
) DEPHandler {
	return &depHandler{
		client:            client,
		authzService:      authzService,
		mdmService:        mdmService,
		depProfileService: depProfileService,
	}
}

// PutTokenPKI godoc
// @Summary Upload and decrypt DEP OAuth1 tokens
// @Description Decrypt the OAuth1 tokens from the Apple ABM/ASM/BE portal and store them.
// @Tags DEP
// @Accept application/pkcs7-mime
// @Produce json
// @Param name path string true "Name of DEP server instance"
// @Param force query string false "Bypass the Consumer Key mismatch check (1 to enable)"
// @Param token body string true "Contents of the .p7m file from Apple"
// @Success 200 {object} response.APIResponse[dto.OAuth1Tokens]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/dep/tokenpki/{name} [put]
func (h *depHandler) PutTokenPKI(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Tham số name là bắt buộc"))
		return
	}

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	result, err := h.mdmService.UploadDEPToken(c.Request.Context(), name, data)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, result, "Token decrypted and saved successfully")
}

// GetTokenPKI godoc
// @Summary Generate and retrieve DEP token PKI certificate
// @Description Generate and store a new X.509 certificate and RSA private key for exchanging encrypted DEP OAuth1 tokens.
// @Tags DEP
// @Produce application/x-pem-file
// @Param name path string true "Name of DEP server instance"
// @Param cn query string false "Common Name"
// @Param validity_days query int false "Validity days"
// @Success 200 {string} string "X.509 certificate PEM"
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/dep/tokenpki/{name} [get]
func (h *depHandler) GetTokenPKI(c *gin.Context) {
	name := c.Param("name")
	cn := c.Query("cn")
	valDaysStr := c.Query("validity_days")
	var valDays int
	if valDaysStr != "" {
		_, _ = fmt.Sscanf(valDaysStr, "%d", &valDays)
	}

	cert, contentDisp, err := h.mdmService.GetDEPTokenPKI(c.Request.Context(), name, cn, valDays)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	if contentDisp != "" {
		c.Header("Content-Disposition", contentDisp)
	}
	c.Data(200, "application/x-pem-file", cert)
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
	if name == "" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Tham số name là bắt buộc"))
		return
	}
	// For compat with apidog, we might want to return the PEM info from NanoMDM instead of local DB
	// But let's check local DB first
	token := h.client.DEPToken.
		Query().
		Where(deptoken.IDEQ(name)).
		OnlyX(c.Request.Context())

	if token == nil {
		response.WriteErrorResponse(c, apperror.ErrNotFound.WithMessage("Token not found"))
		return
	}

	// Maybe call mdmService.GetDEPTokens if we want exactly what apidog shows
	response.OK(c, dto.DEPTokenResponse{
		ID:                token.ID,
		TokenpkiCertPem:   token.TokenpkiCertPem,
		TokenpkiKey_pem:   token.TokenpkiKeyPem,
		AccessTokenExpiry: &token.AccessTokenExpiry,
		CreatedAt:         token.CreatedAt,
		UpdatedAt:         token.UpdatedAt,
	}, "Token retrieved successfully")
}

// SyncDevices godoc
// @Summary Sync DEP devices
// @Description Initiate a sync with Apple DEP servers to fetch new devices via proxy
// @Tags DEP
// @Produce json
// @Param name path string false "DEP name (default: 'default')"
// @Param cursor query string false "Sync cursor"
// @Security BearerAuth
// @Router /v1/dep/proxy/{name}/devices/sync [post]
func (h *depHandler) SyncDevices(c *gin.Context) {
	depName := c.Param("name")
	if depName == "" {
		depName = "default"
	}

	cursor := c.Query("cursor")
	result, err := h.mdmService.SyncDEPDevices(c.Request.Context(), depName, cursor)
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
// @Param request body dto.DEPProfileRequest true "Profile details"
// @Success 200 {object} response.APIResponse[dto.DEPProfileResponse]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/dep/proxy/{name}/profile [post]
func (h *depHandler) DefineProfile(c *gin.Context) {
	depName := c.Param("name")
	if depName == "" {
		depName = "default"
	}

	var req dto.DEPProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	resp, err := h.depProfileService.DefineProfile(c.Request.Context(), depName, &req)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, resp, "Profile defined and saved locally successfully")
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

	resp, err := h.depProfileService.GetProfile(c.Request.Context(), depName, uuid)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, resp, "Profile retrieved successfully")
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
	// For simplicity, we use hardcoded offset/limit or get from query
	profiles, _, err := h.depProfileService.ListProfiles(c.Request.Context(), 0, 100, query.QueryOptions{})
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
// @Accept json
// @Produce json
// @Param name path string true "DEP name"
// @Param request body dto.DEPDevicesRequest true "Devices list to disown"
// @Security BearerAuth
// @Router /v1/dep/proxy/{name}/devices/disown [post]
func (h *depHandler) DisownDevice(c *gin.Context) {
	depName := c.Param("name")
	if depName == "" {
		depName = "default"
	}

	var req dto.DEPDevicesRequest
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
// @Summary Query DEP names
// @Description Query DEP names with optional filters and pagination.
// @Tags DEP
// @Produce json
// @Param dep_name query []string false "Filter by DEP names"
// @Param limit query int false "Limits number of results" default(100)
// @Param offset query int false "Offset results"
// @Param cursor query string false "Pagination cursor"
// @Success 200 {object} response.APIResponse[dto.DEPNamesQueryResponse]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/dep/dep_names [get]
func (h *depHandler) ListNames(c *gin.Context) {
	depNames := c.QueryArray("dep_name")
	limitStr := c.DefaultQuery("limit", "100")
	offsetStr := c.DefaultQuery("offset", "0")
	cursor := c.Query("cursor")

	var limit, offset int
	_, _ = fmt.Sscanf(limitStr, "%d", &limit)
	_, _ = fmt.Sscanf(offsetStr, "%d", &offset)

	result, err := h.mdmService.ListDEPNames(c.Request.Context(), depNames, limit, offset, cursor)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, result, "DEP names retrieved successfully")
}

// GetConfig godoc
// @Summary Return the config for the given DEP name
// @Description Return the config for the given DEP name.
// @Tags DEP
// @Produce json
// @Param name path string true "Name of DEP server instance"
// @Success 200 {object} response.APIResponse[dto.DEPConfig]
// @Failure 401 {object} response.APIResponse[any]
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

// PutConfig godoc
// @Summary Set the config for the given DEP name
// @Description Set the config for the given DEP name.
// @Tags DEP
// @Accept json
// @Produce json
// @Param name path string true "Name of DEP server instance"
// @Param request body dto.DEPConfig true "Config details"
// @Success 200 {object} response.APIResponse[dto.DEPConfig]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/dep/config/{name} [put]
func (h *depHandler) PutConfig(c *gin.Context) {
	name := c.Param("name")
	var req dto.DEPConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	result, err := h.mdmService.SetDEPConfig(c.Request.Context(), name, &req)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, result, "Config updated successfully")
}

// GetAssigner godoc
// @Summary Return the assigner profile UUID
// @Description Return the assigner profile UUID for the given DEP name.
// @Tags DEP
// @Produce json
// @Param name path string true "Name of DEP server instance"
// @Success 200 {object} response.APIResponse[dto.AssignerProfileUUID]
// @Failure 401 {object} response.APIResponse[any]
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
// @Summary Assign a profile UUID for assignment
// @Description Assign a profile UUID for assignment for the given DEP name.
// @Tags DEP
// @Produce json
// @Param name path string true "Name of DEP server instance"
// @Param profile_uuid query string true "Profile UUID to assign"
// @Success 200 {object} response.APIResponse[dto.AssignerProfileUUID]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/dep/assigner/{name} [put]
func (h *depHandler) SetAssigner(c *gin.Context) {
	name := c.Param("name")
	uuid := c.Query("profile_uuid")
	if uuid == "" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Tham số profile_uuid là bắt buộc"))
		return
	}
	result, err := h.mdmService.SetDEPAssigner(c.Request.Context(), name, uuid)
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
	if name == "" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Tham số name là bắt buộc"))
		return
	}
	result, err := h.mdmService.GetDEPAccount(c.Request.Context(), name)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, result, "Account info retrieved successfully")
}

// GetDevices godoc
// @Summary Get DEP devices
// @Description Fetch device list or details from Apple via proxy. If devices array is provided in body, fetches details for those serial numbers.
// @Tags DEP
// @Accept json
// @Produce json
// @Param name path string true "DEP name"
// @Param request body dto.DEPDevicesRequest false "Devices list"
// @Param cursor query string false "Pagination cursor"
// @Security BearerAuth
// @Router /v1/dep/proxy/{name}/devices [post]
func (h *depHandler) GetDevices(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Tham số name là bắt buộc"))
		return
	}

	var req dto.DEPDevicesRequest
	if c.Request.ContentLength > 0 {
		_ = c.ShouldBindJSON(&req)
	}

	cursor := c.Query("cursor")
	result, err := h.mdmService.GetDEPDevices(c.Request.Context(), name, req.Devices, cursor)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, result, "Devices retrieved successfully")
}

// GetTokens godoc
// @Summary Return the DEP OAuth1 tokens
// @Description Return the DEP OAuth1 tokens for the given DEP name.
// @Tags DEP
// @Produce json
// @Param name path string true "Name of DEP server instance"
// @Success 200 {object} response.APIResponse[dto.OAuth1Tokens]
// @Failure 401 {object} response.APIResponse[any]
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

// UpdateTokens godoc
// @Summary Upload and store DEP OAuth1 tokens
// @Description Upload and store DEP OAuth1 tokens for the given DEP Name.
// @Tags DEP
// @Accept json
// @Produce json
// @Param name path string true "Name of DEP server instance"
// @Param tokens body dto.OAuth1Tokens true "OAuth1 tokens"
// @Success 200 {object} response.APIResponse[dto.OAuth1Tokens]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/dep/tokens/{name} [put]
func (h *depHandler) UpdateTokens(c *gin.Context) {
	name := c.Param("name")
	var tokens dto.OAuth1Tokens
	if err := c.ShouldBindJSON(&tokens); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	result, err := h.mdmService.UpdateDEPTokens(c.Request.Context(), name, &tokens)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, result, "Tokens updated successfully")
}

// GetMAIDJWT godoc
// @Summary Generate Managed Apple ID Managed Access JWT
// @Description Generate Managed Apple ID Managed Access JWT.
// @Tags DEP
// @Produce application/jwt
// @Param name path string true "Name of DEP server instance"
// @Param server_uuid query string false "MDM server UUID"
// @Success 200 {string} string "JWT"
// @Security BearerAuth
// @Router /v1/dep/maidjwt/{name} [get]
func (h *depHandler) GetMAIDJWT(c *gin.Context) {
	name := c.Param("name")
	serverUUID := c.Query("server_uuid")

	jwt, serverUuidHeader, jtiHeader, err := h.mdmService.GetMAIDJWT(c.Request.Context(), name, serverUUID)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	c.Header("X-Server-Uuid", serverUuidHeader)
	c.Header("X-Jwt-Jti", jtiHeader)
	c.Data(200, "application/jwt", []byte(jwt))
}

// GetBypassCode godoc
// @Summary Generates or decodes an Activation Lock Bypass Code
// @Description Generates (or decodes) an Activation Lock Bypass Code and returns different forms of it.
// @Tags DEP
// @Produce json
// @Param code query string false "Hex-encoded raw form of bypass code"
// @Param raw query string false "Dash-separated human readable form"
// @Success 200 {object} response.APIResponse[dto.BypassCodeResponse]
// @Security BearerAuth
// @Router /v1/dep/bypasscode [get]
func (h *depHandler) GetBypassCode(c *gin.Context) {
	code := c.Query("code")
	raw := c.Query("raw")

	result, err := h.mdmService.GetBypassCode(c.Request.Context(), code, raw)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, result, "Bypass code generated/decoded successfully")
}

// GetVersion godoc
// @Summary Returns the running NanoDEP version
// @Description Returns the running NanoDEP version.
// @Tags DEP
// @Produce json
// @Success 200 {object} response.APIResponse[dto.NanoDEPVersionResponse]
// @Router /v1/dep/version [get]
func (h *depHandler) GetVersion(c *gin.Context) {
	resp, err := h.mdmService.GetDEPVersion(c.Request.Context())
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	response.OK(c, resp, "Version retrieved successfully")
}
