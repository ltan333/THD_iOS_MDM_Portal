package service

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/pkg/query"
)

type CreateAlertCommand struct {
	Severity string
	Title    string
	Type     string
	DeviceID string
	UserID   *uint
	Details  map[string]interface{}
}

type CreateAlertRuleCommand struct {
	Name        string
	Description string
	Condition   map[string]interface{}
	Actions     map[string]interface{}
	Enabled     bool
}

type UpdateAlertRuleCommand struct {
	ID          uint
	Name        *string
	Description *string
	Condition   map[string]interface{}
	Actions     map[string]interface{}
	Enabled     *bool
}

type AlertService interface {
	// CRUD & Filtering
	List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.Alert, int64, error)
	GetByID(ctx context.Context, id uint) (*ent.Alert, error)
	Create(ctx context.Context, cmd CreateAlertCommand) (*ent.Alert, error)

	// Status Management
	Acknowledge(ctx context.Context, id uint) error
	Resolve(ctx context.Context, id uint) error
	BulkResolve(ctx context.Context, ids []uint) error

	// Stats
	GetStats(ctx context.Context) (*dto.AlertsSummaryResponse, error)

	// Quick Actions
	LockDevice(ctx context.Context, alertID uint) error
	WipeDevice(ctx context.Context, alertID uint) error
	PushPolicy(ctx context.Context, alertID uint, policyID uint) error
	SendMessage(ctx context.Context, alertID uint, message string) error
}

type AlertRuleService interface {
	List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.AlertRule, int64, error)
	GetByID(ctx context.Context, id uint) (*ent.AlertRule, error)
	Create(ctx context.Context, cmd CreateAlertRuleCommand) (*ent.AlertRule, error)
	Update(ctx context.Context, cmd UpdateAlertRuleCommand) (*ent.AlertRule, error)
	Delete(ctx context.Context, id uint) error
	Toggle(ctx context.Context, id uint) error
}
