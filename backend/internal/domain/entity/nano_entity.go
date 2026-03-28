package entity

import (
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
	ID             string
	SerialNumber   *string
	AuthenticateAt time.Time
	TokenUpdateAt  *time.Time
	BootstrapToken *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NanoEnrollment is a read-only view of nanoMDM's `enrollments` table.
// id = UDID; device_id = UDID of the parent device record.
type NanoEnrollment struct {
	ID               string
	DeviceID         string
	Type             string // "Device" | "User"
	Topic            string // APNs topic
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
	Name                string
	AccessToken         *string
	AccessTokenExpiry   *time.Time
	ConfigBaseURL       *string
	SyncerCursor        *string
	AssignerProfileUUID *string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
