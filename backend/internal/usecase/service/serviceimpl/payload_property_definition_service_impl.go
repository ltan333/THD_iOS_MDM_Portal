package serviceimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/thienel/tlog"
	"go.uber.org/zap"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
)

type payloadPropertyDefinitionServiceImpl struct {
	payloadPropertyDefinitionRepo repository.PayloadPropertyDefinitionRepository
}

// NewPayloadPropertyDefinitionService creates a new payload property definition service
func NewPayloadPropertyDefinitionService(repo repository.PayloadPropertyDefinitionRepository) service.PayloadPropertyDefinitionService {
	return &payloadPropertyDefinitionServiceImpl{
		payloadPropertyDefinitionRepo: repo,
	}
}

func (s *payloadPropertyDefinitionServiceImpl) Create(ctx context.Context, cmd service.CreatePayloadPropertyDefinitionCommand) (*ent.PayloadPropertyDefinition, error) {
	if err := validateCreatePayloadPropertyDefinitionCommand(cmd); err != nil {
		return nil, err
	}

	entity := &ent.PayloadPropertyDefinition{
		PayloadType:  strings.TrimSpace(cmd.PayloadType),
		Key:          strings.TrimSpace(cmd.Key),
		ValueType:    strings.TrimSpace(cmd.ValueType),
		DefaultValue: cmd.DefaultValue,
		EnumValues:   cmd.EnumValues,
		Deprecated:   cmd.Deprecated,
		Description:  strings.TrimSpace(cmd.Description),
	}

	if err := s.payloadPropertyDefinitionRepo.Create(ctx, entity); err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *payloadPropertyDefinitionServiceImpl) GetByID(ctx context.Context, id uint) (*ent.PayloadPropertyDefinition, error) {
	if id == 0 {
		return nil, apperror.ErrValidation.WithMessage("id là bắt buộc")
	}

	item, err := s.payloadPropertyDefinitionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *payloadPropertyDefinitionServiceImpl) Update(ctx context.Context, cmd service.UpdatePayloadPropertyDefinitionCommand) (*ent.PayloadPropertyDefinition, error) {
	if cmd.ID == 0 {
		return nil, apperror.ErrValidation.WithMessage("id là bắt buộc")
	}

	if err := validateUpdatePayloadPropertyDefinitionCommand(cmd); err != nil {
		return nil, err
	}

	existing, err := s.payloadPropertyDefinitionRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	existing.PayloadType = strings.TrimSpace(cmd.PayloadType)
	existing.Key = strings.TrimSpace(cmd.Key)
	existing.ValueType = strings.TrimSpace(cmd.ValueType)
	existing.DefaultValue = cmd.DefaultValue
	existing.EnumValues = cmd.EnumValues
	existing.Deprecated = cmd.Deprecated
	existing.Description = strings.TrimSpace(cmd.Description)

	if err := s.payloadPropertyDefinitionRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *payloadPropertyDefinitionServiceImpl) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return apperror.ErrValidation.WithMessage("id là bắt buộc")
	}

	if _, err := s.payloadPropertyDefinitionRepo.FindByID(ctx, id); err != nil {
		return err
	}

	return s.payloadPropertyDefinitionRepo.Delete(ctx, id)
}

func (s *payloadPropertyDefinitionServiceImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.PayloadPropertyDefinition, int64, error) {
	items, total, err := s.payloadPropertyDefinitionRepo.ListWithQuery(ctx, offset, limit, opts)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (s *payloadPropertyDefinitionServiceImpl) ListPayloadTypes(ctx context.Context) ([]string, error) {
	return s.payloadPropertyDefinitionRepo.ListPayloadTypes(ctx)
}

func (s *payloadPropertyDefinitionServiceImpl) ImportFromAppleJSON(ctx context.Context, filename string, data []byte) (*service.ImportPayloadPropertyDefinitionsResult, error) {
	result := &service.ImportPayloadPropertyDefinitionsResult{
		PayloadType: "",
		Total:       0,
		Created:     0,
		Updated:     0,
		Errors:      []string{},
	}

	// Parse the JSON data
	var jsonDoc map[string]interface{}
	if err := json.Unmarshal(data, &jsonDoc); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Lỗi parse JSON: %v", err))
		return result, nil
	}

	// Detect payload type from the document
	payloadType := s.detectPayloadTypeFromDiscussion(jsonDoc)
	if payloadType == "" {
		result.Errors = append(result.Errors, "Không thể xác định payload type từ document")
		return result, nil
	}

	result.PayloadType = payloadType

	// Extract properties from the document
	properties := s.extractPropertiesFromJSON(jsonDoc)
	result.Total = len(properties)

	// Upsert all properties
	for _, prop := range properties {
		prop.PayloadType = payloadType
		_, created, err := s.payloadPropertyDefinitionRepo.UpsertByPayloadTypeAndKey(ctx, prop)
		if err != nil {
			tlog.Error("Failed to upsert payload property definition", zap.String("key", prop.Key), zap.Error(err))
			result.Errors = append(result.Errors, fmt.Sprintf("Lỗi upsert key '%s': %v", prop.Key, err))
			continue
		}

		if created {
			result.Created++
		} else {
			result.Updated++
		}
	}

	tlog.Info("Imported payload property definitions",
		zap.String("payload_type", payloadType),
		zap.Int("total", result.Total),
		zap.Int("created", result.Created),
		zap.Int("updated", result.Updated),
		zap.Int("errors", len(result.Errors)),
	)

	return result, nil
}

// detectPayloadTypeFromDiscussion extracts the payload type from the Discussion section
// of the Apple MDM documentation JSON
func (s *payloadPropertyDefinitionServiceImpl) detectPayloadTypeFromDiscussion(jsonDoc map[string]interface{}) string {
	sections, ok := jsonDoc["primaryContentSections"].([]interface{})
	if !ok {
		return ""
	}

	// Find the content section
	for _, section := range sections {
		sectionMap, ok := section.(map[string]interface{})
		if !ok {
			continue
		}

		kind, ok := sectionMap["kind"].(string)
		if !ok || kind != "content" {
			continue
		}

		content, ok := sectionMap["content"].([]interface{})
		if !ok {
			continue
		}

		// Look for a paragraph with codeVoice containing the payload type
		payloadType := s.extractPayloadTypeFromInlineContent(content)
		if payloadType != "" {
			return payloadType
		}
	}

	return ""
}

// extractPayloadTypeFromInlineContent searches through the content array
// to find a paragraph with inlineContent containing a codeVoice element
func (s *payloadPropertyDefinitionServiceImpl) extractPayloadTypeFromInlineContent(content []interface{}) string {
	for _, item := range content {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		// Check if it's a paragraph
		itemType, ok := itemMap["type"].(string)
		if !ok || itemType != "paragraph" {
			continue
		}

		inlineContent, ok := itemMap["inlineContent"].([]interface{})
		if !ok {
			continue
		}

		// Look for codeVoice in inlineContent
		for _, inline := range inlineContent {
			inlineMap, ok := inline.(map[string]interface{})
			if !ok {
				continue
			}

			inlineType, ok := inlineMap["type"].(string)
			if !ok || inlineType != "codeVoice" {
				continue
			}

			code, ok := inlineMap["code"].(string)
			if ok && strings.HasPrefix(code, "com.apple.") {
				return code
			}
		}
	}

	return ""
}

// extractPropertiesFromJSON extracts property definitions from the properties section
func (s *payloadPropertyDefinitionServiceImpl) extractPropertiesFromJSON(jsonDoc map[string]interface{}) []*ent.PayloadPropertyDefinition {
	properties := []*ent.PayloadPropertyDefinition{}

	sections, ok := jsonDoc["primaryContentSections"].([]interface{})
	if !ok {
		return properties
	}

	// Find the properties section
	for _, section := range sections {
		sectionMap, ok := section.(map[string]interface{})
		if !ok {
			continue
		}

		kind, ok := sectionMap["kind"].(string)
		if !ok || kind != "properties" {
			continue
		}

		items, ok := sectionMap["items"].([]interface{})
		if !ok {
			continue
		}

		// Process each property item
		for _, item := range items {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			prop := s.buildPropertyDefinition(itemMap)
			if prop != nil {
				properties = append(properties, prop)
			}
		}
	}

	return properties
}

// buildPropertyDefinition builds a PayloadPropertyDefinition from a property item
func (s *payloadPropertyDefinitionServiceImpl) buildPropertyDefinition(itemMap map[string]interface{}) *ent.PayloadPropertyDefinition {
	// Extract property name (key)
	name, ok := itemMap["name"].(string)
	if !ok || name == "" {
		return nil
	}

	prop := &ent.PayloadPropertyDefinition{
		Key: strings.TrimSpace(name),
	}

	// Extract required flag
	required := false
	if reqVal, ok := itemMap["required"].(bool); ok {
		required = reqVal
	}

	// Extract deprecated flag
	deprecated := false
	if deprecatedVal, ok := itemMap["deprecated"].(bool); ok {
		deprecated = deprecatedVal
	}

	prop.ValueType = s.extractValueType(itemMap)
	prop.Deprecated = deprecated
	prop.DefaultValue = s.extractDefaultValue(itemMap)
	prop.EnumValues = s.extractEnumValues(itemMap)

	// Extract description from modern Apple docs format (content), fallback to legacy (discussion).
	prop.Description = s.extractDescription(itemMap)

	// Add required flag to description if not empty
	if required {
		if prop.Description != "" {
			prop.Description = fmt.Sprintf("[Required] %s", prop.Description)
		} else {
			prop.Description = "[Required property]"
		}
	}

	return prop
}

func (s *payloadPropertyDefinitionServiceImpl) extractValueType(itemMap map[string]interface{}) string {
	typeInfo := ""
	typeVal, ok := itemMap["type"].([]interface{})
	if !ok || len(typeVal) == 0 {
		return "string"
	}

	first := typeVal[0]
	if typeStr, ok := first.(string); ok {
		typeInfo = strings.TrimSpace(typeStr)
	}
	if typeMap, ok := first.(map[string]interface{}); ok {
		if text, ok := typeMap["text"].(string); ok {
			typeInfo = strings.TrimSpace(text)
		}
	}

	if typeInfo == "" {
		return "string"
	}
	return typeInfo
}

func (s *payloadPropertyDefinitionServiceImpl) extractDefaultValue(itemMap map[string]interface{}) map[string]interface{} {
	attributes, ok := itemMap["attributes"].([]interface{})
	if !ok {
		return map[string]interface{}{}
	}

	for _, attr := range attributes {
		attrMap, ok := attr.(map[string]interface{})
		if !ok {
			continue
		}
		kind, ok := attrMap["kind"].(string)
		if !ok || kind != "default" {
			continue
		}

		if raw, ok := attrMap["value"].(string); ok {
			return map[string]interface{}{"value": parseScalarString(raw)}
		}
		if raw, ok := attrMap["value"]; ok {
			return map[string]interface{}{"value": raw}
		}
	}

	return map[string]interface{}{}
}

func (s *payloadPropertyDefinitionServiceImpl) extractEnumValues(itemMap map[string]interface{}) []interface{} {
	attributes, ok := itemMap["attributes"].([]interface{})
	if !ok {
		return []interface{}{}
	}

	for _, attr := range attributes {
		attrMap, ok := attr.(map[string]interface{})
		if !ok {
			continue
		}
		kind, ok := attrMap["kind"].(string)
		if !ok || kind != "allowedValues" {
			continue
		}

		values, ok := attrMap["values"].([]interface{})
		if !ok {
			continue
		}

		parsed := make([]interface{}, 0, len(values))
		for _, v := range values {
			s, ok := v.(string)
			if ok {
				parsed = append(parsed, parseScalarString(s))
				continue
			}
			parsed = append(parsed, v)
		}
		return parsed
	}

	return []interface{}{}
}

func (s *payloadPropertyDefinitionServiceImpl) extractDescription(itemMap map[string]interface{}) string {
	content, ok := itemMap["content"].([]interface{})
	if ok {
		var parts []string
		for _, block := range content {
			blockMap, ok := block.(map[string]interface{})
			if !ok {
				continue
			}
			blockType, _ := blockMap["type"].(string)
			if blockType == "paragraph" {
				text := extractInlineText(blockMap["inlineContent"])
				if strings.TrimSpace(text) != "" {
					parts = append(parts, text)
				}
			}
		}
		if len(parts) > 0 {
			return strings.TrimSpace(strings.Join(parts, "\n"))
		}
	}

	if discussion, ok := itemMap["discussion"].([]interface{}); ok && len(discussion) > 0 {
		if discMap, ok := discussion[0].(map[string]interface{}); ok {
			text := extractInlineText(discMap["inlineContent"])
			if strings.TrimSpace(text) != "" {
				return strings.TrimSpace(text)
			}
		}
	}

	return ""
}

func extractInlineText(raw interface{}) string {
	inlineContent, ok := raw.([]interface{})
	if !ok {
		return ""
	}
	var b strings.Builder
	for _, inline := range inlineContent {
		inlineMap, ok := inline.(map[string]interface{})
		if !ok {
			continue
		}
		if text, ok := inlineMap["text"].(string); ok {
			b.WriteString(text)
			continue
		}
		if code, ok := inlineMap["code"].(string); ok {
			b.WriteString(code)
		}
	}
	return strings.TrimSpace(b.String())
}

func parseScalarString(raw string) interface{} {
	v := strings.TrimSpace(raw)
	if v == "" {
		return ""
	}

	lower := strings.ToLower(v)
	if lower == "true" {
		return true
	}
	if lower == "false" {
		return false
	}

	if i, err := strconv.ParseInt(v, 10, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(v, 64); err == nil {
		return f
	}

	return v
}

// Validation functions

func validateCreatePayloadPropertyDefinitionCommand(cmd service.CreatePayloadPropertyDefinitionCommand) error {
	if strings.TrimSpace(cmd.PayloadType) == "" {
		return apperror.ErrValidation.WithMessage("payload_type là bắt buộc")
	}
	if strings.TrimSpace(cmd.Key) == "" {
		return apperror.ErrValidation.WithMessage("key là bắt buộc")
	}
	if strings.TrimSpace(cmd.ValueType) == "" {
		return apperror.ErrValidation.WithMessage("value_type là bắt buộc")
	}
	return nil
}

func validateUpdatePayloadPropertyDefinitionCommand(cmd service.UpdatePayloadPropertyDefinitionCommand) error {
	if strings.TrimSpace(cmd.PayloadType) == "" {
		return apperror.ErrValidation.WithMessage("payload_type là bắt buộc")
	}
	if strings.TrimSpace(cmd.Key) == "" {
		return apperror.ErrValidation.WithMessage("key là bắt buộc")
	}
	if strings.TrimSpace(cmd.ValueType) == "" {
		return apperror.ErrValidation.WithMessage("value_type là bắt buộc")
	}
	return nil
}
