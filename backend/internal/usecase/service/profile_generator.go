package service

import (
	"context"
	"github.com/thienel/go-backend-template/internal/ent"
)

type ProfileGenerator interface {
	GenerateXML(ctx context.Context, profile *ent.Profile) ([]byte, error)
}
