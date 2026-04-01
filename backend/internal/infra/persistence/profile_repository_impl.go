package persistence

import (
	"context"
	"strings"
	"time"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/device"
	"github.com/thienel/go-backend-template/internal/ent/devicegroup"
	"github.com/thienel/go-backend-template/internal/ent/profile"
	"github.com/thienel/go-backend-template/internal/ent/profileassignment"
	"github.com/thienel/go-backend-template/internal/ent/profiledeploymentstatus"
	"github.com/thienel/go-backend-template/internal/infra/database"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
)

type profileRepositoryImpl struct {
	client *ent.Client
}

func NewProfileRepository(client *ent.Client) repository.ProfileRepository {
	return &profileRepositoryImpl{client: client}
}

func (r *profileRepositoryImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.Profile, int64, error) {
	q := r.client.Profile.Query()

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

func (r *profileRepositoryImpl) GetByID(ctx context.Context, id uint) (*ent.Profile, error) {
	p, err := r.client.Profile.Query().
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

func (r *profileRepositoryImpl) Create(ctx context.Context, entity *ent.Profile, deviceGroupIDs []uint) (*ent.Profile, error) {
	var p *ent.Profile
	err := database.WithTx(ctx, func(tx *ent.Tx) error {
		create := tx.Profile.Create().
			SetName(entity.Name).
			SetStatus(entity.Status).
			SetVersion(entity.Version)

		if string(entity.Platform) != "" {
			create.SetPlatform(entity.Platform)
		}
		if string(entity.Scope) != "" {
			create.SetScope(entity.Scope)
		}
		if entity.SecuritySettings != nil {
			create.SetSecuritySettings(entity.SecuritySettings)
		}
		if entity.NetworkConfig != nil {
			create.SetNetworkConfig(entity.NetworkConfig)
		}
		if entity.Restrictions != nil {
			create.SetRestrictions(entity.Restrictions)
		}
		if entity.ContentFilter != nil {
			create.SetContentFilter(entity.ContentFilter)
		}
		if entity.ComplianceRules != nil {
			create.SetComplianceRules(entity.ComplianceRules)
		}
		if entity.Payloads != nil {
			create.SetPayloads(entity.Payloads)
		}

		if len(deviceGroupIDs) > 0 {
			create.AddDeviceGroupIDs(deviceGroupIDs...)
		}

		var createErr error
		p, createErr = create.Save(ctx)
		if createErr != nil {
			if ent.IsConstraintError(createErr) {
				return apperror.ErrConflict.WithMessage("Tên profile đã tồn tại")
			}
			return apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo profile").WithError(createErr)
		}
		return nil
	})

	return p, err
}

func (r *profileRepositoryImpl) Update(ctx context.Context, id uint, entity *ent.Profile, deviceGroupIDs []uint) (*ent.Profile, error) {
	var p *ent.Profile
	err := database.WithTx(ctx, func(tx *ent.Tx) error {
		update := tx.Profile.UpdateOneID(id)

		if entity.Name != "" {
			update.SetName(entity.Name)
		}
		if string(entity.Platform) != "" {
			update.SetPlatform(entity.Platform)
		}
		if string(entity.Scope) != "" {
			update.SetScope(entity.Scope)
		}
		if string(entity.Status) != "" {
			update.SetStatus(entity.Status)
		}
		if entity.SecuritySettings != nil {
			update.SetSecuritySettings(entity.SecuritySettings)
		}
		if entity.NetworkConfig != nil {
			update.SetNetworkConfig(entity.NetworkConfig)
		}
		if entity.Restrictions != nil {
			update.SetRestrictions(entity.Restrictions)
		}
		if entity.ContentFilter != nil {
			update.SetContentFilter(entity.ContentFilter)
		}
		if entity.ComplianceRules != nil {
			update.SetComplianceRules(entity.ComplianceRules)
		}
		if entity.Payloads != nil {
			update.SetPayloads(entity.Payloads)
		}

		// Propagate version increment when explicitly set by service layer
		if entity.Version > 0 {
			update.SetVersion(entity.Version)
		}

		if deviceGroupIDs != nil {
			update.ClearDeviceGroups().AddDeviceGroupIDs(deviceGroupIDs...)
		}

		var updateErr error
		p, updateErr = update.Save(ctx)
		if ent.IsConstraintError(updateErr) {
			return apperror.ErrConflict.WithMessage("Tên profile đã tồn tại")
		}
		if ent.IsNotFound(updateErr) {
			return apperror.ErrNotFound.WithMessage("Profile không tồn tại")
		}
		if updateErr != nil {
			return apperror.ErrInternalServerError.WithMessage("Lỗi khi cập nhật profile").WithError(updateErr)
		}
		return nil
	})

	return p, err
}

func (r *profileRepositoryImpl) Delete(ctx context.Context, id uint) error {
	err := r.client.Profile.DeleteOneID(id).Exec(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Profile không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi xóa profile").WithError(err)
	}
	return nil
}

func (r *profileRepositoryImpl) UpdateStatus(ctx context.Context, id uint, status string) error {
	_, err := r.client.Profile.UpdateOneID(id).SetStatus(profile.Status(status)).Save(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Profile không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi cập nhật trạng thái").WithError(err)
	}
	return nil
}

func (r *profileRepositoryImpl) SaveVersion(ctx context.Context, profileID uint, version int, data map[string]any, changeNotes string) error {
	_, err := r.client.ProfileVersion.Create().
		SetProfileID(profileID).
		SetVersion(version).
		SetData(data).
		SetChangeNotes(changeNotes).
		Save(ctx)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo bản sao lưu version").WithError(err)
	}
	return nil
}

// Assignments
func (r *profileRepositoryImpl) Assign(ctx context.Context, cmd service.AssignProfileCommand) (*ent.ProfileAssignment, error) {
	create := r.client.ProfileAssignment.Create().
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

	assignment, err := create.Save(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi gán profile").WithError(err)
	}
	return assignment, nil
}

func (r *profileRepositoryImpl) Unassign(ctx context.Context, profileID uint, assignmentID uint) error {
	err := r.client.ProfileAssignment.DeleteOneID(assignmentID).Exec(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Assignment không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi xóa assignment").WithError(err)
	}
	return nil
}

func (r *profileRepositoryImpl) ListAssignments(ctx context.Context, profileID uint) ([]*ent.ProfileAssignment, error) {
	assignments, err := r.client.ProfileAssignment.Query().
		Where(profileassignment.HasProfileWith(profile.IDEQ(profileID))).
		All(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất assignments").WithError(err)
	}
	return assignments, nil
}

// Versions
func (r *profileRepositoryImpl) ListVersions(ctx context.Context, profileID uint) ([]*ent.ProfileVersion, error) {
	p, err := r.client.Profile.Query().
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

func (r *profileRepositoryImpl) Rollback(ctx context.Context, profileID uint, versionID uint) error {
	version, err := r.client.ProfileVersion.Get(ctx, versionID)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Version không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất version").WithError(err)
	}

	update := r.client.Profile.UpdateOneID(profileID)

	if ss, ok := version.Data["security_settings"].(map[string]any); ok && ss != nil {
		update = update.SetSecuritySettings(ss)
	} else {
		update = update.ClearSecuritySettings()
	}

	if nc, ok := version.Data["network_config"].(map[string]any); ok && nc != nil {
		update = update.SetNetworkConfig(nc)
	} else {
		update = update.ClearNetworkConfig()
	}

	if res, ok := version.Data["restrictions"].(map[string]any); ok && res != nil {
		update = update.SetRestrictions(res)
	} else {
		update = update.ClearRestrictions()
	}

	if cf, ok := version.Data["content_filter"].(map[string]any); ok && cf != nil {
		update = update.SetContentFilter(cf)
	} else {
		update = update.ClearContentFilter()
	}

	if cr, ok := version.Data["compliance_rules"].(map[string]any); ok && cr != nil {
		update = update.SetComplianceRules(cr)
	} else {
		update = update.ClearComplianceRules()
	}

	if pl, ok := version.Data["payloads"].(map[string]any); ok && pl != nil {
		update = update.SetPayloads(pl)
	} else {
		update = update.ClearPayloads()
	}

	if name, ok := version.Data["name"].(string); ok && name != "" {
		update = update.SetName(name)
	}

	// Because we roll back, we must ensure Version decrements or matches the snapshot ID
	if v_id, ok := version.Data["version"].(float64); ok {
		update = update.SetVersion(int(v_id))
	}
	
	_, err = update.Save(ctx)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi rollback").WithError(err)
	}

	return nil
}

// Deployment Status
func (r *profileRepositoryImpl) CreateDeploymentStatus(ctx context.Context, profileID uint, deviceID string, status string) (*ent.ProfileDeploymentStatus, error) {
	ds, err := r.client.ProfileDeploymentStatus.Create().
		SetProfileID(profileID).
		SetDeviceID(deviceID).
		SetStatus(profiledeploymentstatus.Status(status)).
		Save(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo deployment status").WithError(err)
	}
	return ds, nil
}

func (r *profileRepositoryImpl) UpdateDeploymentStatus(ctx context.Context, id uint, status string, errorMessage string) error {
	update := r.client.ProfileDeploymentStatus.UpdateOneID(id).
		SetStatus(profiledeploymentstatus.Status(status))

	if errorMessage != "" {
		update = update.SetErrorMessage(errorMessage)
	}

	if status == "success" {
		now := time.Now()
		update = update.SetAppliedAt(now)
	}

	_, err := update.Save(ctx)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi cập nhật deployment status").WithError(err)
	}
	return nil
}

func (r *profileRepositoryImpl) GetDeploymentStatus(ctx context.Context, profileID uint) ([]*ent.ProfileDeploymentStatus, error) {
	p, err := r.client.Profile.Query().
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

func (r *profileRepositoryImpl) GetProfilesByDevice(ctx context.Context, deviceID string) ([]*ent.Profile, error) {
	profiles, err := r.client.Profile.Query().
		Where(
			profile.Or(
				profile.HasAssignmentsWith(profileassignment.DeviceIDEQ(deviceID)),
				profile.HasDeviceGroupsWith(devicegroup.HasDevicesWith(device.UdidEQ(deviceID))),
			),
		).All(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy vấn profile gán cho thiết bị").WithError(err)
	}
	return profiles, nil
}

func (r *profileRepositoryImpl) GetFlattenedDeviceUDIDsByProfile(ctx context.Context, profileID uint) ([]string, error) {
	// Directly assigned devices:
	directAssignments, err := r.client.ProfileAssignment.Query().
		Where(
			profileassignment.ProfileIDEQ(profileID),
			profileassignment.DeviceIDNotNil(),
		).All(ctx)
	if err != nil {
		return nil, err
	}

	var udids []string
	for _, a := range directAssignments {
		if a.DeviceID != nil {
			udids = append(udids, *a.DeviceID)
		}
	}

    // Group assigned devices
	groupAssignments, err := r.client.ProfileAssignment.Query().
		Where(
			profileassignment.ProfileIDEQ(profileID),
			profileassignment.GroupIDNotNil(),
		).All(ctx)
	if err != nil {
		return nil, err
	}

    if len(groupAssignments) > 0 {
		var groupIDs []uint
		for _, a := range groupAssignments {
			if a.GroupID != nil {
				groupIDs = append(groupIDs, *a.GroupID)
			}
		}

		groupDeviceUDIDs, err := r.client.Device.Query().
			Where(
				device.HasGroupsWith(devicegroup.IDIn(groupIDs...)),
			).
			Select(device.FieldUdid).
			Strings(ctx)
		if err != nil {
			return nil, err
		}

		udids = append(udids, groupDeviceUDIDs...)
    }

	// deduplicate
	unique := make(map[string]bool)
	var final []string
	for _, u := range udids {
		if !unique[u] {
			unique[u] = true
			final = append(final, u)
		}
	}

	return final, nil
}
