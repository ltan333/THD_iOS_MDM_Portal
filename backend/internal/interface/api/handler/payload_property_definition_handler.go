package handler

import (
	"archive/zip"
	"bytes"
	"fmt"
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
	ListVariants(c *gin.Context)
	DeleteAll(c *gin.Context)
	ImportYAMLFolder(c *gin.Context)
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
// @Router /v1/payload-property-definitions/payload-types [get]
func (h *payloadPropertyDefinitionHandlerImpl) ListPayloadTypes(c *gin.Context) {
	types, err := h.service.ListPayloadTypes(c.Request.Context())
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, types, "")
}

// ListVariants godoc
// @Summary List payload variants by payload type
// @Description Get all variants of a payload type with property count
// @Tags Payload Property Definitions
// @Produce json
// @Param payload_type path string true "Payload type"
// @Success 200 {object} response.APIResponse[[]service.PayloadVariantInfo]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/payload-property-definitions/{payload_type}/variants [get]
func (h *payloadPropertyDefinitionHandlerImpl) ListVariants(c *gin.Context) {
	payloadType := strings.TrimSpace(c.Param("payload_type"))
	if payloadType == "" {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Thiếu payload_type"))
		return
	}

	variants, err := h.service.ListVariants(c.Request.Context(), payloadType)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, variants, "Lấy danh sách variants thành công")
}

// DeleteAll godoc
// @Summary Delete all payload property definitions
// @Description Delete all payload property definitions from database
// @Tags Payload Property Definitions
// @Produce json
// @Success 200 {object} response.APIResponse[map[string]int]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/payload-property-definitions [delete]
func (h *payloadPropertyDefinitionHandlerImpl) DeleteAll(c *gin.Context) {
	deletedCount, err := h.service.DeleteAll(c.Request.Context())
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, gin.H{"deleted_count": deletedCount}, "Xóa toàn bộ payload property definitions thành công")
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
// @Router /v1/payload-property-definitions/schema [get]
func (h *payloadPropertyDefinitionHandlerImpl) GetNestedSchema(c *gin.Context) {
	payloadType := c.Query("payload_type") // optional
	payloadVariant := c.Query("variant")   // optional

	schemas, err := h.service.GetNestedSchema(c.Request.Context(), payloadType, payloadVariant)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, schemas, "Lấy nested schema thành công")
}

// ImportYAMLFolder godoc
// @Summary Import payload property definitions from YAML folder
// @Description Upload a folder or zip file containing Apple Device Management YAML files. Supports multipart fields: folder, folder[], files, files[], file.
// @Tags Payload Property Definitions
// @Accept multipart/form-data
// @Produce json
// @Param folder formData file true "Folder zip file or files uploaded from folder picker"
// @Param files formData []file false "Alternative field for folder files"
// @Success 200 {object} response.APIResponse[service.ImportPayloadPropertyDefinitionsResult]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/payload-property-definitions/import-yaml-folder [post]
func (h *payloadPropertyDefinitionHandlerImpl) ImportYAMLFolder(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Không thể parse multipart form").WithError(err))
		return
	}

	var allFiles []*multipart.FileHeader
	for _, fieldName := range []string{"folder", "folder[]", "files", "files[]", "file"} {
		if fhs, ok := form.File[fieldName]; ok {
			allFiles = append(allFiles, fhs...)
		}
	}

	if len(allFiles) == 0 {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage(
			"Thiếu folder/file YAML để import. Vui lòng upload field 'folder' hoặc 'files'"))
		return
	}

	fileMap := make(map[string][]byte)
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

		nameLower := strings.ToLower(filepath.Base(fh.Filename))
		if strings.HasSuffix(nameLower, ".zip") {
			if err := addYAMLFilesFromZip(fileMap, data); err != nil {
				response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Không thể đọc zip: "+fh.Filename).WithError(err))
				return
			}
			continue
		}

		if strings.HasSuffix(nameLower, ".yaml") || strings.HasSuffix(nameLower, ".yml") {
			fileMap[nameLower] = data
		}
	}

	if len(fileMap) == 0 {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Không tìm thấy file .yaml hoặc .yml hợp lệ để import"))
		return
	}

	result, err := h.service.ImportFromAppleYAMLFiles(c.Request.Context(), fileMap)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, result, buildImportPayloadDefinitionsMessage("folder YAML", result))
}

func buildImportPayloadDefinitionsMessage(source string, result *service.ImportPayloadPropertyDefinitionsResult) string {
	if result == nil {
		return fmt.Sprintf("Import %s payload property definitions hoàn tất", source)
	}

	processed := result.Created + result.Updated
	errorCount := len(result.Errors)

	if errorCount == 0 {
		return fmt.Sprintf("Import %s payload property definitions thành công", source)
	}

	if processed == 0 {
		return fmt.Sprintf("Import %s payload property definitions thất bại: %d lỗi", source, errorCount)
	}

	return fmt.Sprintf("Import %s payload property definitions hoàn tất một phần: %d thành công, %d lỗi", source, processed, errorCount)
}

func addYAMLFilesFromZip(fileMap map[string][]byte, zipData []byte) error {
	r, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return err
	}

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}

		nameLower := strings.ToLower(filepath.Base(f.Name))
		if !strings.HasSuffix(nameLower, ".yaml") && !strings.HasSuffix(nameLower, ".yml") {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}
		data, err := io.ReadAll(rc)
		_ = rc.Close()
		if err != nil {
			return err
		}

		fileMap[nameLower] = data
	}

	return nil
}
