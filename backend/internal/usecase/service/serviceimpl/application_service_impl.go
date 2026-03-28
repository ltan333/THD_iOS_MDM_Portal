package serviceimpl

import (
	"context"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/appdeployment"
	"github.com/thienel/go-backend-template/internal/ent/application"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
)

type applicationServiceImpl struct {
	repo repository.ApplicationRepository
}

func NewApplicationService(repo repository.ApplicationRepository) service.ApplicationService {
	return &applicationServiceImpl{repo: repo}
}

func (s *applicationServiceImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.Application, int64, error) {
	return s.repo.List(ctx, offset, limit, opts)
}

func (s *applicationServiceImpl) GetByID(ctx context.Context, id uint) (*ent.Application, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *applicationServiceImpl) Create(ctx context.Context, cmd service.CreateApplicationCommand) (*ent.Application, error) {
	exists, err := s.repo.BundleIDExists(ctx, cmd.BundleID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperror.ErrConflict.WithMessage("Bundle ID đã tồn tại trong hệ thống")
	}

	return s.repo.Create(ctx, &ent.Application{
		Name:        cmd.Name,
		BundleID:    cmd.BundleID,
		Platform:    application.Platform(cmd.Platform),
		Type:        application.Type(cmd.Type),
		Description: cmd.Description,
		IconURL:     cmd.IconURL,
	})
}

func (s *applicationServiceImpl) Update(ctx context.Context, cmd service.UpdateApplicationCommand) (*ent.Application, error) {
	name := ""
	if cmd.Name != nil {
		name = *cmd.Name
	}
	platform := application.Platform("")
	if cmd.Platform != nil {
		platform = application.Platform(*cmd.Platform)
	}
	appType := application.Type("")
	if cmd.Type != nil {
		appType = application.Type(*cmd.Type)
	}
	desc := ""
	if cmd.Description != nil {
		desc = *cmd.Description
	}
	icon := ""
	if cmd.IconURL != nil {
		icon = *cmd.IconURL
	}

	return s.repo.Update(ctx, cmd.ID, &ent.Application{
		Name:        name,
		Platform:    platform,
		Type:        appType,
		Description: desc,
		IconURL:     icon,
	})
}

func (s *applicationServiceImpl) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *applicationServiceImpl) ListVersions(ctx context.Context, appID uint) ([]*ent.AppVersion, error) {
	return s.repo.ListVersions(ctx, appID)
}

func (s *applicationServiceImpl) CreateVersion(ctx context.Context, cmd service.CreateAppVersionCommand) (*ent.AppVersion, error) {
	exists, err := s.repo.AppExists(ctx, cmd.ApplicationID)
	if err != nil || !exists {
		return nil, apperror.ErrNotFound.WithMessage("Không tìm thấy ứng dụng gốc")
	}

	return s.repo.CreateVersion(ctx, &ent.AppVersion{
		ApplicationID:    cmd.ApplicationID,
		Version:          cmd.Version,
		BuildNumber:      cmd.BuildNumber,
		MinimumOsVersion: cmd.MinimumOSVersion,
		FileURL:          cmd.FileURL,
		Size:             cmd.Size,
		Metadata:         cmd.Metadata,
	})
}

func (s *applicationServiceImpl) DeleteVersion(ctx context.Context, id uint) error {
	return s.repo.DeleteVersion(ctx, id)
}

func (s *applicationServiceImpl) Deploy(ctx context.Context, cmd service.CreateAppDeploymentCommand) (*ent.AppDeployment, error) {
	exists, err := s.repo.AppVersionExists(ctx, cmd.AppVersionID)
	if err != nil || !exists {
		return nil, apperror.ErrNotFound.WithMessage("Không tìm thấy phiên bản ứng dụng để deploy")
	}

	return s.repo.CreateDeployment(ctx, &ent.AppDeployment{
		AppVersionID: cmd.AppVersionID,
		TargetType:   appdeployment.TargetType(cmd.TargetType),
		TargetID:     cmd.TargetID,
		Status:       appdeployment.StatusPending,
	})
}

func (s *applicationServiceImpl) ListDeployments(ctx context.Context, appVersionID uint) ([]*ent.AppDeployment, error) {
	return s.repo.ListDeployments(ctx, appVersionID)
}
