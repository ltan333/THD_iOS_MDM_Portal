package service

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
)

// DepDeviceService handles DEP device operations.
type DepDeviceService interface {
	// HandleDEPDeviceEvent processes DEP device events (FetchDevices/SyncDevices).
	// It compares device profile_uuid with assigner profile and reassigns if needed,
	// then upserts devices to the database.
	HandleDEPDeviceEvent(ctx context.Context, depName string, devices []dto.DEPDevice, assignerProfileUUID string, nanomdmSvc NanoMDMService) error

	// ListNeedsManualReassign returns devices that need manual reassignment.
	ListNeedsManualReassign(ctx context.Context) ([]*ent.DepDevice, error)

	// List returns a paginated list of DEP devices.
	List(ctx context.Context, offset, limit int) ([]*ent.DepDevice, int64, error)

	// GetBySerialNumber retrieves a DEP device by its serial number.
	GetBySerialNumber(ctx context.Context, serialNumber string) (*ent.DepDevice, error)
}
