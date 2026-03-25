package serviceimpl

import (
	"context"
	"strings"
	"time"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/alert"
	"github.com/thienel/go-backend-template/internal/ent/alertrule"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
)

type alertServiceImpl struct {
	client *ent.Client
}

func NewAlertService(client *ent.Client) service.AlertService {
	return &alertServiceImpl{client: client}
}

func (s *alertServiceImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.Alert, int64, error) {
	q := s.client.Alert.Query()

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

func (s *alertServiceImpl) GetByID(ctx context.Context, id uint) (*ent.Alert, error) {
	if id == 0 {
		return nil, apperror.ErrValidation.WithMessage("ID alert là bắt buộc")
	}

	a, err := s.client.Alert.Get(ctx, id)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Alert không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất alert").WithError(err)
	}

	return a, nil
}

func (s *alertServiceImpl) Create(ctx context.Context, cmd service.CreateAlertCommand) (*ent.Alert, error) {
	if strings.TrimSpace(cmd.Title) == "" {
		return nil, apperror.ErrValidation.WithMessage("Tiêu đề alert là bắt buộc")
	}

	create := s.client.Alert.Create().
		SetTitle(cmd.Title).
		SetSeverity(alert.Severity(cmd.Severity)).
		SetType(alert.Type(cmd.Type)).
		SetStatus(alert.StatusOpen)

	if cmd.DeviceID != "" {
		create = create.SetDeviceID(cmd.DeviceID)
	}
	if cmd.UserID != nil {
		create = create.SetUserID(*cmd.UserID)
	}
	if cmd.Details != nil {
		create = create.SetDetails(cmd.Details)
	}

	a, err := create.Save(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo alert").WithError(err)
	}

	return a, nil
}

func (s *alertServiceImpl) Acknowledge(ctx context.Context, id uint) error {
	now := time.Now()
	_, err := s.client.Alert.UpdateOneID(id).
		SetStatus(alert.StatusAcknowledged).
		SetAcknowledgedAt(now).
		Save(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Alert không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi acknowledge alert").WithError(err)
	}
	return nil
}

func (s *alertServiceImpl) Resolve(ctx context.Context, id uint) error {
	now := time.Now()
	_, err := s.client.Alert.UpdateOneID(id).
		SetStatus(alert.StatusResolved).
		SetResolvedAt(now).
		Save(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Alert không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi resolve alert").WithError(err)
	}
	return nil
}

func (s *alertServiceImpl) BulkResolve(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return apperror.ErrValidation.WithMessage("Danh sách ID không được rỗng")
	}

	now := time.Now()
	_, err := s.client.Alert.Update().
		Where(alert.IDIn(ids...)).
		SetStatus(alert.StatusResolved).
		SetResolvedAt(now).
		Save(ctx)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi bulk resolve").WithError(err)
	}
	return nil
}

func (s *alertServiceImpl) GetStats(ctx context.Context) (*dto.AlertsSummaryResponse, error) {
	total, _ := s.client.Alert.Query().Count(ctx)
	open, _ := s.client.Alert.Query().Where(alert.StatusEQ(alert.StatusOpen)).Count(ctx)
	acknowledged, _ := s.client.Alert.Query().Where(alert.StatusEQ(alert.StatusAcknowledged)).Count(ctx)
	resolved, _ := s.client.Alert.Query().Where(alert.StatusEQ(alert.StatusResolved)).Count(ctx)

	critical, _ := s.client.Alert.Query().Where(alert.SeverityEQ(alert.SeverityCritical)).Count(ctx)
	high, _ := s.client.Alert.Query().Where(alert.SeverityEQ(alert.SeverityHigh)).Count(ctx)
	medium, _ := s.client.Alert.Query().Where(alert.SeverityEQ(alert.SeverityMedium)).Count(ctx)
	low, _ := s.client.Alert.Query().Where(alert.SeverityEQ(alert.SeverityLow)).Count(ctx)

	security, _ := s.client.Alert.Query().Where(alert.TypeEQ(alert.TypeSecurity)).Count(ctx)
	compliance, _ := s.client.Alert.Query().Where(alert.TypeEQ(alert.TypeCompliance)).Count(ctx)
	connectivity, _ := s.client.Alert.Query().Where(alert.TypeEQ(alert.TypeConnectivity)).Count(ctx)
	application, _ := s.client.Alert.Query().Where(alert.TypeEQ(alert.TypeApplication)).Count(ctx)
	deviceHealth, _ := s.client.Alert.Query().Where(alert.TypeEQ(alert.TypeDeviceHealth)).Count(ctx)

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

func (s *alertServiceImpl) LockDevice(ctx context.Context, alertID uint) error {
	a, err := s.GetByID(ctx, alertID)
	if err != nil {
		return err
	}
	if a.DeviceID == "" {
		return apperror.ErrBadRequest.WithMessage("Alert không có device liên quan")
	}
	// In real implementation, this would call MDM to lock device
	return nil
}

func (s *alertServiceImpl) WipeDevice(ctx context.Context, alertID uint) error {
	a, err := s.GetByID(ctx, alertID)
	if err != nil {
		return err
	}
	if a.DeviceID == "" {
		return apperror.ErrBadRequest.WithMessage("Alert không có device liên quan")
	}
	// In real implementation, this would call MDM to wipe device
	return nil
}

func (s *alertServiceImpl) PushPolicy(ctx context.Context, alertID uint, policyID uint) error {
	_, err := s.GetByID(ctx, alertID)
	if err != nil {
		return err
	}
	// In real implementation, this would push policy to device
	return nil
}

func (s *alertServiceImpl) SendMessage(ctx context.Context, alertID uint, message string) error {
	_, err := s.GetByID(ctx, alertID)
	if err != nil {
		return err
	}
	// In real implementation, this would send message to user
	return nil
}

// AlertRuleService implementation
type alertRuleServiceImpl struct {
	client *ent.Client
}

func NewAlertRuleService(client *ent.Client) service.AlertRuleService {
	return &alertRuleServiceImpl{client: client}
}

func (s *alertRuleServiceImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.AlertRule, int64, error) {
	q := s.client.AlertRule.Query()

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

func (s *alertRuleServiceImpl) GetByID(ctx context.Context, id uint) (*ent.AlertRule, error) {
	r, err := s.client.AlertRule.Get(ctx, id)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Alert rule không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất alert rule").WithError(err)
	}
	return r, nil
}

func (s *alertRuleServiceImpl) Create(ctx context.Context, cmd service.CreateAlertRuleCommand) (*ent.AlertRule, error) {
	if strings.TrimSpace(cmd.Name) == "" {
		return nil, apperror.ErrValidation.WithMessage("Tên rule là bắt buộc")
	}

	r, err := s.client.AlertRule.Create().
		SetName(cmd.Name).
		SetDescription(cmd.Description).
		SetCondition(cmd.Condition).
		SetActions(cmd.Actions).
		SetEnabled(cmd.Enabled).
		Save(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi tạo alert rule").WithError(err)
	}

	return r, nil
}

func (s *alertRuleServiceImpl) Update(ctx context.Context, cmd service.UpdateAlertRuleCommand) (*ent.AlertRule, error) {
	if cmd.ID == 0 {
		return nil, apperror.ErrValidation.WithMessage("ID rule là bắt buộc")
	}

	update := s.client.AlertRule.UpdateOneID(cmd.ID)

	if cmd.Name != nil {
		update = update.SetName(*cmd.Name)
	}
	if cmd.Description != nil {
		update = update.SetDescription(*cmd.Description)
	}
	if cmd.Condition != nil {
		update = update.SetCondition(cmd.Condition)
	}
	if cmd.Actions != nil {
		update = update.SetActions(cmd.Actions)
	}
	if cmd.Enabled != nil {
		update = update.SetEnabled(*cmd.Enabled)
	}

	r, err := update.Save(ctx)
	if ent.IsNotFound(err) {
		return nil, apperror.ErrNotFound.WithMessage("Alert rule không tồn tại")
	}
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Lỗi khi cập nhật alert rule").WithError(err)
	}

	return r, nil
}

func (s *alertRuleServiceImpl) Delete(ctx context.Context, id uint) error {
	err := s.client.AlertRule.DeleteOneID(id).Exec(ctx)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Alert rule không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi xóa alert rule").WithError(err)
	}
	return nil
}

func (s *alertRuleServiceImpl) Toggle(ctx context.Context, id uint) error {
	r, err := s.client.AlertRule.Get(ctx, id)
	if ent.IsNotFound(err) {
		return apperror.ErrNotFound.WithMessage("Alert rule không tồn tại")
	}
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi truy xuất alert rule").WithError(err)
	}

	_, err = s.client.AlertRule.UpdateOneID(id).
		SetEnabled(!r.Enabled).
		Save(ctx)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Lỗi khi toggle alert rule").WithError(err)
	}

	return nil
}
