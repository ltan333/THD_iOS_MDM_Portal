package service

import (
	"context"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
)

type NanoCMDService interface {
	GetVersion(ctx context.Context) (*dto.NanoCMDVersionResponse, error)
	StartWorkflow(ctx context.Context, name string, enrollmentIDs []string, contextStr string) (*dto.NanoCMDWorkflowStartResponse, error)
	GetEvent(ctx context.Context, name string) (*dto.EventSubscription, error)
	PutEvent(ctx context.Context, name string, subscription *dto.EventSubscription) error
	GetFVEnableProfileTemplate(ctx context.Context) ([]byte, error)
	GetProfile(ctx context.Context, name string) ([]byte, error)
	PutProfile(ctx context.Context, name string, profileData []byte) error
	DeleteProfile(ctx context.Context, name string) error
	GetProfiles(ctx context.Context, names []string) (map[string]dto.NanoCMDProfile, error)
	GetCMDPlan(ctx context.Context, name string) (*dto.CMDPlan, error)
	PutCMDPlan(ctx context.Context, name string, plan *dto.CMDPlan) error
	GetInventory(ctx context.Context, enrollmentIDs []string) (dto.NanoCMDInventoryResponse, error)
}
