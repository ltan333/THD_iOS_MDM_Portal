package handler

import (
	"io"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
	"github.com/thienel/go-backend-template/pkg/response"
)

var payloadPropertyDefinitionAllowedQueryFields = map[string]bool{
	"id":           true,
	"payload_type": true,
	"key":          true,
	"value_type":   true,
	"description":  true,
	"created_at":   true,
	"updated_at":   true,
	"search":       true,
}

type PayloadPropertyDefinitionHandler interface {
	List(c *gin.Context)
	GetByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Import(c *gin.Context)
}

type payloadPropertyDefinitionHandlerImpl struct {
	service service.PayloadPropertyDefinitionService
}

func NewPayloadPropertyDefinitionHandler(service service.PayloadPropertyDefinitionService) PayloadPropertyDefinitionHandler {
	return &payloadPropertyDefinitionHandlerImpl{service: service}
}

// List godoc
// @Summary List payload property definitions
// @Description Get a paginated list of payload property definitions with filtering and sorting
// @Tags Payload Property Definitions
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 20)"
// @Param search query string false "Search by payload_type, key, description"
// @Param payload_type query string false "Filter by payload_type"
// @Param key query string false "Filter by key"
// @Param value_type query string false "Filter by value_type"
// @Param sort query string false "Sort by field (id,payload_type,key,value_type,created_at)"
// @Success 200 {object} response.APIResponse[dto.ListResponse[dto.PayloadPropertyDefinitionResponse]]
// @Failure 401 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/payload-property-definitions [get]
func (h *payloadPropertyDefinitionHandlerImpl) List(c *gin.Context) {
	params := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	offset, limit := query.GetPagination(params, 20)
	opts := query.ParseQueryParams(params, payloadPropertyDefinitionAllowedQueryFields)

	items, total, err := h.service.List(c.Request.Context(), offset, limit, opts)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	responses := make([]dto.PayloadPropertyDefinitionResponse, len(items))
	for i, item := range items {
		responses[i] = toPayloadPropertyDefinitionResponse(item)
	}

	page := (offset / limit) + 1
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response.OK(c, dto.ListResponse[dto.PayloadPropertyDefinitionResponse]{
		Items:      responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, "")
}

// GetByID godoc
// @Summary Get payload property definition by ID
// @Description Get details of a payload property definition by ID
// @Tags Payload Property Definitions
// @Produce json
// @Param id path int true "Definition ID"
// @Success 200 {object} response.APIResponse[dto.PayloadPropertyDefinitionResponse]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 404 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/payload-property-definitions/{id} [get]
func (h *payloadPropertyDefinitionHandlerImpl) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("ID không hợp lệ"))
		return
	}

	item, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, toPayloadPropertyDefinitionResponse(item), "")
}

// Create godoc
// @Summary Create payload property definition
// @Description Create a new payload property definition
// @Tags Payload Property Definitions
// @Accept json
// @Produce json
// @Param request body dto.CreatePayloadPropertyDefinitionRequest true "Payload property definition data"
// @Success 201 {object} response.APIResponse[dto.PayloadPropertyDefinitionResponse]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 409 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/payload-property-definitions [post]
func (h *payloadPropertyDefinitionHandlerImpl) Create(c *gin.Context) {
	var req dto.CreatePayloadPropertyDefinitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Dữ liệu không hợp lệ"))
		return
	}

	item, err := h.service.Create(c.Request.Context(), service.CreatePayloadPropertyDefinitionCommand{
		PayloadType:  req.PayloadType,
		Key:          req.Key,
		ValueType:    req.ValueType,
		DefaultValue: req.DefaultValue,
		EnumValues:   req.EnumValues,
		Deprecated:   req.Deprecated,
		Description:  req.Description,
	})
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.Created(c, toPayloadPropertyDefinitionResponse(item), "Tạo định nghĩa thuộc tính payload thành công")
}

// Update godoc
// @Summary Update payload property definition
// @Description Update an existing payload property definition
// @Tags Payload Property Definitions
// @Accept json
// @Produce json
// @Param id path int true "Definition ID"
// @Param request body dto.UpdatePayloadPropertyDefinitionRequest true "Payload property definition data"
// @Success 200 {object} response.APIResponse[dto.PayloadPropertyDefinitionResponse]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 404 {object} response.APIResponse[any]
// @Failure 409 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/payload-property-definitions/{id} [put]
func (h *payloadPropertyDefinitionHandlerImpl) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("ID không hợp lệ"))
		return
	}

	var req dto.UpdatePayloadPropertyDefinitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Dữ liệu không hợp lệ"))
		return
	}

	item, err := h.service.Update(c.Request.Context(), service.UpdatePayloadPropertyDefinitionCommand{
		ID:           uint(id),
		PayloadType:  req.PayloadType,
		Key:          req.Key,
		ValueType:    req.ValueType,
		DefaultValue: req.DefaultValue,
		EnumValues:   req.EnumValues,
		Deprecated:   req.Deprecated,
		Description:  req.Description,
	})
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, toPayloadPropertyDefinitionResponse(item), "Cập nhật định nghĩa thuộc tính payload thành công")
}

// Delete godoc
// @Summary Delete payload property definition
// @Description Delete a payload property definition by ID
// @Tags Payload Property Definitions
// @Produce json
// @Param id path int true "Definition ID"
// @Success 204
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 404 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/payload-property-definitions/{id} [delete]
func (h *payloadPropertyDefinitionHandlerImpl) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("ID không hợp lệ"))
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.NoContent(c)
}

// Import godoc
// @Summary Import payload property definitions from Apple JSON
// @Description Upload an Apple documentation JSON file and import payload property definitions
// @Tags Payload Property Definitions
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Apple payload JSON file"
// @Success 200 {object} response.APIResponse[service.ImportPayloadPropertyDefinitionsResult]
// @Failure 400 {object} response.APIResponse[any]
// @Failure 401 {object} response.APIResponse[any]
// @Failure 500 {object} response.APIResponse[any]
// @Security BearerAuth
// @Router /v1/payload-property-definitions/import [post]
func (h *payloadPropertyDefinitionHandlerImpl) Import(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Thiếu file JSON để import"))
		return
	}

	src, err := file.Open()
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Không thể đọc file upload").WithError(err))
		return
	}
	defer func() { _ = src.Close() }()

	data, err := io.ReadAll(src)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("Không thể đọc nội dung file").WithError(err))
		return
	}

	result, err := h.service.ImportFromAppleJSON(c.Request.Context(), file.Filename, data)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, result, "Import payload property definitions thành công")
}

func toPayloadPropertyDefinitionResponse(item *ent.PayloadPropertyDefinition) dto.PayloadPropertyDefinitionResponse {
	return dto.PayloadPropertyDefinitionResponse{
		ID:           item.ID,
		PayloadType:  item.PayloadType,
		Key:          item.Key,
		ValueType:    item.ValueType,
		DefaultValue: item.DefaultValue,
		EnumValues:   item.EnumValues,
		Deprecated:   item.Deprecated,
		Description:  item.Description,
		CreatedAt:    item.CreatedAt,
		UpdatedAt:    item.UpdatedAt,
	}
}
