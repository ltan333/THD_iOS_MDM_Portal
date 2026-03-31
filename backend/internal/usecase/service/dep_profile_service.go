package service

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
)

// DepProfileService handles DEP profile operations.
type DepProfileService interface {
	// Create creates a new DEP profile locally and registers it with Apple DEP.
	Create(ctx context.Context, depName string, req *dto.DEPProfileRequest) (*ent.DepProfile, error)

	// GetByID retrieves a DEP profile by its ID.
	GetByID(ctx context.Context, id uint) (*ent.DepProfile, error)

	// GetByProfileUUID retrieves a DEP profile by its Apple profile UUID.
	GetByProfileUUID(ctx context.Context, profileUUID string) (*ent.DepProfile, error)

	// List returns all DEP profiles with pagination.
	List(ctx context.Context, offset, limit int) ([]*ent.DepProfile, int64, error)

	// Update updates an existing DEP profile and syncs with Apple DEP.
	Update(ctx context.Context, depName string, id uint, req *dto.DEPProfileRequest) (*ent.DepProfile, error)

	// Delete removes a DEP profile from local DB and optionally from Apple DEP.
	Delete(ctx context.Context, depName string, id uint) error

	// SetAsAssigner sets a profile as the default assigner profile for new devices.
	SetAsAssigner(ctx context.Context, depName string, id uint) error

	// GetAssigner gets the current assigner profile.
	GetAssigner(ctx context.Context, depName string) (*ent.DepProfile, error)
}
