package service

import (
	"context"

	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/pkg/query"
)

// DepProfileService defines the DEP profile service interface
type DepProfileService interface {
	// DefineProfile handles local persistence and upstream sync
	DefineProfile(ctx context.Context, depName string, req *dto.DEPProfileRequest) (*dto.DEPProfileResponse, error)

	// GetProfile retrieves a profile by UUID, with fallback to upstream
	GetProfile(ctx context.Context, depName, uuid string) (*dto.DEPProfileResponse, error)

	// ListProfiles returns all locally stored profiles
	ListProfiles(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*dto.DEPProfileResponse, int64, error)
}
