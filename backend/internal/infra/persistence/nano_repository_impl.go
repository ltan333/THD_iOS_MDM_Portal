package persistence

import (
	"context"
	"database/sql"

	"github.com/thienel/go-backend-template/internal/domain/entity"
	"github.com/thienel/go-backend-template/internal/domain/repository"
)

// nanoRepositoryImpl implements repository.NanoRepository
type nanoRepositoryImpl struct {
	db *sql.DB
}

func NewNanoRepository(db *sql.DB) repository.NanoRepository {
	return &nanoRepositoryImpl{db: db}
}

func (r *nanoRepositoryImpl) GetEnrollmentByUDID(ctx context.Context, udid string) (*entity.NanoEnrollment, error) {
	const q = `
		SELECT id, device_id, type, topic, push_magic, token_hex,
		       enabled, token_update_tally, last_seen_at, created_at, updated_at
		FROM enrollments
		WHERE id = $1 AND enabled = true
		LIMIT 1`

	row := r.db.QueryRowContext(ctx, q, udid)
	var e entity.NanoEnrollment
	err := row.Scan(
		&e.ID, &e.DeviceID, &e.Type, &e.Topic, &e.PushMagic, &e.TokenHex,
		&e.Enabled, &e.TokenUpdateTally, &e.LastSeenAt, &e.CreatedAt, &e.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *nanoRepositoryImpl) GetEnrollmentBySerialNumber(ctx context.Context, sn string) (*entity.NanoEnrollment, error) {
	const q = `
		SELECT e.id, e.device_id, e.type, e.topic, e.push_magic, e.token_hex,
		       e.enabled, e.token_update_tally, e.last_seen_at, e.created_at, e.updated_at
		FROM enrollments e
		JOIN devices d ON d.id = e.device_id
		WHERE d.serial_number = $1 AND e.enabled = true
		LIMIT 1`

	row := r.db.QueryRowContext(ctx, q, sn)
	var e entity.NanoEnrollment
	err := row.Scan(
		&e.ID, &e.DeviceID, &e.Type, &e.Topic, &e.PushMagic, &e.TokenHex,
		&e.Enabled, &e.TokenUpdateTally, &e.LastSeenAt, &e.CreatedAt, &e.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *nanoRepositoryImpl) GetNanoDeviceByUDID(ctx context.Context, udid string) (*entity.NanoDevice, error) {
	const q = `
		SELECT id, serial_number, authenticate_at, token_update_at,
		       bootstrap_token_b64, created_at, updated_at
		FROM devices
		WHERE id = $1
		LIMIT 1`

	row := r.db.QueryRowContext(ctx, q, udid)
	var d entity.NanoDevice
	err := row.Scan(
		&d.ID, &d.SerialNumber, &d.AuthenticateAt, &d.TokenUpdateAt,
		&d.BootstrapToken, &d.CreatedAt, &d.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *nanoRepositoryImpl) GetCommandResults(ctx context.Context, udid string, limit int) ([]entity.NanoCommandResult, error) {
	if limit <= 0 {
		limit = 20
	}
	const q = `
		SELECT id, command_uuid, status, result, created_at, updated_at
		FROM command_results
		WHERE id = $1
		ORDER BY updated_at DESC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, q, udid, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	var results []entity.NanoCommandResult
	for rows.Next() {
		var cr entity.NanoCommandResult
		if err := rows.Scan(
			&cr.EnrollmentID, &cr.CommandUUID, &cr.Status,
			&cr.Result, &cr.CreatedAt, &cr.UpdatedAt,
		); err != nil {
			return nil, err
		}
		results = append(results, cr)
	}
	return results, rows.Err()
}

func (r *nanoRepositoryImpl) GetDepNames(ctx context.Context) ([]entity.NanoDepName, error) {
	const q = `
		SELECT name, access_token, access_token_expiry, config_base_url,
		       syncer_cursor, assigner_profile_uuid, created_at, updated_at
		FROM dep_names
		ORDER BY name`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	var names []entity.NanoDepName
	for rows.Next() {
		var n entity.NanoDepName
		if err := rows.Scan(
			&n.Name, &n.AccessToken, &n.AccessTokenExpiry, &n.ConfigBaseURL,
			&n.SyncerCursor, &n.AssignerProfileUUID, &n.CreatedAt, &n.UpdatedAt,
		); err != nil {
			return nil, err
		}
		names = append(names, n)
	}
	return names, rows.Err()
}

func (r *nanoRepositoryImpl) IsEnrolled(ctx context.Context, udid string) (bool, error) {
	const q = `SELECT EXISTS(SELECT 1 FROM enrollments WHERE id = $1 AND enabled = true)`
	var exists bool
	err := r.db.QueryRowContext(ctx, q, udid).Scan(&exists)
	return exists, err
}
