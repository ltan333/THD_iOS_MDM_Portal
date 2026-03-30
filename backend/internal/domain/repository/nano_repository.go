package repository

import (
	"context"

	"github.com/thienel/go-backend-template/internal/domain/entity"
)

// NanoRepository — read-only interface against nano-owned tables.
type NanoRepository interface {
	// GetEnrollmentByUDID returns the active enrollment record for a UDID.
	// Returns nil, nil if the device is not enrolled.
	GetEnrollmentByUDID(ctx context.Context, udid string) (*entity.NanoEnrollment, error)

	// GetEnrollmentBySerialNumber returns the active enrollment for a device
	// identified by serial number.
	GetEnrollmentBySerialNumber(ctx context.Context, sn string) (*entity.NanoEnrollment, error)

	// GetNanoDeviceByUDID reads a device record from nanoMDM's `devices` table.
	GetNanoDeviceByUDID(ctx context.Context, udid string) (*entity.NanoDevice, error)

	// GetCommandResults returns the most recent command results for a UDID,
	// ordered by most recently updated. Limit controls how many to fetch.
	GetCommandResults(ctx context.Context, udid string, limit int) ([]entity.NanoCommandResult, error)

	// GetDepNames returns all DEP provider configurations from nanoDEP's table.
	GetDepNames(ctx context.Context) ([]entity.NanoDepName, error)

	// IsEnrolled is a lightweight check — returns true if nanoMDM has an active
	// enrollment record for the UDID.
	IsEnrolled(ctx context.Context, udid string) (bool, error)
}
