package persistence

import (
	"context"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/depprofile"
	apperror "github.com/thienel/go-backend-template/pkg/error"
)

type depProfileRepositoryImpl struct {
	client *ent.Client
}

func NewDepProfileRepository(client *ent.Client) repository.DepProfileRepository {
	return &depProfileRepositoryImpl{client: client}
}

func (r *depProfileRepositoryImpl) Create(ctx context.Context, profile *ent.DepProfile) (*ent.DepProfile, error) {
	create := r.client.DepProfile.Create().
		SetProfileName(profile.ProfileName).
		SetAllowPairing(profile.AllowPairing).
		SetAutoAdvanceSetup(profile.AutoAdvanceSetup).
		SetAwaitDeviceConfigured(profile.AwaitDeviceConfigured).
		SetDoNotUseProfileFromBackup(profile.DoNotUseProfileFromBackup).
		SetIsReturnToService(profile.IsReturnToService).
		SetIsMandatory(profile.IsMandatory).
		SetIsMdmRemovable(profile.IsMdmRemovable).
		SetIsMultiUser(profile.IsMultiUser).
		SetIsSupervised(profile.IsSupervised)

	if profile.ProfileUUID != "" {
		create = create.SetProfileUUID(profile.ProfileUUID)
	}
	if len(profile.AnchorCerts) > 0 {
		create = create.SetAnchorCerts(profile.AnchorCerts)
	}
	if profile.ConfigurationWebURL != "" {
		create = create.SetConfigurationWebURL(profile.ConfigurationWebURL)
	}
	if profile.Department != "" {
		create = create.SetDepartment(profile.Department)
	}
	if len(profile.Devices) > 0 {
		create = create.SetDevices(profile.Devices)
	}
	if profile.Language != "" {
		create = create.SetLanguage(profile.Language)
	}
	if profile.OrgMagic != "" {
		create = create.SetOrgMagic(profile.OrgMagic)
	}
	if profile.Region != "" {
		create = create.SetRegion(profile.Region)
	}
	if len(profile.SkipSetupItems) > 0 {
		create = create.SetSkipSetupItems(profile.SkipSetupItems)
	}
	if len(profile.SupervisingHostCerts) > 0 {
		create = create.SetSupervisingHostCerts(profile.SupervisingHostCerts)
	}
	if profile.SupportEmailAddress != "" {
		create = create.SetSupportEmailAddress(profile.SupportEmailAddress)
	}
	if profile.SupportPhoneNumber != "" {
		create = create.SetSupportPhoneNumber(profile.SupportPhoneNumber)
	}
	if profile.URL != "" {
		create = create.SetURL(profile.URL)
	}
	if len(profile.ProfileData) > 0 {
		create = create.SetProfileData(profile.ProfileData)
	}

	result, err := create.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, apperror.ErrConflict.WithMessage("Profile with this name already exists")
		}
		return nil, apperror.ErrInternalServerError.WithMessage("Failed to create DEP profile").WithError(err)
	}
	return result, nil
}

func (r *depProfileRepositoryImpl) GetByID(ctx context.Context, id uint) (*ent.DepProfile, error) {
	profile, err := r.client.DepProfile.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperror.ErrNotFound.WithMessage("DEP profile not found")
		}
		return nil, apperror.ErrInternalServerError.WithMessage("Failed to get DEP profile").WithError(err)
	}
	return profile, nil
}

func (r *depProfileRepositoryImpl) GetByProfileUUID(ctx context.Context, profileUUID string) (*ent.DepProfile, error) {
	profile, err := r.client.DepProfile.
		Query().
		Where(depprofile.ProfileUUIDEQ(profileUUID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperror.ErrNotFound.WithMessage("DEP profile not found")
		}
		return nil, apperror.ErrInternalServerError.WithMessage("Failed to get DEP profile").WithError(err)
	}
	return profile, nil
}

func (r *depProfileRepositoryImpl) GetByName(ctx context.Context, profileName string) (*ent.DepProfile, error) {
	profile, err := r.client.DepProfile.
		Query().
		Where(depprofile.ProfileNameEQ(profileName)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperror.ErrNotFound.WithMessage("DEP profile not found")
		}
		return nil, apperror.ErrInternalServerError.WithMessage("Failed to get DEP profile").WithError(err)
	}
	return profile, nil
}

func (r *depProfileRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*ent.DepProfile, int64, error) {
	q := r.client.DepProfile.Query()

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, apperror.ErrInternalServerError.WithMessage("Failed to count DEP profiles").WithError(err)
	}

	profiles, err := q.
		Offset(offset).
		Limit(limit).
		Order(ent.Desc(depprofile.FieldUpdatedAt)).
		All(ctx)
	if err != nil {
		return nil, 0, apperror.ErrInternalServerError.WithMessage("Failed to list DEP profiles").WithError(err)
	}

	return profiles, int64(total), nil
}

func (r *depProfileRepositoryImpl) Update(ctx context.Context, id uint, profile *ent.DepProfile) (*ent.DepProfile, error) {
	update := r.client.DepProfile.UpdateOneID(id).
		SetProfileName(profile.ProfileName).
		SetAllowPairing(profile.AllowPairing).
		SetAutoAdvanceSetup(profile.AutoAdvanceSetup).
		SetAwaitDeviceConfigured(profile.AwaitDeviceConfigured).
		SetDoNotUseProfileFromBackup(profile.DoNotUseProfileFromBackup).
		SetIsReturnToService(profile.IsReturnToService).
		SetIsMandatory(profile.IsMandatory).
		SetIsMdmRemovable(profile.IsMdmRemovable).
		SetIsMultiUser(profile.IsMultiUser).
		SetIsSupervised(profile.IsSupervised)

	if profile.ProfileUUID != "" {
		update = update.SetProfileUUID(profile.ProfileUUID)
	}
	if len(profile.AnchorCerts) > 0 {
		update = update.SetAnchorCerts(profile.AnchorCerts)
	}
	if profile.ConfigurationWebURL != "" {
		update = update.SetConfigurationWebURL(profile.ConfigurationWebURL)
	}
	if profile.Department != "" {
		update = update.SetDepartment(profile.Department)
	}
	if len(profile.Devices) > 0 {
		update = update.SetDevices(profile.Devices)
	}
	if profile.Language != "" {
		update = update.SetLanguage(profile.Language)
	}
	if profile.OrgMagic != "" {
		update = update.SetOrgMagic(profile.OrgMagic)
	}
	if profile.Region != "" {
		update = update.SetRegion(profile.Region)
	}
	if len(profile.SkipSetupItems) > 0 {
		update = update.SetSkipSetupItems(profile.SkipSetupItems)
	}
	if len(profile.SupervisingHostCerts) > 0 {
		update = update.SetSupervisingHostCerts(profile.SupervisingHostCerts)
	}
	if profile.SupportEmailAddress != "" {
		update = update.SetSupportEmailAddress(profile.SupportEmailAddress)
	}
	if profile.SupportPhoneNumber != "" {
		update = update.SetSupportPhoneNumber(profile.SupportPhoneNumber)
	}
	if profile.URL != "" {
		update = update.SetURL(profile.URL)
	}
	if len(profile.ProfileData) > 0 {
		update = update.SetProfileData(profile.ProfileData)
	}

	result, err := update.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperror.ErrNotFound.WithMessage("DEP profile not found")
		}
		if ent.IsConstraintError(err) {
			return nil, apperror.ErrConflict.WithMessage("Profile with this name already exists")
		}
		return nil, apperror.ErrInternalServerError.WithMessage("Failed to update DEP profile").WithError(err)
	}
	return result, nil
}

func (r *depProfileRepositoryImpl) Delete(ctx context.Context, id uint) error {
	err := r.client.DepProfile.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return apperror.ErrNotFound.WithMessage("DEP profile not found")
		}
		return apperror.ErrInternalServerError.WithMessage("Failed to delete DEP profile").WithError(err)
	}
	return nil
}

func (r *depProfileRepositoryImpl) SetProfileUUID(ctx context.Context, id uint, profileUUID string) error {
	_, err := r.client.DepProfile.UpdateOneID(id).
		SetProfileUUID(profileUUID).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return apperror.ErrNotFound.WithMessage("DEP profile not found")
		}
		return apperror.ErrInternalServerError.WithMessage("Failed to update DEP profile UUID").WithError(err)
	}
	return nil
}
