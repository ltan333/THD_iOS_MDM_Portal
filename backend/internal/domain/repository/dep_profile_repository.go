package repository

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
)

// DepProfileRepository extends BaseRepository for DepProfile entity
type DepProfileRepository interface {
	BaseRepository[ent.DepProfile]

	// Additional DEP profile specific methods
	FindByProfileUUID(ctx context.Context, uuid string) (*ent.DepProfile, error)
	FindByProfileName(ctx context.Context, name string) (*ent.DepProfile, error)
}
