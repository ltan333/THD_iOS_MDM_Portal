package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// DeviceGroup holds the schema definition for the DeviceGroup entity.
type DeviceGroup struct {
	ent.Schema
}

// Annotations of the DeviceGroup.
func (DeviceGroup) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "device_groups"},
	}
}

// Fields of the DeviceGroup.
func (DeviceGroup) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id"),
		field.String("name").NotEmpty().MaxLen(100),
		field.String("description").Optional().MaxLen(500),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the DeviceGroup.
func (DeviceGroup) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("devices", Device.Type),
		edge.To("profiles", Profile.Type),
	}
}
