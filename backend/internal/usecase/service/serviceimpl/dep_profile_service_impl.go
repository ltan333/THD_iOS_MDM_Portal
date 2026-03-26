package serviceimpl

import (
	"context"

	"github.com/thienel/tlog"
	"go.uber.org/zap"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
)

type depProfileServiceImpl struct {
	depProfileRepo repository.DepProfileRepository
	mdmService     service.NanoMDMService
}

// NewDepProfileService creates a new DEP profile service
func NewDepProfileService(depProfileRepo repository.DepProfileRepository, mdmService service.NanoMDMService) service.DepProfileService {
	return &depProfileServiceImpl{
		depProfileRepo: depProfileRepo,
		mdmService:     mdmService,
	}
}

func (s *depProfileServiceImpl) DefineProfile(ctx context.Context, depName string, req *dto.DEPProfileRequest) (*dto.DEPProfileResponse, error) {
	// Call upstream service
	profileUUID, err := s.mdmService.DefineDEPProfile(ctx, depName, req)
	if err != nil {
		tlog.Error("Failed to define DEP profile upstream", zap.Error(err))
		// Continue even if upstream fails, if possible, but keep uuid empty
		// In production, we might want to fail here if we really need upstream UUID
	}

	// Check if profile exists locally
	existing, err := s.depProfileRepo.FindByProfileName(ctx, req.ProfileName)
	var profile *ent.DepProfile
	if err == nil {
		// Update existing
		profile = existing
		mapRequestToEntity(req, profile)
		profile.ProfileUUID = profileUUID
		if err := s.depProfileRepo.Update(ctx, profile); err != nil {
			return nil, err
		}
	} else {
		// Create new
		profile = &ent.DepProfile{
			ProfileName: req.ProfileName,
			ProfileUUID: profileUUID,
		}
		mapRequestToEntity(req, profile)
		if err := s.depProfileRepo.Create(ctx, profile); err != nil {
			return nil, err
		}
	}

	return mapEntityToResponse(profile), nil
}

func (s *depProfileServiceImpl) GetProfile(ctx context.Context, depName, uuid string) (*dto.DEPProfileResponse, error) {
	// Try local first
	profile, err := s.depProfileRepo.FindByProfileUUID(ctx, uuid)
	if err == nil {
		return mapEntityToResponse(profile), nil
	}

	// Fallback to upstream
	tlog.Info("DEP profile not found locally, falling back to upstream", zap.String("uuid", uuid))
	upstreamResp, err := s.mdmService.GetDEPProfile(ctx, depName, uuid)
	if err != nil {
		return nil, err
	}

	// Type assertion if possible, or return as is if implementation allows
	if resp, ok := upstreamResp.(*dto.DEPProfileResponse); ok {
		return resp, nil
	}

	// If it's another type (e.g. map), we might need more complex mapping
	// but for now let's assume it returns the DO or we handle it
	return nil, apperror.ErrInternalServerError.WithMessage("Unexpected upstream response type")
}

func (s *depProfileServiceImpl) ListProfiles(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*dto.DEPProfileResponse, int64, error) {
	profiles, total, err := s.depProfileRepo.List(ctx, offset, limit, opts)
	if err != nil {
		return nil, 0, err
	}

	resp := make([]*dto.DEPProfileResponse, len(profiles))
	for i, p := range profiles {
		resp[i] = mapEntityToResponse(&p)
	}

	return resp, total, nil
}

// Helper mappings

func mapRequestToEntity(req *dto.DEPProfileRequest, e *ent.DepProfile) {
	if req.AllowPairing != nil {
		e.AllowPairing = *req.AllowPairing
	}
	if req.AnchorCerts != nil {
		e.AnchorCerts = req.AnchorCerts
	}
	if req.AutoAdvanceSetup != nil {
		e.AutoAdvanceSetup = *req.AutoAdvanceSetup
	}
	if req.AwaitDeviceConfigured != nil {
		e.AwaitDeviceConfigured = *req.AwaitDeviceConfigured
	}
	if req.ConfigurationWebURL != "" {
		e.ConfigurationWebURL = req.ConfigurationWebURL
	}
	if req.Department != "" {
		e.Department = req.Department
	}
	if req.Devices != nil {
		e.Devices = req.Devices
	}
	if req.DoNotUseProfileFromBackup != nil {
		e.DoNotUseProfileFromBackup = *req.DoNotUseProfileFromBackup
	}
	if req.IsReturnToService != nil {
		e.IsReturnToService = *req.IsReturnToService
	}
	if req.IsMandatory != nil {
		e.IsMandatory = *req.IsMandatory
	}
	if req.IsMDMRemovable != nil {
		e.IsMdmRemovable = *req.IsMDMRemovable
	}
	if req.IsMultiUser != nil {
		e.IsMultiUser = *req.IsMultiUser
	}
	if req.IsSupervised != nil {
		e.IsSupervised = *req.IsSupervised
	}
	if req.Language != "" {
		e.Language = req.Language
	}
	if req.OrgMagic != "" {
		e.OrgMagic = req.OrgMagic
	}
	if req.Region != "" {
		e.Region = req.Region
	}
	if req.SkipSetupItems != nil {
		e.SkipSetupItems = req.SkipSetupItems
	}
	if req.SupervisingHostCerts != nil {
		e.SupervisingHostCerts = req.SupervisingHostCerts
	}
	if req.SupportEmailAddress != "" {
		e.SupportEmailAddress = req.SupportEmailAddress
	}
	if req.SupportPhoneNumber != "" {
		e.SupportPhoneNumber = req.SupportPhoneNumber
	}
	if req.URL != "" {
		e.URL = req.URL
	}
}

func mapEntityToResponse(p *ent.DepProfile) *dto.DEPProfileResponse {
	return &dto.DEPProfileResponse{
		Name:                      p.ProfileName,
		ProfileUUID:               p.ProfileUUID,
		AllowPairing:              p.AllowPairing,
		AnchorCerts:               p.AnchorCerts,
		AutoAdvanceSetup:          p.AutoAdvanceSetup,
		AwaitDeviceConfigured:     p.AwaitDeviceConfigured,
		ConfigurationWebURL:       p.ConfigurationWebURL,
		Department:                p.Department,
		Devices:                   p.Devices,
		DoNotUseProfileFromBackup: p.DoNotUseProfileFromBackup,
		IsReturnToService:         p.IsReturnToService,
		IsMandatory:               p.IsMandatory,
		IsMDMRemovable:            p.IsMdmRemovable,
		IsMultiUser:               p.IsMultiUser,
		IsSupervised:              p.IsSupervised,
		Language:                  p.Language,
		OrgMagic:                  p.OrgMagic,
		Region:                    p.Region,
		SkipSetupItems:            p.SkipSetupItems,
		SupervisingHostCerts:      p.SupervisingHostCerts,
		SupportEmailAddress:       p.SupportEmailAddress,
		SupportPhoneNumber:        p.SupportPhoneNumber,
		URL:                       p.URL,
	}
}
