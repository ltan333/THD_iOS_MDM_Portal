package service

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
)

type CreateSettingCommand struct {
	Key         string
	Value       string
	Description string
}

type UpdateSettingCommand struct {
	Key         string
	Value       *string
	Description *string
}

type SettingService interface {
	List(ctx context.Context) ([]*ent.Setting, error)
	GetByKey(ctx context.Context, key string) (*ent.Setting, error)
	Create(ctx context.Context, cmd CreateSettingCommand) (*ent.Setting, error)
	Update(ctx context.Context, cmd UpdateSettingCommand) (*ent.Setting, error)
	Delete(ctx context.Context, key string) error
}
