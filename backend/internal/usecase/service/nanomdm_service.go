package service

import (
	"context"
)

type NanoMDMService interface {
	// DEP related (proxied or direct)
	DefineDEPProfile(ctx context.Context, depName string, profile interface{}) (string, error)
	GetDEPProfile(ctx context.Context, depName, profileUUID string) (interface{}, error)
	
	// Push Certificate related
	UploadPushCert(ctx context.Context, certData []byte) error
	GetPushCert(ctx context.Context) (interface{}, error)
}
