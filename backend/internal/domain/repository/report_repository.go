package repository

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/pkg/query"
)

type ReportRepository interface {
	GetDevicesForExport(ctx context.Context, opts query.QueryOptions) ([]*ent.Device, error)
	GetAlertsForExport(ctx context.Context, opts query.QueryOptions) ([]*ent.Alert, error)
	GetApplicationsForExport(ctx context.Context, opts query.QueryOptions) ([]*ent.Application, error)
}
