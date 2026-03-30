package repository

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/pkg/query"
)

type UniqueFieldConflict struct {
	Field string
	Value string
}

type MobileConfigRepository interface {
	List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.MobileConfig, int64, error)
	GetByID(ctx context.Context, id uint) (*ent.MobileConfig, error)
	GetByIDWithPayloads(ctx context.Context, id uint) (*ent.MobileConfig, error)
	GetFullForExport(ctx context.Context, id uint) (*ent.MobileConfig, error)
	List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.MobileConfig, int64, error)
	Create(ctx context.Context, entity *ent.MobileConfig, payload []*ent.Payload) (*ent.MobileConfig, error)
	Update(ctx context.Context, id uint, entity *ent.MobileConfig, payload []*ent.Payload) (*ent.MobileConfig, error)
	Delete(ctx context.Context, id uint) error
	FindCreateUniqueFieldConflict(ctx context.Context, name string, payloadIdentifier string, payloadIdentifiers []string) (*UniqueFieldConflict, error)
	FindUpdateUniqueFieldConflict(ctx context.Context, id uint, name string, payloadIdentifier string, payloadIdentifiers []string) (*UniqueFieldConflict, error)
}
