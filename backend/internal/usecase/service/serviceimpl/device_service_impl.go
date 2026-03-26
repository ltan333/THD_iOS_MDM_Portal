package serviceimpl

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"strings"
	"time"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/device"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
	"github.com/thienel/tlog"
	"go.uber.org/zap"
)

type deviceServiceImpl struct {
	client         *ent.Client
	profileService service.ProfileService
}

func NewDeviceService(client *ent.Client, profileService service.ProfileService) service.DeviceService {
	return &deviceServiceImpl{
		client:         client,
		profileService: profileService,
	}
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

	if cmd.MacAddress != "" {
		create = create.SetMACAddress(cmd.MacAddress)
	}
	if cmd.IpAddress != "" {
		create = create.SetIPAddress(cmd.IpAddress)
	}
	if cmd.BatteryLevel > 0 {
		create = create.SetBatteryLevel(cmd.BatteryLevel)
	}
	if cmd.StorageCapacity > 0 {
		create = create.SetStorageCapacity(cmd.StorageCapacity)
	}
	if cmd.StorageUsed > 0 {
		create = create.SetStorageUsed(cmd.StorageUsed)
	}
	create = create.SetIsJailbroken(cmd.IsJailbroken)

	if cmd.EnrollmentType != "" {
		create = create.SetEnrollmentType(device.EnrollmentType(cmd.EnrollmentType))
	}

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
	if cmd.MacAddress != nil {
		update = update.SetMACAddress(*cmd.MacAddress)
	}
	if cmd.IpAddress != nil {
		update = update.SetIPAddress(*cmd.IpAddress)
	}
	if cmd.BatteryLevel != nil {
		update = update.SetBatteryLevel(*cmd.BatteryLevel)
	}
	if cmd.StorageCapacity != nil {
		update = update.SetStorageCapacity(*cmd.StorageCapacity)
	}
	if cmd.StorageUsed != nil {
		update = update.SetStorageUsed(*cmd.StorageUsed)
	}
	if cmd.IsJailbroken != nil {
		update = update.SetIsJailbroken(*cmd.IsJailbroken)
	}
	if cmd.EnrollmentType != nil {
		update = update.SetEnrollmentType(device.EnrollmentType(*cmd.EnrollmentType))
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
		_ = writer.Write([]string{"ID", "Serial Number", "Model", "Name", "Platform", "Status", "Is Enrolled"})
		for _, d := range devices {
			_ = writer.Write([]string{
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

func (s *deviceServiceImpl) HandleWebhook(ctx context.Context, payload *dto.NanoCMDWebhook) error {
	if payload.Topic == "" {
		return nil
	}

	// MDM Check-in events usually have the UDID in the checkin_event map
	var udid string
	if payload.Checkin_event != nil {
		if u, ok := payload.Checkin_event["udid"].(string); ok {
			udid = u
		}
	}

	if udid == "" {
		return nil
	}

	switch payload.Topic {
	case "mdm.Authenticate":
		tlog.Info("Device authenticating", zap.String("udid", udid))
		// Create or update device record
		exists, _ := s.client.Device.Query().Where(device.IDEQ(udid)).Exist(ctx)
		if !exists {
			return s.client.Device.Create().
				SetID(udid).
				SetStatus(device.StatusPending).
				SetIsEnrolled(false).
				Exec(ctx)
		}
		return nil
	case "mdm.TokenUpdate":
		tlog.Info("Device tokens updated/enrolled", zap.String("udid", udid))
		
		var finalID = udid
		sn, _ := payload.Checkin_event["serial_number"].(string)

		// IDENTITY MIGRATION: Check if device exists as "dep-SN"
		if sn != "" {
			depID := "dep-" + sn
			depDev, err := s.client.Device.Query().
				Where(device.IDEQ(depID)).
				WithGroups().
				Only(ctx)
			
			if err == nil && depDev != nil {
				tlog.Info("Migrating identity from DEP SN to UDID", zap.String("sn", sn), zap.String("udid", udid))
				
				// Copy data and groups to new UDID record
				create := s.client.Device.Create().
					SetID(udid).
					SetSerialNumber(sn).
					SetIsEnrolled(true).
					SetStatus(device.StatusActive).
					SetEnrolledAt(time.Now()).
					SetLastSeen(time.Now())

				if depDev.OwnerID != 0 { create.SetOwnerID(depDev.OwnerID) }
				if depDev.Model != "" { create.SetModel(depDev.Model) }
				if depDev.Name != "" { create.SetName(depDev.Name) }

				// Add groups
				groupIDs := make([]uint, len(depDev.Edges.Groups))
				for i, g := range depDev.Edges.Groups {
					groupIDs[i] = g.ID
				}
				create.AddGroupIDs(groupIDs...)

				_, err = create.Save(ctx)
				if err == nil {
					// Delete old DEP record
					_ = s.client.Device.DeleteOneID(depID).Exec(ctx)
				} else {
					tlog.Error("Failed to create migrated device record", zap.Error(err))
				}
			} else {
				// Regular update if no DEP record found
				updater := s.client.Device.UpdateOneID(udid).
					SetIsEnrolled(true).
					SetStatus(device.StatusActive).
					SetEnrolledAt(time.Now()).
					SetLastSeen(time.Now())

				if sn != "" { updater.SetSerialNumber(sn) }
				if model, ok := payload.Checkin_event["model"].(string); ok && model != "" {
					updater.SetModel(model)
				}
				if osVer, ok := payload.Checkin_event["os_version"].(string); ok && osVer != "" {
					updater.SetOsVersion(osVer)
				}
				
				_, err = updater.Save(ctx)
				if err != nil {
					// If update fails because it doesn't exist, create it
					if ent.IsNotFound(err) {
						_, _ = s.client.Device.Create().
							SetID(udid).
							SetSerialNumber(sn).
							SetIsEnrolled(true).
							SetStatus(device.StatusActive).
							SetEnrolledAt(time.Now()).
							SetLastSeen(time.Now()).
							Save(ctx)
					}
				}
			}
		}

		// TRIGGER AUTO-DEPLOY: Push profiles assigned to this device or its groups
		if err := s.profileService.DeployToDevice(ctx, finalID); err != nil {
			tlog.Error("Failed to trigger auto-deploy after enrollment", zap.String("udid", finalID), zap.Error(err))
		}
		return nil

	case "mdm.CheckOut":
		tlog.Info("Device checked out/unenrolled", zap.String("udid", udid))
		return s.client.Device.UpdateOneID(udid).
			SetIsEnrolled(false).
			SetStatus(device.StatusInactive).
			Exec(ctx)
	}

	return nil
}

func (s *deviceServiceImpl) UpsertFromDEP(ctx context.Context, devices []any) error {
	for _, devAny := range devices {
		devMap, ok := devAny.(map[string]any)
		if !ok {
			continue
		}

		sn, _ := devMap["serial_number"].(string)
		if sn == "" {
			continue
		}

		// Use serial number as ID if UDID is not available yet (DEP devices)
		// Or try to find by serial number
		existing, _ := s.client.Device.Query().Where(device.SerialNumberEQ(sn)).Only(ctx)
		
		model, _ := devMap["model"].(string)
		desc, _ := devMap["description"].(string)
		
		if existing != nil {
			updater := s.client.Device.UpdateOne(existing)
			if model != "" { updater.SetModel(model) }
			if desc != "" { updater.SetName(desc) }
			updater.SetEnrollmentType(device.EnrollmentTypeDep)
			if err := updater.Exec(ctx); err != nil {
				tlog.Error("Failed to update DEP device", zap.String("sn", sn), zap.Error(err))
			}
		} else {
			// Create new record using serial number as ID temporarily if no UDID
			// Most DEP devices will eventually enroll and update their ID to UDID
			// For now, we use a prefixed SN or just SN to represent the "pre-enrolled" state
			err := s.client.Device.Create().
				SetID("dep-" + sn).
				SetSerialNumber(sn).
				SetModel(model).
				SetName(desc).
				SetStatus(device.StatusPending).
				SetEnrollmentType(device.EnrollmentTypeDep).
				Exec(ctx)
			if err != nil {
				tlog.Error("Failed to create DEP device", zap.String("sn", sn), zap.Error(err))
			}
		}
	}
	return nil
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
