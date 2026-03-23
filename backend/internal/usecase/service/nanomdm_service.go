package service

import (
	"context"
	"net/http"

	"github.com/thienel/go-backend-template/internal/interface/api/dto"
)

type NanoMDMService interface {
	// DEP related (proxied or direct)
	DefineDEPProfile(ctx context.Context, depName string, profile interface{}) (string, error)
	GetDEPProfile(ctx context.Context, depName, profileUUID string) (interface{}, error)
	SyncDEPDevices(ctx context.Context, depName string, cursor string) (interface{}, error)
	DisownDEPDevices(ctx context.Context, depName string, devices []string) (interface{}, error)
	UploadDEPToken(ctx context.Context, depName string, tokenData []byte) (interface{}, error)
	ListDEPProfiles(ctx context.Context, depName string) (interface{}, error)

	// New methods from apidog
	ListDEPNames(ctx context.Context) (interface{}, error)
	GetDEPConfig(ctx context.Context, depName string) (interface{}, error)
	GetDEPAssigner(ctx context.Context, depName string) (interface{}, error)
	SetDEPAssigner(ctx context.Context, depName string, assigner interface{}) (interface{}, error)
	GetDEPAccount(ctx context.Context, depName string) (interface{}, error)
	GetDEPDevices(ctx context.Context, depName string, devices []string, cursor string) (interface{}, error)
	GetDEPTokens(ctx context.Context, depName string) (interface{}, error)
	EnqueueCommand(ctx context.Context, udid string, cmdData []byte) (*dto.APIResult, error)
	Push(ctx context.Context, enrollments []string) (*dto.APIResult, error)
	EscrowKeyUnlock(ctx context.Context, req *dto.EscrowKeyUnlockRequest) ([]byte, http.Header, int, error)
	GetVersion(ctx context.Context) (*dto.NanoMDMVersionResponse, error)

	// Push Certificate related
	UploadPushCert(ctx context.Context, certData []byte) (*dto.PushCertResponse, error)
	GetPushCert(ctx context.Context, topic string) (*dto.PushCertResponse, error)
}
