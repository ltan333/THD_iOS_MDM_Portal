package persistence

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/device"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
	"github.com/thienel/tlog"
	"go.uber.org/zap"
)

type deviceRepositoryImpl struct {
	client *ent.Client
}

func NewDeviceRepository(client *ent.Client) repository.DeviceRepository {
	return &deviceRepositoryImpl{client: client}
}

func (r *deviceRepositoryImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.Device, int64, error) {
	q := r.client.Device.Query()

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

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, apperror.ErrInternalServerError.WithMessage("Lỗi khi đếm thiết bị").WithError(err)
	}

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

func (r *deviceRepositoryImpl) GetByID(ctx context.Context, id string) (*ent.Device, error) {
	d, err := r.client.Device.Get(ctx, id)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Thiết bị không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất thiết bị").WithError(err)
	}
	return d, nil
}

func (r *deviceRepositoryImpl) Create(ctx context.Context, entity *ent.Device) (*ent.Device, error) {
	exists, err := r.client.Device.Query().Where(device.IDEQ(entity.ID)).Exist(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi kiểm tra thiết bị").WithError(err)
	}
	if exists {
		return nil, apperror.ErrConflict.WithMessage("Thiết bị đã tồn tại")
	}

	create := r.client.Device.Create().
		SetID(entity.ID).
		SetNillableSerialNumber(&entity.SerialNumber).
		SetNillableModel(&entity.Model).
		SetNillableName(&entity.Name)

	if entity.MACAddress != "" {
		create = create.SetMACAddress(entity.MACAddress)
	}
	if entity.IPAddress != "" {
		create = create.SetIPAddress(entity.IPAddress)
	}
	if entity.BatteryLevel > 0 {
		create = create.SetBatteryLevel(entity.BatteryLevel)
	}
	if entity.StorageCapacity > 0 {
		create = create.SetStorageCapacity(entity.StorageCapacity)
	}
	if entity.StorageUsed > 0 {
		create = create.SetStorageUsed(entity.StorageUsed)
	}
	create = create.SetIsJailbroken(entity.IsJailbroken)

	if string(entity.EnrollmentType) != "" {
		create = create.SetEnrollmentType(entity.EnrollmentType)
	}
	if string(entity.Platform) != "" {
		create = create.SetPlatform(entity.Platform)
	}
	if entity.OwnerID != 0 { // Assume 0 means not set
		create = create.SetOwnerID(entity.OwnerID)
	}

	d, err := create.Save(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo thiết bị").WithError(err)
	}

	return d, nil
}

func (r *deviceRepositoryImpl) Update(ctx context.Context, id string, entity *ent.Device) (*ent.Device, error) {
	update := r.client.Device.UpdateOneID(id)

	if entity.SerialNumber != "" {
		update = update.SetSerialNumber(entity.SerialNumber)
	}
	if entity.Model != "" {
		update = update.SetModel(entity.Model)
	}
	if entity.Name != "" {
		update = update.SetName(entity.Name)
	}
	if string(entity.Platform) != "" {
		update = update.SetPlatform(entity.Platform)
	}
	if string(entity.Status) != "" {
		update = update.SetStatus(entity.Status)
	}
	if string(entity.ComplianceStatus) != "" {
		update = update.SetComplianceStatus(entity.ComplianceStatus)
	}
	
	update = update.SetIsEnrolled(entity.IsEnrolled)

	if entity.OwnerID != 0 {
		update = update.SetOwnerID(entity.OwnerID)
	}
	if entity.OsVersion != "" {
		update = update.SetOsVersion(entity.OsVersion)
	}
	if entity.DeviceType != "" {
		update = update.SetDeviceType(entity.DeviceType)
	}
	if entity.MACAddress != "" {
		update = update.SetMACAddress(entity.MACAddress)
	}
	if entity.IPAddress != "" {
		update = update.SetIPAddress(entity.IPAddress)
	}
	if entity.BatteryLevel > 0 {
		update = update.SetBatteryLevel(entity.BatteryLevel)
	}
	if entity.StorageCapacity > 0 {
		update = update.SetStorageCapacity(entity.StorageCapacity)
	}
	if entity.StorageUsed > 0 {
		update = update.SetStorageUsed(entity.StorageUsed)
	}
	update = update.SetIsJailbroken(entity.IsJailbroken)

	if string(entity.EnrollmentType) != "" {
		update = update.SetEnrollmentType(entity.EnrollmentType)
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

func (r *deviceRepositoryImpl) Delete(ctx context.Context, id string) error {
	err := r.client.Device.DeleteOneID(id).Exec(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Thiết bị không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi xóa thiết bị").WithError(err)
	}
	return nil
}

func (r *deviceRepositoryImpl) FindByUDID(ctx context.Context, udid string) (*ent.Device, error) {
	d, err := r.client.Device.Query().Where(device.UdidEQ(udid)).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Thiết bị không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tìm thiết bị bằng UDID").WithError(err)
	}
	return d, nil
}

func (r *deviceRepositoryImpl) FindBySerialNumber(ctx context.Context, sn string) (*ent.Device, error) {
	d, err := r.client.Device.Query().Where(device.SerialNumberEQ(sn)).Only(ctx)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Thiết bị không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tìm bằng Serial Number").WithError(err)
	}
	return d, nil
}

func (r *deviceRepositoryImpl) GetStats(ctx context.Context) (*repository.DeviceStats, error) {
	total, _ := r.client.Device.Query().Count(ctx)
	active, _ := r.client.Device.Query().Where(device.StatusEQ(device.StatusActive)).Count(ctx)
	inactive, _ := r.client.Device.Query().Where(device.StatusEQ(device.StatusInactive)).Count(ctx)
	enrolled, _ := r.client.Device.Query().Where(device.IsEnrolledEQ(true)).Count(ctx)
	compliant, _ := r.client.Device.Query().Where(device.ComplianceStatusEQ(device.ComplianceStatusCompliant)).Count(ctx)
	nonCompliant, _ := r.client.Device.Query().Where(device.ComplianceStatusEQ(device.ComplianceStatusNonCompliant)).Count(ctx)

	byPlatform := map[string]int64{
		"ios":     0,
		"android": 0,
		"windows": 0,
		"macos":   0,
		"other":   0,
	}
	ios, _ := r.client.Device.Query().Where(device.PlatformEQ(device.PlatformIos)).Count(ctx)
	android, _ := r.client.Device.Query().Where(device.PlatformEQ(device.PlatformAndroid)).Count(ctx)
	windows, _ := r.client.Device.Query().Where(device.PlatformEQ(device.PlatformWindows)).Count(ctx)
	macos, _ := r.client.Device.Query().Where(device.PlatformEQ(device.PlatformMacos)).Count(ctx)
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

	return &repository.DeviceStats{
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

func (r *deviceRepositoryImpl) GetAll(ctx context.Context) ([]*ent.Device, error) {
	devices, err := r.client.Device.Query().All(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất thiết bị").WithError(err)
	}
	return devices, nil
}

func (r *deviceRepositoryImpl) EnsureMinimalByUDID(ctx context.Context, udid string, sn string) error {
	create := r.client.Device.Create().
		SetID(uuid.NewString()).
		SetUdid(udid).
		SetStatus(device.StatusPending).
		SetIsEnrolled(false)

	if sn != "" {
		create = create.SetSerialNumber(sn)
	}

	err := create.Exec(ctx)
	// If a concurrent thread created the device just 1ms before us, 
	// ignore the constraint error. The device is safely there.
	if err != nil && ent.IsConstraintError(err) {
		return nil
	}
	return err
}

func (r *deviceRepositoryImpl) UpdateTokenEnrolledBySN(ctx context.Context, udid string, sn string, model string, osVer string) error {
	existing, err := r.client.Device.Query().
		Where(device.SerialNumberEQ(sn)).
		Only(ctx)

	if err == nil && existing != nil {
		now := time.Now()
		updater := r.client.Device.UpdateOneID(existing.ID).
			SetUdid(udid).
			SetIsEnrolled(true).
			SetStatus(device.StatusActive).
			SetEnrolledAt(now).
			SetLastSeen(now).
			SetEnrollmentType(device.EnrollmentTypeDep)

		if model != "" {
			updater = updater.SetModel(model)
		}
		if osVer != "" {
			updater = updater.SetOsVersion(osVer)
		}
		_, err = updater.Save(ctx)
		return err
	}
	return apperror.ErrNotFound.WithMessage("Không tìm thấy bằng sn")
}

func (r *deviceRepositoryImpl) UpdateTokenEnrolledByUDID(ctx context.Context, udid string, sn string, model string, osVer string) error {
	existing, err := r.client.Device.Query().Where(device.UdidEQ(udid)).Only(ctx)
	if err == nil && existing != nil {
		now := time.Now()
		updater := r.client.Device.UpdateOneID(existing.ID).
			SetIsEnrolled(true).
			SetStatus(device.StatusActive).
			SetEnrolledAt(now).
			SetLastSeen(now)
		if sn != "" {
			updater = updater.SetSerialNumber(sn)
		}
		if model != "" {
			updater = updater.SetModel(model)
		}
		if osVer != "" {
			updater = updater.SetOsVersion(osVer)
		}
		_, err = updater.Save(ctx)
		return err
	}
	return apperror.ErrNotFound.WithMessage("Không tìm thấy bằng udid")
}

func (r *deviceRepositoryImpl) CreateEnrolledDevice(ctx context.Context, udid string, sn string, model string, osVer string) error {
	now := time.Now()
	create := r.client.Device.Create().
		SetID(uuid.NewString()).
		SetUdid(udid).
		SetIsEnrolled(true).
		SetStatus(device.StatusActive).
		SetEnrolledAt(now).
		SetLastSeen(now)
	if sn != "" {
		create = create.SetSerialNumber(sn)
	}
	if model != "" {
		create = create.SetModel(model)
	}
	if osVer != "" {
		create = create.SetOsVersion(osVer)
	}
	_, err := create.Save(ctx)
	return err
}

func (r *deviceRepositoryImpl) UpdateCheckOut(ctx context.Context, udid string) error {
	_, err := r.client.Device.Update().
		Where(device.UdidEQ(udid)).
		SetIsEnrolled(false).
		SetStatus(device.StatusInactive).
		Save(ctx)
	return err
}

func (r *deviceRepositoryImpl) ApplyDeviceInformation(ctx context.Context, udid string, qr map[string]any) error {
	updater := r.client.Device.Update().
		Where(device.UdidEQ(udid)).
		SetLastSeen(time.Now())

	if v, ok := qr["OSVersion"].(string); ok && v != "" {
		updater = updater.SetOsVersion(v)
	}
	if v, ok := qr["ModelName"].(string); ok && v != "" {
		updater = updater.SetModel(v)
	}
	if v, ok := qr["DeviceName"].(string); ok && v != "" {
		updater = updater.SetName(v)
	}
	if v, ok := qr["WiFiMAC"].(string); ok && v != "" {
		updater = updater.SetMACAddress(v)
	}
	if v, ok := qr["BatteryLevel"].(float64); ok {
		updater = updater.SetBatteryLevel(v * 100)
	}

	const bytesInGB = 1024 * 1024 * 1024

	if v, ok := qr["DeviceCapacity"].(float64); ok {
		updater = updater.SetStorageCapacity(uint64(v * float64(bytesInGB)))
	}
	if avail, ok := qr["AvailableDeviceCapacity"].(float64); ok {
		if cap, ok := qr["DeviceCapacity"].(float64); ok && cap > 0 {
			updater = updater.SetStorageUsed(uint64((cap - avail) * float64(bytesInGB)))
		}
	}

	_, err := updater.Save(ctx)
	if err != nil {
		tlog.Error("Failed to apply DeviceInformation", zap.String("udid", udid), zap.Error(err))
		return err
	}
	return nil
}

func (r *deviceRepositoryImpl) UpsertFromDEP(ctx context.Context, devices []map[string]any) error {
	for _, devMap := range devices {
		sn, _ := devMap["serial_number"].(string)
		if sn == "" {
			continue
		}

		model, _ := devMap["model"].(string)
		desc, _ := devMap["description"].(string)

		existing, _ := r.client.Device.Query().Where(device.SerialNumberEQ(sn)).Only(ctx)

		if existing != nil {
			updater := r.client.Device.UpdateOne(existing)
			if model != "" {
				updater = updater.SetModel(model)
			}
			if desc != "" {
				updater = updater.SetName(desc)
			}
			updater = updater.SetEnrollmentType(device.EnrollmentTypeDep)
			if err := updater.Exec(ctx); err != nil {
				tlog.Error("Failed to update DEP device", zap.String("sn", sn), zap.Error(err))
			}
			continue
		}

		portalID := uuid.NewString()
		err := r.client.Device.Create().
			SetID(portalID).
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
	return nil
}
