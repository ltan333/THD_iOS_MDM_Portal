package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Payload struct {
	ent.Schema
}

func (Payload) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id"),
		field.String("payload_description").Optional(),
		field.String("payload_display_name").MaxLen(255).NotEmpty(),
		field.String("payload_identifier").Unique().NotEmpty(),
		field.String("payload_organization").MaxLen(255).Optional(),
		field.String("payload_type").NotEmpty(),
		field.String("payload_uuid").NotEmpty(),
		field.Int("payload_version").Default(1),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

func (Payload) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("mobile_config", MobileConfig.Type).Ref("payloads").Unique().Required(),
		edge.To("properties", PayloadProperty.Type),
	}
}
