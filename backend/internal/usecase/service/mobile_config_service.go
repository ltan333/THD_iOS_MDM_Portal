package service

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/pkg/query"
)

type CreateMobileConfigPropertyCommand struct {
	Key       string
	ValueJSON map[string]interface{}
}

type CreateMobileConfigPayloadCommand struct {
	PayloadDescription  string
	PayloadDisplayName  string
	PayloadIdentifier   string
	PayloadOrganization string
	PayloadType         string
	PayloadVersion      int
	Properties          []CreateMobileConfigPropertyCommand
}

type CreateMobileConfigCommand struct {
	Name                     string
	PayloadIdentifier        string
	PayloadType              string
	PayloadDisplayName       string
	PayloadDescription       string
	PayloadOrganization      string
	PayloadVersion           int
	PayloadRemovalDisallowed bool
	Payloads                 []CreateMobileConfigPayloadCommand
}

type UpdateMobileConfigCommand struct {
	ID                       uint
	Name                     string
	PayloadIdentifier        string
	PayloadType              string
	PayloadDisplayName       string
	PayloadDescription       string
	PayloadOrganization      string
	PayloadVersion           int
	PayloadRemovalDisallowed bool
	Payloads                 []CreateMobileConfigPayloadCommand
}

type GenerateMobileConfigXMLCommand struct {
	ID uint
}

type MobileConfigService interface {
	List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.MobileConfig, int64, error)
	GetByID(ctx context.Context, id uint) (*ent.MobileConfig, error)
	Create(ctx context.Context, cmd CreateMobileConfigCommand) (*ent.MobileConfig, error)
	Update(ctx context.Context, cmd UpdateMobileConfigCommand) (*ent.MobileConfig, error)
	Delete(ctx context.Context, id uint) error
	GenerateXML(ctx context.Context, cmd GenerateMobileConfigXMLCommand) ([]byte, error)
}
