package repository

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/pkg/query"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
)

type AlertRepository interface {
	List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.Alert, int64, error)
	GetByID(ctx context.Context, id uint) (*ent.Alert, error)
	Create(ctx context.Context, entity *ent.Alert) (*ent.Alert, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
	BulkUpdateStatus(ctx context.Context, ids []uint, status string) error
	GetStats(ctx context.Context) (*dto.AlertsSummaryResponse, error)
}

type AlertRuleRepository interface {
	List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.AlertRule, int64, error)
	GetByID(ctx context.Context, id uint) (*ent.AlertRule, error)
	Create(ctx context.Context, entity *ent.AlertRule) (*ent.AlertRule, error)
	Update(ctx context.Context, id uint, entity *ent.AlertRule) (*ent.AlertRule, error)
	Delete(ctx context.Context, id uint) error
	SetEnabled(ctx context.Context, id uint, enabled bool) error
}
