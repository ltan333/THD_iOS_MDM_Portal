package repository

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
)

type UniqueFieldConflict struct {
	Field string
	Value string
}

type MobileConfigRepository interface {
	GetByID(ctx context.Context, id uint) (*ent.MobileConfig, error)
	GetFullForExport(ctx context.Context, id uint) (*ent.MobileConfig, error)
	Create(ctx context.Context, entity *ent.MobileConfig, payload []*ent.Payload) (*ent.MobileConfig, error)
	Update(ctx context.Context, id uint, entity *ent.MobileConfig, payload []*ent.Payload) (*ent.MobileConfig, error)
	Delete(ctx context.Context, id uint) error
	FindCreateUniqueFieldConflict(ctx context.Context, name string, payloadIdentifier string, payloadIdentifiers []string) (*UniqueFieldConflict, error)
	FindUpdateUniqueFieldConflict(ctx context.Context, id uint, name string, payloadIdentifier string, payloadIdentifiers []string) (*UniqueFieldConflict, error)
}
