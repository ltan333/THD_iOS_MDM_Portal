package persistence

import (
	"context"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/device"
	"github.com/thienel/go-backend-template/internal/ent/user"
)

type dashboardRepositoryImpl struct {
	client *ent.Client
}

func NewDashboardRepository(client *ent.Client) repository.DashboardRepository {
	return &dashboardRepositoryImpl{client: client}
}

func (r *dashboardRepositoryImpl) CountDevices(ctx context.Context) (int, error) {
	return r.client.Device.Query().Count(ctx)
}

func (r *dashboardRepositoryImpl) CountEnrolledDevices(ctx context.Context) (int, error) {
	return r.client.Device.Query().Where(device.IsEnrolledEQ(true)).Count(ctx)
}

func (r *dashboardRepositoryImpl) CountUsers(ctx context.Context) (int, error) {
	return r.client.User.Query().Where(user.DeletedAtIsNil()).Count(ctx)
}

func (r *dashboardRepositoryImpl) CountActiveUsers(ctx context.Context) (int, error) {
	return r.client.User.Query().Where(user.DeletedAtIsNil(), user.StatusEQ("ACTIVE")).Count(ctx)
}

func (r *dashboardRepositoryImpl) GetDevicePlatformCounts(ctx context.Context) (map[string]int, error) {
	var platformStats []struct {
		Platform string `json:"platform"`
		Count    int    `json:"count"`
	}

	result := make(map[string]int)

	if err := r.client.Device.Query().
		GroupBy(device.FieldPlatform).
		Aggregate(ent.Count()).
		Scan(ctx, &platformStats); err != nil {
		return result, err
	}

	for _, stat := range platformStats {
		if stat.Platform != "" {
			result[stat.Platform] = stat.Count
		}
	}
	return result, nil
}
