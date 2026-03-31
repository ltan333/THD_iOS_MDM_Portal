package service

import (
	"context"
	"net/http"

	"github.com/thienel/go-backend-template/internal/interface/api/dto"
)

type NanoMDMService interface {
	// DEP related (proxied or direct)
	GetDEPProfile(ctx context.Context, depName, profileUUID string) (any, error)
	SyncDEPDevices(ctx context.Context, depName string, cursor string) (any, error)
	DisownDEPDevices(ctx context.Context, depName string, devices []string) (any, error)
	UploadDEPToken(ctx context.Context, depName string, tokenData []byte) (any, error)

	// New methods from apidog / NanoDEP spec
	ListDEPNames(ctx context.Context, depNames []string, limit, offset int, cursor string) (*dto.DEPNamesQueryResponse, error)
	GetDEPVersion(ctx context.Context) (*dto.NanoDEPVersionResponse, error)
	GetDEPConfig(ctx context.Context, depName string) (*dto.DEPConfig, error)
	SetDEPConfig(ctx context.Context, depName string, config *dto.DEPConfig) (*dto.DEPConfig, error)
	GetDEPAssigner(ctx context.Context, depName string) (*dto.AssignerProfileUUID, error)
	SetDEPAssigner(ctx context.Context, depName string, profileUUID string) (*dto.AssignerProfileUUID, error)
	GetDEPAccount(ctx context.Context, depName string) (any, error)
	GetDEPDevices(ctx context.Context, depName string, devices []string, cursor string) (any, error)
	GetDEPTokens(ctx context.Context, depName string) (*dto.OAuth1Tokens, error)
	UpdateDEPTokens(ctx context.Context, depName string, tokens *dto.OAuth1Tokens) (*dto.OAuth1Tokens, error)
	GetDEPTokenPKI(ctx context.Context, depName string, cn string, validityDays int) ([]byte, string, error)
	GetMAIDJWT(ctx context.Context, depName string, serverUUID string) (string, string, string, error)
	GetBypassCode(ctx context.Context, code, raw string) (*dto.BypassCodeResponse, error)
	EnqueueCommand(ctx context.Context, udid string, cmdData []byte) (*dto.APIResult, error)
	Push(ctx context.Context, enrollments []string) (*dto.APIResult, error)
	EscrowKeyUnlock(ctx context.Context, req *dto.EscrowKeyUnlockRequest) ([]byte, http.Header, int, error)
	GetVersion(ctx context.Context) (*dto.NanoMDMVersionResponse, error)

	// Push Certificate related
	UploadPushCert(ctx context.Context, certData []byte) (*dto.PushCertResponse, error)
	GetPushCert(ctx context.Context, topic string) (*dto.PushCertResponse, error)

	// DEP Syncer management
	ReloadDEPSyncer(ctx context.Context) error
}
