package service

import (
	"context"

	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/interface/api/dto"
	"github.com/thienel/go-backend-template/pkg/query"
)

type CreateDeviceCommand struct {
	ID             string
	SerialNumber   string
	Model          string
	Name           string
	Platform       string
	OwnerID        *uint
	MacAddress     string
	IpAddress      string
	BatteryLevel   float64
	StorageCapacity uint64
	StorageUsed    uint64
	IsJailbroken   bool
	EnrollmentType string
}

type UpdateDeviceCommand struct {
	ID               string
	SerialNumber     *string
	Model            *string
	Name             *string
	Platform         *string
	Status           *string
	ComplianceStatus *string
	IsEnrolled       *bool
	OwnerID          *uint
	OsVersion        *string
	DeviceType       *string
	MacAddress       *string
	IpAddress        *string
	BatteryLevel     *float64
	StorageCapacity  *uint64
	StorageUsed      *uint64
	IsJailbroken     *bool
	EnrollmentType   *string
}

type DeviceService interface {
	List(ctx context.Context, offset, limit int, opts query.QueryOptions) ([]*ent.Device, int64, error)
	GetByID(ctx context.Context, id string) (*ent.Device, error)
	Create(ctx context.Context, cmd CreateDeviceCommand) (*ent.Device, error)
	Update(ctx context.Context, cmd UpdateDeviceCommand) (*ent.Device, error)
	Delete(ctx context.Context, id string) error
	GetStats(ctx context.Context) (*dto.DeviceStatsResponse, error)
	Export(ctx context.Context, format string) ([]byte, error)
}
