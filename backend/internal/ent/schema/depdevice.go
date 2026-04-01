package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// DepDevice holds the schema definition for the DepDevice entity.
// This entity stores devices from Apple DEP (Device Enrollment Program).
type DepDevice struct {
	ent.Schema
}

// Annotations of the DepDevice.
func (DepDevice) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "dep_devices"},
	}
}

// Fields of the DepDevice.
func (DepDevice) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().NotEmpty(),
		field.String("serial_number").NotEmpty(),
		field.String("dep_name").NotEmpty(),
		field.String("model").Optional(),
		field.String("description").Optional(),
		field.String("color").Optional(),
		field.String("asset_tag").Optional(),
		field.String("os").Optional(),
		field.String("device_family").Optional(),

		// Profile tracking
		field.String("profile_uuid").Optional(),
		field.String("profile_status").Optional().
			Comment("empty | assigned | pushed | removed"),
		field.Time("profile_assign_time").Optional().Nillable(),
		field.Time("profile_push_time").Optional().Nillable(),

		// Assignment info
		field.String("device_assigned_by").Optional(),
		field.Time("device_assigned_date").Optional().Nillable(),

		// Sync metadata
		field.String("op_type").Optional().
			Comment("added | modified | deleted"),
		field.Time("op_date").Optional().Nillable(),

		// Reassign status
		field.Bool("needs_manual_reassign").Default(false).
			Comment("true when Apple returns NOT_ACCESSIBLE"),
		field.String("reassign_error").Optional().
			Comment("Reason for reassign failure"),

		field.Bool("is_active").Default(true),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the DepDevice.
func (DepDevice) Edges() []ent.Edge {
	return nil
}

// Indexes of the DepDevice.
func (DepDevice) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("serial_number").Unique(),
		index.Fields("dep_name"),
		index.Fields("profile_uuid"),
		index.Fields("needs_manual_reassign"),
		index.Fields("op_type"),
	}
}
