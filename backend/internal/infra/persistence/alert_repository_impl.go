package persistence

import (
	"context"
	"strings"
	"time"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/alert"
	"github.com/thienel/go-backend-template/internal/ent/alertrule"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
)

// Alert Repository Implementation

type alertRepositoryImpl struct {
	client *ent.Client
}

func NewAlertRepository(client *ent.Client) repository.AlertRepository {
	return &alertRepositoryImpl{client: client}
}

func (r *alertRepositoryImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.Alert, int64, error) {
	q := r.client.Alert.Query()

	for field, filter := range opts.Filters {
		switch field {
		case "search":
			if val, ok := filter.Value.(string); ok && val != "" {
				q = q.Where(alert.TitleContainsFold(val))
			}
		case "severity":
			if val, ok := filter.Value.(string); ok && val != "" {
				q = q.Where(alert.SeverityEQ(alert.Severity(val)))
			}
		case "status":
			if val, ok := filter.Value.(string); ok && val != "" {
				q = q.Where(alert.StatusEQ(alert.Status(val)))
			}
		case "type":
			if val, ok := filter.Value.(string); ok && val != "" {
				q = q.Where(alert.TypeEQ(alert.Type(val)))
			}
		case "device_id":
			if val, ok := filter.Value.(string); ok && val != "" {
				q = q.Where(alert.DeviceIDEQ(val))
			}
		}
	}

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, apperror.ErrInternalServerError.WithMessage("Lỗi khi đếm alert").WithError(err)
	}

	if len(opts.Sort) > 0 {
		for _, sortField := range opts.Sort {
			switch strings.ToLower(sortField.Field) {
			case "created_at":
				if sortField.Desc {
					q = q.Order(ent.Desc(alert.FieldCreatedAt))
				} else {
					q = q.Order(ent.Asc(alert.FieldCreatedAt))
				}
			case "severity":
				if sortField.Desc {
					q = q.Order(ent.Desc(alert.FieldSeverity))
				} else {
					q = q.Order(ent.Asc(alert.FieldSeverity))
				}
			}
		}
	} else {
		q = q.Order(ent.Desc(alert.FieldCreatedAt))
	}

	alerts, err := q.Offset(offset).Limit(limit).All(ctx)
	if err != nil {
		return nil, 0, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất alerts").WithError(err)
	}

	return alerts, int64(total), nil
}

func (r *alertRepositoryImpl) GetByID(ctx context.Context, id uint) (*ent.Alert, error) {
	a, err := r.client.Alert.Get(ctx, id)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Alert không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất alert").WithError(err)
	}
	return a, nil
}

func (r *alertRepositoryImpl) Create(ctx context.Context, entity *ent.Alert) (*ent.Alert, error) {
	create := r.client.Alert.Create().
		SetTitle(entity.Title).
		SetSeverity(entity.Severity).
		SetType(entity.Type).
		SetStatus(entity.Status)

	if entity.DeviceID != "" {
		create = create.SetDeviceID(entity.DeviceID)
	}
	if entity.UserID != nil {
		create = create.SetUserID(*entity.UserID)
	}
	if entity.Details != nil {
		create = create.SetDetails(entity.Details)
	}

	a, err := create.Save(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo alert").WithError(err)
	}

	return a, nil
}

func (r *alertRepositoryImpl) UpdateStatus(ctx context.Context, id uint, status string) error {
	now := time.Now()
	update := r.client.Alert.UpdateOneID(id).SetStatus(alert.Status(status))
	
	if status == string(alert.StatusAcknowledged) {
		update.SetAcknowledgedAt(now)
	} else if status == string(alert.StatusResolved) {
		update.SetResolvedAt(now)
	}

	_, err := update.Save(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Alert không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi cập nhật alert").WithError(err)
	}
	return nil
}

func (r *alertRepositoryImpl) BulkUpdateStatus(ctx context.Context, ids []uint, status string) error {
	now := time.Now()
	update := r.client.Alert.Update().
		Where(alert.IDIn(ids...)).
		SetStatus(alert.Status(status))

	if status == string(alert.StatusAcknowledged) {
		update.SetAcknowledgedAt(now)
	} else if status == string(alert.StatusResolved) {
		update.SetResolvedAt(now)
	}

	_, err := update.Save(ctx)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi bulk update alert").WithError(err)
	}
	return nil
}

func (r *alertRepositoryImpl) GetStats(ctx context.Context) (*dto.AlertsSummaryResponse, error) {
	total, _ := r.client.Alert.Query().Count(ctx)
	open, _ := r.client.Alert.Query().Where(alert.StatusEQ(alert.StatusOpen)).Count(ctx)
	acknowledged, _ := r.client.Alert.Query().Where(alert.StatusEQ(alert.StatusAcknowledged)).Count(ctx)
	resolved, _ := r.client.Alert.Query().Where(alert.StatusEQ(alert.StatusResolved)).Count(ctx)

	critical, _ := r.client.Alert.Query().Where(alert.SeverityEQ(alert.SeverityCritical)).Count(ctx)
	high, _ := r.client.Alert.Query().Where(alert.SeverityEQ(alert.SeverityHigh)).Count(ctx)
	medium, _ := r.client.Alert.Query().Where(alert.SeverityEQ(alert.SeverityMedium)).Count(ctx)
	low, _ := r.client.Alert.Query().Where(alert.SeverityEQ(alert.SeverityLow)).Count(ctx)

	security, _ := r.client.Alert.Query().Where(alert.TypeEQ(alert.TypeSecurity)).Count(ctx)
	compliance, _ := r.client.Alert.Query().Where(alert.TypeEQ(alert.TypeCompliance)).Count(ctx)
	connectivity, _ := r.client.Alert.Query().Where(alert.TypeEQ(alert.TypeConnectivity)).Count(ctx)
	application, _ := r.client.Alert.Query().Where(alert.TypeEQ(alert.TypeApplication)).Count(ctx)
	deviceHealth, _ := r.client.Alert.Query().Where(alert.TypeEQ(alert.TypeDeviceHealth)).Count(ctx)

	return &dto.AlertsSummaryResponse{
		Total:        int64(total),
		Open:         int64(open),
		Acknowledged: int64(acknowledged),
		Resolved:     int64(resolved),
		BySeverity: map[string]int64{
			"critical": int64(critical),
			"high":     int64(high),
			"medium":   int64(medium),
			"low":      int64(low),
		},
		ByType: map[string]int64{
			"security":      int64(security),
			"compliance":    int64(compliance),
			"connectivity":  int64(connectivity),
			"application":   int64(application),
			"device_health": int64(deviceHealth),
		},
	}, nil
}

// Alert Rule Repository Implementation

type alertRuleRepositoryImpl struct {
	client *ent.Client
}

func NewAlertRuleRepository(client *ent.Client) repository.AlertRuleRepository {
	return &alertRuleRepositoryImpl{client: client}
}

func (r *alertRuleRepositoryImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.AlertRule, int64, error) {
	q := r.client.AlertRule.Query()

	for field, filter := range opts.Filters {
		switch field {
		case "search":
			if val, ok := filter.Value.(string); ok && val != "" {
				q = q.Where(alertrule.NameContainsFold(val))
			}
		case "enabled":
			if val, ok := filter.Value.(string); ok {
				enabled := val == "true"
				q = q.Where(alertrule.EnabledEQ(enabled))
			}
		}
	}

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, apperror.ErrInternalServerError.WithMessage("Lỗi khi đếm alert rules").WithError(err)
	}

	rules, err := q.Offset(offset).Limit(limit).Order(ent.Desc(alertrule.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, 0, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất alert rules").WithError(err)
	}

	return rules, int64(total), nil
}

func (r *alertRuleRepositoryImpl) GetByID(ctx context.Context, id uint) (*ent.AlertRule, error) {
	rule, err := r.client.AlertRule.Get(ctx, id)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Alert rule không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất alert rule").WithError(err)
	}
	return rule, nil
}

func (r *alertRuleRepositoryImpl) Create(ctx context.Context, entity *ent.AlertRule) (*ent.AlertRule, error) {
	rule, err := r.client.AlertRule.Create().
		SetName(entity.Name).
		SetDescription(entity.Description).
		SetCondition(entity.Condition).
		SetActions(entity.Actions).
		SetEnabled(entity.Enabled).
		Save(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo alert rule").WithError(err)
	}
	return rule, nil
}

func (r *alertRuleRepositoryImpl) Update(ctx context.Context, id uint, entity *ent.AlertRule) (*ent.AlertRule, error) {
	update := r.client.AlertRule.UpdateOneID(id)

	if entity.Name != "" {
		update.SetName(entity.Name)
	}
	if entity.Description != "" {
		update.SetDescription(entity.Description)
	}
	if entity.Condition != nil {
		update.SetCondition(entity.Condition)
	}
	if entity.Actions != nil {
		update.SetActions(entity.Actions)
	}
	
	// Because Enabled is boolean, we might want to let SetEnabled explicitly handle it 
	// but the service maps it, so we'll just always apply if called from service with new values.
	// We'll trust the domain model passed in has the correct final values if they were meant to be updated.
	update.SetEnabled(entity.Enabled)

	rule, err := update.Save(ctx)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Alert rule không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi cập nhật alert rule").WithError(err)
	}
	return rule, nil
}

func (r *alertRuleRepositoryImpl) Delete(ctx context.Context, id uint) error {
	err := r.client.AlertRule.DeleteOneID(id).Exec(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Alert rule không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi xóa alert rule").WithError(err)
	}
	return nil
}

func (r *alertRuleRepositoryImpl) SetEnabled(ctx context.Context, id uint, enabled bool) error {
	_, err := r.client.AlertRule.UpdateOneID(id).SetEnabled(enabled).Save(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Alert rule không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi toggle alert rule").WithError(err)
	}
	return nil
}
