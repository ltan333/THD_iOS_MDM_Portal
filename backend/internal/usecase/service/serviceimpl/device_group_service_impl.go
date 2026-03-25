package serviceimpl

import (
	"context"
	"strings"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/devicegroup"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
)

type deviceGroupServiceImpl struct {
	client *ent.Client
}

func NewDeviceGroupService(client *ent.Client) service.DeviceGroupService {
	return &deviceGroupServiceImpl{client: client}
}

func (s *deviceGroupServiceImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.DeviceGroup, int64, error) {
	q := s.client.DeviceGroup.Query()

	// Apply filters
	for field, filter := range opts.Filters {
		switch field {
		case "search":
			if searchVal, ok := filter.Value.(string); ok && searchVal != "" {
				q = q.Where(
					devicegroup.Or(
						devicegroup.NameContainsFold(searchVal),
						devicegroup.DescriptionContainsFold(searchVal),
					),
				)
			}
		case "name":
			if val, ok := filter.Value.(string); ok && val != "" {
				q = q.Where(devicegroup.NameContainsFold(val))
			}
		}
	}

	// Count total
	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, apperror.ErrInternalServerError.WithMessage("Lỗi khi đếm nhóm thiết bị").WithError(err)
	}

	// Apply sorting
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

	groups, err := q.Offset(offset).Limit(limit).WithDevices().All(ctx)
	if err != nil {
		return nil, 0, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất danh sách nhóm").WithError(err)
	}

	return groups, int64(total), nil
}

func (s *deviceGroupServiceImpl) GetByID(ctx context.Context, id uint) (*ent.DeviceGroup, error) {
	if id == 0 {
		return nil, apperror.ErrValidation.WithMessage("ID nhóm là bắt buộc")
	}

	g, err := s.client.DeviceGroup.Query().
		Where(devicegroup.IDEQ(id)).
		WithDevices().
		WithProfiles().
		First(ctx)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Nhóm thiết bị không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất nhóm").WithError(err)
	}

	return g, nil
}

func (s *deviceGroupServiceImpl) Create(ctx context.Context, cmd service.CreateDeviceGroupCommand) (*ent.DeviceGroup, error) {
	if strings.TrimSpace(cmd.Name) == "" {
		return nil, apperror.ErrValidation.WithMessage("Tên nhóm là bắt buộc")
	}

	// Check name uniqueness
	exists, err := s.client.DeviceGroup.Query().
		Where(devicegroup.NameEQ(cmd.Name)).
		Exist(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi kiểm tra tên nhóm").WithError(err)
	}
	if exists {
		return nil, apperror.ErrConflict.WithMessage("Tên nhóm đã tồn tại")
	}

	g, err := s.client.DeviceGroup.Create().
		SetName(cmd.Name).
		SetDescription(cmd.Description).
		Save(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo nhóm").WithError(err)
	}

	return g, nil
}

func (s *deviceGroupServiceImpl) Update(ctx context.Context, cmd service.UpdateDeviceGroupCommand) (*ent.DeviceGroup, error) {
	if cmd.ID == 0 {
		return nil, apperror.ErrValidation.WithMessage("ID nhóm là bắt buộc")
	}

	update := s.client.DeviceGroup.UpdateOneID(cmd.ID)

	if cmd.Name != nil && strings.TrimSpace(*cmd.Name) != "" {
		// Check name uniqueness for update
		exists, err := s.client.DeviceGroup.Query().
			Where(devicegroup.NameEQ(*cmd.Name), devicegroup.IDNEQ(cmd.ID)).
			Exist(ctx)
		if err != nil {
			return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi kiểm tra tên nhóm").WithError(err)
		}
		if exists {
			return nil, apperror.ErrConflict.WithMessage("Tên nhóm đã tồn tại")
		}
		update = update.SetName(*cmd.Name)
	}
	if cmd.Description != nil {
		update = update.SetDescription(*cmd.Description)
	}

	g, err := update.Save(ctx)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Nhóm thiết bị không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi cập nhật nhóm").WithError(err)
	}

	return g, nil
}

func (s *deviceGroupServiceImpl) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return apperror.ErrValidation.WithMessage("ID nhóm là bắt buộc")
	}

	err := s.client.DeviceGroup.DeleteOneID(id).Exec(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Nhóm thiết bị không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi xóa nhóm").WithError(err)
	}

	return nil
}

func (s *deviceGroupServiceImpl) AddDevices(ctx context.Context, groupID uint, deviceIDs []string) error {
	if groupID == 0 {
		return apperror.ErrValidation.WithMessage("ID nhóm là bắt buộc")
	}
	if len(deviceIDs) == 0 {
		return apperror.ErrValidation.WithMessage("Danh sách thiết bị không được rỗng")
	}

	_, err := s.client.DeviceGroup.UpdateOneID(groupID).
		AddDeviceIDs(deviceIDs...).
		Save(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Nhóm thiết bị không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi thêm thiết bị vào nhóm").WithError(err)
	}

	return nil
}

func (s *deviceGroupServiceImpl) RemoveDevice(ctx context.Context, groupID uint, deviceID string) error {
	if groupID == 0 {
		return apperror.ErrValidation.WithMessage("ID nhóm là bắt buộc")
	}
	if deviceID == "" {
		return apperror.ErrValidation.WithMessage("ID thiết bị là bắt buộc")
	}

	_, err := s.client.DeviceGroup.UpdateOneID(groupID).
		RemoveDeviceIDs(deviceID).
		Save(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Nhóm thiết bị không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi xóa thiết bị khỏi nhóm").WithError(err)
	}

	return nil
}

func (s *deviceGroupServiceImpl) AssignProfile(ctx context.Context, groupID uint, profileID uint) error {
	if groupID == 0 {
		return apperror.ErrValidation.WithMessage("ID nhóm là bắt buộc")
	}
	if profileID == 0 {
		return apperror.ErrValidation.WithMessage("ID profile là bắt buộc")
	}

	_, err := s.client.DeviceGroup.UpdateOneID(groupID).
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
