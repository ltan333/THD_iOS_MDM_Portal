package serviceimpl

import (
	"context"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/tlog"
	"go.uber.org/zap"
)

type depProfileServiceImpl struct {
	repo       repository.DepProfileRepository
	nanomdmSvc service.NanoMDMService
}

func NewDepProfileService(repo repository.DepProfileRepository, nanomdmSvc service.NanoMDMService) service.DepProfileService {
	return &depProfileServiceImpl{
		repo:       repo,
		nanomdmSvc: nanomdmSvc,
	}
}

// Create creates a new DEP profile locally and registers it with Apple DEP.
func (s *depProfileServiceImpl) Create(ctx context.Context, depName string, req *dto.DEPProfileRequest) (*ent.DepProfile, error) {
	// 1. Create profile in Apple DEP first to get the profile_uuid
	appleResp, err := s.nanomdmSvc.CreateDEPProfile(ctx, depName, req)
	if err != nil {
		tlog.Error("Failed to create DEP profile in Apple",
			zap.String("profile_name", req.ProfileName),
			zap.Error(err))
		return nil, apperror.ErrInternalServerError.WithMessage("Failed to create profile in Apple DEP").WithError(err)
	}

	tlog.Info("Created DEP profile in Apple",
		zap.String("profile_name", req.ProfileName),
		zap.String("profile_uuid", appleResp.ProfileUUID))

	// 2. Create profile in local database with the Apple profile_uuid
	profile := mapDEPProfileRequestToEnt(req)
	profile.ProfileUUID = appleResp.ProfileUUID

	result, err := s.repo.Create(ctx, profile)
	if err != nil {
		tlog.Error("Failed to save DEP profile to database",
			zap.String("profile_name", req.ProfileName),
			zap.Error(err))
		return nil, err
	}

	tlog.Info("Created DEP profile in database",
		zap.Uint("id", result.ID),
		zap.String("profile_uuid", result.ProfileUUID))

	return result, nil
}

func (s *depProfileServiceImpl) GetByID(ctx context.Context, id uint) (*ent.DepProfile, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *depProfileServiceImpl) GetByProfileUUID(ctx context.Context, profileUUID string) (*ent.DepProfile, error) {
	return s.repo.GetByProfileUUID(ctx, profileUUID)
}

func (s *depProfileServiceImpl) List(ctx context.Context, offset, limit int) ([]*ent.DepProfile, int64, error) {
	return s.repo.List(ctx, offset, limit)
}

// Update updates an existing DEP profile and syncs with Apple DEP.
func (s *depProfileServiceImpl) Update(ctx context.Context, depName string, id uint, req *dto.DEPProfileRequest) (*ent.DepProfile, error) {
	// 1. Get existing profile to get the Apple profile_uuid
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 2. Update profile in Apple DEP (this creates a new version with the same name)
	appleResp, err := s.nanomdmSvc.CreateDEPProfile(ctx, depName, req)
	if err != nil {
		tlog.Error("Failed to update DEP profile in Apple",
			zap.Uint("id", id),
			zap.String("profile_name", req.ProfileName),
			zap.Error(err))
		return nil, apperror.ErrInternalServerError.WithMessage("Failed to update profile in Apple DEP").WithError(err)
	}

	tlog.Info("Updated DEP profile in Apple",
		zap.Uint("id", id),
		zap.String("old_uuid", existing.ProfileUUID),
		zap.String("new_uuid", appleResp.ProfileUUID))

	// 3. Update profile in local database
	profile := mapDEPProfileRequestToEnt(req)
	profile.ProfileUUID = appleResp.ProfileUUID

	result, err := s.repo.Update(ctx, id, profile)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Delete removes a DEP profile from local DB and optionally from Apple DEP.
func (s *depProfileServiceImpl) Delete(ctx context.Context, depName string, id uint) error {
	// 1. Get existing profile
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. Try to remove from Apple DEP (optional - may fail if devices are assigned)
	if existing.ProfileUUID != "" {
		if err := s.nanomdmSvc.RemoveDEPProfile(ctx, depName, existing.ProfileUUID); err != nil {
			tlog.Warn("Failed to remove DEP profile from Apple (continuing with local delete)",
				zap.Uint("id", id),
				zap.String("profile_uuid", existing.ProfileUUID),
				zap.Error(err))
			// Continue with local delete even if Apple API fails
		}
	}

	// 3. Delete from local database
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	tlog.Info("Deleted DEP profile",
		zap.Uint("id", id),
		zap.String("profile_uuid", existing.ProfileUUID))

	return nil
}

// SetAsAssigner sets a profile as the default assigner profile for new devices.
func (s *depProfileServiceImpl) SetAsAssigner(ctx context.Context, depName string, id uint) error {
	// 1. Get profile to get the Apple profile_uuid
	profile, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if profile.ProfileUUID == "" {
		return apperror.ErrBadRequest.WithMessage("Profile has not been registered with Apple DEP yet")
	}

	// 2. Set as assigner in NanoDEP
	_, err = s.nanomdmSvc.SetDEPAssigner(ctx, depName, profile.ProfileUUID)
	if err != nil {
		tlog.Error("Failed to set DEP assigner profile",
			zap.Uint("id", id),
			zap.String("profile_uuid", profile.ProfileUUID),
			zap.Error(err))
		return apperror.ErrInternalServerError.WithMessage("Failed to set assigner profile").WithError(err)
	}

	tlog.Info("Set DEP assigner profile",
		zap.Uint("id", id),
		zap.String("profile_uuid", profile.ProfileUUID),
		zap.String("dep_name", depName))

	return nil
}

// GetAssigner gets the current assigner profile.
func (s *depProfileServiceImpl) GetAssigner(ctx context.Context, depName string) (*ent.DepProfile, error) {
	// 1. Get assigner profile UUID from NanoDEP
	assigner, err := s.nanomdmSvc.GetDEPAssigner(ctx, depName)
	if err != nil {
		return nil, err
	}

	if assigner.ProfileUUID == "" {
		return nil, apperror.ErrNotFound.WithMessage("No assigner profile configured")
	}

	// 2. Get profile from local database
	profile, err := s.repo.GetByProfileUUID(ctx, assigner.ProfileUUID)
	if err != nil {
		// Profile exists in Apple but not in local DB
		tlog.Warn("Assigner profile not found in local database",
			zap.String("profile_uuid", assigner.ProfileUUID))
		return nil, err
	}

	return profile, nil
}

// mapDEPProfileRequestToEnt converts a DEPProfileRequest DTO to an ent.DepProfile entity.
func mapDEPProfileRequestToEnt(req *dto.DEPProfileRequest) *ent.DepProfile {
	profile := &ent.DepProfile{
		ProfileName:          req.ProfileName,
		ConfigurationWebURL:  req.ConfigurationWebURL,
		Department:           req.Department,
		Language:             req.Language,
		OrgMagic:             req.OrgMagic,
		Region:               req.Region,
		SupportEmailAddress:  req.SupportEmailAddress,
		SupportPhoneNumber:   req.SupportPhoneNumber,
		URL:                  req.URL,
		AnchorCerts:          req.AnchorCerts,
		Devices:              req.Devices,
		SkipSetupItems:       req.SkipSetupItems,
		SupervisingHostCerts: req.SupervisingHostCerts,
		ProfileData:          req.ProfileData,
	}

	// Handle boolean pointers
	if req.AllowPairing != nil {
		profile.AllowPairing = *req.AllowPairing
	} else {
		profile.AllowPairing = true // default
	}
	if req.AutoAdvanceSetup != nil {
		profile.AutoAdvanceSetup = *req.AutoAdvanceSetup
	}
	if req.AwaitDeviceConfigured != nil {
		profile.AwaitDeviceConfigured = *req.AwaitDeviceConfigured
	}
	if req.DoNotUseProfileFromBackup != nil {
		profile.DoNotUseProfileFromBackup = *req.DoNotUseProfileFromBackup
	}
	if req.IsReturnToService != nil {
		profile.IsReturnToService = *req.IsReturnToService
	}
	if req.IsMandatory != nil {
		profile.IsMandatory = *req.IsMandatory
	}
	if req.IsMDMRemovable != nil {
		profile.IsMdmRemovable = *req.IsMDMRemovable
	} else {
		profile.IsMdmRemovable = true // default
	}
	if req.IsMultiUser != nil {
		profile.IsMultiUser = *req.IsMultiUser
	}
	if req.IsSupervised != nil {
		profile.IsSupervised = *req.IsSupervised
	}

	return profile
}
