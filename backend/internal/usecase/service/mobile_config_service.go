package service

import "context"

type GenerateMobileConfigXMLCommand struct {
	ID uint
}

type MobileConfigService interface {
	GenerateXML(ctx context.Context, cmd GenerateMobileConfigXMLCommand) ([]byte, error)
}
