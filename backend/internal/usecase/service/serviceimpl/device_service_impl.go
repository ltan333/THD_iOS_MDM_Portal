package serviceimpl

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"strings"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/device"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
)

type deviceServiceImpl struct {
	client *ent.Client
}

func NewDeviceService(client *ent.Client) service.DeviceService {
	return &deviceServiceImpl{client: client}
}

func (s *deviceServiceImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.Device, int64, error) {
	q := s.client.Device.Query()

	// Apply filters
	for field, filter := range opts.Filters {
		switch field {
		case "search":
			if searchVal, ok := filter.Value.(string); ok && searchVal != "" {
				q = q.Where(
					device.Or(
						device.NameContainsFold(searchVal),
						device.SerialNumberContainsFold(searchVal),
						device.ModelContainsFold(searchVal),
					),
				)
			}
		case "platform":
			if val, ok := filter.Value.(string); ok && val != "" {
				q = q.Where(device.PlatformEQ(device.Platform(val)))
			}
		case "status":
			if val, ok := filter.Value.(string); ok && val != "" {
				q = q.Where(device.StatusEQ(device.Status(val)))
			}
		case "is_enrolled":
			if val, ok := filter.Value.(string); ok {
				enrolled := val == "true"
				q = q.Where(device.IsEnrolledEQ(enrolled))
			}
		case "compliance_status":
			if val, ok := filter.Value.(string); ok && val != "" {
				q = q.Where(device.ComplianceStatusEQ(device.ComplianceStatus(val)))
			}
		}
	}

	// Count total
	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, apperror.ErrInternalServerError.WithMessage("Lỗi khi đếm thiết bị").WithError(err)
	}

	// Apply sorting
	if len(opts.Sort) > 0 {
		for _, sortField := range opts.Sort {
			switch strings.ToLower(sortField.Field) {
			case "name":
				if sortField.Desc {
					q = q.Order(ent.Desc(device.FieldName))
				} else {
					q = q.Order(ent.Asc(device.FieldName))
				}
			case "created_at":
				if sortField.Desc {
					q = q.Order(ent.Desc(device.FieldCreatedAt))
				} else {
					q = q.Order(ent.Asc(device.FieldCreatedAt))
				}
			case "last_seen":
				if sortField.Desc {
					q = q.Order(ent.Desc(device.FieldLastSeen))
				} else {
					q = q.Order(ent.Asc(device.FieldLastSeen))
				}
			}
		}
	} else {
		q = q.Order(ent.Desc(device.FieldCreatedAt))
	}

	devices, err := q.Offset(offset).Limit(limit).All(ctx)
	if err != nil {
		return nil, 0, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất danh sách thiết bị").WithError(err)
	}

	return devices, int64(total), nil
}

func (s *deviceServiceImpl) GetByID(ctx context.Context, id string) (*ent.Device, error) {
	if id == "" {
		return nil, apperror.ErrValidation.WithMessage("ID thiết bị là bắt buộc")
	}

	d, err := s.client.Device.Get(ctx, id)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Thiết bị không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất thiết bị").WithError(err)
	}

	return d, nil
}

func (s *deviceServiceImpl) Create(ctx context.Context, cmd service.CreateDeviceCommand) (*ent.Device, error) {
	if cmd.ID == "" {
		return nil, apperror.ErrValidation.WithMessage("ID thiết bị là bắt buộc")
	}

	// Check if device exists
	exists, err := s.client.Device.Query().Where(device.IDEQ(cmd.ID)).Exist(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi kiểm tra thiết bị").WithError(err)
	}
	if exists {
		return nil, apperror.ErrConflict.WithMessage("Thiết bị đã tồn tại")
	}

	create := s.client.Device.Create().
		SetID(cmd.ID).
		SetNillableSerialNumber(&cmd.SerialNumber).
		SetNillableModel(&cmd.Model).
		SetNillableName(&cmd.Name)

	if cmd.Platform != "" {
		create = create.SetPlatform(device.Platform(cmd.Platform))
	}
	if cmd.OwnerID != nil {
		create = create.SetOwnerID(*cmd.OwnerID)
	}

	d, err := create.Save(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo thiết bị").WithError(err)
	}

	return d, nil
}

func (s *deviceServiceImpl) Update(ctx context.Context, cmd service.UpdateDeviceCommand) (*ent.Device, error) {
	if cmd.ID == "" {
		return nil, apperror.ErrValidation.WithMessage("ID thiết bị là bắt buộc")
	}

	update := s.client.Device.UpdateOneID(cmd.ID)

	if cmd.SerialNumber != nil {
		update = update.SetSerialNumber(*cmd.SerialNumber)
	}
	if cmd.Model != nil {
		update = update.SetModel(*cmd.Model)
	}
	if cmd.Name != nil {
		update = update.SetName(*cmd.Name)
	}
	if cmd.Platform != nil {
		update = update.SetPlatform(device.Platform(*cmd.Platform))
	}
	if cmd.Status != nil {
		update = update.SetStatus(device.Status(*cmd.Status))
	}
	if cmd.ComplianceStatus != nil {
		update = update.SetComplianceStatus(device.ComplianceStatus(*cmd.ComplianceStatus))
	}
	if cmd.IsEnrolled != nil {
		update = update.SetIsEnrolled(*cmd.IsEnrolled)
	}
	if cmd.OwnerID != nil {
		update = update.SetOwnerID(*cmd.OwnerID)
	}
	if cmd.OsVersion != nil {
		update = update.SetOsVersion(*cmd.OsVersion)
	}
	if cmd.DeviceType != nil {
		update = update.SetDeviceType(*cmd.DeviceType)
	}

	d, err := update.Save(ctx)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Thiết bị không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi cập nhật thiết bị").WithError(err)
	}

	return d, nil
}

func (s *deviceServiceImpl) Delete(ctx context.Context, id string) error {
	if id == "" {
		return apperror.ErrValidation.WithMessage("ID thiết bị là bắt buộc")
	}

	err := s.client.Device.DeleteOneID(id).Exec(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Thiết bị không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi xóa thiết bị").WithError(err)
	}

	return nil
}

func (s *deviceServiceImpl) GetStats(ctx context.Context) (*dto.DeviceStatsResponse, error) {
	total, _ := s.client.Device.Query().Count(ctx)
	active, _ := s.client.Device.Query().Where(device.StatusEQ(device.StatusActive)).Count(ctx)
	inactive, _ := s.client.Device.Query().Where(device.StatusEQ(device.StatusInactive)).Count(ctx)
	enrolled, _ := s.client.Device.Query().Where(device.IsEnrolledEQ(true)).Count(ctx)
	compliant, _ := s.client.Device.Query().Where(device.ComplianceStatusEQ(device.ComplianceStatusCompliant)).Count(ctx)
	nonCompliant, _ := s.client.Device.Query().Where(device.ComplianceStatusEQ(device.ComplianceStatusNonCompliant)).Count(ctx)

	// Count by platform
	byPlatform := map[string]int64{
		"ios":     0,
		"android": 0,
		"windows": 0,
		"macos":   0,
		"other":   0,
	}
	ios, _ := s.client.Device.Query().Where(device.PlatformEQ(device.PlatformIos)).Count(ctx)
	android, _ := s.client.Device.Query().Where(device.PlatformEQ(device.PlatformAndroid)).Count(ctx)
	windows, _ := s.client.Device.Query().Where(device.PlatformEQ(device.PlatformWindows)).Count(ctx)
	macos, _ := s.client.Device.Query().Where(device.PlatformEQ(device.PlatformMacos)).Count(ctx)
	byPlatform["ios"] = int64(ios)
	byPlatform["android"] = int64(android)
	byPlatform["windows"] = int64(windows)
	byPlatform["macos"] = int64(macos)
	byPlatform["other"] = int64(total - ios - android - windows - macos)

	byStatus := map[string]int64{
		"active":   int64(active),
		"inactive": int64(inactive),
		"pending":  int64(total - active - inactive),
	}

	return &dto.DeviceStatsResponse{
		Total:        int64(total),
		Active:       int64(active),
		Inactive:     int64(inactive),
		Enrolled:     int64(enrolled),
		ByPlatform:   byPlatform,
		ByStatus:     byStatus,
		Compliant:    int64(compliant),
		NonCompliant: int64(nonCompliant),
	}, nil
}

func (s *deviceServiceImpl) Export(ctx context.Context, format string) ([]byte, error) {
	devices, err := s.client.Device.Query().All(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất thiết bị").WithError(err)
	}

	switch format {
	case "json":
		return json.Marshal(devices)
	case "csv":
		var buf strings.Builder
		writer := csv.NewWriter(&buf)
		writer.Write([]string{"ID", "Serial Number", "Model", "Name", "Platform", "Status", "Is Enrolled"})
		for _, d := range devices {
			writer.Write([]string{
				d.ID,
				d.SerialNumber,
				d.Model,
				d.Name,
				string(d.Platform),
				string(d.Status),
				boolToString(d.IsEnrolled),
			})
		}
		writer.Flush()
		return []byte(buf.String()), nil
	default:
		return nil, apperror.ErrBadRequest.WithMessage("Format không hỗ trợ: " + format)
	}
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
