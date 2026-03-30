package service

import (
	"context"
)

type ImportPayloadPropertyDefinitionsResult struct {
	PayloadType       string   `json:"payload_type"`
	Total             int      `json:"total"`
	Created           int      `json:"created"`
	Updated           int      `json:"updated"`
	UpdatedProperties []string `json:"updated_properties,omitempty"`
	Errors            []string `json:"errors,omitempty"`
}

// NestedProperty represents a single property in the nested schema tree.
type NestedProperty struct {
	ValueType       string                     `json:"value_type"`
	Title           *string                    `json:"title,omitempty"`
	Description     string                     `json:"description,omitempty"`
	Presence        string                     `json:"presence,omitempty"`
	DefaultValue    map[string]interface{}     `json:"default_value,omitempty"`
	EnumValues      []interface{}              `json:"enum_values,omitempty"`
	Deprecated      bool                       `json:"deprecated,omitempty"`
	NestedReference *string                    `json:"nested_reference,omitempty"`
	ItemsType       *string                    `json:"items_type,omitempty"`
	ItemsReference  *string                    `json:"items_reference,omitempty"`
	SupportedOS     map[string]interface{}     `json:"supported_os,omitempty"`
	Conditions      map[string]interface{}     `json:"conditions,omitempty"`
	ItemsSchema     map[string]*NestedProperty `json:"items_schema,omitempty"`
	Properties      map[string]*NestedProperty `json:"properties,omitempty"`
}

type NestedPayloadSchema struct {
	PayloadType    string                     `json:"payload_type"`
	PayloadVariant string                     `json:"payload_variant,omitempty"`
	Properties     map[string]*NestedProperty `json:"properties"`
}

type PayloadVariantInfo struct {
	Variant       string `json:"variant,omitempty"`
	DisplayName   string `json:"display_name"`
	Description   string `json:"description,omitempty"`
	PropertyCount int    `json:"property_count"`
}

type PayloadPropertyDefinitionService interface {
	ListPayloadTypes(ctx context.Context) ([]string, error)
	ListVariants(ctx context.Context, payloadType string) ([]*PayloadVariantInfo, error)
	DeleteAll(ctx context.Context) (int, error)
	ImportFromAppleJSONFiles(ctx context.Context, fileMap map[string][]byte) (*ImportPayloadPropertyDefinitionsResult, error)
	ImportFromAppleYAMLFiles(ctx context.Context, fileMap map[string][]byte) (*ImportPayloadPropertyDefinitionsResult, error)
	GetNestedSchema(ctx context.Context, payloadType, payloadVariant string) ([]*NestedPayloadSchema, error)
}
