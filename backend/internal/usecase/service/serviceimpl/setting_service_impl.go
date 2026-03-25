package serviceimpl

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/setting"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
)

type settingServiceImpl struct {
	client *ent.Client
}

func NewSettingService(client *ent.Client) service.SettingService {
	return &settingServiceImpl{client: client}
}

func (s *settingServiceImpl) List(ctx context.Context) ([]*ent.Setting, error) {
	settings, err := s.client.Setting.Query().
		Order(ent.Asc(setting.FieldKey)).
		All(ctx)

	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi lấy danh sách cài đặt").WithError(err)
	}

	return settings, nil
}

func (s *settingServiceImpl) GetByKey(ctx context.Context, key string) (*ent.Setting, error) {
	st, err := s.client.Setting.Query().Where(setting.KeyEQ(key)).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Không tìm thấy cấu hình với key này")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi lấy cấu hình").WithError(err)
	}

	return st, nil
}

func (s *settingServiceImpl) Create(ctx context.Context, cmd service.CreateSettingCommand) (*ent.Setting, error) {
	exists, err := s.client.Setting.Query().Where(setting.KeyEQ(cmd.Key)).Exist(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi kiểm tra key").WithError(err)
	}
	if exists {
		return nil, apperror.ErrConflict.WithMessage("Cấu hình với key này đã tồn tại")
	}

	st, err := s.client.Setting.Create().
		SetKey(cmd.Key).
		SetValue(cmd.Value).
		SetNillableDescription(&cmd.Description).
		Save(ctx)

	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo cấu hình").WithError(err)
	}

	return st, nil
}

func (s *settingServiceImpl) Update(ctx context.Context, cmd service.UpdateSettingCommand) (*ent.Setting, error) {
	st, err := s.client.Setting.Query().Where(setting.KeyEQ(cmd.Key)).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Không tìm thấy cấu hình với key này")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi lấy cấu hình").WithError(err)
	}

	update := st.Update()
	if cmd.Value != nil {
		update.SetValue(*cmd.Value)
	}
	if cmd.Description != nil {
		update.SetDescription(*cmd.Description)
	}

	updatedSt, err := update.Save(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi cập nhật cấu hình").WithError(err)
	}

	return updatedSt, nil
}

func (s *settingServiceImpl) Delete(ctx context.Context, key string) error {
	_, err := s.client.Setting.Delete().Where(setting.KeyEQ(key)).Exec(ctx)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi xóa cấu hình").WithError(err)
	}
	return nil
}
