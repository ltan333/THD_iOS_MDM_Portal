package persistence

import (
	"context"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/setting"
	apperror "github.com/thienel/go-backend-template/pkg/error"
)

type settingRepositoryImpl struct {
	client *ent.Client
}

func NewSettingRepository(client *ent.Client) repository.SettingRepository {
	return &settingRepositoryImpl{client: client}
}

func (r *settingRepositoryImpl) List(ctx context.Context) ([]*ent.Setting, error) {
	return r.client.Setting.Query().
		Order(ent.Asc(setting.FieldKey)).
		All(ctx)
}

func (r *settingRepositoryImpl) GetByKey(ctx context.Context, key string) (*ent.Setting, error) {
	st, err := r.client.Setting.Query().Where(setting.KeyEQ(key)).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Không tìm thấy cấu hình với key này")
	}
	return st, err
}

func (r *settingRepositoryImpl) Create(ctx context.Context, entity *ent.Setting) (*ent.Setting, error) {
	create := r.client.Setting.Create().
		SetKey(entity.Key).
		SetValue(entity.Value)
	
	if entity.Description != "" {
		create = create.SetDescription(entity.Description)
	}

	return create.Save(ctx)
}

func (r *settingRepositoryImpl) Update(ctx context.Context, key string, entity *ent.Setting) (*ent.Setting, error) {
	update := r.client.Setting.Update().Where(setting.KeyEQ(key))

	if entity.Value != "" {
		update.SetValue(entity.Value)
	}
	if entity.Description != "" {
		update.SetDescription(entity.Description)
	}

	_, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}
	return r.GetByKey(ctx, key)
}

func (r *settingRepositoryImpl) Delete(ctx context.Context, key string) error {
	_, err := r.client.Setting.Delete().Where(setting.KeyEQ(key)).Exec(ctx)
	return err
}

func (r *settingRepositoryImpl) KeyExists(ctx context.Context, key string) (bool, error) {
	return r.client.Setting.Query().Where(setting.KeyEQ(key)).Exist(ctx)
}
