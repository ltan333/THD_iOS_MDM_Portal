package serviceimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/thienel/tlog"
	"go.uber.org/zap"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	"github.com/thienel/go-backend-template/pkg/query"
)

// isProfileSpecificFile reports whether a filename is the profile-specific-payload-keys aggregate file.
// This file is an aggregate index and is never a payload definition document.
func isProfileSpecificFile(name string) bool {
	return strings.Contains(strings.ToLower(name), "profile-specific-payload-keys")
}

// normalizeFileName returns a lowercase, base-name-only key used as the key in a fileMap.
func normalizeFileName(name string) string {
	return strings.ToLower(filepath.Base(name))
}

type payloadPropertyDefinitionServiceImpl struct {
	payloadPropertyDefinitionRepo repository.PayloadPropertyDefinitionRepository
}

// NewPayloadPropertyDefinitionService creates a new payload property definition service
func NewPayloadPropertyDefinitionService(repo repository.PayloadPropertyDefinitionRepository) service.PayloadPropertyDefinitionService {
	return &payloadPropertyDefinitionServiceImpl{
		payloadPropertyDefinitionRepo: repo,
	}
}

func (s *payloadPropertyDefinitionServiceImpl) ListPayloadTypes(ctx context.Context) ([]string, error) {
	return s.payloadPropertyDefinitionRepo.ListPayloadTypes(ctx)
}

// ImportFromAppleJSONFiles imports from multiple uploaded files (in-memory).
// fileMap keys are normalized (lowercase) filenames, values are file contents.
//
// Classification is content-based (per HANDLING_NESTED_PAYLOAD_PROPERTIES v1.1):
//   - Top-level payload doc: has 'com.apple.' payload type in its Discussion section.
//   - Nested dictionary / array element: no payload type detected → skipped as top-level,
//     used on-demand via resolveNestedFileFromMap when a parent property references it.
//   - Aggregate file (profile-specific-payload-keys): explicitly skipped by filename.
func (s *payloadPropertyDefinitionServiceImpl) ImportFromAppleJSONFiles(ctx context.Context, fileMap map[string][]byte) (*service.ImportPayloadPropertyDefinitionsResult, error) {
	result := &service.ImportPayloadPropertyDefinitionsResult{
		PayloadType: "",
		Total:       0,
		Created:     0,
		Updated:     0,
		Errors:      []string{},
	}

	payloadTypes := []string{}

	for filename, data := range fileMap {
		// Only skip the known aggregate index file; everything else is classified by content.
		if isProfileSpecificFile(filename) {
			continue
		}

		r := s.importSingleDocFromMap(ctx, filename, data, fileMap)
		result.Total += r.Total
		result.Created += r.Created
		result.Updated += r.Updated
		result.Errors = append(result.Errors, r.Errors...)
		if r.PayloadType != "" {
			payloadTypes = append(payloadTypes, r.PayloadType)
		}
	}

	switch len(payloadTypes) {
	case 0:
		result.PayloadType = ""
	case 1:
		result.PayloadType = payloadTypes[0]
	default:
		result.PayloadType = strings.Join(payloadTypes, ", ")
	}

	return result, nil
}

// importSingleDocFromMap parses and imports one payload document using in-memory fileMap for nested resolution.
func (s *payloadPropertyDefinitionServiceImpl) importSingleDocFromMap(ctx context.Context, filename string, data []byte, fileMap map[string][]byte) *service.ImportPayloadPropertyDefinitionsResult {
	result := &service.ImportPayloadPropertyDefinitionsResult{Errors: []string{}}

	var jsonDoc map[string]interface{}
	if err := json.Unmarshal(data, &jsonDoc); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("[%s] Lỗi parse JSON: %v", filename, err))
		return result
	}

	payloadType := s.detectPayloadTypeFromDiscussion(jsonDoc)
	if payloadType == "" {
		// Not a payload doc (e.g. it might be a dictionary sub-object uploaded alone) — skip silently
		return result
	}

	result.PayloadType = payloadType

	visited := map[string]bool{}
	properties := s.extractPropertiesRecursivelyFromMap(jsonDoc, fileMap, "", visited)
	result.Total = len(properties)

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

	return result
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

// extractPropertiesRecursivelyFromMap extracts property definitions from the properties section,
// and recursively resolves nested dictionary types using the in-memory fileMap.
// keyPrefix is the dot-notation prefix for nested properties (e.g. "CommunicationServiceRules.").
// visited prevents infinite loops when circular references exist.
//
// Recursion triggers when a property has nested_reference set AND:
//   - value_type == "dictionary" (direct nested object), or
//   - value_type == "array"      (array of nested objects — imports element properties too)
func (s *payloadPropertyDefinitionServiceImpl) extractPropertiesRecursivelyFromMap(
	jsonDoc map[string]interface{},
	fileMap map[string][]byte,
	keyPrefix string,
	visited map[string]bool,
) []*ent.PayloadPropertyDefinition {
	properties := []*ent.PayloadPropertyDefinition{}

	sections, ok := jsonDoc["primaryContentSections"].([]interface{})
	if !ok {
		return properties
	}

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

		for i, item := range items {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			prop := s.buildPropertyDefinitionWithPrefix(itemMap, keyPrefix, i)
			if prop == nil {
				continue
			}
			properties = append(properties, prop)

			// Recurse into nested reference when:
			//   - value_type == "dictionary" → nested_reference holds the identifier
			//   - value_type == "array"      → items_reference holds the identifier (items must be dictionaries)
			var nestedRef string
			switch prop.ValueType {
			case "dictionary":
				if prop.NestedReference == nil {
					continue
				}
				nestedRef = *prop.NestedReference
			case "array":
				if prop.ItemsReference == nil {
					continue
				}
				nestedRef = *prop.ItemsReference
			default:
				continue
			}

			if visited[nestedRef] {
				continue
			}
			visited[nestedRef] = true

			nestedData, found := s.resolveNestedFileFromMap(fileMap, itemMap)
			if !found {
				continue
			}

			var nestedDoc map[string]interface{}
			if err := json.Unmarshal(nestedData, &nestedDoc); err != nil {
				tlog.Warn("Không thể parse nested JSON",
					zap.String("nested_ref", nestedRef),
					zap.Error(err),
				)
				continue
			}

			nestedPrefix := prop.Key + "."
			nestedProps := s.extractPropertiesRecursivelyFromMap(nestedDoc, fileMap, nestedPrefix, visited)
			properties = append(properties, nestedProps...)
		}
	}

	return properties
}

// resolveNestedFileFromMap looks up a nested dictionary file in the in-memory fileMap.
// It extracts the last segment of the identifier URL and lowercases it to form the lookup key.
//
// Example:
//
//	identifier = "doc://com.apple.devicemanagement/.../CardDAV/CommunicationServiceRules-data.dictionary"
//	lookup key = "communicationservicerules-data.dictionary.json"
func (s *payloadPropertyDefinitionServiceImpl) resolveNestedFileFromMap(fileMap map[string][]byte, itemMap map[string]interface{}) ([]byte, bool) {
	typeVal, ok := itemMap["type"].([]interface{})
	if !ok {
		return nil, false
	}

	for _, t := range typeVal {
		typeMap, ok := t.(map[string]interface{})
		if !ok {
			continue
		}
		kind, _ := typeMap["kind"].(string)
		if kind != "typeIdentifier" {
			continue
		}

		identifier, _ := typeMap["identifier"].(string)
		if identifier == "" {
			continue
		}

		// Extract the last path segment from the identifier URL and build a filename
		// e.g. ".../CardDAV/CommunicationServiceRules-data.dictionary" → "communicationservicerules-data.dictionary.json"
		parts := strings.Split(identifier, "/")
		if len(parts) == 0 {
			continue
		}
		lastSegment := strings.ToLower(parts[len(parts)-1])
		if !strings.HasSuffix(lastSegment, ".json") {
			lastSegment += ".json"
		}

		if data, ok := fileMap[lastSegment]; ok {
			return data, true
		}
	}

	return nil, false
}

// buildPropertyDefinitionWithPrefix builds a PayloadPropertyDefinition from a property item,
// prepending keyPrefix (dot-notation) to the property key.
// orderIndex is the 0-based position of this item in the properties list.
func (s *payloadPropertyDefinitionServiceImpl) buildPropertyDefinitionWithPrefix(itemMap map[string]interface{}, keyPrefix string, orderIndex int) *ent.PayloadPropertyDefinition {
	// Extract property name (key)
	name, ok := itemMap["name"].(string)
	if !ok || name == "" {
		return nil
	}

	prop := &ent.PayloadPropertyDefinition{
		Key:        keyPrefix + strings.TrimSpace(name),
		IsNested:   keyPrefix != "",
		OrderIndex: orderIndex,
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
	prop.Description = s.extractDescription(itemMap)

	// Populate reference fields based on value_type
	switch prop.ValueType {
	case "dictionary":
		// nested_reference: identifier of the nested dictionary
		prop.NestedReference = s.extractNestedReference(itemMap)
	case "array":
		// items_type: type of array elements; items_reference: identifier if elements are dictionaries
		prop.ItemsType = s.extractItemsType(itemMap)
		if prop.ItemsType != nil && *prop.ItemsType == "dictionary" {
			prop.ItemsReference = s.extractItemsReference(itemMap)
		}
	}

	// Prepend [Required] to description
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
	typeVal, ok := itemMap["type"].([]interface{})
	if !ok || len(typeVal) == 0 {
		return "string"
	}

	// Scan all tokens to characterize the type:
	//   - bracketOpen / bracketClose: text tokens "[" and "]"
	//   - hasTypeIdentifier: any token with kind == "typeIdentifier" (reference to another object)
	//   - primitiveType: last non-bracket text token value
	var primitiveType string
	hasTypeIdentifier := false
	bracketOpen := false
	bracketClose := false

	for _, t := range typeVal {
		typeMap, ok := t.(map[string]interface{})
		if !ok {
			continue
		}
		kind, _ := typeMap["kind"].(string)
		switch kind {
		case "typeIdentifier":
			hasTypeIdentifier = true
		case "text":
			text, _ := typeMap["text"].(string)
			text = strings.TrimSpace(text)
			switch text {
			case "[":
				bracketOpen = true
			case "]":
				bracketClose = true
			default:
				if text != "" {
					primitiveType = text
				}
			}
		}
	}

	switch {
	case hasTypeIdentifier && bracketOpen && bracketClose:
		// Pattern: [ TypeIdentifier ] → array whose elements are a nested object
		return "array"
	case hasTypeIdentifier:
		// Direct reference to a nested object → dictionary
		return "dictionary"
	case primitiveType != "":
		return primitiveType
	default:
		return "string"
	}
}

// extractNestedReference extracts the preciseIdentifier from a typeIdentifier type entry.
// Only returns a value when value_type is "dictionary" (not array — use extractItemsReference for that).
func (s *payloadPropertyDefinitionServiceImpl) extractNestedReference(itemMap map[string]interface{}) *string {
	// Only applicable for direct dictionary references (no bracket pattern)
	if s.extractValueType(itemMap) != "dictionary" {
		return nil
	}
	typeVal, ok := itemMap["type"].([]interface{})
	if !ok {
		return nil
	}
	for _, t := range typeVal {
		typeMap, ok := t.(map[string]interface{})
		if !ok {
			continue
		}
		kind, _ := typeMap["kind"].(string)
		if kind != "typeIdentifier" {
			continue
		}
		if preciseID, ok := typeMap["preciseIdentifier"].(string); ok && preciseID != "" {
			v := preciseID
			return &v
		}
	}
	return nil
}

// extractItemsType returns the type of elements in an array property.
// Returns nil if the property is not an array.
func (s *payloadPropertyDefinitionServiceImpl) extractItemsType(itemMap map[string]interface{}) *string {
	if s.extractValueType(itemMap) != "array" {
		return nil
	}
	typeVal, ok := itemMap["type"].([]interface{})
	if !ok {
		return nil
	}
	// Inside bracket pattern [ ... ], determine element type
	for _, t := range typeVal {
		typeMap, ok := t.(map[string]interface{})
		if !ok {
			continue
		}
		kind, _ := typeMap["kind"].(string)
		switch kind {
		case "typeIdentifier":
			v := "dictionary"
			return &v
		case "text":
			text := strings.TrimSpace(typeMap["text"].(string))
			if text != "" && text != "[" && text != "]" {
				v := text
				return &v
			}
		}
	}
	v := "string"
	return &v
}

// extractItemsReference extracts the preciseIdentifier of an array's element dictionary type.
// Returns nil if elements are not a typeIdentifier (i.e. primitive array).
func (s *payloadPropertyDefinitionServiceImpl) extractItemsReference(itemMap map[string]interface{}) *string {
	typeVal, ok := itemMap["type"].([]interface{})
	if !ok {
		return nil
	}
	for _, t := range typeVal {
		typeMap, ok := t.(map[string]interface{})
		if !ok {
			continue
		}
		kind, _ := typeMap["kind"].(string)
		if kind != "typeIdentifier" {
			continue
		}
		if preciseID, ok := typeMap["preciseIdentifier"].(string); ok && preciseID != "" {
			v := preciseID
			return &v
		}
	}
	return nil
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

// GetNestedSchema returns nested schemas. If payloadType is empty, returns schemas for ALL payload types.
// Results are sorted by payload_type then order_index within each payload type.
func (s *payloadPropertyDefinitionServiceImpl) GetNestedSchema(ctx context.Context, payloadType string) ([]*service.NestedPayloadSchema, error) {
	opts := query.NewQueryOptions()
	if strings.TrimSpace(payloadType) != "" {
		opts.AddFilter("payload_type", "eq", payloadType)
	}
	// Sort by payload_type first so groups are contiguous, then by order_index for correctness within group.
	opts.AddSort("payload_type", false)
	opts.AddSort("order_index", false)

	flatList, _, err := s.payloadPropertyDefinitionRepo.List(ctx, 0, 100000, opts)
	if err != nil {
		return nil, err
	}

	// Group by payload_type preserving insertion order
	typeOrder := []string{}
	groups := map[string][]*ent.PayloadPropertyDefinition{}
	for i := range flatList {
		pt := flatList[i].PayloadType
		if _, exists := groups[pt]; !exists {
			typeOrder = append(typeOrder, pt)
			groups[pt] = []*ent.PayloadPropertyDefinition{}
		}
		groups[pt] = append(groups[pt], &flatList[i])
	}

	schemas := make([]*service.NestedPayloadSchema, 0, len(typeOrder))
	for _, pt := range typeOrder {
		schemas = append(schemas, &service.NestedPayloadSchema{
			PayloadType: pt,
			Properties:  buildNestedSchema(groups[pt]),
		})
	}

	return schemas, nil
}

// buildNestedSchema converts a flat, dot-notation list of properties into a nested map.
//
// Two-pass approach:
//
//	Pass 1 – build a lookup table: fullKey → *NestedProperty (populated from DB).
//	Pass 2 – for each top-level key (depth 0), recursively attach children into the
//	          correct sub-map (Properties for dictionary, ItemsSchema for array).
//
// This guarantees that parent ValueType is always known before we try to route children,
// even when the flat list is not in any particular depth order.
func buildNestedSchema(flat []*ent.PayloadPropertyDefinition) map[string]*service.NestedProperty {
	// Pass 1: build the node map keyed by full dot-notation key.
	nodeMap := make(map[string]*service.NestedProperty, len(flat))
	for _, p := range flat {
		nodeMap[p.Key] = &service.NestedProperty{
			ValueType:       p.ValueType,
			Description:     p.Description,
			DefaultValue:    p.DefaultValue,
			EnumValues:      p.EnumValues,
			Deprecated:      p.Deprecated,
			NestedReference: p.NestedReference,
			ItemsType:       p.ItemsType,
			ItemsReference:  p.ItemsReference,
		}
	}

	// Pass 2: wire up parent → children relationships.
	for _, p := range flat {
		parts := strings.Split(p.Key, ".")
		if len(parts) < 2 {
			continue // top-level, nothing to wire
		}
		parentKey := strings.Join(parts[:len(parts)-1], ".")
		childKey := parts[len(parts)-1]

		parent, ok := nodeMap[parentKey]
		if !ok {
			// Parent missing from DB — create a placeholder so the child is not lost.
			parent = &service.NestedProperty{ValueType: "dictionary"}
			nodeMap[parentKey] = parent
		}

		childNode := nodeMap[p.Key]

		if parent.ValueType == "array" {
			if parent.ItemsSchema == nil {
				parent.ItemsSchema = map[string]*service.NestedProperty{}
			}
			parent.ItemsSchema[childKey] = childNode
		} else {
			// dictionary or unknown → put in Properties
			if parent.Properties == nil {
				parent.Properties = map[string]*service.NestedProperty{}
			}
			parent.Properties[childKey] = childNode
		}
	}

	// Collect root-level nodes (no dot in key).
	root := map[string]*service.NestedProperty{}
	for _, p := range flat {
		if !strings.Contains(p.Key, ".") {
			root[p.Key] = nodeMap[p.Key]
		}
	}
	return root
}
