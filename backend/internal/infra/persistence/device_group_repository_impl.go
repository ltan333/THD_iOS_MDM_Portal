package persistence

import (
	"context"
	"strings"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/devicegroup"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
)

type deviceGroupRepositoryImpl struct {
	client *ent.Client
}

func NewDeviceGroupRepository(client *ent.Client) repository.DeviceGroupRepository {
	return &deviceGroupRepositoryImpl{client: client}
}

func (r *deviceGroupRepositoryImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.DeviceGroup, int64, error) {
	q := r.client.DeviceGroup.Query()

	for field, filter := range opts.Filters {
		switch field {
		case "search":
			if val, ok := filter.Value.(string); ok && val != "" {
				q = q.Where(
					devicegroup.Or(
						devicegroup.NameContainsFold(val),
						devicegroup.DescriptionContainsFold(val),
					),
				)
			}
		}
	}

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, apperror.ErrInternalServerError.WithMessage("Lỗi khi đếm device group").WithError(err)
	}

	if len(opts.Sort) > 0 {
		for _, sortField := range opts.Sort {
			switch strings.ToLower(sortField.Field) {
			case "name":
				if sortField.Desc {
					q = q.Order(ent.Desc(devicegroup.FieldName))
				} else {
					q = q.Order(ent.Asc(devicegroup.FieldName))
				}
			case "created_at":
				if sortField.Desc {
					q = q.Order(ent.Desc(devicegroup.FieldCreatedAt))
				} else {
					q = q.Order(ent.Asc(devicegroup.FieldCreatedAt))
				}
			}
		}
	} else {
		q = q.Order(ent.Desc(devicegroup.FieldCreatedAt))
	}

	groups, err := q.Offset(offset).Limit(limit).All(ctx)
	if err != nil {
		return nil, 0, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất device groups").WithError(err)
	}

	return groups, int64(total), nil
}

func (r *deviceGroupRepositoryImpl) GetByID(ctx context.Context, id uint) (*ent.DeviceGroup, error) {
	group, err := r.client.DeviceGroup.Query().
		Where(devicegroup.IDEQ(id)).
		WithDevices().
		Only(ctx)

	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Nhóm thiết bị không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi lấy thông tin nhóm thiết bị").WithError(err)
	}

	return group, nil
}

func (r *deviceGroupRepositoryImpl) Create(ctx context.Context, entity *ent.DeviceGroup) (*ent.DeviceGroup, error) {
	exists, err := r.client.DeviceGroup.Query().Where(devicegroup.NameEQ(entity.Name)).Exist(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi kiểm tra tên nhóm thiết bị").WithError(err)
	}
	if exists {
		return nil, apperror.ErrConflict.WithMessage("Tên nhóm thiết bị đã tồn tại")
	}

	create := r.client.DeviceGroup.Create().
		SetName(entity.Name).
		SetDescription(entity.Description)

	group, err := create.Save(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo nhóm thiết bị").WithError(err)
	}

	return group, nil
}

func (r *deviceGroupRepositoryImpl) Update(ctx context.Context, id uint, entity *ent.DeviceGroup) (*ent.DeviceGroup, error) {
	if entity.Name != "" {
		exists, err := r.client.DeviceGroup.Query().
			Where(
				devicegroup.NameEQ(entity.Name),
				devicegroup.IDNEQ(id),
			).
			Exist(ctx)
		if err != nil {
			return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi kiểm tra tên nhóm thiết bị").WithError(err)
		}
		if exists {
			return nil, apperror.ErrConflict.WithMessage("Tên nhóm thiết bị đã tồn tại")
		}
	}

	update := r.client.DeviceGroup.UpdateOneID(id)

	if entity.Name != "" {
		update.SetName(entity.Name)
	}
	if entity.Description != "" {
		update.SetDescription(entity.Description)
	}

	group, err := update.Save(ctx)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Nhóm thiết bị không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi cập nhật nhóm thiết bị").WithError(err)
	}

	return group, nil
}

func (r *deviceGroupRepositoryImpl) Delete(ctx context.Context, id uint) error {
	err := r.client.DeviceGroup.DeleteOneID(id).Exec(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Nhóm thiết bị không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi xoá nhóm thiết bị").WithError(err)
	}
	return nil
}

func (r *deviceGroupRepositoryImpl) AddDevices(ctx context.Context, groupID uint, deviceIDs []string) error {
	_, err := r.client.DeviceGroup.UpdateOneID(groupID).AddDeviceIDs(deviceIDs...).Save(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Nhóm thiết bị không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi thêm thiết bị vào nhóm").WithError(err)
	}
	return nil
}

func (r *deviceGroupRepositoryImpl) RemoveDevice(ctx context.Context, groupID uint, deviceID string) error {
	_, err := r.client.DeviceGroup.UpdateOneID(groupID).RemoveDeviceIDs(deviceID).Save(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Nhóm thiết bị không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi xóa thiết bị khỏi nhóm").WithError(err)
	}
	return nil
}

func (r *deviceGroupRepositoryImpl) AssignProfile(ctx context.Context, groupID uint, profileID uint) error {
	_, err := r.client.DeviceGroup.UpdateOneID(groupID).
		AddProfileIDs(profileID).
		Save(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Nhóm thiết bị không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi gán profile").WithError(err)
	}
	return nil
}
