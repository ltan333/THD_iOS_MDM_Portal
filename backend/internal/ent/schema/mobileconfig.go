package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type MobileConfig struct {
	ent.Schema
}

func (MobileConfig) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id"),
		field.String("name").Unique().MaxLen(255).NotEmpty(),
		field.String("payload_identifier").Unique().NotEmpty(),
		field.String("payload_type").NotEmpty(),
		field.String("payload_display_name").MaxLen(255).NotEmpty(),
		field.String("payload_description").Optional(),
		field.String("payload_organization").MaxLen(255).Optional(),
		field.String("payload_uuid").NotEmpty(),
		field.Int("payload_version").Default(1),
		field.Bool("payload_removal_disallowed").Default(false),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

func (MobileConfig) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("payloads", Payload.Type),
	}
}

func (MobileConfig) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("payload_identifier").Unique(),
	}
}
