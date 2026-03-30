package persistence

import (
	"context"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/alert"
	"github.com/thienel/go-backend-template/internal/ent/application"
	"github.com/thienel/go-backend-template/internal/ent/device"
	"github.com/thienel/go-backend-template/pkg/query"
)

type reportRepositoryImpl struct {
	client *ent.Client
}

func NewReportRepository(client *ent.Client) repository.ReportRepository {
	return &reportRepositoryImpl{client: client}
}

func (r *reportRepositoryImpl) GetDevicesForExport(ctx context.Context, opts query.QueryOptions) ([]*ent.Device, error) {
	q := r.client.Device.Query()

	if searchFilter, ok := opts.Filters["search"]; ok {
		searchStr, _ := searchFilter.Value.(string)
		if searchStr != "" {
			q = q.Where(
				device.Or(
					device.NameContainsFold(searchStr),
					device.SerialNumberContainsFold(searchStr),
					device.ModelContainsFold(searchStr),
				),
			)
		}
	}

	return q.Order(ent.Desc(device.FieldCreatedAt)).All(ctx)
}

func (r *reportRepositoryImpl) GetAlertsForExport(ctx context.Context, opts query.QueryOptions) ([]*ent.Alert, error) {
	q := r.client.Alert.Query()

	if searchFilter, ok := opts.Filters["search"]; ok {
		searchStr, _ := searchFilter.Value.(string)
		if searchStr != "" {
			q = q.Where(alert.TitleContainsFold(searchStr))
		}
	}

	return q.Order(ent.Desc(alert.FieldCreatedAt)).All(ctx)
}

func (r *reportRepositoryImpl) GetApplicationsForExport(ctx context.Context, opts query.QueryOptions) ([]*ent.Application, error) {
	q := r.client.Application.Query()

	if searchFilter, ok := opts.Filters["search"]; ok {
		searchStr, _ := searchFilter.Value.(string)
		if searchStr != "" {
			q = q.Where(
				application.Or(
					application.NameContainsFold(searchStr),
					application.BundleIDContainsFold(searchStr),
				),
			)
		}
	}

	return q.Order(ent.Desc(application.FieldCreatedAt)).All(ctx)
}
