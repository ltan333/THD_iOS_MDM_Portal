package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/response"
)

type MobileConfigHandler interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetXML(c *gin.Context)
}

type mobileConfigHandlerImpl struct {
	mobileConfigService service.MobileConfigService
}

func NewMobileConfigHandler(mobileConfigService service.MobileConfigService) MobileConfigHandler {
	return &mobileConfigHandlerImpl{mobileConfigService: mobileConfigService}
}

// Create godoc
// @Summary Create mobile config
// @Description Create a new Apple mobileconfig with payloads and payload properties
// @Tags Mobile Config
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateMobileConfigRequest true "Mobile config payload"
// @Success 201 {object} MobileConfigSuccessResponse
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 409 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /v1/mobile-configs [post]
func (m *mobileConfigHandlerImpl) Create(c *gin.Context) {
	var req dto.CreateMobileConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if fieldErrors := buildCreateValidationFieldErrors(err); len(fieldErrors) > 0 {
			response.ValidationError(c, fieldErrors)
			return
		}
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Dữ liệu không hợp lệ"))
		return
	}

	if fieldErrors := validateCreateMobileConfigRequest(req); len(fieldErrors) > 0 {
		response.ValidationError(c, fieldErrors)
		return
	}

	payloads := make([]service.CreateMobileConfigPayloadCommand, 0, len(req.Payloads))
	for _, payloadReq := range req.Payloads {
		properties := make([]service.CreateMobileConfigPropertyCommand, 0, len(payloadReq.Properties))
		for _, propReq := range payloadReq.Properties {
			properties = append(properties, service.CreateMobileConfigPropertyCommand{
				Key:       propReq.Key,
				ValueJSON: propReq.ValueJSON,
			})
		}

		payloads = append(payloads, service.CreateMobileConfigPayloadCommand{
			PayloadDescription:  payloadReq.PayloadDescription,
			PayloadDisplayName:  payloadReq.PayloadDisplayName,
			PayloadIdentifier:   payloadReq.PayloadIdentifier,
			PayloadOrganization: payloadReq.PayloadOrganization,
			PayloadType:         payloadReq.PayloadType,
			PayloadVersion:      payloadReq.PayloadVersion,
			Properties:          properties,
		})
	}

	created, err := m.mobileConfigService.Create(c.Request.Context(), service.CreateMobileConfigCommand{
		Name:                     req.Name,
		PayloadIdentifier:        req.PayloadIdentifier,
		PayloadType:              req.PayloadType,
		PayloadDisplayName:       req.PayloadDisplayName,
		PayloadDescription:       req.PayloadDescription,
		PayloadOrganization:      req.PayloadOrganization,
		PayloadVersion:           req.PayloadVersion,
		PayloadRemovalDisallowed: req.PayloadRemovalDisallowed,
		Payloads:                 payloads,
	})
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}
	if created == nil {
		response.WriteErrorResponse(c, apperror.ErrInternalServerError.WithMessage("Không thể tạo mobile config"))
		return
	}

	response.Created(c, toMobileConfigResponse(created), "Tạo mobile config thành công")
}

// Update godoc
// @Summary Update mobile config
// @Description Update an existing Apple mobileconfig with payloads and payload properties
// @Tags Mobile Config
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Mobile config ID"
// @Param request body dto.UpdateMobileConfigRequest true "Mobile config payload"
// @Success 200 {object} MobileConfigSuccessResponse
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 409 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /v1/mobile-configs/{id} [put]
func (m *mobileConfigHandlerImpl) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("ID không hợp lệ"))
		return
	}

	var req dto.UpdateMobileConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if fieldErrors := buildValidationFieldErrors(err, "UpdateMobileConfigRequest."); len(fieldErrors) > 0 {
			response.ValidationError(c, fieldErrors)
			return
		}
		response.WriteErrorResponse(c, apperror.ErrValidation.WithMessage("Dữ liệu không hợp lệ"))
		return
	}

	createReq := dto.CreateMobileConfigRequest(req)
	if fieldErrors := validateCreateMobileConfigRequest(createReq); len(fieldErrors) > 0 {
		response.ValidationError(c, fieldErrors)
		return
	}

	updated, err := m.mobileConfigService.Update(c.Request.Context(), service.UpdateMobileConfigCommand{
		ID:                       uint(id),
		Name:                     req.Name,
		PayloadIdentifier:        req.PayloadIdentifier,
		PayloadType:              req.PayloadType,
		PayloadDisplayName:       req.PayloadDisplayName,
		PayloadDescription:       req.PayloadDescription,
		PayloadOrganization:      req.PayloadOrganization,
		PayloadVersion:           req.PayloadVersion,
		PayloadRemovalDisallowed: req.PayloadRemovalDisallowed,
		Payloads:                 toCreatePayloadCommands(req.Payloads),
	})
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK(c, toMobileConfigResponse(updated), "Cập nhật mobile config thành công")
}

// Delete godoc
// @Summary Delete mobile config
// @Description Delete an existing Apple mobileconfig by ID
// @Tags Mobile Config
// @Produce json
// @Security BearerAuth
// @Param id path int true "Mobile config ID"
// @Success 200 {object} EmptySuccessResponse
// @Failure 400 {object} APIErrorResponse
// @Failure 401 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /v1/mobile-configs/{id} [delete]
func (m *mobileConfigHandlerImpl) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("ID không hợp lệ"))
		return
	}

	if err := m.mobileConfigService.Delete(c.Request.Context(), uint(id)); err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	response.OK[any](c, nil, "Xóa mobile config thành công")
}

func buildCreateValidationFieldErrors(err error) []response.FieldError {
	return buildValidationFieldErrors(err, "CreateMobileConfigRequest.")
}

func buildValidationFieldErrors(err error, namespacePrefix string) []response.FieldError {
	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	fields := make([]response.FieldError, 0, len(validationErrs))
	for _, fieldErr := range validationErrs {
		fieldName := toRequestFieldName(fieldErr.Namespace(), namespacePrefix)
		fields = append(fields, response.FieldError{
			Field:   fieldName,
			Message: validationMessageFromTag(fieldErr.Tag(), fieldErr.Param()),
		})
	}

	return fields
}

func toRequestFieldName(namespace string, namespacePrefix string) string {
	path := strings.TrimPrefix(namespace, namespacePrefix)
	replacer := strings.NewReplacer(
		"Name", "name",
		"PayloadIdentifier", "payload_identifier",
		"PayloadType", "payload_type",
		"PayloadDisplayName", "payload_display_name",
		"PayloadDescription", "payload_description",
		"PayloadOrganization", "payload_organization",
		"PayloadVersion", "payload_version",
		"PayloadRemovalDisallowed", "payload_removal_disallowed",
		"Payloads", "payloads",
		"Properties", "properties",
		"Key", "key",
		"ValueJSON", "value_json",
	)

	return replacer.Replace(path)
}

func toCreatePayloadCommands(payloadReqs []dto.UpdateMobileConfigPayloadRequest) []service.CreateMobileConfigPayloadCommand {
	payloads := make([]service.CreateMobileConfigPayloadCommand, 0, len(payloadReqs))
	for _, payloadReq := range payloadReqs {
		properties := make([]service.CreateMobileConfigPropertyCommand, 0, len(payloadReq.Properties))
		for _, propReq := range payloadReq.Properties {
			properties = append(properties, service.CreateMobileConfigPropertyCommand{
				Key:       propReq.Key,
				ValueJSON: propReq.ValueJSON,
			})
		}

		payloads = append(payloads, service.CreateMobileConfigPayloadCommand{
			PayloadDescription:  payloadReq.PayloadDescription,
			PayloadDisplayName:  payloadReq.PayloadDisplayName,
			PayloadIdentifier:   payloadReq.PayloadIdentifier,
			PayloadOrganization: payloadReq.PayloadOrganization,
			PayloadType:         payloadReq.PayloadType,
			PayloadVersion:      payloadReq.PayloadVersion,
			Properties:          properties,
		})
	}

	return payloads
}

func validationMessageFromTag(tag string, param string) string {
	switch tag {
	case "required":
		return "Trường này là bắt buộc"
	case "min":
		return "Giá trị phải lớn hơn hoặc bằng " + param
	default:
		return "Dữ liệu không hợp lệ"
	}
}

func validateCreateMobileConfigRequest(req dto.CreateMobileConfigRequest) []response.FieldError {
	fieldErrors := make([]response.FieldError, 0)
	if strings.TrimSpace(req.Name) == "" {
		fieldErrors = append(fieldErrors, response.FieldError{Field: "name", Message: "Trường này là bắt buộc"})
	}
	if strings.TrimSpace(req.PayloadIdentifier) == "" {
		fieldErrors = append(fieldErrors, response.FieldError{Field: "payload_identifier", Message: "Trường này là bắt buộc"})
	}
	if strings.TrimSpace(req.PayloadType) == "" {
		fieldErrors = append(fieldErrors, response.FieldError{Field: "payload_type", Message: "Trường này là bắt buộc"})
	}
	if strings.TrimSpace(req.PayloadDisplayName) == "" {
		fieldErrors = append(fieldErrors, response.FieldError{Field: "payload_display_name", Message: "Trường này là bắt buộc"})
	}

	for i, payloadReq := range req.Payloads {
		if strings.TrimSpace(payloadReq.PayloadDisplayName) == "" {
			fieldErrors = append(fieldErrors, response.FieldError{Field: "payloads[" + strconv.Itoa(i) + "].payload_display_name", Message: "Trường này là bắt buộc"})
		}
		if strings.TrimSpace(payloadReq.PayloadIdentifier) == "" {
			fieldErrors = append(fieldErrors, response.FieldError{Field: "payloads[" + strconv.Itoa(i) + "].payload_identifier", Message: "Trường này là bắt buộc"})
		}
		if strings.TrimSpace(payloadReq.PayloadType) == "" {
			fieldErrors = append(fieldErrors, response.FieldError{Field: "payloads[" + strconv.Itoa(i) + "].payload_type", Message: "Trường này là bắt buộc"})
		}

		for j, propReq := range payloadReq.Properties {
			if strings.TrimSpace(propReq.Key) == "" {
				fieldErrors = append(fieldErrors, response.FieldError{Field: "payloads[" + strconv.Itoa(i) + "].properties[" + strconv.Itoa(j) + "].key", Message: "Trường này là bắt buộc"})
			}
			if propReq.ValueJSON == nil {
				fieldErrors = append(fieldErrors, response.FieldError{Field: "payloads[" + strconv.Itoa(i) + "].properties[" + strconv.Itoa(j) + "].value_json", Message: "Trường này là bắt buộc"})
			}
		}
	}

	return fieldErrors
}

// GetXML godoc
// @Summary Export mobile config XML
// @Description Generate and return raw Apple mobileconfig XML content by ID
// @Tags Mobile Config
// @Produce xml
// @Param id path int true "Mobile config ID"
// @Success 200 {string} string "Raw XML"
// @Failure 400 {object} APIErrorResponse
// @Failure 404 {object} APIErrorResponse
// @Failure 500 {object} APIErrorResponse
// @Router /v1/mobile-configs/{id}/xml [get]
func (m *mobileConfigHandlerImpl) GetXML(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.WriteErrorResponse(c, apperror.ErrBadRequest.WithMessage("ID không hợp lệ"))
		return
	}

	cmd := service.GenerateMobileConfigXMLCommand{ID: uint(id)}
	xmlBytes, err := m.mobileConfigService.GenerateXML(c.Request.Context(), cmd)
	if err != nil {
		response.WriteErrorResponse(c, err)
		return
	}

	// Return raw XML
	c.Data(http.StatusOK, "text/xml; charset=utf-8", xmlBytes)
}

func toMobileConfigResponse(mc *ent.MobileConfig) dto.MobileConfigResponse {
	payloadResponses := make([]dto.MobileConfigPayloadResponse, 0, len(mc.Edges.Payloads))
	for _, p := range mc.Edges.Payloads {
		propertyResponses := make([]dto.MobileConfigPropertyResponse, 0, len(p.Edges.Properties))
		for _, prop := range p.Edges.Properties {
			key := ""
			if prop.Edges.Definition != nil {
				key = prop.Edges.Definition.Key
			}

			propertyResponses = append(propertyResponses, dto.MobileConfigPropertyResponse{
				ID:        prop.ID,
				Key:       key,
				ValueJSON: prop.ValueJSON,
			})
		}

		payloadResponses = append(payloadResponses, dto.MobileConfigPayloadResponse{
			ID:                  p.ID,
			PayloadDescription:  p.PayloadDescription,
			PayloadDisplayName:  p.PayloadDisplayName,
			PayloadIdentifier:   p.PayloadIdentifier,
			PayloadOrganization: p.PayloadOrganization,
			PayloadType:         p.PayloadType,
			PayloadUUID:         p.PayloadUUID,
			PayloadVersion:      p.PayloadVersion,
			Properties:          propertyResponses,
		})
	}

	return dto.MobileConfigResponse{
		ID:                       mc.ID,
		Name:                     mc.Name,
		PayloadIdentifier:        mc.PayloadIdentifier,
		PayloadType:              mc.PayloadType,
		PayloadDisplayName:       mc.PayloadDisplayName,
		PayloadDescription:       mc.PayloadDescription,
		PayloadOrganization:      mc.PayloadOrganization,
		PayloadUUID:              mc.PayloadUUID,
		PayloadVersion:           mc.PayloadVersion,
		PayloadRemovalDisallowed: mc.PayloadRemovalDisallowed,
		Payloads:                 payloadResponses,
		CreatedAt:                mc.CreatedAt,
		UpdatedAt:                mc.UpdatedAt,
	}
}
