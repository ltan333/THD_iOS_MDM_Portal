package serviceimpl

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/device"
	"github.com/thienel/go-backend-template/internal/ent/devicegroup"
	"github.com/thienel/go-backend-template/internal/ent/profile"
	"github.com/thienel/go-backend-template/internal/ent/profileassignment"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
	"github.com/thienel/tlog"
	"go.uber.org/zap"
)

type profileServiceImpl struct {
	client    *ent.Client
	generator service.ProfileGenerator
	mdmService service.NanoMDMService
}

func NewProfileService(client *ent.Client, generator service.ProfileGenerator, mdmService service.NanoMDMService) service.ProfileService {
	return &profileServiceImpl{
		client:    client,
		generator: generator,
		mdmService: mdmService,
	}
}

func (s *profileServiceImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.Profile, int64, error) {
	q := s.client.Profile.Query()

	for field, filter := range opts.Filters {
		switch field {
		case "search":
			if val, ok := filter.Value.(string); ok && val != "" {
				q = q.Where(profile.NameContainsFold(val))
			}
		case "platform":
			if val, ok := filter.Value.(string); ok && val != "" {
				q = q.Where(profile.PlatformEQ(profile.Platform(val)))
			}
		case "status":
			if val, ok := filter.Value.(string); ok && val != "" {
				q = q.Where(profile.StatusEQ(profile.Status(val)))
			}
		case "scope":
			if val, ok := filter.Value.(string); ok && val != "" {
				q = q.Where(profile.ScopeEQ(profile.Scope(val)))
			}
		}
	}

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, apperror.ErrInternalServerError.WithMessage("Lỗi khi đếm profile").WithError(err)
	}

	if len(opts.Sort) > 0 {
		for _, sortField := range opts.Sort {
			switch strings.ToLower(sortField.Field) {
			case "name":
				if sortField.Desc {
					q = q.Order(ent.Desc(profile.FieldName))
				} else {
					q = q.Order(ent.Asc(profile.FieldName))
				}
			case "created_at":
				if sortField.Desc {
					q = q.Order(ent.Desc(profile.FieldCreatedAt))
				} else {
					q = q.Order(ent.Asc(profile.FieldCreatedAt))
				}
			}
		}
	} else {
		q = q.Order(ent.Desc(profile.FieldCreatedAt))
	}

	profiles, err := q.Offset(offset).Limit(limit).All(ctx)
	if err != nil {
		return nil, 0, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất profile").WithError(err)
	}

	return profiles, int64(total), nil
}

func (s *profileServiceImpl) GetByID(ctx context.Context, id uint) (*ent.Profile, error) {
	if id == 0 {
		return nil, apperror.ErrValidation.WithMessage("ID profile là bắt buộc")
	}

	p, err := s.client.Profile.Query().
		Where(profile.IDEQ(id)).
		WithAssignments().
		WithVersions().
		WithDeploymentStatuses().
		First(ctx)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Profile không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất profile").WithError(err)
	}

	return p, nil
}

func (s *profileServiceImpl) Create(ctx context.Context, cmd service.CreateProfileCommand) (*ent.Profile, error) {
	if strings.TrimSpace(cmd.Name) == "" {
		return nil, apperror.ErrValidation.WithMessage("Tên profile là bắt buộc")
	}

	exists, _ := s.client.Profile.Query().Where(profile.NameEQ(cmd.Name)).Exist(ctx)
	if exists {
		return nil, apperror.ErrConflict.WithMessage("Tên profile đã tồn tại")
	}

	create := s.client.Profile.Create().
		SetName(cmd.Name).
		SetStatus(profile.StatusDraft).
		SetVersion(1)

	if cmd.Platform != "" {
		create = create.SetPlatform(profile.Platform(cmd.Platform))
	}
	if cmd.Scope != "" {
		create = create.SetScope(profile.Scope(cmd.Scope))
	}
	if cmd.SecuritySettings != nil {
		create = create.SetSecuritySettings(cmd.SecuritySettings)
	}
	if cmd.NetworkConfig != nil {
		create = create.SetNetworkConfig(cmd.NetworkConfig)
	}
	if cmd.Restrictions != nil {
		create = create.SetRestrictions(cmd.Restrictions)
	}
	if cmd.ContentFilter != nil {
		create = create.SetContentFilter(cmd.ContentFilter)
	}
	if cmd.ComplianceRules != nil {
		create = create.SetComplianceRules(cmd.ComplianceRules)
	}
	if cmd.Payloads != nil {
		create = create.SetPayloads(cmd.Payloads)
	}

	p, err := create.Save(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo profile").WithError(err)
	}

	return p, nil
}

func (s *profileServiceImpl) Update(ctx context.Context, cmd service.UpdateProfileCommand) (*ent.Profile, error) {
	if cmd.ID == 0 {
		return nil, apperror.ErrValidation.WithMessage("ID profile là bắt buộc")
	}

	// 1. Lấy dữ liệu hiện tại để tạo bản sao lưu (Version Snapshot)
	old, err := s.client.Profile.Get(ctx, cmd.ID)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Profile không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất dữ liệu cũ").WithError(err)
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

	_, err = s.client.ProfileVersion.Create().
		SetProfileID(old.ID).
		SetVersion(old.Version).
		SetData(snapshotData).
		SetChangeNotes("Cập nhật tự động tăng version").
		Save(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo bản sao lưu version").WithError(err)
	}

	// 3. Cập nhật dữ liệu mới và tăng version
	update := s.client.Profile.UpdateOneID(cmd.ID).
		SetVersion(old.Version + 1)

	if cmd.Name != nil && strings.TrimSpace(*cmd.Name) != "" {
		exists, _ := s.client.Profile.Query().
			Where(profile.NameEQ(*cmd.Name), profile.IDNEQ(cmd.ID)).
			Exist(ctx)
		if exists {
			return nil, apperror.ErrConflict.WithMessage("Tên profile đã tồn tại")
		}
		update = update.SetName(*cmd.Name)
	}
	if cmd.Platform != nil {
		update = update.SetPlatform(profile.Platform(*cmd.Platform))
	}
	if cmd.Scope != nil {
		update = update.SetScope(profile.Scope(*cmd.Scope))
	}
	if cmd.SecuritySettings != nil {
		update = update.SetSecuritySettings(cmd.SecuritySettings)
	}
	if cmd.NetworkConfig != nil {
		update = update.SetNetworkConfig(cmd.NetworkConfig)
	}
	if cmd.Restrictions != nil {
		update = update.SetRestrictions(cmd.Restrictions)
	}
	if cmd.ContentFilter != nil {
		update = update.SetContentFilter(cmd.ContentFilter)
	}
	if cmd.ComplianceRules != nil {
		update = update.SetComplianceRules(cmd.ComplianceRules)
	}
	if cmd.Payloads != nil {
		update = update.SetPayloads(cmd.Payloads)
	}

	p, err := update.Save(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi cập nhật profile").WithError(err)
	}

	return p, nil
}

func (s *profileServiceImpl) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return apperror.ErrValidation.WithMessage("ID profile là bắt buộc")
	}

	err := s.client.Profile.DeleteOneID(id).Exec(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Profile không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi xóa profile").WithError(err)
	}

	return nil
}

func (s *profileServiceImpl) UpdateStatus(ctx context.Context, id uint, status string) error {
	if id == 0 {
		return apperror.ErrValidation.WithMessage("ID profile là bắt buộc")
	}

	_, err := s.client.Profile.UpdateOneID(id).
		SetStatus(profile.Status(status)).
		Save(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Profile không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi cập nhật trạng thái").WithError(err)
	}

	return nil
}

func (s *profileServiceImpl) UpdateSecuritySettings(ctx context.Context, id uint, settings map[string]any) error {
	_, err := s.client.Profile.UpdateOneID(id).
		SetSecuritySettings(settings).
		Save(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Profile không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi cập nhật security settings").WithError(err)
	}
	return nil
}

func (s *profileServiceImpl) UpdateNetworkConfig(ctx context.Context, id uint, config map[string]any) error {
	_, err := s.client.Profile.UpdateOneID(id).
		SetNetworkConfig(config).
		Save(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Profile không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi cập nhật network config").WithError(err)
	}
	return nil
}

func (s *profileServiceImpl) UpdateRestrictions(ctx context.Context, id uint, restrictions map[string]any) error {
	_, err := s.client.Profile.UpdateOneID(id).
		SetRestrictions(restrictions).
		Save(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Profile không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi cập nhật restrictions").WithError(err)
	}
	return nil
}

func (s *profileServiceImpl) UpdateContentFilter(ctx context.Context, id uint, filter map[string]any) error {
	_, err := s.client.Profile.UpdateOneID(id).
		SetContentFilter(filter).
		Save(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Profile không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi cập nhật content filter").WithError(err)
	}
	return nil
}

func (s *profileServiceImpl) UpdateComplianceRules(ctx context.Context, id uint, rules map[string]any) error {
	_, err := s.client.Profile.UpdateOneID(id).
		SetComplianceRules(rules).
		Save(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Profile không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi cập nhật compliance rules").WithError(err)
	}
	return nil
}

func (s *profileServiceImpl) Assign(ctx context.Context, cmd service.AssignProfileCommand) error {
	create := s.client.ProfileAssignment.Create().
		SetProfileID(cmd.ProfileID).
		SetTargetType(profileassignment.TargetType(cmd.TargetType)).
		SetScheduleType(profileassignment.ScheduleType(cmd.ScheduleType)).
		SetNillableScheduledAt(cmd.ScheduledAt)

	if cmd.DeviceID != nil {
		create = create.SetDeviceID(*cmd.DeviceID)
	}
	if cmd.GroupID != nil {
		create = create.SetGroupID(*cmd.GroupID)
	}

	_, err := create.Save(ctx)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi gán profile").WithError(err)
	}
	return nil
}

func (s *profileServiceImpl) Unassign(ctx context.Context, profileID uint, assignmentID uint) error {
	err := s.client.ProfileAssignment.DeleteOneID(assignmentID).Exec(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Assignment không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi xóa assignment").WithError(err)
	}
	return nil
}

func (s *profileServiceImpl) ListAssignments(ctx context.Context, profileID uint) ([]*ent.ProfileAssignment, error) {
	assignments, err := s.client.ProfileAssignment.Query().
		Where(profileassignment.HasProfileWith(profile.IDEQ(profileID))).
		All(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất assignments").WithError(err)
	}
	return assignments, nil
}

func (s *profileServiceImpl) ListVersions(ctx context.Context, profileID uint) ([]*ent.ProfileVersion, error) {
	p, err := s.client.Profile.Query().
		Where(profile.IDEQ(profileID)).
		WithVersions().
		First(ctx)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Profile không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất versions").WithError(err)
	}
	return p.Edges.Versions, nil
}

func (s *profileServiceImpl) Rollback(ctx context.Context, profileID uint, versionID uint) error {
	version, err := s.client.ProfileVersion.Get(ctx, versionID)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Version không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất version").WithError(err)
	}

	_, err = s.client.Profile.UpdateOneID(profileID).
		SetSecuritySettings(version.Data["security_settings"].(map[string]any)).
		SetNetworkConfig(version.Data["network_config"].(map[string]any)).
		SetRestrictions(version.Data["restrictions"].(map[string]any)).
		Save(ctx)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi rollback").WithError(err)
	}

	return nil
}

func (s *profileServiceImpl) GetDeploymentStatus(ctx context.Context, profileID uint) ([]*ent.ProfileDeploymentStatus, error) {
	p, err := s.client.Profile.Query().
		Where(profile.IDEQ(profileID)).
		WithDeploymentStatuses().
		First(ctx)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Profile không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất deployment status").WithError(err)
	}
	return p.Edges.DeploymentStatuses, nil
}

func (s *profileServiceImpl) Repush(ctx context.Context, profileID uint) error {
	// 1. Fetch profile
	p, err := s.client.Profile.Get(ctx, profileID)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Profile không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất profile").WithError(err)
	}

	// 2. Generate XML
	xmlData, err := s.generator.GenerateXML(ctx, p)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo XML profile").WithError(err)
	}

	// 3. Find all assigned devices
	assignments, err := s.client.ProfileAssignment.Query().
		Where(profileassignment.ProfileIDEQ(profileID), profileassignment.DeviceIDNEQ("")).
		All(ctx)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất danh sách máy gán").WithError(err)
	}

	// 4. Enqueue InstallProfile command for each device
	installProfileXML := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Command</key>
	<dict>
		<key>Payload</key>
		<data>%s</data>
		<key>RequestType</key>
		<string>InstallProfile</string>
	</dict>
	<key>CommandUUID</key>
	<string>InstallProfile-%d-%s</string>
</dict>
</plist>`

	encodedProfile := base64.StdEncoding.EncodeToString(xmlData)

	for _, as := range assignments {
		if as.DeviceID == nil {
			continue
		}
		udid := *as.DeviceID
		finalXML := fmt.Sprintf(installProfileXML, encodedProfile, p.ID, udid)
		_, err := s.mdmService.EnqueueCommand(ctx, udid, []byte(finalXML))
		if err != nil {
			tlog.Error("Failed to enqueue InstallProfile", zap.String("udid", udid), zap.Error(err))
			continue
		}
		
		// 5. Trigger Push
		_, _ = s.mdmService.Push(ctx, []string{udid})
	}

	return nil
}

func (s *profileServiceImpl) DeployToDevice(ctx context.Context, deviceID string) error {
	// 1. Find all profiles assigned to this device (directly or via group)
	profiles, err := s.client.Profile.Query().
		Where(
			profile.Or(
				profile.HasAssignmentsWith(profileassignment.DeviceIDEQ(deviceID)),
				profile.HasDeviceGroupsWith(devicegroup.HasDevicesWith(device.IDEQ(deviceID))),
			),
		).All(ctx)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi truy vấn profile gán cho thiết bị").WithError(err)
	}

	for _, p := range profiles {
		// 2. Generate XML
		xmlData, err := s.generator.GenerateXML(ctx, p)
		if err != nil {
			tlog.Error("Failed to generate XML for auto-deploy", zap.Uint("profile_id", p.ID), zap.Error(err))
			continue
		}

		// 3. Enqueue InstallProfile command
		installProfileXML := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Command</key>
	<dict>
		<key>Payload</key>
		<data>%s</data>
		<key>RequestType</key>
		<string>InstallProfile</string>
	</dict>
	<key>CommandUUID</key>
	<string>InstallProfile-%d-%s</string>
</dict>
</plist>`

		encodedProfile := base64.StdEncoding.EncodeToString(xmlData)
		finalXML := fmt.Sprintf(installProfileXML, encodedProfile, p.ID, deviceID)

		_, err = s.mdmService.EnqueueCommand(ctx, deviceID, []byte(finalXML))
		if err != nil {
			tlog.Error("Failed to enqueue auto-InstallProfile", zap.String("udid", deviceID), zap.Error(err))
			continue
		}

		// 4. Trigger Push
		_, _ = s.mdmService.Push(ctx, []string{deviceID})
	}

	return nil
}

func (s *profileServiceImpl) Duplicate(ctx context.Context, profileID uint) (*ent.Profile, error) {
	original, err := s.client.Profile.Get(ctx, profileID)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Profile không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất profile").WithError(err)
	}

	newName := original.Name + " (Copy) " + time.Now().Format("20060102150405")

	duplicate, err := s.client.Profile.Create().
		SetName(newName).
		SetPlatform(original.Platform).
		SetScope(original.Scope).
		SetStatus(profile.StatusDraft).
		SetSecuritySettings(original.SecuritySettings).
		SetNetworkConfig(original.NetworkConfig).
		SetRestrictions(original.Restrictions).
		SetContentFilter(original.ContentFilter).
		SetComplianceRules(original.ComplianceRules).
		SetPayloads(original.Payloads).
		SetVersion(1).
		Save(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi duplicate profile").WithError(err)
	}

	return duplicate, nil
}
