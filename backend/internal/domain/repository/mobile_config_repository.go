package repository

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
)

type MobileConfigRepository interface {
	BaseRepository[ent.MobileConfig]

	FindByID(ctx context.Context, id uint) (*ent.MobileConfig, error)
	GetFullForExport(ctx context.Context, id uint) (*ent.MobileConfig, error)
}
