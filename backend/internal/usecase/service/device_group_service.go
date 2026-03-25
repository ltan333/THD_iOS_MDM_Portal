package service

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/pkg/query"
)

type CreateDeviceGroupCommand struct {
	Name        string
	Description string
}

type UpdateDeviceGroupCommand struct {
	ID          uint
	Name        *string
	Description *string
}

type DeviceGroupService interface {
	List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.DeviceGroup, int64, error)
	GetByID(ctx context.Context, id uint) (*ent.DeviceGroup, error)
	Create(ctx context.Context, cmd CreateDeviceGroupCommand) (*ent.DeviceGroup, error)
	Update(ctx context.Context, cmd UpdateDeviceGroupCommand) (*ent.DeviceGroup, error)
	Delete(ctx context.Context, id uint) error
	AddDevices(ctx context.Context, groupID uint, deviceIDs []string) error
	RemoveDevice(ctx context.Context, groupID uint, deviceID string) error
	AssignProfile(ctx context.Context, groupID uint, profileID uint) error
}
