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
	SecuritySettings map[string]interface{}
	NetworkConfig    map[string]interface{}
	Restrictions     map[string]interface{}
	ContentFilter    map[string]interface{}
	ComplianceRules  map[string]interface{}
	Payloads         map[string]interface{}
}

type UpdateProfileCommand struct {
	ID               uint
	Name             *string
	Platform         *string
	Scope            *string
	SecuritySettings map[string]interface{}
	NetworkConfig    map[string]interface{}
	Restrictions     map[string]interface{}
	ContentFilter    map[string]interface{}
	ComplianceRules  map[string]interface{}
	Payloads         map[string]interface{}
}

type AssignProfileCommand struct {
	ProfileID    uint
	TargetType   string // device, group, user
	TargetID     string
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
	UpdateSecuritySettings(ctx context.Context, id uint, settings map[string]interface{}) error
	UpdateNetworkConfig(ctx context.Context, id uint, config map[string]interface{}) error
	UpdateRestrictions(ctx context.Context, id uint, restrictions map[string]interface{}) error
	UpdateContentFilter(ctx context.Context, id uint, filter map[string]interface{}) error
	UpdateComplianceRules(ctx context.Context, id uint, rules map[string]interface{}) error

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
