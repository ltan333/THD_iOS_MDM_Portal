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
	FindByPayloadTypeAndKey(ctx context.Context, payloadType, key string) (*ent.PayloadPropertyDefinition, error)
	UpsertByPayloadTypeAndKey(ctx context.Context, entity *ent.PayloadPropertyDefinition) (*ent.PayloadPropertyDefinition, bool, error)
}
