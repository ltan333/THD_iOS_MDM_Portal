package repository

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/pkg/query"
)

// PayloadPropertyDefinitionRepository extends BaseRepository for PayloadPropertyDefinition entity.
type PayloadPropertyDefinitionRepository interface {
	BaseRepository[ent.PayloadPropertyDefinition]

	ListWithQuery(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.PayloadPropertyDefinition, int64, error)
	ListPayloadTypes(ctx context.Context) ([]string, error)
	DeleteAll(ctx context.Context) (int, error)
	FindByPayloadTypeVariantAndKey(ctx context.Context, payloadType, payloadVariant, key string) (*ent.PayloadPropertyDefinition, error)
	UpsertByPayloadTypeVariantAndKey(ctx context.Context, entity *ent.PayloadPropertyDefinition) (*ent.PayloadPropertyDefinition, bool, error)
}
