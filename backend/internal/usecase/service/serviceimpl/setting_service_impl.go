package serviceimpl

import (
	"context"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
)

type settingServiceImpl struct {
	repo repository.SettingRepository
}

func NewSettingService(repo repository.SettingRepository) service.SettingService {
	return &settingServiceImpl{repo: repo}
}

func (s *settingServiceImpl) List(ctx context.Context) ([]*ent.Setting, error) {
	settings, err := s.repo.List(ctx)

	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi lấy danh sách cài đặt").WithError(err)
	}

	return settings, nil
}

func (s *settingServiceImpl) GetByKey(ctx context.Context, key string) (*ent.Setting, error) {
	st, err := s.repo.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}

	return st, nil
}

func (s *settingServiceImpl) Create(ctx context.Context, cmd service.CreateSettingCommand) (*ent.Setting, error) {
	exists, err := s.repo.KeyExists(ctx, cmd.Key)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi kiểm tra key").WithError(err)
	}
	if exists {
		return nil, apperror.ErrConflict.WithMessage("Cấu hình với key này đã tồn tại")
	}

	st, err := s.repo.Create(ctx, &ent.Setting{
		Key:         cmd.Key,
		Value:       cmd.Value,
		Description: cmd.Description,
	})

	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo cấu hình").WithError(err)
	}

	return st, nil
}

func (s *settingServiceImpl) Update(ctx context.Context, cmd service.UpdateSettingCommand) (*ent.Setting, error) {
	exists, err := s.repo.KeyExists(ctx, cmd.Key)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi kết nối dữ liệu").WithError(err)
	}
	if !exists {
		return nil, apperror.ErrNotFound.WithMessage("Không tìm thấy cấu hình với key này")
	}

	val := ""
	if cmd.Value != nil {
		val = *cmd.Value
	}
	desc := ""
	if cmd.Description != nil {
		desc = *cmd.Description
	}

	updatedSt, err := s.repo.Update(ctx, cmd.Key, &ent.Setting{
		Value: val,
		Description: desc,
	})
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi cập nhật cấu hình").WithError(err)
	}

	return updatedSt, nil
}

func (s *settingServiceImpl) Delete(ctx context.Context, key string) error {
	err := s.repo.Delete(ctx, key)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi xóa cấu hình").WithError(err)
	}
	return nil
}
