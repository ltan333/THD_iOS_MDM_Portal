package repository

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/pkg/query"
)

// DeviceStats represents the aggregated statistics for devices
type DeviceStats struct {
	Total        int64
	Active       int64
	Inactive     int64
	Enrolled     int64
	ByPlatform   map[string]int64
	ByStatus     map[string]int64
	Compliant    int64
	NonCompliant int64
}

type DeviceRepository interface {
	List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.Device, int64, error)
	GetByID(ctx context.Context, id string) (*ent.Device, error)
	Create(ctx context.Context, entity *ent.Device) (*ent.Device, error)
	Update(ctx context.Context, id string, entity *ent.Device) (*ent.Device, error)
	Delete(ctx context.Context, id string) error

	// Specific queries
	FindByUDID(ctx context.Context, udid string) (*ent.Device, error)
	FindBySerialNumber(ctx context.Context, sn string) (*ent.Device, error)
	GetStats(ctx context.Context) (*DeviceStats, error)
	GetAll(ctx context.Context) ([]*ent.Device, error)

	// Deprecated / Kept for backwards compatibility if needed specifically
	// GetUDID(ctx context.Context, id string) (string, error)
	// UpdateUDID(ctx context.Context, id string, udid string) error
	// UpdateStatusByUDID(ctx context.Context, udid string, status string) error
	// UpdateInventoryByUDID(ctx context.Context, udid string, osVersion, modelName, deviceName, batteryLevel, availableCapacity string) error

	// For specific Webhook/DEP updates
	UpsertFromDEP(ctx context.Context, devices []map[string]any) error
	EnsureMinimalByUDID(ctx context.Context, udid string, sn string) error
	UpdateTokenEnrolledBySN(ctx context.Context, udid string, sn string, model string, osVer string) error
	UpdateTokenEnrolledByUDID(ctx context.Context, udid string, sn string, model string, osVer string) error
	CreateEnrolledDevice(ctx context.Context, udid string, sn string, model string, osVer string) error
	UpdateCheckOut(ctx context.Context, udid string) error
	ApplyDeviceInformation(ctx context.Context, udid string, qr map[string]any) error
	ReconcileBySerialAndUDID(ctx context.Context, serial string, udid string) error
}
