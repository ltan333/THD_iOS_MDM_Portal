package serviceimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	"github.com/thienel/tlog"
	"go.uber.org/zap"
)

type depDeviceServiceImpl struct {
	repo repository.DepDeviceRepository
}

func NewDepDeviceService(repo repository.DepDeviceRepository) service.DepDeviceService {
	return &depDeviceServiceImpl{repo: repo}
}

func (s *depDeviceServiceImpl) HandleDEPDeviceEvent(
	ctx context.Context,
	depName string,
	devices []dto.DEPDevice,
	assignerProfileUUID string,
	nanomdmSvc service.NanoMDMService,
) error {
	// Identify devices that need profile reassignment
	var needsAssign []string
	for _, d := range devices {
		if d.OpType != "deleted" &&
			assignerProfileUUID != "" &&
			d.ProfileUUID != assignerProfileUUID {
			needsAssign = append(needsAssign, d.SerialNumber)
			tlog.Debug("Device needs profile reassignment",
				zap.String("serial", d.SerialNumber),
				zap.String("current_profile", d.ProfileUUID),
				zap.String("assigner_profile", assignerProfileUUID))
		}
	}

	// Batch reassign profiles if needed
	assignResults := map[string]string{}
	if len(needsAssign) > 0 && nanomdmSvc != nil {
		results, err := nanomdmSvc.AssignDEPProfile(ctx, depName, assignerProfileUUID, needsAssign)
		if err != nil {
			tlog.Error("Failed to assign DEP profile batch",
				zap.String("dep_name", depName),
				zap.Error(err))
			// Don't return error - continue to upsert devices with current status
		} else {
			assignResults = results
		}
	}

	// Map webhook devices to ent.DepDevice and apply assignment results
	depDevices := make([]*ent.DepDevice, 0, len(devices))
	for _, d := range devices {
		depDev := mapWebhookDeviceToEnt(d, depName)

		// Handle op_type = deleted
		if d.OpType == "deleted" {
			depDev.IsActive = false
			depDevices = append(depDevices, depDev)
			continue
		}

		// Apply assignment results if this device was reassigned
		if result, ok := assignResults[d.SerialNumber]; ok {
			switch result {
			case "SUCCESS":
				depDev.ProfileUUID = assignerProfileUUID
				depDev.NeedsManualReassign = false
				depDev.ReassignError = ""
				tlog.Info("DEP profile assigned successfully",
					zap.String("serial", d.SerialNumber))
			case "NOT_ACCESSIBLE":
				depDev.NeedsManualReassign = true
				depDev.ReassignError = "NOT_ACCESSIBLE: device belongs to another MDM server"
				tlog.Warn("Device not accessible for reassign — needs manual action",
					zap.String("serial", d.SerialNumber))
			case "NOT_FOUND":
				depDev.NeedsManualReassign = false
				depDev.ReassignError = "NOT_FOUND: serial not in ABM"
				tlog.Warn("Device serial not found in ABM",
					zap.String("serial", d.SerialNumber))
			default: // FAILED or other
				depDev.NeedsManualReassign = false
				depDev.ReassignError = fmt.Sprintf("FAILED: %s", result)
				tlog.Error("Failed to assign DEP profile",
					zap.String("serial", d.SerialNumber),
					zap.String("result", result))
			}
		}

		depDevices = append(depDevices, depDev)
	}

	// Batch upsert to database
	if err := s.repo.UpsertBatch(ctx, depDevices); err != nil {
		return fmt.Errorf("upsert dep devices: %w", err)
	}

	tlog.Info("DEP device event processed",
		zap.String("dep_name", depName),
		zap.Int("total_devices", len(devices)),
		zap.Int("reassigned", len(needsAssign)))

	return nil
}

func (s *depDeviceServiceImpl) ListNeedsManualReassign(ctx context.Context) ([]*ent.DepDevice, error) {
	return s.repo.ListNeedsManualReassign(ctx)
}

func (s *depDeviceServiceImpl) List(ctx context.Context, offset, limit int) ([]*ent.DepDevice, int64, error) {
	return s.repo.List(ctx, offset, limit)
}

func (s *depDeviceServiceImpl) GetBySerialNumber(ctx context.Context, serialNumber string) (*ent.DepDevice, error) {
	return s.repo.GetBySerialNumber(ctx, serialNumber)
}

// mapWebhookDeviceToEnt converts a DEPDevice DTO to an ent.DepDevice entity.
func mapWebhookDeviceToEnt(d dto.DEPDevice, depName string) *ent.DepDevice {
	dev := &ent.DepDevice{
		SerialNumber:        d.SerialNumber,
		DepName:             depName,
		Model:               d.Model,
		Description:         d.Description,
		Color:               d.Color,
		AssetTag:            d.AssetTag,
		Os:                  d.OS,
		DeviceFamily:        d.DeviceFamily,
		ProfileUUID:         d.ProfileUUID,
		ProfileStatus:       d.ProfileStatus,
		DeviceAssignedBy:    d.DeviceAssignedBy,
		OpType:              d.OpType,
		IsActive:            true,
		NeedsManualReassign: false,
		ReassignError:       "",
	}

	// Parse and set nullable time fields
	if t := parseTime(d.ProfileAssignTime); t != nil {
		dev.ProfileAssignTime = t
	}
	if t := parseTime(d.ProfilePushTime); t != nil {
		dev.ProfilePushTime = t
	}
	if t := parseTime(d.DeviceAssignedDate); t != nil {
		dev.DeviceAssignedDate = t
	}
	if t := parseTime(d.OpDate); t != nil {
		dev.OpDate = t
	}

	return dev
}

// parseTime parses a time string from DEP webhook payload.
// Returns nil if the string is empty or parsing fails.
func parseTime(s string) *time.Time {
	if s == "" || s == "0001-01-01T00:00:00Z" {
		return nil
	}
	// Try multiple formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02",
	}
	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return &t
		}
	}
	return nil
}
