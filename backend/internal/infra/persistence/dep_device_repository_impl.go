package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/thienel/go-backend-template/internal/domain/repository"
	"github.com/thienel/go-backend-template/internal/ent"
	"github.com/thienel/go-backend-template/internal/ent/depdevice"
	apperror "github.com/thienel/go-backend-template/pkg/error"
)

type depDeviceRepositoryImpl struct {
	client *ent.Client
}

func NewDepDeviceRepository(client *ent.Client) repository.DepDeviceRepository {
	return &depDeviceRepositoryImpl{client: client}
}

func (r *depDeviceRepositoryImpl) UpsertBatch(ctx context.Context, devices []*ent.DepDevice) error {
	// Use transaction for batch upsert
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Failed to start transaction").WithError(err)
	}

	for _, d := range devices {
		if err := r.upsertOne(ctx, tx.Client(), d); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return apperror.ErrInternalServerError.WithMessage("Failed to commit transaction").WithError(err)
	}
	return nil
}

func (r *depDeviceRepositoryImpl) upsertOne(ctx context.Context, client *ent.Client, d *ent.DepDevice) error {
	// Check if device already exists
	existing, err := client.DepDevice.
		Query().
		Where(depdevice.SerialNumberEQ(d.SerialNumber)).
		Only(ctx)

	if ent.IsNotFound(err) {
		// Insert new device
		create := client.DepDevice.Create().
			SetID(uuid.New().String()).
			SetSerialNumber(d.SerialNumber).
			SetDepName(d.DepName).
			SetModel(d.Model).
			SetDescription(d.Description).
			SetColor(d.Color).
			SetAssetTag(d.AssetTag).
			SetOs(d.Os).
			SetDeviceFamily(d.DeviceFamily).
			SetProfileUUID(d.ProfileUUID).
			SetProfileStatus(d.ProfileStatus).
			SetDeviceAssignedBy(d.DeviceAssignedBy).
			SetOpType(d.OpType).
			SetIsActive(d.IsActive).
			SetNeedsManualReassign(d.NeedsManualReassign).
			SetReassignError(d.ReassignError)

		if d.ProfileAssignTime != nil {
			create = create.SetProfileAssignTime(*d.ProfileAssignTime)
		}
		if d.ProfilePushTime != nil {
			create = create.SetProfilePushTime(*d.ProfilePushTime)
		}
		if d.DeviceAssignedDate != nil {
			create = create.SetDeviceAssignedDate(*d.DeviceAssignedDate)
		}
		if d.OpDate != nil {
			create = create.SetOpDate(*d.OpDate)
		}

		_, err = create.Save(ctx)
		if err != nil {
			return apperror.ErrInternalServerError.WithMessage("Failed to create DEP device").WithError(err)
		}
		return nil
	}

	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Failed to query DEP device").WithError(err)
	}

	// Update existing device
	update := client.DepDevice.UpdateOne(existing).
		SetModel(d.Model).
		SetDescription(d.Description).
		SetColor(d.Color).
		SetAssetTag(d.AssetTag).
		SetOs(d.Os).
		SetDeviceFamily(d.DeviceFamily).
		SetProfileUUID(d.ProfileUUID).
		SetProfileStatus(d.ProfileStatus).
		SetDeviceAssignedBy(d.DeviceAssignedBy).
		SetOpType(d.OpType).
		SetIsActive(d.IsActive).
		SetNeedsManualReassign(d.NeedsManualReassign).
		SetReassignError(d.ReassignError)

	if d.ProfileAssignTime != nil {
		update = update.SetProfileAssignTime(*d.ProfileAssignTime)
	} else {
		update = update.ClearProfileAssignTime()
	}
	if d.ProfilePushTime != nil {
		update = update.SetProfilePushTime(*d.ProfilePushTime)
	} else {
		update = update.ClearProfilePushTime()
	}
	if d.DeviceAssignedDate != nil {
		update = update.SetDeviceAssignedDate(*d.DeviceAssignedDate)
	} else {
		update = update.ClearDeviceAssignedDate()
	}
	if d.OpDate != nil {
		update = update.SetOpDate(*d.OpDate)
	} else {
		update = update.ClearOpDate()
	}

	_, err = update.Save(ctx)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Failed to update DEP device").WithError(err)
	}
	return nil
}

func (r *depDeviceRepositoryImpl) GetBySerialNumber(ctx context.Context, serialNumber string) (*ent.DepDevice, error) {
	device, err := r.client.DepDevice.
		Query().
		Where(depdevice.SerialNumberEQ(serialNumber)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperror.ErrNotFound.WithMessage("DEP device not found")
		}
		return nil, apperror.ErrInternalServerError.WithMessage("Failed to query DEP device").WithError(err)
	}
	return device, nil
}

func (r *depDeviceRepositoryImpl) ListNeedsManualReassign(ctx context.Context) ([]*ent.DepDevice, error) {
	devices, err := r.client.DepDevice.
		Query().
		Where(
			depdevice.NeedsManualReassignEQ(true),
			depdevice.IsActiveEQ(true),
		).
		All(ctx)
	if err != nil {
		return nil, apperror.ErrInternalServerError.WithMessage("Failed to list DEP devices needing reassign").WithError(err)
	}
	return devices, nil
}

func (r *depDeviceRepositoryImpl) MarkInactive(ctx context.Context, serialNumber string) error {
	_, err := r.client.DepDevice.
		Update().
		Where(depdevice.SerialNumberEQ(serialNumber)).
		SetIsActive(false).
		Save(ctx)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Failed to mark DEP device inactive").WithError(err)
	}
	return nil
}

func (r *depDeviceRepositoryImpl) UpdateReassignStatus(ctx context.Context, serialNumber string, needsReassign bool, errMsg string) error {
	_, err := r.client.DepDevice.
		Update().
		Where(depdevice.SerialNumberEQ(serialNumber)).
		SetNeedsManualReassign(needsReassign).
		SetReassignError(errMsg).
		Save(ctx)
	if err != nil {
		return apperror.ErrInternalServerError.WithMessage("Failed to update DEP device reassign status").WithError(err)
	}
	return nil
}

func (r *depDeviceRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*ent.DepDevice, int64, error) {
	q := r.client.DepDevice.Query()

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, apperror.ErrInternalServerError.WithMessage("Failed to count DEP devices").WithError(err)
	}

	devices, err := q.
		Offset(offset).
		Limit(limit).
		Order(ent.Desc(depdevice.FieldUpdatedAt)).
		All(ctx)
	if err != nil {
		return nil, 0, apperror.ErrInternalServerError.WithMessage("Failed to list DEP devices").WithError(err)
	}

	return devices, int64(total), nil
}
