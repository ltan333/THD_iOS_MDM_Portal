package persistence

import (
	"context"
	"database/sql"
	"time"
)

// ─────────────────────────────────────────────────────────────────────────────
// Domain types — read-only mirrors of nanoMDM / nanoDEP tables.
//
// These structs are NOT managed by Ent. They map directly to the tables that
// the nano servers own. The portal treats them as read-only; all writes go
// through the nano server APIs (HTTP).
// ─────────────────────────────────────────────────────────────────────────────

// NanoDevice is a read-only view of nanoMDM's `devices` table.
// id = UDID (Apple MDM enrollment identifier).
type NanoDevice struct {
	ID              string
	SerialNumber    *string
	AuthenticateAt  time.Time
	TokenUpdateAt   *time.Time
	BootstrapToken  *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NanoEnrollment is a read-only view of nanoMDM's `enrollments` table.
// id = UDID; device_id = UDID of the parent device record.
type NanoEnrollment struct {
	ID               string
	DeviceID         string
	Type             string   // "Device" | "User"
	Topic            string   // APNs topic
	PushMagic        string
	TokenHex         string
	Enabled          bool
	TokenUpdateTally int
	LastSeenAt       time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NanoCommandResult is a read-only view of nanoMDM's `command_results` table.
type NanoCommandResult struct {
	EnrollmentID string
	CommandUUID  string
	Status       string // "Acknowledged" | "Error" | "CommandFormatError" | "Idle" | "NotNow"
	Result       string // raw XML plist response
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NanoDepName is a read-only view of nanoDEP's `dep_names` table.
type NanoDepName struct {
	Name                 string
	AccessToken          *string
	AccessTokenExpiry    *time.Time
	ConfigBaseURL        *string
	SyncerCursor         *string
	AssignerProfileUUID  *string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// ─────────────────────────────────────────────────────────────────────────────
// NanoRepository — read-only queries against nano-owned tables.
// ─────────────────────────────────────────────────────────────────────────────

type NanoRepository struct {
	db *sql.DB
}

func NewNanoRepository(db *sql.DB) *NanoRepository {
	return &NanoRepository{db: db}
}

// GetEnrollmentByUDID returns the active enrollment record for a UDID.
// Returns nil, nil if the device is not enrolled.
func (r *NanoRepository) GetEnrollmentByUDID(ctx context.Context, udid string) (*NanoEnrollment, error) {
	const q = `
		SELECT id, device_id, type, topic, push_magic, token_hex,
		       enabled, token_update_tally, last_seen_at, created_at, updated_at
		FROM enrollments
		WHERE id = $1 AND enabled = true
		LIMIT 1`

	row := r.db.QueryRowContext(ctx, q, udid)
	var e NanoEnrollment
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

// GetEnrollmentBySerialNumber returns the active enrollment for a device
// identified by serial number. Joins devices → enrollments.
func (r *NanoRepository) GetEnrollmentBySerialNumber(ctx context.Context, sn string) (*NanoEnrollment, error) {
	const q = `
		SELECT e.id, e.device_id, e.type, e.topic, e.push_magic, e.token_hex,
		       e.enabled, e.token_update_tally, e.last_seen_at, e.created_at, e.updated_at
		FROM enrollments e
		JOIN devices d ON d.id = e.device_id
		WHERE d.serial_number = $1 AND e.enabled = true
		LIMIT 1`

	row := r.db.QueryRowContext(ctx, q, sn)
	var e NanoEnrollment
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

// GetNanoDeviceByUDID reads a device record from nanoMDM's `devices` table.
func (r *NanoRepository) GetNanoDeviceByUDID(ctx context.Context, udid string) (*NanoDevice, error) {
	const q = `
		SELECT id, serial_number, authenticate_at, token_update_at,
		       bootstrap_token_b64, created_at, updated_at
		FROM devices
		WHERE id = $1
		LIMIT 1`

	row := r.db.QueryRowContext(ctx, q, udid)
	var d NanoDevice
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

// GetCommandResults returns the most recent command results for a UDID,
// ordered by most recently updated. Limit controls how many to fetch.
func (r *NanoRepository) GetCommandResults(ctx context.Context, udid string, limit int) ([]NanoCommandResult, error) {
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
	defer rows.Close()

	var results []NanoCommandResult
	for rows.Next() {
		var cr NanoCommandResult
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

// GetDepNames returns all DEP provider configurations from nanoDEP's table.
// Sensitive fields (keys, secrets) are intentionally excluded.
func (r *NanoRepository) GetDepNames(ctx context.Context) ([]NanoDepName, error) {
	const q = `
		SELECT name, access_token, access_token_expiry, config_base_url,
		       syncer_cursor, assigner_profile_uuid, created_at, updated_at
		FROM dep_names
		ORDER BY name`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []NanoDepName
	for rows.Next() {
		var n NanoDepName
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

// IsEnrolled is a lightweight check — returns true if nanoMDM has an active
// enrollment record for the UDID. Cheaper than fetching the full enrollment.
func (r *NanoRepository) IsEnrolled(ctx context.Context, udid string) (bool, error) {
	const q = `SELECT EXISTS(SELECT 1 FROM enrollments WHERE id = $1 AND enabled = true)`
	var exists bool
	err := r.db.QueryRowContext(ctx, q, udid).Scan(&exists)
	return exists, err
}
