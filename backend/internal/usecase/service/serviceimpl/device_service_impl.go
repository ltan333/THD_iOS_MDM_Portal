package serviceimpl

import (
	"context"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/device"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/event"
	"github.com/thienel/go-backend-template/pkg/query"
	"github.com/thienel/tlog"
	"howett.net/plist"
)

type deviceServiceImpl struct {
	repo     repository.DeviceRepository
	eventBus *event.Bus
}

func NewDeviceService(repo repository.DeviceRepository, eventBus *event.Bus) service.DeviceService {
	return &deviceServiceImpl{
		repo:     repo,
		eventBus: eventBus,
	}
}

func (s *deviceServiceImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.Device, int64, error) {
	return s.repo.List(ctx, offset, limit, opts)
}

func (s *deviceServiceImpl) GetByID(ctx context.Context, id string) (*ent.Device, error) {
	if id == "" {
		return nil, apperror.ErrValidation.WithMessage("ID thiết bị là bắt buộc")
	}

	return s.repo.GetByID(ctx, id)
}

// GetUDID resolves the portal device ID to the MDM enrollment UDID.
// Handlers call this before passing an ID to nanoMDM's EnqueueCommand or Push.
func (s *deviceServiceImpl) GetUDID(ctx context.Context, id string) (string, error) {
	if id == "" {
		return "", apperror.ErrValidation.WithMessage("ID thiết bị là bắt buộc")
	}

	d, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return "", err
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

	ownerID := uint(0)
	if cmd.OwnerID != nil {
		ownerID = *cmd.OwnerID
	}

	return s.repo.Create(ctx, &ent.Device{
		ID:              cmd.ID,
		SerialNumber:    cmd.SerialNumber,
		Model:           cmd.Model,
		Name:            cmd.Name,
		MACAddress:      cmd.MacAddress,
		IPAddress:       cmd.IpAddress,
		BatteryLevel:    cmd.BatteryLevel,
		StorageCapacity: cmd.StorageCapacity,
		StorageUsed:     cmd.StorageUsed,
		IsJailbroken:    cmd.IsJailbroken,
		EnrollmentType:  device.EnrollmentType(cmd.EnrollmentType),
		Platform:        device.Platform(cmd.Platform),
		OwnerID:         ownerID,
	})
}

func (s *deviceServiceImpl) Update(ctx context.Context, cmd service.UpdateDeviceCommand) (*ent.Device, error) {
	if cmd.ID == "" {
		return nil, apperror.ErrValidation.WithMessage("ID thiết bị là bắt buộc")
	}

	updates := &ent.Device{}
	if cmd.SerialNumber != nil {
		updates.SerialNumber = *cmd.SerialNumber
	}
	if cmd.Model != nil {
		updates.Model = *cmd.Model
	}
	if cmd.Name != nil {
		updates.Name = *cmd.Name
	}
	if cmd.Platform != nil {
		updates.Platform = device.Platform(*cmd.Platform)
	}
	if cmd.Status != nil {
		updates.Status = device.Status(*cmd.Status)
	}
	if cmd.ComplianceStatus != nil {
		updates.ComplianceStatus = device.ComplianceStatus(*cmd.ComplianceStatus)
	}
	if cmd.IsEnrolled != nil {
		updates.IsEnrolled = *cmd.IsEnrolled
	}
	if cmd.OwnerID != nil {
		updates.OwnerID = *cmd.OwnerID
	}
	if cmd.OsVersion != nil {
		updates.OsVersion = *cmd.OsVersion
	}
	if cmd.DeviceType != nil {
		updates.DeviceType = *cmd.DeviceType
	}
	if cmd.MacAddress != nil {
		updates.MACAddress = *cmd.MacAddress
	}
	if cmd.IpAddress != nil {
		updates.IPAddress = *cmd.IpAddress
	}
	if cmd.BatteryLevel != nil {
		updates.BatteryLevel = *cmd.BatteryLevel
	}
	if cmd.StorageCapacity != nil {
		updates.StorageCapacity = *cmd.StorageCapacity
	}
	if cmd.StorageUsed != nil {
		updates.StorageUsed = *cmd.StorageUsed
	}
	if cmd.IsJailbroken != nil {
		updates.IsJailbroken = *cmd.IsJailbroken
	}
	if cmd.EnrollmentType != nil {
		updates.EnrollmentType = device.EnrollmentType(*cmd.EnrollmentType)
	}

	return s.repo.Update(ctx, cmd.ID, updates)
}

func (s *deviceServiceImpl) Delete(ctx context.Context, id string) error {
	if id == "" {
		return apperror.ErrValidation.WithMessage("ID thiết bị là bắt buộc")
	}
	return s.repo.Delete(ctx, id)
}

func (s *deviceServiceImpl) GetStats(ctx context.Context) (*dto.DeviceStatsResponse, error) {
	stats, err := s.repo.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.DeviceStatsResponse{
		Total:        stats.Total,
		Active:       stats.Active,
		Inactive:     stats.Inactive,
		Enrolled:     stats.Enrolled,
		ByPlatform:   stats.ByPlatform,
		ByStatus:     stats.ByStatus,
		Compliant:    stats.Compliant,
		NonCompliant: stats.NonCompliant,
	}, nil
}

func (s *deviceServiceImpl) Export(ctx context.Context, format string) ([]byte, error) {
	devices, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
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
			d, _ := s.repo.FindBySerialNumber(ctx, sn)
			if d != nil {
				return nil
			}
		}

		err := s.repo.EnsureMinimalByUDID(ctx, udid, sn)
		return err

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

		return s.repo.UpdateCheckOut(ctx, udid)

	// ------------------------------------------------------------------ //
	// mdm.Connect — device responded to a command (Acknowledge)           //
	// ------------------------------------------------------------------ //
	case "mdm.Connect":
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
	if sn != "" {
		err := s.repo.UpdateTokenEnrolledBySN(ctx, udid, sn, model, osVer)
		if err == nil {
			return nil
		}
		if !ent.IsNotFound(err) {
			return err
		}
	}

	err := s.repo.UpdateTokenEnrolledByUDID(ctx, udid, sn, model, osVer)
	if err == nil {
		return nil
	}
	if !ent.IsNotFound(err) {
		return err
	}

	return s.repo.CreateEnrolledDevice(ctx, udid, sn, model, osVer)
}

// handleAcknowledge processes mdm.Acknowledge events and dispatches
// DeviceInformation responses to update the device record.
func (s *deviceServiceImpl) handleAcknowledge(ctx context.Context, payload *dto.NanoCMDWebhook) {
	if payload.AcknowledgeEvent == nil && payload.Checkin_event == nil {
		return
	}

	ack := payload.AcknowledgeEvent
	if ack == nil {
		// Some deployments emit mdm.Connect payload fields under checkin_event.
		ack = payload.Checkin_event
	}

	requestType := deepFindString(ack, "request_type", "RequestType")
	rawRequestType, rawUDID, rawQueryResponses, rawErr := extractFromRawPayload(ack)
	if rawErr != nil {
		tlog.Warn("Failed to parse mdm.Connect raw payload",
			zap.Error(rawErr),
			zap.Any("ack_keys", mapKeys(ack)))
	}
	if requestType == "" {
		requestType = rawRequestType
	}

	queryResponses, _ := deepFindMap(ack, "query_responses", "QueryResponses")
	if len(queryResponses) == 0 {
		queryResponses = rawQueryResponses
	}

	if requestType == "" {
		// Some MDM command responses omit RequestType and only include QueryResponses.
		if len(queryResponses) == 0 {
			tlog.Debug("mdm.Connect payload missing request type",
				zap.Any("ack_keys", mapKeys(ack)))
			return
		}
		requestType = "DeviceInformation"
	}
	if requestType != "DeviceInformation" {
		return
	}

	udid := deepFindString(ack, "udid", "UDID")
	if udid == "" {
		udid = rawUDID
	}
	if udid == "" {
		udid = stringFromCheckin(payload, "udid")
	}
	if udid == "" {
		tlog.Warn("DeviceInformation acknowledge missing UDID",
			zap.Any("ack_keys", mapKeys(ack)))
		return
	}

	if len(queryResponses) == 0 {
		tlog.Warn("DeviceInformation acknowledge missing query responses",
			zap.String("udid", udid),
			zap.Any("ack_keys", mapKeys(ack)))
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
	_ = s.repo.ApplyDeviceInformation(ctx, udid, qr)
	serial := stringFromMapAny(qr, "SerialNumber", "serial_number")
	if serial != "" {
		_ = s.repo.ReconcileBySerialAndUDID(ctx, serial, udid)
	}
}

// UpsertFromDEP syncs devices discovered via the DEP API.
//
// Fix 6: each DEP device gets a stable portal UUID as its primary key.
// The `serial_number` is the hardware identity. The `udid` field starts as nil
// and is populated when the device enrolls via mdm.TokenUpdate — at which point
// the same record is updated in place without any deletion.
func (s *deviceServiceImpl) UpsertFromDEP(ctx context.Context, devices []any) error {
	var formatted []map[string]any
	for _, devAny := range devices {
		if devMap, ok := devAny.(map[string]any); ok {
			formatted = append(formatted, devMap)
		}
	}
	return s.repo.UpsertFromDEP(ctx, formatted)
}

// ------------------------------------------------------------------ //
// Helpers                                                             //
// ------------------------------------------------------------------ //

func stringFromCheckin(payload *dto.NanoCMDWebhook, key string) string {
	if payload.Checkin_event == nil {
		return ""
	}

	// 1. Try exact match (snake_case if nanoMDM ever converts them)
	if v, ok := payload.Checkin_event[key].(string); ok && v != "" {
		return v
	}

	// 2. Fallback to Apple's Raw Plist Keys (PascalCase)
	appleKeys := map[string]string{
		"udid":          "UDID",
		"serial_number": "SerialNumber",
		"model":         "Model",
		"os_version":    "OSVersion",
	}

	if appleKey, exists := appleKeys[key]; exists {
		if v, ok := payload.Checkin_event[appleKey].(string); ok && v != "" {
			return v
		}
	}

	return ""
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func stringFromMapAny(m map[string]any, keys ...string) string {
	for _, key := range keys {
		if v, ok := m[key].(string); ok && v != "" {
			return v
		}
	}
	return ""
}

func mapFromMapAny(m map[string]any, keys ...string) (map[string]any, bool) {
	for _, key := range keys {
		if v, ok := m[key]; ok {
			if mv, ok := v.(map[string]any); ok {
				return mv, true
			}
		}
	}
	return nil, false
}

func deepFindString(m map[string]any, keys ...string) string {
	for _, key := range keys {
		if v, ok := m[key].(string); ok && v != "" {
			return v
		}
	}

	for _, value := range m {
		switch vv := value.(type) {
		case map[string]any:
			if result := deepFindString(vv, keys...); result != "" {
				return result
			}
		case []any:
			for _, item := range vv {
				if nested, ok := item.(map[string]any); ok {
					if result := deepFindString(nested, keys...); result != "" {
						return result
					}
				}
			}
		}
	}

	return ""
}

func deepFindMap(m map[string]any, keys ...string) (map[string]any, bool) {
	for _, key := range keys {
		if v, ok := m[key].(map[string]any); ok {
			return v, true
		}
	}

	for _, value := range m {
		switch vv := value.(type) {
		case map[string]any:
			if result, ok := deepFindMap(vv, keys...); ok {
				return result, true
			}
		case []any:
			for _, item := range vv {
				if nested, ok := item.(map[string]any); ok {
					if result, found := deepFindMap(nested, keys...); found {
						return result, true
					}
				}
			}
		}
	}

	return nil, false
}

func mapKeys(m map[string]any) []string {
	if len(m) == 0 {
		return nil
	}

	keys := make([]string, 0, len(m))
	for k, v := range m {
		keys = append(keys, fmt.Sprintf("%s(%T)", k, v))
	}
	return keys
}

func extractFromRawPayload(ack map[string]any) (string, string, map[string]any, error) {
	rawPayload := deepFindString(ack, "raw_payload", "rawPayload", "RawPayload")
	if rawPayload == "" {
		return "", "", nil, nil
	}

	plistMap, err := parseRawPayloadPlist(rawPayload)
	if err != nil {
		return "", "", nil, err
	}

	requestType := deepFindString(plistMap, "request_type", "RequestType")
	udid := deepFindString(plistMap, "udid", "UDID")
	queryResponses, _ := deepFindMap(plistMap, "query_responses", "QueryResponses")

	return requestType, udid, queryResponses, nil
}

func parseRawPayloadPlist(rawPayload string) (map[string]any, error) {
	if m, err := unmarshalPlistToMap([]byte(rawPayload)); err == nil {
		return m, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(rawPayload)
	if err == nil {
		if m, decodeErr := unmarshalPlistToMap(decoded); decodeErr == nil {
			return m, nil
		}
	}

	return nil, fmt.Errorf("unsupported raw_payload plist format")
}

func unmarshalPlistToMap(data []byte) (map[string]any, error) {
	var out any
	if _, err := plist.Unmarshal(data, &out); err != nil {
		return nil, err
	}

	result, ok := normalizeAnyToStringMap(out).(map[string]any)
	if !ok {
		return nil, fmt.Errorf("plist root is not a dictionary")
	}

	return result, nil
}

func normalizeAnyToStringMap(v any) any {
	switch vv := v.(type) {
	case map[string]any:
		m := make(map[string]any, len(vv))
		for k, val := range vv {
			m[k] = normalizeAnyToStringMap(val)
		}
		return m
	case map[any]any:
		m := make(map[string]any, len(vv))
		for k, val := range vv {
			m[fmt.Sprint(k)] = normalizeAnyToStringMap(val)
		}
		return m
	case []any:
		arr := make([]any, 0, len(vv))
		for _, item := range vv {
			arr = append(arr, normalizeAnyToStringMap(item))
		}
		return arr
	default:
		return vv
	}
}
