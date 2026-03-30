package service

import (
	"context"

	"github.com/thienel/go-backend-template/internal/interface/api/dto"
)

type DashboardService interface {
	GetStats(ctx context.Context) (*dto.DashboardStatsResponse, error)
	GetDeviceStats(ctx context.Context) (*dto.DeviceStatsResponse, error)
	GetAlertsSummary(ctx context.Context) (*dto.AlertsSummaryResponse, error)
	GetChartData(ctx context.Context, chartType string) (*dto.ChartDataResponse, error)
}
