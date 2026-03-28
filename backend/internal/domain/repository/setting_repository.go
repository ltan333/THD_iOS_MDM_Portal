package repository

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
)

type SettingRepository interface {
	List(ctx context.Context) ([]*ent.Setting, error)
	GetByKey(ctx context.Context, key string) (*ent.Setting, error)
	Create(ctx context.Context, entity *ent.Setting) (*ent.Setting, error)
	Update(ctx context.Context, key string, entity *ent.Setting) (*ent.Setting, error)
	Delete(ctx context.Context, key string) error
	KeyExists(ctx context.Context, key string) (bool, error)
}
