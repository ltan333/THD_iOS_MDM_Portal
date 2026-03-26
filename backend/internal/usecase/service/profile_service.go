package service

import (
	"context"
	"time"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/pkg/query"
)

type CreateProfileCommand struct {
	Name             string
	Platform         string
	Scope            string
	SecuritySettings map[string]any
	NetworkConfig    map[string]any
	Restrictions     map[string]any
	ContentFilter    map[string]any
	ComplianceRules  map[string]any
	Payloads         map[string]any
}

type UpdateProfileCommand struct {
	ID               uint
	Name             *string
	Platform         *string
	Scope            *string
	SecuritySettings map[string]any
	NetworkConfig    map[string]any
	Restrictions     map[string]any
	ContentFilter    map[string]any
	ComplianceRules  map[string]any
	Payloads         map[string]any
}

type AssignProfileCommand struct {
	ProfileID    uint
	TargetType   string // device, group
	DeviceID     *string
	GroupID      *uint
	ScheduleType string // immediate, scheduled
	ScheduledAt  *time.Time
}

type ProfileService interface {
	// CRUD
	List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.Profile, int64, error)
	GetByID(ctx context.Context, id uint) (*ent.Profile, error)
	Create(ctx context.Context, cmd CreateProfileCommand) (*ent.Profile, error)
	Update(ctx context.Context, cmd UpdateProfileCommand) (*ent.Profile, error)
	Delete(ctx context.Context, id uint) error

	// Status
	UpdateStatus(ctx context.Context, id uint, status string) error

	// Settings sections
	UpdateSecuritySettings(ctx context.Context, id uint, settings map[string]any) error
	UpdateNetworkConfig(ctx context.Context, id uint, config map[string]any) error
	UpdateRestrictions(ctx context.Context, id uint, restrictions map[string]any) error
	UpdateContentFilter(ctx context.Context, id uint, filter map[string]any) error
	UpdateComplianceRules(ctx context.Context, id uint, rules map[string]any) error

	// Assignment
	Assign(ctx context.Context, cmd AssignProfileCommand) error
	Unassign(ctx context.Context, profileID uint, assignmentID uint) error
	ListAssignments(ctx context.Context, profileID uint) ([]*ent.ProfileAssignment, error)

	// Version
	ListVersions(ctx context.Context, profileID uint) ([]*ent.ProfileVersion, error)
	Rollback(ctx context.Context, profileID uint, versionID uint) error

	// Deployment Status
	GetDeploymentStatus(ctx context.Context, profileID uint) ([]*ent.ProfileDeploymentStatus, error)
	Repush(ctx context.Context, profileID uint) error
	Duplicate(ctx context.Context, profileID uint) (*ent.Profile, error)
}
