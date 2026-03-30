package service

import (
	"context"
)

type ImportPayloadPropertyDefinitionsResult struct {
	PayloadType string   `json:"payload_type"`
	Total       int      `json:"total"`
	Created     int      `json:"created"`
	Updated     int      `json:"updated"`
	Errors      []string `json:"errors,omitempty"`
}

type NestedProperty struct {
	ValueType       string                     `json:"value_type"`
	Description     string                     `json:"description,omitempty"`
	DefaultValue    map[string]interface{}     `json:"default_value,omitempty"`
	EnumValues      []interface{}              `json:"enum_values,omitempty"`
	Deprecated      bool                       `json:"deprecated,omitempty"`
	NestedReference *string                    `json:"nested_reference,omitempty"`
	ItemsType       *string                    `json:"items_type,omitempty"`
	ItemsReference  *string                    `json:"items_reference,omitempty"`
	ItemsSchema     map[string]*NestedProperty `json:"items_schema,omitempty"`
	Properties      map[string]*NestedProperty `json:"properties,omitempty"`
}

type NestedPayloadSchema struct {
	PayloadType string                     `json:"payload_type"`
	Properties  map[string]*NestedProperty `json:"properties"`
}

type PayloadPropertyDefinitionService interface {
	ListPayloadTypes(ctx context.Context) ([]string, error)
	ImportFromAppleJSONFiles(ctx context.Context, fileMap map[string][]byte) (*ImportPayloadPropertyDefinitionsResult, error)
	GetNestedSchema(ctx context.Context, payloadType string) ([]*NestedPayloadSchema, error)
}
