package serviceimpl

import (
	"context"
	"strings"
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
	generator  service.ProfileGenerator
	mdmService service.NanoMDMService
}

func NewProfileService(repo repository.ProfileRepository, generator service.ProfileGenerator, mdmService service.NanoMDMService) service.ProfileService {
	return &profileServiceImpl{
		repo:       repo,
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
	return s.repo.Delete(ctx, id)
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

func (s *profileServiceImpl) Assign(ctx context.Context, cmd service.AssignProfileCommand) error {
	return s.repo.Assign(ctx, cmd)
}

func (s *profileServiceImpl) Unassign(ctx context.Context, profileID uint, assignmentID uint) error {
	return s.repo.Unassign(ctx, profileID, assignmentID)
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

	xmlData, err := s.generator.GenerateXML(ctx, p)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo XML profile").WithError(err)
	}

	cmdBuilder := mdmcmd.NewBuilder("")
	cmdData, _, err := cmdBuilder.InstallProfile(xmlData)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi build InstallProfile command").WithError(err)
	}

	udids, err := s.repo.GetFlattenedDeviceUDIDsByProfile(ctx, profileID)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi fetch device IDs").WithError(err)
	}

	for _, udid := range udids {
		if _, err := s.mdmService.EnqueueCommand(ctx, udid, cmdData); err != nil {
			tlog.Error("Failed to enqueue InstallProfile", zap.String("udid", udid), zap.Error(err))
			continue
		}
		_, _ = s.mdmService.Push(ctx, []string{udid})
	}

	return nil
}

func (s *profileServiceImpl) DeployToDevice(ctx context.Context, deviceID string) error {
	profiles, err := s.repo.GetProfilesByDevice(ctx, deviceID)
	if err != nil {
		return err
	}

	cmdBuilder := mdmcmd.NewBuilder("")

	for _, p := range profiles {
		xmlData, err := s.generator.GenerateXML(ctx, p)
		if err != nil {
			tlog.Error("Failed to generate XML for auto-deploy", zap.Uint("profile_id", p.ID), zap.Error(err))
			continue
		}

		cmdData, _, err := cmdBuilder.InstallProfile(xmlData)
		if err != nil {
			tlog.Error("Failed to build InstallProfile command", zap.Uint("profile_id", p.ID), zap.Error(err))
			continue
		}

		if _, err = s.mdmService.EnqueueCommand(ctx, deviceID, cmdData); err != nil {
			tlog.Error("Failed to enqueue auto-InstallProfile", zap.String("udid", deviceID), zap.Error(err))
			continue
		}

		_, _ = s.mdmService.Push(ctx, []string{deviceID})
	}

	return nil
}

func (s *profileServiceImpl) InstallOnDevice(ctx context.Context, profileID uint, deviceID string) error {
	p, err := s.repo.GetByID(ctx, profileID)
	if err != nil {
		return err
	}

	xmlData, err := s.generator.GenerateXML(ctx, p)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo XML profile").WithError(err)
	}

	cmdBuilder := mdmcmd.NewBuilder("")
	cmdData, _, err := cmdBuilder.InstallProfile(xmlData)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi build InstallProfile command").WithError(err)
	}

	if _, err = s.mdmService.EnqueueCommand(ctx, deviceID, cmdData); err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi enqueue InstallProfile").WithError(err)
	}

	// 4. Trigger APNs push
	_, _ = s.mdmService.Push(ctx, []string{deviceID})

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
