package handler

import (
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/response"
)

type PayloadPropertyDefinitionHandler interface {
	ListPayloadTypes(c *gin.Context)
	Import(c *gin.Context)
	GetNestedSchema(c *gin.Context)
}

type payloadPropertyDefinitionHandlerImpl struct {
	service service.PayloadPropertyDefinitionService
}

func NewPayloadPropertyDefinitionHandler(service service.PayloadPropertyDefinitionService) PayloadPropertyDefinitionHandler {
	return &payloadPropertyDefinitionHandlerImpl{service: service}
}

// ListPayloadTypes godoc
// @Summary List payload types
// @Description Get all distinct payload types from payload property definitions
// @Tags Payload Property Definitions
// @Produce json
// @Success 200 {object} response.APIResponse[[]string]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /api/v1/payload-property-definitions/payload-types [get]
func (h *payloadPropertyDefinitionHandlerImpl) ListPayloadTypes(c *gin.Context) {
	types, err := h.service.ListPayloadTypes(c.Request.Context())
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, types, "")
}

// Import godoc
// @Summary Import payload property definitions from Apple JSON files
// @Description Upload one or more Apple documentation JSON files. The system auto-detects top-level payload files vs nested dictionary files and resolves nested properties recursively.
// @Tags Payload Property Definitions
// @Accept multipart/form-data
// @Produce json
// @Param files formData []file true "Apple payload JSON files (upload nhiều files; bao gồm cả nested dictionary files để resolve đệ quy)"
// @Success 200 {object} response.APIResponse[service.ImportPayloadPropertyDefinitionsResult]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /api/v1/payload-property-definitions/import [post]
func (h *payloadPropertyDefinitionHandlerImpl) Import(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Không thể parse multipart form").WithError(err))
		return
	}

	// Collect files from all common field names for compatibility:
	// "files"   – multiple files (recommended)
	// "files[]" – some clients append [] for array fields
	// "file"    – legacy single-file field name
	var allFiles []*multipart.FileHeader
	for _, fieldName := range []string{"files", "files[]", "file"} {
		if fhs, ok := form.File[fieldName]; ok {
			allFiles = append(allFiles, fhs...)
		}
	}

	if len(allFiles) == 0 {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage(
			"Thiếu file JSON để import. Vui lòng upload với field name: 'files' (hỗ trợ nhiều files)"))
		return
	}

	// Build in-memory fileMap from all uploaded files
	fileMap := make(map[string][]byte, len(allFiles))
	for _, fh := range allFiles {
		src, err := fh.Open()
		if err != nil {
			response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Không thể mở file: "+fh.Filename).WithError(err))
			return
		}
		data, err := io.ReadAll(src)
		_ = src.Close()
		if err != nil {
			response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Không thể đọc file: "+fh.Filename).WithError(err))
			return
		}
		// Normalize key: lowercase basename (strips any directory prefix from browser paths)
		key := strings.ToLower(filepath.Base(fh.Filename))
		fileMap[key] = data
	}

	result, err := h.service.ImportFromAppleJSONFiles(c.Request.Context(), fileMap)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, result, "Import payload property definitions thành công")
}

// GetNestedSchema godoc
// @Summary Get nested schema for one or all payload types
// @Description Returns payload properties as nested JSON tree. If payload_type query param is provided, only that type is returned. Otherwise all payload types are returned.
// @Tags Payload Property Definitions
// @Produce json
// @Param payload_type query string false "Filter by payload type (e.g. com.apple.carddav.account). Omit to get all."
// @Success 200 {object} response.APIResponse[[]service.NestedPayloadSchema]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
// @Router /api/v1/payload-property-definitions/schema [get]
func (h *payloadPropertyDefinitionHandlerImpl) GetNestedSchema(c *gin.Context) {
	payloadType := c.Query("payload_type") // optional

	schemas, err := h.service.GetNestedSchema(c.Request.Context(), payloadType)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, schemas, "Lấy nested schema thành công")
}
