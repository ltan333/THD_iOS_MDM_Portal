package serviceimpl

import (
	"context"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
)

type deviceGroupServiceImpl struct {
	repo repository.DeviceGroupRepository
}

func NewDeviceGroupService(repo repository.DeviceGroupRepository) service.DeviceGroupService {
	return &deviceGroupServiceImpl{repo: repo}
}

func (s *deviceGroupServiceImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.DeviceGroup, int64, error) {
	return s.repo.List(ctx, offset, limit, opts)
}

func (s *deviceGroupServiceImpl) GetByID(ctx context.Context, id uint) (*ent.DeviceGroup, error) {
	if id == 0 {
		return nil, apperror.ErrValidation.WithMessage("ID nhóm thiết bị là bắt buộc")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *deviceGroupServiceImpl) Create(ctx context.Context, cmd service.CreateDeviceGroupCommand) (*ent.DeviceGroup, error) {
	if cmd.Name == "" {
		return nil, apperror.ErrValidation.WithMessage("Tên nhóm thiết bị là bắt buộc")
	}

	return s.repo.Create(ctx, &ent.DeviceGroup{
		Name:        cmd.Name,
		Description: cmd.Description,
	})
}

func (s *deviceGroupServiceImpl) Update(ctx context.Context, cmd service.UpdateDeviceGroupCommand) (*ent.DeviceGroup, error) {
	if cmd.ID == 0 {
		return nil, apperror.ErrValidation.WithMessage("ID nhóm thiết bị là bắt buộc")
	}

	name := ""
	if cmd.Name != nil {
		name = *cmd.Name
	}
	desc := ""
	if cmd.Description != nil {
		desc = *cmd.Description
	}

	return s.repo.Update(ctx, cmd.ID, &ent.DeviceGroup{
		Name:        name,
		Description: desc,
	})
}

func (s *deviceGroupServiceImpl) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return apperror.ErrValidation.WithMessage("ID nhóm thiết bị là bắt buộc")
	}

	return s.repo.Delete(ctx, id)
}

func (s *deviceGroupServiceImpl) AddDevices(ctx context.Context, groupID uint, deviceIDs []string) error {
	if groupID == 0 {
		return apperror.ErrValidation.WithMessage("ID nhóm thiết bị là bắt buộc")
	}
	if len(deviceIDs) == 0 {
		return apperror.ErrValidation.WithMessage("Danh sách thiết bị không được rỗng")
	}

	return s.repo.AddDevices(ctx, groupID, deviceIDs)
}

func (s *deviceGroupServiceImpl) RemoveDevice(ctx context.Context, groupID uint, deviceID string) error {
	if groupID == 0 {
		return apperror.ErrValidation.WithMessage("ID nhóm thiết bị là bắt buộc")
	}
	if deviceID == "" {
		return apperror.ErrValidation.WithMessage("ID thiết bị là bắt buộc")
	}

	return s.repo.RemoveDevice(ctx, groupID, deviceID)
}

func (s *deviceGroupServiceImpl) AssignProfile(ctx context.Context, groupID uint, profileID uint) error {
	if groupID == 0 {
		return apperror.ErrValidation.WithMessage("ID nhóm là bắt buộc")
	}
	if profileID == 0 {
		return apperror.ErrValidation.WithMessage("ID profile là bắt buộc")
	}

	return s.repo.AssignProfile(ctx, groupID, profileID)
}
