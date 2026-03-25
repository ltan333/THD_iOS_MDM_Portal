package serviceimpl

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
	"howett.net/plist"
)

type mobileConfigServiceImpl struct {
	mobileConfigRepo repository.MobileConfigRepository
}

func NewMobileConfigService(mobileConfigRepo repository.MobileConfigRepository) service.MobileConfigService {
	return &mobileConfigServiceImpl{mobileConfigRepo: mobileConfigRepo}
}

func (m *mobileConfigServiceImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.MobileConfig, int64, error) {
	return m.mobileConfigRepo.List(ctx, offset, limit, opts)
}

func (m *mobileConfigServiceImpl) GetByID(ctx context.Context, id uint) (*ent.MobileConfig, error) {
	if id == 0 {
		return nil, apperror.ErrValidation.WithMessage("id là bắt buộc")
	}
	return m.mobileConfigRepo.GetByIDWithPayloads(ctx, id)
}

func (m *mobileConfigServiceImpl) Create(ctx context.Context, cmd service.CreateMobileConfigCommand) (*ent.MobileConfig, error) {
	if err := validateCreateMobileConfigCommand(cmd); err != nil {
		return nil, err
	}
	if err := validateDuplicateUniqueFieldsInRequest(cmd); err != nil {
		return nil, err
	}

	name := strings.TrimSpace(cmd.Name)
	payloadIdentifier := strings.TrimSpace(cmd.PayloadIdentifier)
	payloadIdentifiers := buildTrimmedPayloadIdentifiers(cmd.Payloads)

	conflict, err := m.mobileConfigRepo.FindCreateUniqueFieldConflict(ctx, name, payloadIdentifier, payloadIdentifiers)
	if err != nil {
		return nil, err
	}
	if conflict != nil {
		return nil, apperror.ErrConflict.WithMessage(fmt.Sprintf("%s đã tồn tại: %s", conflict.Field, conflict.Value))
	}

	entity, payloads := buildCreateEntities(cmd)

	created, err := m.mobileConfigRepo.Create(ctx, entity, payloads)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (m *mobileConfigServiceImpl) Update(ctx context.Context, cmd service.UpdateMobileConfigCommand) (*ent.MobileConfig, error) {
	if cmd.ID == 0 {
		return nil, apperror.ErrValidation.WithMessage("id là bắt buộc")
	}

	createCmd := service.CreateMobileConfigCommand{
		Name:                     cmd.Name,
		PayloadIdentifier:        cmd.PayloadIdentifier,
		PayloadType:              cmd.PayloadType,
		PayloadDisplayName:       cmd.PayloadDisplayName,
		PayloadDescription:       cmd.PayloadDescription,
		PayloadOrganization:      cmd.PayloadOrganization,
		PayloadVersion:           cmd.PayloadVersion,
		PayloadRemovalDisallowed: cmd.PayloadRemovalDisallowed,
		Payloads:                 cmd.Payloads,
	}

	if err := validateCreateMobileConfigCommand(createCmd); err != nil {
		return nil, err
	}
	if err := validateDuplicateUniqueFieldsInRequest(createCmd); err != nil {
		return nil, err
	}

	if _, err := m.mobileConfigRepo.GetByID(ctx, cmd.ID); err != nil {
		return nil, err
	}

	name := strings.TrimSpace(cmd.Name)
	payloadIdentifier := strings.TrimSpace(cmd.PayloadIdentifier)
	payloadIdentifiers := buildTrimmedPayloadIdentifiers(cmd.Payloads)

	conflict, err := m.mobileConfigRepo.FindUpdateUniqueFieldConflict(ctx, cmd.ID, name, payloadIdentifier, payloadIdentifiers)
	if err != nil {
		return nil, err
	}
	if conflict != nil {
		return nil, apperror.ErrConflict.WithMessage(fmt.Sprintf("%s đã tồn tại: %s", conflict.Field, conflict.Value))
	}

	entity, payloads := buildCreateEntities(createCmd)
	updated, err := m.mobileConfigRepo.Update(ctx, cmd.ID, entity, payloads)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (m *mobileConfigServiceImpl) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return apperror.ErrValidation.WithMessage("id là bắt buộc")
	}

	if _, err := m.mobileConfigRepo.GetByID(ctx, id); err != nil {
		return err
	}

	return m.mobileConfigRepo.Delete(ctx, id)
}

func buildCreateEntities(cmd service.CreateMobileConfigCommand) (*ent.MobileConfig, []*ent.Payload) {

	payloadVersion := cmd.PayloadVersion
	if payloadVersion <= 0 {
		payloadVersion = 1
	}

	entity := &ent.MobileConfig{
		Name:                     strings.TrimSpace(cmd.Name),
		PayloadIdentifier:        strings.TrimSpace(cmd.PayloadIdentifier),
		PayloadType:              strings.TrimSpace(cmd.PayloadType),
		PayloadDisplayName:       strings.TrimSpace(cmd.PayloadDisplayName),
		PayloadDescription:       cmd.PayloadDescription,
		PayloadOrganization:      cmd.PayloadOrganization,
		PayloadUUID:              uuid.NewString(),
		PayloadVersion:           payloadVersion,
		PayloadRemovalDisallowed: cmd.PayloadRemovalDisallowed,
	}

	payloads := make([]*ent.Payload, 0, len(cmd.Payloads))
	for _, payloadCmd := range cmd.Payloads {
		payloadItemVersion := payloadCmd.PayloadVersion
		if payloadItemVersion <= 0 {
			payloadItemVersion = 1
		}

		properties := make([]*ent.PayloadProperty, 0, len(payloadCmd.Properties))
		for _, propCmd := range payloadCmd.Properties {
			valueJSON := propCmd.ValueJSON
			if valueJSON == nil {
				valueJSON = map[string]interface{}{}
			}

			properties = append(properties, &ent.PayloadProperty{
				ValueJSON: valueJSON,
				Edges: ent.PayloadPropertyEdges{
					Definition: &ent.PayloadPropertyDefinition{Key: strings.TrimSpace(propCmd.Key)},
				},
			})
		}

		payloads = append(payloads, &ent.Payload{
			PayloadDescription:  payloadCmd.PayloadDescription,
			PayloadDisplayName:  strings.TrimSpace(payloadCmd.PayloadDisplayName),
			PayloadIdentifier:   strings.TrimSpace(payloadCmd.PayloadIdentifier),
			PayloadOrganization: payloadCmd.PayloadOrganization,
			PayloadType:         strings.TrimSpace(payloadCmd.PayloadType),
			PayloadUUID:         uuid.NewString(),
			PayloadVersion:      payloadItemVersion,
			Edges: ent.PayloadEdges{
				Properties: properties,
			},
		})
	}

	return entity, payloads
}

func validateCreateMobileConfigCommand(cmd service.CreateMobileConfigCommand) error {
	if strings.TrimSpace(cmd.Name) == "" {
		return apperror.ErrValidation.WithMessage("name là bắt buộc")
	}
	if strings.TrimSpace(cmd.PayloadIdentifier) == "" {
		return apperror.ErrValidation.WithMessage("payload_identifier là bắt buộc")
	}
	if strings.TrimSpace(cmd.PayloadType) == "" {
		return apperror.ErrValidation.WithMessage("payload_type là bắt buộc")
	}
	if strings.TrimSpace(cmd.PayloadDisplayName) == "" {
		return apperror.ErrValidation.WithMessage("payload_display_name là bắt buộc")
	}
	if len(cmd.Payloads) == 0 {
		return apperror.ErrValidation.WithMessage("Mobile config phải có ít nhất một payload")
	}

	for payloadIdx, payload := range cmd.Payloads {
		if strings.TrimSpace(payload.PayloadDisplayName) == "" {
			return apperror.ErrValidation.WithMessage(fmt.Sprintf("payloads[%d].payload_display_name là bắt buộc", payloadIdx))
		}
		if strings.TrimSpace(payload.PayloadIdentifier) == "" {
			return apperror.ErrValidation.WithMessage(fmt.Sprintf("payloads[%d].payload_identifier là bắt buộc", payloadIdx))
		}
		if strings.TrimSpace(payload.PayloadType) == "" {
			return apperror.ErrValidation.WithMessage(fmt.Sprintf("payloads[%d].payload_type là bắt buộc", payloadIdx))
		}

		for propIdx, property := range payload.Properties {
			if strings.TrimSpace(property.Key) == "" {
				return apperror.ErrValidation.WithMessage(fmt.Sprintf("payloads[%d].properties[%d].key là bắt buộc", payloadIdx, propIdx))
			}
		}
	}

	return nil
}

func validateDuplicateUniqueFieldsInRequest(cmd service.CreateMobileConfigCommand) error {
	seenPayloadIdentifiers := make(map[string]struct{}, len(cmd.Payloads))
	for idx, payloadCmd := range cmd.Payloads {
		identifier := strings.TrimSpace(payloadCmd.PayloadIdentifier)
		if _, exists := seenPayloadIdentifiers[identifier]; exists {
			return apperror.ErrConflict.WithMessage(fmt.Sprintf("payloads[%d].payload_identifier bị trùng trong request: %s", idx, identifier))
		}
		seenPayloadIdentifiers[identifier] = struct{}{}
	}

	return nil
}

func buildTrimmedPayloadIdentifiers(payloads []service.CreateMobileConfigPayloadCommand) []string {
	if len(payloads) == 0 {
		return nil
	}

	identifiers := make([]string, 0, len(payloads))
	for _, payloadCmd := range payloads {
		identifiers = append(identifiers, strings.TrimSpace(payloadCmd.PayloadIdentifier))
	}

	return identifiers
}

func (m *mobileConfigServiceImpl) GenerateXML(ctx context.Context, cmd service.GenerateMobileConfigXMLCommand) ([]byte, error) {
	mc, err := m.mobileConfigRepo.GetFullForExport(ctx, cmd.ID)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("MobileConfig không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất MobileConfig").WithError(err)
	}

	xmlBytes, err := buildMobileConfigXML(mc)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo XML").WithError(err)
	}

	return xmlBytes, nil
}

func buildMobileConfigXML(mc *ent.MobileConfig) ([]byte, error) {
	root := map[string]interface{}{
		"PayloadDescription":       mc.PayloadDescription,
		"PayloadDisplayName":       mc.PayloadDisplayName,
		"PayloadIdentifier":        mc.PayloadIdentifier,
		"PayloadOrganization":      mc.PayloadOrganization,
		"PayloadRemovalDisallowed": mc.PayloadRemovalDisallowed,
		"PayloadType":              mc.PayloadType,
		"PayloadUUID":              mc.PayloadUUID,
		"PayloadVersion":           mc.PayloadVersion,
		"PayloadContent":           buildPayloadContent(mc.Edges.Payloads),
	}

	return plist.MarshalIndent(root, plist.XMLFormat, "\t")
}

func buildPayloadContent(payloads []*ent.Payload) []map[string]interface{} {
	var content []map[string]interface{}

	for _, p := range payloads {
		payloadDict := map[string]interface{}{
			"PayloadType":         p.PayloadType,
			"PayloadVersion":      p.PayloadVersion,
			"PayloadIdentifier":   p.PayloadIdentifier,
			"PayloadUUID":         p.PayloadUUID,
			"PayloadDisplayName":  p.PayloadDisplayName,
			"PayloadDescription":  p.PayloadDescription,
			"PayloadOrganization": p.PayloadOrganization,
		}

		// Convert properties to key-value pairs in the payload content
		for _, prop := range p.Edges.Properties {
			if prop.Edges.Definition == nil {
				continue
			}
			key := prop.Edges.Definition.Key
			val := prop.ValueJSON

			// Normalize value based on the defined value type
			payloadDict[key] = normalizeValue(val, prop.Edges.Definition.ValueType)
		}

		content = append(content, payloadDict)
	}

	return content
}

func normalizeValue(val interface{}, valueType string) interface{} {
	raw := extractRawValue(val)

	switch strings.ToLower(valueType) {
	case "bool":
		if b, ok := raw.(bool); ok {
			return b
		}
		return false
	case "integer", "number":
		switch v := raw.(type) {
		case float64:
			return int(v)
		case int:
			return v
		case int64:
			return int(v)
		}
		return 0
	case "string":
		if s, ok := raw.(string); ok {
			return s
		}
		return ""
	case "array", "dict":
		return raw
	default:
		return raw
	}
}

func extractRawValue(val interface{}) interface{} {
	if m, ok := val.(map[string]interface{}); ok {
		if unwrapped, exists := m["value"]; exists {
			return unwrapped
		}
	}
	return val
}
