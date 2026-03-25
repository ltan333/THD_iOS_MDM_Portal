package service

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/pkg/query"
)

type CreateApplicationCommand struct {
	Name        string
	BundleID    string
	Platform    string
	Type        string
	Description string
	IconURL     string
}

type UpdateApplicationCommand struct {
	ID          uint
	Name        *string
	Platform    *string
	Type        *string
	Description *string
	IconURL     *string
}

type CreateAppVersionCommand struct {
	ApplicationID    uint
	Version          string
	BuildNumber      string
	MinimumOSVersion string
	FileURL          string
	Size             int64
	Metadata         map[string]interface{}
}

type CreateAppDeploymentCommand struct {
	AppVersionID uint
	TargetType   string
	TargetID     string
}

type ApplicationService interface {
	// Applications
	List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.Application, int64, error)
	GetByID(ctx context.Context, id uint) (*ent.Application, error)
	Create(ctx context.Context, cmd CreateApplicationCommand) (*ent.Application, error)
	Update(ctx context.Context, cmd UpdateApplicationCommand) (*ent.Application, error)
	Delete(ctx context.Context, id uint) error

	// App Versions
	ListVersions(ctx context.Context, appID uint) ([]*ent.AppVersion, error)
	CreateVersion(ctx context.Context, cmd CreateAppVersionCommand) (*ent.AppVersion, error)
	DeleteVersion(ctx context.Context, id uint) error

	// App Deployments
	Deploy(ctx context.Context, cmd CreateAppDeploymentCommand) (*ent.AppDeployment, error)
	ListDeployments(ctx context.Context, appVersionID uint) ([]*ent.AppDeployment, error)
}
