package persistence

import (
	"context"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/appdeployment"
	"github.com/thienel/go-backend-template/internal/ent/application"
	"github.com/thienel/go-backend-template/internal/ent/appversion"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
)

type applicationRepositoryImpl struct {
	client *ent.Client
}

func NewApplicationRepository(client *ent.Client) repository.ApplicationRepository {
	return &applicationRepositoryImpl{client: client}
}

func (r *applicationRepositoryImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.Application, int64, error) {
	q := r.client.Application.Query()

	// Apply search
	if searchFilter, ok := opts.Filters["search"]; ok {
		searchStr, _ := searchFilter.Value.(string)
		if searchStr != "" {
			q = q.Where(
				application.Or(
					application.NameContainsFold(searchStr),
					application.BundleIDContainsFold(searchStr),
				),
			)
		}
	}

	// Apply filters
	for field, val := range opts.Filters {
		switch field {
		case "platform":
			if strVal, ok := val.Value.(string); ok {
				q = q.Where(application.PlatformEQ(application.Platform(strVal)))
			}
		case "type":
			if strVal, ok := val.Value.(string); ok {
				q = q.Where(application.TypeEQ(application.Type(strVal)))
			}
		case "search":
			// handled above
		}
	}

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, apperror.ErrInternalServerError.WithMessage("Lỗi khi đếm số lượng ứng dụng").WithError(err)
	}

	apps, err := q.
		WithVersions(func(vq *ent.AppVersionQuery) {
			vq.Order(ent.Desc(appversion.FieldCreatedAt))
		}).
		Limit(limit).
		Offset(offset).
		Order(ent.Desc(application.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, 0, apperror.ErrInternalServerError.WithMessage("Lỗi khi lấy danh sách ứng dụng").WithError(err)
	}

	return apps, int64(total), nil
}

func (r *applicationRepositoryImpl) GetByID(ctx context.Context, id uint) (*ent.Application, error) {
	app, err := r.client.Application.Query().
		Where(application.IDEQ(id)).
		WithVersions(func(vq *ent.AppVersionQuery) {
			vq.Order(ent.Desc(appversion.FieldCreatedAt))
		}).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Không tìm thấy ứng dụng này")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy vấn ứng dụng").WithError(err)
	}
	return app, nil
}

func (r *applicationRepositoryImpl) Create(ctx context.Context, entity *ent.Application) (*ent.Application, error) {
	create := r.client.Application.Create().
		SetName(entity.Name).
		SetBundleID(entity.BundleID).
		SetPlatform(entity.Platform).
		SetType(entity.Type)

	if entity.Description != "" {
		create = create.SetDescription(entity.Description)
	}
	if entity.IconURL != "" {
		create = create.SetIconURL(entity.IconURL)
	}

	app, err := create.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, apperror.ErrConflict.WithMessage("Bundle ID đã tồn tại trong hệ thống")
		}
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo ứng dụng").WithError(err)
	}

	return app, nil
}

func (r *applicationRepositoryImpl) Update(ctx context.Context, id uint, entity *ent.Application) (*ent.Application, error) {
	update := r.client.Application.UpdateOneID(id)

	if entity.Name != "" {
		update.SetName(entity.Name)
	}
	if string(entity.Platform) != "" {
		update.SetPlatform(entity.Platform)
	}
	if string(entity.Type) != "" {
		update.SetType(entity.Type)
	}
	if entity.Description != "" {
		update.SetDescription(entity.Description)
	}
	if entity.IconURL != "" {
		update.SetIconURL(entity.IconURL)
	}

	app, err := update.Save(ctx)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Không tìm thấy ứng dụng để cập nhật")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi cập nhật ứng dụng").WithError(err)
	}

	return app, nil
}

func (r *applicationRepositoryImpl) Delete(ctx context.Context, id uint) error {
	err := r.client.Application.DeleteOneID(id).Exec(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Không tìm thấy ứng dụng để xóa")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi xóa ứng dụng").WithError(err)
	}
	return nil
}

func (r *applicationRepositoryImpl) BundleIDExists(ctx context.Context, bundleID string) (bool, error) {
	exists, err := r.client.Application.Query().
		Where(application.BundleIDEQ(bundleID)).
		Exist(ctx)
	if err != nil {
		return false, apperror.ErrInternalServerError.WithMessage("Lỗi khi kiểm tra Bundle ID").WithError(err)
	}
	return exists, nil
}

// App Version Methods

func (r *applicationRepositoryImpl) ListVersions(ctx context.Context, appID uint) ([]*ent.AppVersion, error) {
	versions, err := r.client.AppVersion.Query().
		Where(appversion.ApplicationIDEQ(appID)).
		Order(ent.Desc(appversion.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi lấy danh sách phiên bản ứng dụng").WithError(err)
	}
	return versions, nil
}

func (r *applicationRepositoryImpl) CreateVersion(ctx context.Context, entity *ent.AppVersion) (*ent.AppVersion, error) {
	create := r.client.AppVersion.Create().
		SetApplicationID(entity.ApplicationID).
		SetVersion(entity.Version).
		SetBuildNumber(entity.BuildNumber)

	if entity.MinimumOsVersion != "" {
		create = create.SetMinimumOsVersion(entity.MinimumOsVersion)
	}
	if entity.FileURL != "" {
		create = create.SetFileURL(entity.FileURL)
	}
	create = create.SetSize(entity.Size)
	if entity.Metadata != nil {
		create = create.SetMetadata(entity.Metadata)
	}

	version, err := create.Save(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo phiên bản ứng dụng").WithError(err)
	}

	return version, nil
}

func (r *applicationRepositoryImpl) DeleteVersion(ctx context.Context, id uint) error {
	err := r.client.AppVersion.DeleteOneID(id).Exec(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Không tìm thấy phiên bản để xóa")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi xóa phiên bản").WithError(err)
	}
	return nil
}

func (r *applicationRepositoryImpl) AppExists(ctx context.Context, id uint) (bool, error) {
	exists, err := r.client.Application.Query().Where(application.IDEQ(id)).Exist(ctx)
	if err != nil {
		return false, apperror.ErrInternalServerError.WithMessage("Lỗi khi kiểm tra ứng dụng gốc").WithError(err)
	}
	return exists, nil
}

// App Deployment Methods

func (r *applicationRepositoryImpl) CreateDeployment(ctx context.Context, entity *ent.AppDeployment) (*ent.AppDeployment, error) {
	create := r.client.AppDeployment.Create().
		SetAppVersionID(entity.AppVersionID).
		SetTargetType(entity.TargetType).
		SetTargetID(entity.TargetID).
		SetStatus(entity.Status)

	deployment, err := create.Save(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi thu xếp triển khai ứng dụng").WithError(err)
	}

	return deployment, nil
}

func (r *applicationRepositoryImpl) ListDeployments(ctx context.Context, appVersionID uint) ([]*ent.AppDeployment, error) {
	deployments, err := r.client.AppDeployment.Query().
		Where(appdeployment.AppVersionIDEQ(appVersionID)).
		Order(ent.Desc(appdeployment.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi lấy danh sách thiết bị triển khai").WithError(err)
	}
	return deployments, nil
}

func (r *applicationRepositoryImpl) AppVersionExists(ctx context.Context, id uint) (bool, error) {
	exists, err := r.client.AppVersion.Query().Where(appversion.IDEQ(id)).Exist(ctx)
	if err != nil {
		return false, apperror.ErrInternalServerError.WithMessage("Lỗi khi kiểm tra phiên bản ứng dụng").WithError(err)
	}
	return exists, nil
}

func (r *applicationRepositoryImpl) GetVersionByID(ctx context.Context, id uint) (*ent.AppVersion, error) {
	version, err := r.client.AppVersion.Query().
		Where(appversion.IDEQ(id)).
		WithApplication().
		First(ctx)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Không tìm thấy phiên bản ứng dụng")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất phiên bản ứng dụng").WithError(err)
	}
	return version, nil
}

func (r *applicationRepositoryImpl) UpdateDeploymentStatus(ctx context.Context, id uint, status string, errorMessage string) error {
	update := r.client.AppDeployment.UpdateOneID(id).
		SetStatus(appdeployment.Status(status))
	
	if errorMessage != "" {
		update = update.SetErrorMessage(errorMessage)
	}
	
	if err := update.Exec(ctx); err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi cập nhật trạng thái triển khai").WithError(err)
	}
	return nil
}

func (r *applicationRepositoryImpl) CountDeployments(ctx context.Context) (int, error) {
	count, err := r.client.AppDeployment.Query().Count(ctx)
	if err != nil {
		return 0, apperror.ErrInternalServerError.WithMessage("Lỗi khi đếm số lượng triển khai").WithError(err)
	}
	return count, nil
}
