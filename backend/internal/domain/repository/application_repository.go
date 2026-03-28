package repository

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/pkg/query"
)

type ApplicationRepository interface {
	List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.Application, int64, error)
	GetByID(ctx context.Context, id uint) (*ent.Application, error)
	Create(ctx context.Context, entity *ent.Application) (*ent.Application, error)
	Update(ctx context.Context, id uint, entity *ent.Application) (*ent.Application, error)
	Delete(ctx context.Context, id uint) error
	BundleIDExists(ctx context.Context, bundleID string) (bool, error)

	// AppVersion
	ListVersions(ctx context.Context, appID uint) ([]*ent.AppVersion, error)
	CreateVersion(ctx context.Context, entity *ent.AppVersion) (*ent.AppVersion, error)
	DeleteVersion(ctx context.Context, id uint) error
	AppExists(ctx context.Context, id uint) (bool, error)

	// AppDeployment
	CreateDeployment(ctx context.Context, entity *ent.AppDeployment) (*ent.AppDeployment, error)
	ListDeployments(ctx context.Context, appVersionID uint) ([]*ent.AppDeployment, error)
	AppVersionExists(ctx context.Context, id uint) (bool, error)
}
