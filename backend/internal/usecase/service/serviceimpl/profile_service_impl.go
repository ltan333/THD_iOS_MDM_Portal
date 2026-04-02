package serviceimpl

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/profile"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/mdmcmd"
	"github.com/thienel/go-backend-template/pkg/query"
	"github.com/thienel/tlog"
	"go.uber.org/zap"
)

type profileServiceImpl struct {
	repo       repository.ProfileRepository
	deviceRepo repository.DeviceRepository
	generator  service.ProfileGenerator
	mdmService service.NanoMDMService
	// pendingCmds maps commandUUID → deploymentStatusID for in-flight InstallProfile commands.
	// Used to correlate mdm.Connect ACKs back to the correct deployment status record
	// when NanoMDM omits request_type from the webhook.
	pendingCmds sync.Map
}

func NewProfileService(repo repository.ProfileRepository, deviceRepo repository.DeviceRepository, generator service.ProfileGenerator, mdmService service.NanoMDMService) service.ProfileService {
	return &profileServiceImpl{
		repo:       repo,
		deviceRepo: deviceRepo,
		generator:  generator,
		mdmService: mdmService,
	}
}

func (s *profileServiceImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.Profile, int64, error) {
	return s.repo.List(ctx, offset, limit, opts)
}

func (s *profileServiceImpl) GetByID(ctx context.Context, id uint) (*ent.Profile, error) {
	if id == 0 {
		return nil, apperror.ErrValidation.WithMessage("ID profile là bắt buộc")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *profileServiceImpl) Create(ctx context.Context, cmd service.CreateProfileCommand) (*ent.Profile, error) {
	if strings.TrimSpace(cmd.Name) == "" {
		return nil, apperror.ErrValidation.WithMessage("Tên profile là bắt buộc")
	}

	return s.repo.Create(ctx, &ent.Profile{
		Name:             cmd.Name,
		Platform:         profile.Platform(cmd.Platform),
		Scope:            profile.Scope(cmd.Scope),
		Status:           profile.StatusDraft,
		Version:          1,
		SecuritySettings: cmd.SecuritySettings,
		NetworkConfig:    cmd.NetworkConfig,
		Restrictions:     cmd.Restrictions,
		ContentFilter:    cmd.ContentFilter,
		ComplianceRules:  cmd.ComplianceRules,
		Payloads:         cmd.Payloads,
	}, nil) // Removed cmd.DeviceGroupIDs since CreateProfileCommand doesn't have it
}

func (s *profileServiceImpl) Update(ctx context.Context, cmd service.UpdateProfileCommand) (*ent.Profile, error) {
	if cmd.ID == 0 {
		return nil, apperror.ErrValidation.WithMessage("ID profile là bắt buộc")
	}

	// 1. Lấy dữ liệu hiện tại để tạo bản sao lưu (Version Snapshot)
	old, err := s.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	// 2. Tạo bản ghi version mới
	snapshotData := map[string]any{
		"name":              old.Name,
		"platform":          string(old.Platform),
		"scope":             string(old.Scope),
		"security_settings": old.SecuritySettings,
		"network_config":    old.NetworkConfig,
		"restrictions":      old.Restrictions,
		"content_filter":    old.ContentFilter,
		"compliance_rules":  old.ComplianceRules,
		"payloads":          old.Payloads,
	}

	err = s.repo.SaveVersion(ctx, old.ID, old.Version, snapshotData, "Cập nhật tự động tăng version")
	if err != nil {
		return nil, err
	}

	// 3. Cập nhật dữ liệu mới và tăng version
	name := ""
	if cmd.Name != nil && strings.TrimSpace(*cmd.Name) != "" {
		name = *cmd.Name
	}

	platform := profile.Platform("")
	if cmd.Platform != nil {
		platform = profile.Platform(*cmd.Platform)
	}
	scope := profile.Scope("")
	if cmd.Scope != nil {
		scope = profile.Scope(*cmd.Scope)
	}

	return s.repo.Update(ctx, cmd.ID, &ent.Profile{
		Name:             name,
		Platform:         platform,
		Scope:            scope,
		SecuritySettings: cmd.SecuritySettings,
		NetworkConfig:    cmd.NetworkConfig,
		Restrictions:     cmd.Restrictions,
		ContentFilter:    cmd.ContentFilter,
		ComplianceRules:  cmd.ComplianceRules,
		Payloads:         cmd.Payloads,
	}, nil) // Assume deviceGroupIDs doesn't change here unless explicitly specified in cmd (it wasn't in original Update logic)
}

func (s *profileServiceImpl) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return apperror.ErrValidation.WithMessage("ID profile là bắt buộc")
	}

	udidsBefore, err := s.repo.GetFlattenedDeviceUDIDsByProfile(ctx, id)
	if err != nil {
		tlog.Error("Failed to fetch UDIDs before delete", zap.Error(err))
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	if len(udidsBefore) > 0 {
		identifier := s.generator.GetProfileIdentifier(id)
		cmdBuilder := mdmcmd.NewBuilder("")
		if cmdData, _, err := cmdBuilder.RemoveProfile(identifier); err == nil {
			for _, udid := range udidsBefore {
				if _, err := s.mdmService.EnqueueCommand(ctx, udid, cmdData); err != nil {
					tlog.Error("Failed to enqueue RemoveProfile on delete", zap.String("udid", udid), zap.Error(err))
				}
			}
			_, _ = s.mdmService.Push(ctx, udidsBefore)
		} else {
			tlog.Error("Failed to build RemoveProfile MDM command on delete", zap.Error(err))
		}
	}

	return nil
}

func (s *profileServiceImpl) UpdateStatus(ctx context.Context, id uint, status string) error {
	if id == 0 {
		return apperror.ErrValidation.WithMessage("ID profile là bắt buộc")
	}
	return s.repo.UpdateStatus(ctx, id, status)
}

// snapshotAndIncrement fetches the current profile state and saves a version snapshot
// before any granular mutator overwrites it. Ensures every field-level update is reversible.
func (s *profileServiceImpl) snapshotAndIncrement(ctx context.Context, id uint) (*ent.Profile, error) {
	old, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	snapshotData := map[string]any{
		"name":              old.Name,
		"platform":          string(old.Platform),
		"scope":             string(old.Scope),
		"security_settings": old.SecuritySettings,
		"network_config":    old.NetworkConfig,
		"restrictions":      old.Restrictions,
		"content_filter":    old.ContentFilter,
		"compliance_rules":  old.ComplianceRules,
		"payloads":          old.Payloads,
	}
	if err = s.repo.SaveVersion(ctx, old.ID, old.Version, snapshotData, "Cập nhật qua Granular Mutator"); err != nil {
		return nil, err
	}
	return old, nil
}

func (s *profileServiceImpl) UpdateSecuritySettings(ctx context.Context, id uint, settings map[string]any) error {
	old, err := s.snapshotAndIncrement(ctx, id)
	if err != nil {
		return err
	}
	_, err = s.repo.Update(ctx, id, &ent.Profile{SecuritySettings: settings, Version: old.Version + 1}, nil)
	return err
}

func (s *profileServiceImpl) UpdateNetworkConfig(ctx context.Context, id uint, config map[string]any) error {
	old, err := s.snapshotAndIncrement(ctx, id)
	if err != nil {
		return err
	}
	_, err = s.repo.Update(ctx, id, &ent.Profile{NetworkConfig: config, Version: old.Version + 1}, nil)
	return err
}

func (s *profileServiceImpl) UpdateRestrictions(ctx context.Context, id uint, restrictions map[string]any) error {
	old, err := s.snapshotAndIncrement(ctx, id)
	if err != nil {
		return err
	}
	_, err = s.repo.Update(ctx, id, &ent.Profile{Restrictions: restrictions, Version: old.Version + 1}, nil)
	return err
}

func (s *profileServiceImpl) UpdateContentFilter(ctx context.Context, id uint, filter map[string]any) error {
	old, err := s.snapshotAndIncrement(ctx, id)
	if err != nil {
		return err
	}
	_, err = s.repo.Update(ctx, id, &ent.Profile{ContentFilter: filter, Version: old.Version + 1}, nil)
	return err
}

func (s *profileServiceImpl) UpdateComplianceRules(ctx context.Context, id uint, rules map[string]any) error {
	old, err := s.snapshotAndIncrement(ctx, id)
	if err != nil {
		return err
	}
	_, err = s.repo.Update(ctx, id, &ent.Profile{ComplianceRules: rules, Version: old.Version + 1}, nil)
	return err
}

func (s *profileServiceImpl) Assign(ctx context.Context, cmd service.AssignProfileCommand) (*ent.ProfileAssignment, error) {
	assignment, err := s.repo.Assign(ctx, cmd)
	if err != nil {
		return nil, err
	}
	// Immediately deploy if schedule_type is "immediate" and a direct device is targeted.
	if cmd.ScheduleType == "immediate" {
		if cmd.DeviceID != nil && *cmd.DeviceID != "" {
			if err := s.InstallOnDevice(ctx, cmd.ProfileID, *cmd.DeviceID); err != nil {
				// Non-fatal: log but do not roll back the assignment record.
				tlog.Error("Immediate deploy on assign failed",
					zap.Uint("profile_id", cmd.ProfileID),
					zap.String("device_id", *cmd.DeviceID),
					zap.Error(err))
			}
		}
	}
	return assignment, nil
}

func (s *profileServiceImpl) Unassign(ctx context.Context, profileID uint, assignmentID uint) error {
	udidsBefore, err := s.repo.GetFlattenedDeviceUDIDsByProfile(ctx, profileID)
	if err != nil {
		tlog.Error("Failed to fetch UDIDs before unassign", zap.Error(err))
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi load danh sách device").WithError(err)
	}

	if err := s.repo.Unassign(ctx, profileID, assignmentID); err != nil {
		return err
	}

	udidsAfter, err := s.repo.GetFlattenedDeviceUDIDsByProfile(ctx, profileID)
	if err != nil {
		tlog.Error("Failed to fetch UDIDs after unassign", zap.Error(err))
		return nil
	}

	afterMap := make(map[string]struct{}, len(udidsAfter))
	for _, u := range udidsAfter {
		afterMap[u] = struct{}{}
	}

	var udidsToRemove []string
	for _, u := range udidsBefore {
		if _, exists := afterMap[u]; !exists {
			udidsToRemove = append(udidsToRemove, u)
		}
	}

	if len(udidsToRemove) > 0 {
		identifier := s.generator.GetProfileIdentifier(profileID)
		cmdBuilder := mdmcmd.NewBuilder("")
		if cmdData, _, err := cmdBuilder.RemoveProfile(identifier); err == nil {
			for _, udid := range udidsToRemove {
				if _, err := s.mdmService.EnqueueCommand(ctx, udid, cmdData); err != nil {
					tlog.Error("Failed to enqueue RemoveProfile on unassign", zap.String("udid", udid), zap.Error(err))
				}
			}
			_, _ = s.mdmService.Push(ctx, udidsToRemove)
		} else {
			tlog.Error("Failed to build RemoveProfile MDM command on unassign", zap.Error(err))
		}
	}

	return nil
}

func (s *profileServiceImpl) ListAssignments(ctx context.Context, profileID uint) ([]*ent.ProfileAssignment, error) {
	return s.repo.ListAssignments(ctx, profileID)
}

func (s *profileServiceImpl) ListVersions(ctx context.Context, profileID uint) ([]*ent.ProfileVersion, error) {
	return s.repo.ListVersions(ctx, profileID)
}

func (s *profileServiceImpl) Rollback(ctx context.Context, profileID uint, versionID uint) error {
	return s.repo.Rollback(ctx, profileID, versionID)
}

func (s *profileServiceImpl) GetDeploymentStatus(ctx context.Context, profileID uint) ([]*ent.ProfileDeploymentStatus, error) {
	return s.repo.GetDeploymentStatus(ctx, profileID)
}

func (s *profileServiceImpl) Repush(ctx context.Context, profileID uint) error {
	p, err := s.repo.GetByID(ctx, profileID)
	if err != nil {
		return err
	}

	if p.Status != profile.StatusActive {
		return apperror.ErrBadRequest.WithMessage("Chỉ có thể repush profile đang active")
	}

	xmlData, err := s.generator.GenerateXML(ctx, p)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo XML profile").WithError(err)
	}

	udids, err := s.repo.GetFlattenedDeviceUDIDsByProfile(ctx, profileID)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi fetch device IDs").WithError(err)
	}

	cmdBuilder := mdmcmd.NewBuilder("")
	var failed int
	for _, udid := range udids {
		// Each device gets its own command with a unique CommandUUID
		cmdData, commandUUID, err := cmdBuilder.InstallProfile(xmlData)
		if err != nil {
			tlog.Error("Failed to build InstallProfile command", zap.String("udid", udid), zap.Error(err))
			failed++
			continue
		}

		// Resolve UDID → portal device ID to create deployment status
		device, err := s.deviceRepo.FindByUDID(ctx, udid)
		if err != nil {
			tlog.Warn("Repush: device not found for UDID", zap.String("udid", udid), zap.Error(err))
			failed++
			continue
		}

		ds, err := s.repo.CreateDeploymentStatus(ctx, profileID, device.ID, "pending")
		if err != nil {
			tlog.Error("Repush: failed to create deployment status", zap.String("udid", udid), zap.Error(err))
		}

		if _, err := s.mdmService.EnqueueCommand(ctx, udid, cmdData); err != nil {
			tlog.Error("Failed to enqueue InstallProfile", zap.String("udid", udid), zap.Error(err))
			if ds != nil {
				_ = s.repo.UpdateDeploymentStatus(ctx, ds.ID, "failed", err.Error())
			}
			failed++
			continue
		}

		if ds != nil && commandUUID != "" {
			s.pendingCmds.Store(commandUUID, ds.ID)
		}

		_, _ = s.mdmService.Push(ctx, []string{udid})
	}

	if len(udids) > 0 && failed == len(udids) {
		return apperror.ErrInternalServerError.WithMessage(
			fmt.Sprintf("Repush failed for all %d devices", failed))
	}
	return nil
}

func (s *profileServiceImpl) DeployToDevice(ctx context.Context, deviceID string) error {
	// deviceID is the MDM UDID (called from enrollment worker with ev.DeviceID)
	profiles, err := s.repo.GetProfilesByDevice(ctx, deviceID)
	if err != nil {
		return err
	}

	// Resolve UDID → portal device ID once for all profiles
	device, err := s.deviceRepo.FindByUDID(ctx, deviceID)
	if err != nil {
		tlog.Warn("DeployToDevice: device not found for UDID", zap.String("udid", deviceID), zap.Error(err))
		return nil
	}

	cmdBuilder := mdmcmd.NewBuilder("")

	for _, p := range profiles {
		// Only deploy active profiles
		if p.Status != profile.StatusActive {
			continue
		}

		xmlData, err := s.generator.GenerateXML(ctx, p)
		if err != nil {
			tlog.Error("Failed to generate XML for auto-deploy", zap.Uint("profile_id", p.ID), zap.Error(err))
			continue
		}

		// Each profile gets its own command with a unique CommandUUID
		cmdData, commandUUID, err := cmdBuilder.InstallProfile(xmlData)
		if err != nil {
			tlog.Error("Failed to build InstallProfile command", zap.Uint("profile_id", p.ID), zap.Error(err))
			continue
		}

		ds, err := s.repo.CreateDeploymentStatus(ctx, p.ID, device.ID, "pending")
		if err != nil {
			tlog.Error("DeployToDevice: failed to create deployment status",
				zap.Uint("profile_id", p.ID), zap.String("udid", deviceID), zap.Error(err))
		}

		if _, err = s.mdmService.EnqueueCommand(ctx, deviceID, cmdData); err != nil {
			tlog.Error("Failed to enqueue auto-InstallProfile", zap.String("udid", deviceID), zap.Error(err))
			if ds != nil {
				_ = s.repo.UpdateDeploymentStatus(ctx, ds.ID, "failed", err.Error())
			}
			continue
		}

		if ds != nil && commandUUID != "" {
			s.pendingCmds.Store(commandUUID, ds.ID)
		}

		_, _ = s.mdmService.Push(ctx, []string{deviceID})
	}

	return nil
}

// InstallOnDevice installs a profile on the device identified by its portal device ID.
// It resolves the MDM UDID internally before enqueuing commands to NanoMDM.
func (s *profileServiceImpl) InstallOnDevice(ctx context.Context, profileID uint, portalDeviceID string) error {
	p, err := s.repo.GetByID(ctx, profileID)
	if err != nil {
		return err
	}

	// Resolve portal device ID → MDM UDID
	device, err := s.deviceRepo.GetByID(ctx, portalDeviceID)
	if err != nil {
		return apperror.ErrNotFound.WithMessage("Không tìm thấy device").WithError(err)
	}
	if device.Udid == nil || *device.Udid == "" {
		return apperror.ErrBadRequest.WithMessage("Device chưa enroll vào MDM, không có UDID")
	}
	udid := *device.Udid

	// 1. Create deployment status record (FK references portal_devices.id)
	ds, err := s.repo.CreateDeploymentStatus(ctx, profileID, portalDeviceID, "pending")
	if err != nil {
		tlog.Error("Failed to create deployment status", zap.Error(err))
		// Continue anyway - don't block installation
	}

	// 2. Generate XML
	xmlData, err := s.generator.GenerateXML(ctx, p)
	if err != nil {
		if ds != nil {
			_ = s.repo.UpdateDeploymentStatus(ctx, ds.ID, "failed", "Failed to generate XML: "+err.Error())
		}
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo XML profile").WithError(err)
	}

	// 3. Build command — capture commandUUID to correlate the device ACK later
	cmdBuilder := mdmcmd.NewBuilder("")
	cmdData, commandUUID, err := cmdBuilder.InstallProfile(xmlData)
	if err != nil {
		if ds != nil {
			_ = s.repo.UpdateDeploymentStatus(ctx, ds.ID, "failed", "Failed to build command: "+err.Error())
		}
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi build InstallProfile command").WithError(err)
	}

	// 4. Enqueue command via MDM UDID
	if _, err = s.mdmService.EnqueueCommand(ctx, udid, cmdData); err != nil {
		if ds != nil {
			_ = s.repo.UpdateDeploymentStatus(ctx, ds.ID, "failed", "Failed to enqueue command: "+err.Error())
		}
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi enqueue InstallProfile").WithError(err)
	}

	// 5. Store commandUUID → deploymentStatusID so the ACK webhook can update the record
	if ds != nil && commandUUID != "" {
		s.pendingCmds.Store(commandUUID, ds.ID)
	}

	// 6. Trigger APNs push
	_, _ = s.mdmService.Push(ctx, []string{udid})

	if ds != nil {
		tlog.Info("Profile installation queued",
			zap.Uint("profile_id", profileID),
			zap.String("portal_device_id", portalDeviceID),
			zap.String("udid", udid),
			zap.String("command_uuid", commandUUID),
			zap.Uint("deployment_status_id", ds.ID))
	}

	return nil
}

// HandleInstallAck processes an InstallProfile ACK from a device. It resolves
// HandleInstallAck processes an InstallProfile command ACK.
// It first tries to match by commandUUID (precise). If commandUUID is not in
// the pending map (e.g. after a server restart), it falls back to updating all
// pending deployment statuses for the device identified by UDID.
func (s *profileServiceImpl) HandleInstallAck(ctx context.Context, udid string, commandUUID string, ackStatus string, errMsg string) error {
	// NotNow means the device is busy and will retry later — leave status as pending
	// and keep commandUUID in the map so the eventual real ACK is matched correctly.
	if ackStatus == "NotNow" {
		tlog.Info("InstallProfile deferred by device (NotNow), waiting for retry",
			zap.String("udid", udid),
			zap.String("command_uuid", commandUUID))
		return nil
	}

	repoStatus := "failed"
	if ackStatus == "Acknowledged" {
		repoStatus = "success"
	}

	// Primary path: look up exact deployment status by commandUUID
	if commandUUID != "" {
		if val, ok := s.pendingCmds.LoadAndDelete(commandUUID); ok {
			dsID := val.(uint)
			if err := s.repo.UpdateDeploymentStatus(ctx, dsID, repoStatus, errMsg); err != nil {
				tlog.Error("HandleInstallAck: failed to update deployment status by commandUUID",
					zap.String("command_uuid", commandUUID),
					zap.Uint("deployment_status_id", dsID),
					zap.Error(err))
				return err
			}
			tlog.Info("Deployment status updated from ACK",
				zap.String("command_uuid", commandUUID),
				zap.Uint("deployment_status_id", dsID),
				zap.String("status", repoStatus),
				zap.String("error_msg", errMsg))
			return nil
		}
	}

	// Fallback: commandUUID not in pending map (server restarted between enqueue and ACK).
	// Update all pending statuses for this device.
	if udid == "" {
		return nil
	}
	device, err := s.deviceRepo.FindByUDID(ctx, udid)
	if err != nil {
		tlog.Warn("HandleInstallAck: device not found for UDID",
			zap.String("udid", udid), zap.Error(err))
		return nil
	}
	if err := s.repo.UpdateDeploymentStatusByDevice(ctx, device.ID, repoStatus, errMsg); err != nil {
		tlog.Error("HandleInstallAck: fallback update failed",
			zap.String("udid", udid),
			zap.String("portal_device_id", device.ID),
			zap.Error(err))
		return err
	}
	tlog.Info("Deployment status updated from ACK (fallback by device)",
		zap.String("udid", udid),
		zap.String("portal_device_id", device.ID),
		zap.String("status", repoStatus))
	return nil
}

func (s *profileServiceImpl) Duplicate(ctx context.Context, id uint) (*ent.Profile, error) {
	if id == 0 {
		return nil, apperror.ErrValidation.WithMessage("ID profile là bắt buộc")
	}

	original, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	newName := original.Name + " (Copy) " + time.Now().Format("20060102150405")

	newProfile := &ent.Profile{
		Name:             newName,
		Platform:         original.Platform,
		Scope:            original.Scope,
		Status:           profile.StatusDraft,
		Version:          1,
		SecuritySettings: original.SecuritySettings,
		NetworkConfig:    original.NetworkConfig,
		Restrictions:     original.Restrictions,
		ContentFilter:    original.ContentFilter,
		ComplianceRules:  original.ComplianceRules,
		Payloads:         original.Payloads,
	}

	duplicate, err := s.repo.Create(ctx, newProfile, nil)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi duplicate profile").WithError(err)
	}

	return duplicate, nil
}
