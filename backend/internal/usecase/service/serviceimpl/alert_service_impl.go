package serviceimpl

import (
	"context"
	"strings"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/alert"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	apperror "github.com/thienel/go-backend-template/pkg/error"
	"github.com/thienel/go-backend-template/pkg/query"
)

type alertServiceImpl struct {
	repo repository.AlertRepository
}

func NewAlertService(repo repository.AlertRepository) service.AlertService {
	return &alertServiceImpl{repo: repo}
}

func (s *alertServiceImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.Alert, int64, error) {
	return s.repo.List(ctx, offset, limit, opts)
}

func (s *alertServiceImpl) GetByID(ctx context.Context, id uint) (*ent.Alert, error) {
	if id == 0 {
		return nil, apperror.ErrValidation.WithMessage("ID alert là bắt buộc")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *alertServiceImpl) Create(ctx context.Context, cmd service.CreateAlertCommand) (*ent.Alert, error) {
	if strings.TrimSpace(cmd.Title) == "" {
		return nil, apperror.ErrValidation.WithMessage("Tiêu đề alert là bắt buộc")
	}
	
	return s.repo.Create(ctx, &ent.Alert{
		Title:    cmd.Title,
		Severity: alert.Severity(cmd.Severity),
		Type:     alert.Type(cmd.Type),
		Status:   alert.StatusOpen,
		DeviceID: cmd.DeviceID,
		UserID:   cmd.UserID,
		Details:  cmd.Details,
	})
}

func (s *alertServiceImpl) Acknowledge(ctx context.Context, id uint) error {
	return s.repo.UpdateStatus(ctx, id, string(alert.StatusAcknowledged))
}

func (s *alertServiceImpl) Resolve(ctx context.Context, id uint) error {
	return s.repo.UpdateStatus(ctx, id, string(alert.StatusResolved))
}

func (s *alertServiceImpl) BulkResolve(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return apperror.ErrValidation.WithMessage("Danh sách ID không được rỗng")
	}
	return s.repo.BulkUpdateStatus(ctx, ids, string(alert.StatusResolved))
}

func (s *alertServiceImpl) GetStats(ctx context.Context) (*dto.AlertsSummaryResponse, error) {
	return s.repo.GetStats(ctx)
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
	repo repository.AlertRuleRepository
}

func NewAlertRuleService(repo repository.AlertRuleRepository) service.AlertRuleService {
	return &alertRuleServiceImpl{repo: repo}
}

func (s *alertRuleServiceImpl) List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.AlertRule, int64, error) {
	return s.repo.List(ctx, offset, limit, opts)
}

func (s *alertRuleServiceImpl) GetByID(ctx context.Context, id uint) (*ent.AlertRule, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *alertRuleServiceImpl) Create(ctx context.Context, cmd service.CreateAlertRuleCommand) (*ent.AlertRule, error) {
	if strings.TrimSpace(cmd.Name) == "" {
		return nil, apperror.ErrValidation.WithMessage("Tên rule là bắt buộc")
	}

	return s.repo.Create(ctx, &ent.AlertRule{
		Name:        cmd.Name,
		Description: cmd.Description,
		Condition:   cmd.Condition,
		Actions:     cmd.Actions,
		Enabled:     cmd.Enabled,
	})
}

func (s *alertRuleServiceImpl) Update(ctx context.Context, cmd service.UpdateAlertRuleCommand) (*ent.AlertRule, error) {
	if cmd.ID == 0 {
		return nil, apperror.ErrValidation.WithMessage("ID rule là bắt buộc")
	}

	r, err := s.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, apperror.ErrNotFound.WithMessage("Alert rule không tồn tại")
	}

	name := r.Name
	if cmd.Name != nil {
		name = *cmd.Name
	}
	desc := r.Description
	if cmd.Description != nil {
		desc = *cmd.Description
	}
	cond := r.Condition
	if cmd.Condition != nil {
		cond = cmd.Condition
	}
	act := r.Actions
	if cmd.Actions != nil {
		act = cmd.Actions
	}
	enab := r.Enabled
	if cmd.Enabled != nil {
		enab = *cmd.Enabled
	}

	return s.repo.Update(ctx, cmd.ID, &ent.AlertRule{
		Name:        name,
		Description: desc,
		Condition:   cond,
		Actions:     act,
		Enabled:     enab,
	})
}

func (s *alertRuleServiceImpl) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *alertRuleServiceImpl) Toggle(ctx context.Context, id uint) error {
	r, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.SetEnabled(ctx, id, !r.Enabled)
}
