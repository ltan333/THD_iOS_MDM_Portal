package serviceimpl

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/device"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/event"
	"github.com/thienel/go-backend-template/pkg/query"
	"github.com/thienel/tlog"
	"go.uber.org/zap"
)

type deviceServiceImpl struct {
	client   *ent.Client
	eventBus *event.Bus
}

// NewDeviceService creates a device service. The eventBus is used to publish
// enrollment/checkout events so that other services (profile deployment,
// inventory sync) can react without creating circular dependencies.
func NewDeviceService(client *ent.Client, eventBus *event.Bus) service.DeviceService {
	return &deviceServiceImpl{
		client:   client,
		eventBus: eventBus,
	}
}

func (s *deviceServiceImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.Device, int64, error) {
	q := s.client.Device.Query()

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

// GetUDID resolves the portal device ID to the MDM enrollment UDID.
// Handlers call this before passing an ID to nanoMDM's EnqueueCommand or Push.
func (s *deviceServiceImpl) GetUDID(ctx context.Context, id string) (string, error) {
	if id == "" {
		return "", apperror.ErrValidation.WithMessage("ID thiết bị là bắt buộc")
	}

	d, err := s.client.Device.Get(ctx, id)
	if ent.IsNotFound(err) {
		return "", apperror.ErrNotFound.WithMessage("Thiết bị không tồn tại")
	}
	if err != nil {
		return "", apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất thiết bị").WithError(err)
	}

	if d.Udid == nil || *d.Udid == "" {
		return "", apperror.ErrBadRequest.WithMessage("Thiết bị chưa enroll, không có UDID MDM")
	}

	return *d.Udid, nil
}

func (s *deviceServiceImpl) Create(ctx context.Context, cmd service.CreateDeviceCommand) (*ent.Device, error) {
	if cmd.ID == "" {
		return nil, apperror.ErrValidation.WithMessage("ID thiết bị là bắt buộc")
	}

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
		_ = writer.Write([]string{"ID", "UDID", "Serial Number", "Model", "Name", "Platform", "Status", "Is Enrolled"})
		for _, d := range devices {
			udidVal := ""
			if d.Udid != nil {
				udidVal = *d.Udid
			}
			_ = writer.Write([]string{
				d.ID,
				udidVal,
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

// HandleWebhook processes MDM check-in and acknowledge events from nanoMDM/nanoCMD.
//
// Key design after Fix 6:
//   - DEP devices are stored with a portal UUID as their ID; `udid` is nil until enrollment.
//   - mdm.TokenUpdate sets the `udid` field on the EXISTING portal record keyed by serial_number.
//     No record deletion or migration is required.
//   - Events publish the MDM UDID (not the portal UUID) so nanoMDM callers receive the right identifier.
func (s *deviceServiceImpl) HandleWebhook(ctx context.Context, payload *dto.NanoCMDWebhook) error {
	if payload.Topic == "" {
		return nil
	}

	switch payload.Topic {
	// ------------------------------------------------------------------ //
	// mdm.Authenticate — device starts enrollment handshake               //
	// ------------------------------------------------------------------ //
	case "mdm.Authenticate":
		udid := stringFromCheckin(payload, "udid")
		if udid == "" {
			return nil
		}
		sn := stringFromCheckin(payload, "serial_number")
		tlog.Info("Device authenticating", zap.String("udid", udid))

		// If a portal record keyed by serial_number already exists (DEP device),
		// do nothing — TokenUpdate will handle the UDID assignment.
		if sn != "" {
			exists, _ := s.client.Device.Query().Where(device.SerialNumberEQ(sn)).Exist(ctx)
			if exists {
				return nil
			}
		}

		// No existing record — create a minimal one keyed by a portal UUID.
		exists, _ := s.client.Device.Query().Where(device.UdidEQ(udid)).Exist(ctx)
		if !exists {
			return s.client.Device.Create().
				SetID(uuid.NewString()).
				SetUdid(udid).
				SetNillableSerialNumber(nilIfEmpty(sn)).
				SetStatus(device.StatusPending).
				SetIsEnrolled(false).
				Exec(ctx)
		}
		return nil

	// ------------------------------------------------------------------ //
	// mdm.TokenUpdate — device enrolled / token refreshed                 //
	// ------------------------------------------------------------------ //
	case "mdm.TokenUpdate":
		udid := stringFromCheckin(payload, "udid")
		if udid == "" {
			return nil
		}
		sn := stringFromCheckin(payload, "serial_number")
		model := stringFromCheckin(payload, "model")
		osVer := stringFromCheckin(payload, "os_version")

		tlog.Info("Device tokens updated/enrolled", zap.String("udid", udid))

		err := s.handleTokenUpdate(ctx, udid, sn, model, osVer)
		if err != nil {
			return err
		}

		// Publish using the UDID — all downstream MDM callers need the UDID,
		// not the portal's internal UUID.
		if s.eventBus != nil {
			s.eventBus.PublishEnrolled(event.DeviceEnrolledEvent{
				DeviceID:     udid,
				SerialNumber: sn,
			})
		}
		return nil

	// ------------------------------------------------------------------ //
	// mdm.CheckOut — device unenrolled                                    //
	// ------------------------------------------------------------------ //
	case "mdm.CheckOut":
		udid := stringFromCheckin(payload, "udid")
		if udid == "" {
			return nil
		}
		tlog.Info("Device checked out/unenrolled", zap.String("udid", udid))
		if s.eventBus != nil {
			s.eventBus.PublishCheckedOut(event.DeviceCheckedOutEvent{DeviceID: udid})
		}

		// Find the device by UDID and update enrollment state.
		d, err := s.client.Device.Query().Where(device.UdidEQ(udid)).Only(ctx)
		if err != nil {
			return nil // Device not found in portal — no-op
		}
		return s.client.Device.UpdateOneID(d.ID).
			SetIsEnrolled(false).
			SetStatus(device.StatusInactive).
			Exec(ctx)

	// ------------------------------------------------------------------ //
	// mdm.Acknowledge — device responded to a command                     //
	// ------------------------------------------------------------------ //
	case "mdm.Acknowledge":
		s.handleAcknowledge(ctx, payload)
		return nil
	}

	return nil
}

// handleTokenUpdate applies the enrollment state change to the portal DB.
//
// Fix 6 logic:
//  1. If a portal record with the matching serial_number exists → update its `udid`
//     field and mark it as enrolled. No record deletion, no migration.
//  2. If no serial_number match → look for an existing record by UDID (already has udid set).
//  3. If nothing exists → create a new record keyed by UDID (non-DEP direct enrollment).
func (s *deviceServiceImpl) handleTokenUpdate(ctx context.Context, udid, sn, model, osVer string) error {
	now := time.Now()

	// ── Case 1: DEP pre-enrollment record found by serial number ───────────
	if sn != "" {
		existing, err := s.client.Device.Query().
			Where(device.SerialNumberEQ(sn)).
			Only(ctx)

		if err == nil && existing != nil {
			// The device record already exists (created by UpsertFromDEP).
			// Simply set the UDID and update enrollment state — no deletion needed.
			updater := s.client.Device.UpdateOneID(existing.ID).
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
	}

	// ── Case 2: Already have a record with this UDID ────────────────────────
	existing, err := s.client.Device.Query().Where(device.UdidEQ(udid)).Only(ctx)
	if err == nil && existing != nil {
		updater := s.client.Device.UpdateOneID(existing.ID).
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

	// ── Case 3: No existing record — direct non-DEP enrollment ─────────────
	// All devices get a stable portal UUID regardless of enrollment method.
	create := s.client.Device.Create().
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
	_, err = create.Save(ctx)
	return err
}

// handleAcknowledge processes mdm.Acknowledge events and dispatches
// DeviceInformation responses to update the device record.
func (s *deviceServiceImpl) handleAcknowledge(ctx context.Context, payload *dto.NanoCMDWebhook) {
	if payload.AcknowledgeEvent == nil {
		return
	}

	ack := payload.AcknowledgeEvent
	requestType, _ := ack["request_type"].(string)
	if requestType != "DeviceInformation" {
		return
	}

	udid, _ := ack["udid"].(string)
	if udid == "" {
		udid = stringFromCheckin(payload, "udid")
	}
	if udid == "" {
		return
	}

	queryResponses, _ := ack["query_responses"].(map[string]any)
	if len(queryResponses) == 0 {
		return
	}

	tlog.Info("Received DeviceInformation response", zap.String("udid", udid))
	s.applyDeviceInformation(ctx, udid, queryResponses)

	if s.eventBus != nil {
		s.eventBus.PublishDeviceInformation(event.DeviceInformationReceivedEvent{
			DeviceID:       udid,
			QueryResponses: queryResponses,
		})
	}
}

// applyDeviceInformation writes DeviceInformation query responses to the device
// record identified by UDID.
func (s *deviceServiceImpl) applyDeviceInformation(ctx context.Context, udid string, qr map[string]any) {
	// Look up the portal record by UDID.
	d, err := s.client.Device.Query().Where(device.UdidEQ(udid)).Only(ctx)
	if ent.IsNotFound(err) {
		return
	}
	if err != nil {
		tlog.Error("applyDeviceInformation: device lookup failed", zap.String("udid", udid), zap.Error(err))
		return
	}

	updater := s.client.Device.UpdateOneID(d.ID).SetLastSeen(time.Now())

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
	// Apple returns BatteryLevel as 0.0–1.0; store as 0–100 percentage.
	if v, ok := qr["BatteryLevel"].(float64); ok {
		updater = updater.SetBatteryLevel(v * 100)
	}
	// Storage fields are in MB.
	if v, ok := qr["DeviceCapacity"].(float64); ok {
		updater = updater.SetStorageCapacity(uint64(v * 1024 * 1024))
	}
	if avail, ok := qr["AvailableDeviceCapacity"].(float64); ok {
		if cap, ok := qr["DeviceCapacity"].(float64); ok && cap > 0 {
			updater = updater.SetStorageUsed(uint64((cap - avail) * 1024 * 1024))
		}
	}

	if _, err := updater.Save(ctx); err != nil {
		tlog.Error("Failed to apply DeviceInformation", zap.String("udid", udid), zap.Error(err))
	}
}

// UpsertFromDEP syncs devices discovered via the DEP API.
//
// Fix 6: each DEP device gets a stable portal UUID as its primary key.
// The `serial_number` is the hardware identity. The `udid` field starts as nil
// and is populated when the device enrolls via mdm.TokenUpdate — at which point
// the same record is updated in place without any deletion.
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

		model, _ := devMap["model"].(string)
		desc, _ := devMap["description"].(string)

		// Find by serial_number — this covers both:
		//   (a) previously enrolled devices (portal record already has udid set)
		//   (b) DEP-placeholder records from a prior sync
		existing, _ := s.client.Device.Query().Where(device.SerialNumberEQ(sn)).Only(ctx)

		if existing != nil {
			// Device already known — refresh metadata from the DEP response.
			updater := s.client.Device.UpdateOne(existing)
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

		// New device — create a portal record with a stable UUID.
		// udid is left nil; it will be set when the device enrolls.
		portalID := uuid.NewString()
		err := s.client.Device.Create().
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

// ------------------------------------------------------------------ //
// Helpers                                                             //
// ------------------------------------------------------------------ //

func stringFromCheckin(payload *dto.NanoCMDWebhook, key string) string {
	if payload.Checkin_event == nil {
		return ""
	}
	v, _ := payload.Checkin_event[key].(string)
	return v
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
