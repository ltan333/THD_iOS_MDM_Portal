package repository

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	"github.com/thienel/go-backend-template/pkg/query"
)

type ProfileRepository interface {
	List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.Profile, int64, error)
	GetByID(ctx context.Context, id uint) (*ent.Profile, error)
	Create(ctx context.Context, entity *ent.Profile, deviceGroupIDs []uint) (*ent.Profile, error)
	Update(ctx context.Context, id uint, entity *ent.Profile, deviceGroupIDs []uint) (*ent.Profile, error)
	Delete(ctx context.Context, id uint) error
	UpdateStatus(ctx context.Context, id uint, status string) error
	SaveVersion(ctx context.Context, profileID uint, version int, data map[string]any, changeNotes string) error

	// Assignments
	Assign(ctx context.Context, cmd service.AssignProfileCommand) (*ent.ProfileAssignment, error)
	Unassign(ctx context.Context, profileID uint, assignmentID uint) error
	ListAssignments(ctx context.Context, profileID uint) ([]*ent.ProfileAssignment, error)

	// Versions
	ListVersions(ctx context.Context, profileID uint) ([]*ent.ProfileVersion, error)
	Rollback(ctx context.Context, profileID uint, versionID uint) error

	// Deployment Status
	CreateDeploymentStatus(ctx context.Context, profileID uint, deviceID string, status string) (*ent.ProfileDeploymentStatus, error)
	UpdateDeploymentStatus(ctx context.Context, id uint, status string, errorMessage string) error
	GetDeploymentStatus(ctx context.Context, profileID uint) ([]*ent.ProfileDeploymentStatus, error)
	GetProfilesByDevice(ctx context.Context, deviceID string) ([]*ent.Profile, error)
	GetFlattenedDeviceUDIDsByProfile(ctx context.Context, profileID uint) ([]string, error)
}
