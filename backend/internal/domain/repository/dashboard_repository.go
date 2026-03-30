package repository

import (
	"context"
)

type DashboardRepository interface {
	CountDevices(ctx context.Context) (int, error)
	CountEnrolledDevices(ctx context.Context) (int, error)
	CountUsers(ctx context.Context) (int, error)
	CountActiveUsers(ctx context.Context) (int, error)
	GetDevicePlatformCounts(ctx context.Context) (map[string]int, error)
}
