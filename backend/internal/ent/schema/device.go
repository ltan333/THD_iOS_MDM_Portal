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
		entsql.Annotation{Table: "devices"},
	}
}

// Fields of the Device.
func (Device) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().NotEmpty().SchemaType(map[string]string{
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
	}
}
