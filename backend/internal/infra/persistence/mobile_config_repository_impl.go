package persistence

import (
	"context"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/mobileconfig"
	"github.com/thienel/go-backend-template/pkg/query"
)

type mobileConfigRepositoryImpl struct {
	client *ent.Client
}

func NewMobileConfigRepository(client *ent.Client) repository.MobileConfigRepository {
	return &mobileConfigRepositoryImpl{client: client}
}

func (m *mobileConfigRepositoryImpl) Create(ctx context.Context, entity *ent.MobileConfig) error {
	panic("unimplemented")
}

func (m *mobileConfigRepositoryImpl) Delete(ctx context.Context, id uint) error {
	panic("unimplemented")
}

func (m *mobileConfigRepositoryImpl) Exists(ctx context.Context, id uint) (bool, error) {
	panic("unimplemented")
}

func (m *mobileConfigRepositoryImpl) FindByID(ctx context.Context, id uint) (*ent.MobileConfig, error) {
	panic("unimplemented")
}

func (m *mobileConfigRepositoryImpl) List(ctx context.Context, offset int, limit int, opts query.QueryOptions) ([]ent.MobileConfig, int64, error) {
	panic("unimplemented")
}

func (m *mobileConfigRepositoryImpl) Update(ctx context.Context, entity *ent.MobileConfig) error {
	panic("unimplemented")
}

func (m *mobileConfigRepositoryImpl) GetFullForExport(ctx context.Context, id uint) (*ent.MobileConfig, error) {
	return m.client.MobileConfig.Query().
		Where(mobileconfig.IDEQ(id)).
		WithPayloads(func(q *ent.PayloadQuery) {
			q.WithProperties(func(pq *ent.PayloadPropertyQuery) {
				pq.WithDefinition()
			})
		}).
		First(ctx)
}
