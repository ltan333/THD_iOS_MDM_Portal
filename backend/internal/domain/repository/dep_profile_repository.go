package repository

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
)

// DepProfileRepository defines the interface for DEP profile persistence operations.
type DepProfileRepository interface {
	// Create creates a new DEP profile in the database.
	Create(ctx context.Context, profile *ent.DepProfile) (*ent.DepProfile, error)

	// GetByID retrieves a DEP profile by its ID.
	GetByID(ctx context.Context, id uint) (*ent.DepProfile, error)

	// GetByProfileUUID retrieves a DEP profile by its Apple profile UUID.
	GetByProfileUUID(ctx context.Context, profileUUID string) (*ent.DepProfile, error)

	// GetByName retrieves a DEP profile by its profile name.
	GetByName(ctx context.Context, profileName string) (*ent.DepProfile, error)

	// List returns all DEP profiles with pagination.
	List(ctx context.Context, offset, limit int) ([]*ent.DepProfile, int64, error)

	// Update updates an existing DEP profile.
	Update(ctx context.Context, id uint, profile *ent.DepProfile) (*ent.DepProfile, error)

	// Delete removes a DEP profile from the database.
	Delete(ctx context.Context, id uint) error

	// SetProfileUUID sets the Apple profile UUID for a profile (after creating in Apple DEP).
	SetProfileUUID(ctx context.Context, id uint, profileUUID string) error
}
