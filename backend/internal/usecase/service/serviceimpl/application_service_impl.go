package serviceimpl

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/appdeployment"
	"github.com/thienel/go-backend-template/internal/ent/application"
	"github.com/thienel/go-backend-template/internal/ent/appversion"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
)

type applicationServiceImpl struct {
	client *ent.Client
}

func NewApplicationService(client *ent.Client) service.ApplicationService {
	return &applicationServiceImpl{client: client}
}

func (s *applicationServiceImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.Application, int64, error) {
	q := s.client.Application.Query()

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

	total, err := q.Count(ctx)
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

func (s *applicationServiceImpl) GetByID(ctx context.Context, id uint) (*ent.Application, error) {
	app, err := s.client.Application.Query().
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

func (s *applicationServiceImpl) Create(ctx context.Context, cmd service.CreateApplicationCommand) (*ent.Application, error) {
	exists, err := s.client.Application.Query().Where(application.BundleIDEQ(cmd.BundleID)).Exist(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi kiểm tra Bundle ID").WithError(err)
	}
	if exists {
		return nil, apperror.ErrConflict.WithMessage("Bundle ID đã tồn tại trong hệ thống")
	}

	app, err := s.client.Application.Create().
		SetName(cmd.Name).
		SetBundleID(cmd.BundleID).
		SetPlatform(application.Platform(cmd.Platform)).
		SetType(application.Type(cmd.Type)).
		SetNillableDescription(&cmd.Description).
		SetNillableIconURL(&cmd.IconURL).
		Save(ctx)

	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo ứng dụng").WithError(err)
	}

	return app, nil
}

func (s *applicationServiceImpl) Update(ctx context.Context, cmd service.UpdateApplicationCommand) (*ent.Application, error) {
	update := s.client.Application.UpdateOneID(cmd.ID)

	if cmd.Name != nil {
		update.SetName(*cmd.Name)
	}
	if cmd.Platform != nil {
		update.SetPlatform(application.Platform(*cmd.Platform))
	}
	if cmd.Type != nil {
		update.SetType(application.Type(*cmd.Type))
	}
	if cmd.Description != nil {
		update.SetDescription(*cmd.Description)
	}
	if cmd.IconURL != nil {
		update.SetIconURL(*cmd.IconURL)
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

func (s *applicationServiceImpl) Delete(ctx context.Context, id uint) error {
	err := s.client.Application.DeleteOneID(id).Exec(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Không tìm thấy ứng dụng để xóa")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi xóa ứng dụng").WithError(err)
	}
	return nil
}

func (s *applicationServiceImpl) ListVersions(ctx context.Context, appID uint) ([]*ent.AppVersion, error) {
	versions, err := s.client.AppVersion.Query().
		Where(appversion.ApplicationIDEQ(appID)).
		Order(ent.Desc(appversion.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi lấy danh sách phiên bản ứng dụng").WithError(err)
	}

	return versions, nil
}

func (s *applicationServiceImpl) CreateVersion(ctx context.Context, cmd service.CreateAppVersionCommand) (*ent.AppVersion, error) {
	exists, err := s.client.Application.Query().Where(application.IDEQ(cmd.ApplicationID)).Exist(ctx)
	if err != nil || !exists {
		return nil, apperror.ErrNotFound.WithMessage("Không tìm thấy ứng dụng gốc")
	}

	version, err := s.client.AppVersion.Create().
		SetApplicationID(cmd.ApplicationID).
		SetVersion(cmd.Version).
		SetBuildNumber(cmd.BuildNumber).
		SetNillableMinimumOsVersion(&cmd.MinimumOSVersion).
		SetNillableFileURL(&cmd.FileURL).
		SetNillableSize(&cmd.Size).
		SetMetadata(cmd.Metadata).
		Save(ctx)

	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo phiên bản ứng dụng").WithError(err)
	}

	return version, nil
}

func (s *applicationServiceImpl) DeleteVersion(ctx context.Context, id uint) error {
	err := s.client.AppVersion.DeleteOneID(id).Exec(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Không tìm thấy phiên bản để xóa")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi xóa phiên bản").WithError(err)
	}
	return nil
}

func (s *applicationServiceImpl) Deploy(ctx context.Context, cmd service.CreateAppDeploymentCommand) (*ent.AppDeployment, error) {
	exists, err := s.client.AppVersion.Query().Where(appversion.IDEQ(cmd.AppVersionID)).Exist(ctx)
	if err != nil || !exists {
		return nil, apperror.ErrNotFound.WithMessage("Không tìm thấy phiên bản ứng dụng để deploy")
	}

	deployment, err := s.client.AppDeployment.Create().
		SetAppVersionID(cmd.AppVersionID).
		SetTargetType(appdeployment.TargetType(cmd.TargetType)).
		SetTargetID(cmd.TargetID).
		SetStatus(appdeployment.StatusPending).
		Save(ctx)

	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi thu xếp triển khai ứng dụng").WithError(err)
	}

	return deployment, nil
}

func (s *applicationServiceImpl) ListDeployments(ctx context.Context, appVersionID uint) ([]*ent.AppDeployment, error) {
	deployments, err := s.client.AppDeployment.Query().
		Where(appdeployment.AppVersionIDEQ(appVersionID)).
		Order(ent.Desc(appdeployment.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi lấy danh sách thiết bị triển khai").WithError(err)
	}

	return deployments, nil
}
