package dto

import "time"

type CreatePayloadPropertyDefinitionRequest struct {
	PayloadType  string                 `json:"payload_type" binding:"required"`
	Key          string                 `json:"key" binding:"required"`
	ValueType    string                 `json:"value_type" binding:"required"`
	DefaultValue map[string]interface{} `json:"default_value,omitempty"`
	EnumValues   []interface{}          `json:"enum_values,omitempty"`
	Deprecated   bool                   `json:"deprecated"`
	Description  string                 `json:"description,omitempty"`
}

type UpdatePayloadPropertyDefinitionRequest = CreatePayloadPropertyDefinitionRequest

type ImportDirectoryRequest struct {
	DocsDir string `json:"docs_dir" binding:"required"`
}

type PayloadPropertyDefinitionResponse struct {
	ID              int                    `json:"id"`
	PayloadType     string                 `json:"payload_type"`
	Key             string                 `json:"key"`
	ValueType       string                 `json:"value_type"`
	DefaultValue    map[string]interface{} `json:"default_value,omitempty"`
	EnumValues      []interface{}          `json:"enum_values,omitempty"`
	Deprecated      bool                   `json:"deprecated"`
	Description     string                 `json:"description,omitempty"`
	NestedReference *string                `json:"nested_reference,omitempty"`
	ItemsType       *string                `json:"items_type,omitempty"`
	ItemsReference  *string                `json:"items_reference,omitempty"`
	IsNested        bool                   `json:"is_nested"`
	OrderIndex      int                    `json:"order_index"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}
