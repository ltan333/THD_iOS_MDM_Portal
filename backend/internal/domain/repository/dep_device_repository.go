package repository

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
)

// DepDeviceRepository defines the interface for DEP device persistence operations.
type DepDeviceRepository interface {
	// UpsertBatch upserts multiple devices from DEP webhook.
	// If a device exists (by serial_number), it will be updated; otherwise, it will be created.
	UpsertBatch(ctx context.Context, devices []*ent.DepDevice) error

	// GetBySerialNumber retrieves a DEP device by its serial number.
	GetBySerialNumber(ctx context.Context, serialNumber string) (*ent.DepDevice, error)

	// ListNeedsManualReassign returns devices that need manual reassignment.
	ListNeedsManualReassign(ctx context.Context) ([]*ent.DepDevice, error)

	// MarkInactive marks a device as inactive (when op_type = deleted).
	MarkInactive(ctx context.Context, serialNumber string) error

	// UpdateReassignStatus updates the reassignment status for a device.
	UpdateReassignStatus(ctx context.Context, serialNumber string, needsReassign bool, errMsg string) error

	// List returns a paginated list of DEP devices.
	List(ctx context.Context, offset, limit int) ([]*ent.DepDevice, int64, error)
}
