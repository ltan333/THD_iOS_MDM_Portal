package persistence

import (
	"context"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/depprofile"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
)

type depProfileRepositoryImpl struct {
	client *ent.Client
}

// NewDepProfileRepository creates a new DEP profile repository
func NewDepProfileRepository(client *ent.Client) repository.DepProfileRepository {
	return &depProfileRepositoryImpl{client: client}
}

// --- BaseRepository Methods ---

func (r *depProfileRepositoryImpl) Create(ctx context.Context, e *ent.DepProfile) error {
	create := r.client.DepProfile.Create().
		SetProfileName(e.ProfileName).
		SetAllowPairing(e.AllowPairing).
		SetAnchorCerts(e.AnchorCerts).
		SetAutoAdvanceSetup(e.AutoAdvanceSetup).
		SetAwaitDeviceConfigured(e.AwaitDeviceConfigured).
		SetConfigurationWebURL(e.ConfigurationWebURL).
		SetDepartment(e.Department).
		SetDevices(e.Devices).
		SetDoNotUseProfileFromBackup(e.DoNotUseProfileFromBackup).
		SetIsReturnToService(e.IsReturnToService).
		SetIsMandatory(e.IsMandatory).
		SetIsMdmRemovable(e.IsMdmRemovable).
		SetIsMultiUser(e.IsMultiUser).
		SetIsSupervised(e.IsSupervised).
		SetLanguage(e.Language).
		SetOrgMagic(e.OrgMagic).
		SetRegion(e.Region).
		SetSkipSetupItems(e.SkipSetupItems).
		SetSupervisingHostCerts(e.SupervisingHostCerts).
		SetSupportEmailAddress(e.SupportEmailAddress).
		SetSupportPhoneNumber(e.SupportPhoneNumber).
		SetURL(e.URL).
		SetProfileData(e.ProfileData)

	if e.ProfileUUID != "" {
		create.SetProfileUUID(e.ProfileUUID)
	}

	u, err := create.Save(ctx)
	if err != nil {
		return wrapCreateError(err, "DEP profile")
	}
	e.ID = u.ID
	e.CreatedAt = u.CreatedAt
	e.UpdatedAt = u.UpdatedAt
	return nil
}

func (r *depProfileRepositoryImpl) FindByID(ctx context.Context, id uint) (*ent.DepProfile, error) {
	u, err := r.client.DepProfile.Query().Where(depprofile.IDEQ(id)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperror.ErrNotFound.WithMessage("Không tìm thấy DEP profile")
		}
		return nil, wrapFindError(err, "DEP profile")
	}
	return u, nil
}

func (r *depProfileRepositoryImpl) Update(ctx context.Context, e *ent.DepProfile) error {
	update := r.client.DepProfile.UpdateOneID(e.ID).
		SetProfileName(e.ProfileName).
		SetAllowPairing(e.AllowPairing).
		SetAnchorCerts(e.AnchorCerts).
		SetAutoAdvanceSetup(e.AutoAdvanceSetup).
		SetAwaitDeviceConfigured(e.AwaitDeviceConfigured).
		SetConfigurationWebURL(e.ConfigurationWebURL).
		SetDepartment(e.Department).
		SetDevices(e.Devices).
		SetDoNotUseProfileFromBackup(e.DoNotUseProfileFromBackup).
		SetIsReturnToService(e.IsReturnToService).
		SetIsMandatory(e.IsMandatory).
		SetIsMdmRemovable(e.IsMdmRemovable).
		SetIsMultiUser(e.IsMultiUser).
		SetIsSupervised(e.IsSupervised).
		SetLanguage(e.Language).
		SetOrgMagic(e.OrgMagic).
		SetRegion(e.Region).
		SetSkipSetupItems(e.SkipSetupItems).
		SetSupervisingHostCerts(e.SupervisingHostCerts).
		SetSupportEmailAddress(e.SupportEmailAddress).
		SetSupportPhoneNumber(e.SupportPhoneNumber).
		SetURL(e.URL).
		SetProfileData(e.ProfileData)

	if e.ProfileUUID != "" {
		update.SetProfileUUID(e.ProfileUUID)
	}

	u, err := update.Save(ctx)
	if err != nil {
		return wrapUpdateError(err, "DEP profile")
	}
	e.UpdatedAt = u.UpdatedAt
	return nil
}

func (r *depProfileRepositoryImpl) Delete(ctx context.Context, id uint) error {
	err := r.client.DepProfile.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return wrapDeleteError(err, "DEP profile")
	}
	return nil
}

func (r *depProfileRepositoryImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]ent.DepProfile, int64, error) {
	q := r.client.DepProfile.Query()

	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, wrapListError(err, "DEP profile")
	}

	// Simple sort
	q = q.Order(ent.Desc(depprofile.FieldCreatedAt))

	entProfiles, err := q.Offset(offset).Limit(limit).All(ctx)
	if err != nil {
		return nil, 0, wrapListError(err, "DEP profile")
	}

	res := make([]ent.DepProfile, len(entProfiles))
	for i, p := range entProfiles {
		res[i] = *p
	}
	return res, int64(total), nil
}

func (r *depProfileRepositoryImpl) Exists(ctx context.Context, id uint) (bool, error) {
	count, err := r.client.DepProfile.Query().Where(depprofile.IDEQ(id)).Count(ctx)
	if err != nil {
		return false, wrapFindError(err, "DEP profile")
	}
	return count > 0, nil
}

// --- DepProfileRepository Methods ---

func (r *depProfileRepositoryImpl) FindByProfileUUID(ctx context.Context, uuid string) (*ent.DepProfile, error) {
	u, err := r.client.DepProfile.Query().Where(depprofile.ProfileUUIDEQ(uuid)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperror.ErrNotFound.WithMessage("Không tìm thấy DEP profile với UUID: " + uuid)
		}
		return nil, wrapFindError(err, "DEP profile")
	}
	return u, nil
}

func (r *depProfileRepositoryImpl) FindByProfileName(ctx context.Context, name string) (*ent.DepProfile, error) {
	u, err := r.client.DepProfile.Query().Where(depprofile.ProfileNameEQ(name)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperror.ErrNotFound.WithMessage("Không tìm thấy DEP profile với tên: " + name)
		}
		return nil, wrapFindError(err, "DEP profile")
	}
	return u, nil
}
