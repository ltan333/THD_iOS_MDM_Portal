package service

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/pkg/query"
)

type CreatePayloadPropertyDefinitionCommand struct {
	PayloadType  string
	Key          string
	ValueType    string
	DefaultValue map[string]interface{}
	EnumValues   []interface{}
	Deprecated   bool
	Description  string
}

type UpdatePayloadPropertyDefinitionCommand struct {
	ID           uint
	PayloadType  string
	Key          string
	ValueType    string
	DefaultValue map[string]interface{}
	EnumValues   []interface{}
	Deprecated   bool
	Description  string
}

type ImportPayloadPropertyDefinitionsResult struct {
	PayloadType string   `json:"payload_type"`
	Total       int      `json:"total"`
	Created     int      `json:"created"`
	Updated     int      `json:"updated"`
	Errors      []string `json:"errors,omitempty"`
}

// PayloadPropertyDefinitionService defines the payload property definition service interface.
type PayloadPropertyDefinitionService interface {
	Create(ctx context.Context, cmd CreatePayloadPropertyDefinitionCommand) (*ent.PayloadPropertyDefinition, error)
	GetByID(ctx context.Context, id uint) (*ent.PayloadPropertyDefinition, error)
	Update(ctx context.Context, cmd UpdatePayloadPropertyDefinitionCommand) (*ent.PayloadPropertyDefinition, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.PayloadPropertyDefinition, int64, error)

	ImportFromAppleJSON(ctx context.Context, filename string, data []byte) (*ImportPayloadPropertyDefinitionsResult, error)
}
