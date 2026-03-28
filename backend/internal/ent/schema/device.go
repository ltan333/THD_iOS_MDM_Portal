package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Device holds the schema definition for the Device entity.
type Device struct {
	ent.Schema
}

// Annotations of the Device.
func (Device) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "portal_devices"},
	}
}

// Fields of the Device.
func (Device) Fields() []ent.Field {
	return []ent.Field{
		// id is the portal's stable primary key.
		// For DEP-synced devices it is a generated UUID; for directly-enrolled
		// devices it defaults to the UDID for backward compatibility.
		field.String("id").Unique().NotEmpty().SchemaType(map[string]string{
			"postgres": "character varying(255)",
		}),
		// udid is the Apple MDM enrollment identifier (UDID).
		// NULL for DEP devices that have not enrolled yet.
		// Always set on mdm.TokenUpdate (enrollment/token refresh).
		field.String("udid").Unique().Optional().Nillable().SchemaType(map[string]string{
			"postgres": "character varying(255)",
		}),
		field.String("serial_number").Unique().Optional().SchemaType(map[string]string{
			"postgres": "character varying(127)",
		}),
		field.String("model").Optional().SchemaType(map[string]string{
			"postgres": "character varying(255)",
		}),
		field.Uint("owner_id").Optional(),
		field.Bool("is_enrolled").Default(false),
		field.String("name").Optional(),
		field.Time("last_sync").Optional(),
		// New fields for Phase 1
		field.Enum("platform").Values("ios", "android", "windows", "macos", "other").Default("other").Optional(),
		field.Enum("status").Values("active", "inactive", "pending", "lost", "wiped").Default("pending").Optional(),
		field.Enum("compliance_status").Values("compliant", "non_compliant", "unknown").Default("unknown").Optional(),
		field.String("os_version").Optional(),
		field.String("device_type").Optional(), // iphone, ipad, android_phone, etc.
		field.Time("last_seen").Optional(),
		field.Time("enrolled_at").Optional(),
		// Newly added fields for comprehensive device health & network
		field.String("mac_address").Optional(),
		field.String("ip_address").Optional(),
		field.Float("battery_level").Optional().Comment("Battery level percentage (0-100)"),
		field.Uint64("storage_capacity").Optional().Comment("Total storage in bytes"),
		field.Uint64("storage_used").Optional().Comment("Used storage in bytes"),
		field.Bool("is_jailbroken").Default(false).Comment("True if device is jailbroken/rooted"),
		field.Enum("enrollment_type").Values("dep", "qr", "manual", "unknown").Default("unknown").Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Device.
func (Device) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("devices").
			Unique().
			Field("owner_id"),
		edge.From("groups", DeviceGroup.Type).
			Ref("devices"),
	}
}
