package service

import (
	"context"

	"github.com/thienel/go-backend-template/pkg/query"
)

type ReportService interface {
	ExportDevicesCSV(ctx context.Context, opts query.QueryOptions) ([]byte, error)
	ExportAlertsCSV(ctx context.Context, opts query.QueryOptions) ([]byte, error)
	ExportApplicationsCSV(ctx context.Context, opts query.QueryOptions) ([]byte, error)
}
