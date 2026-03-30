package schema

import (
	"encoding/json"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type PayloadProperty struct {
	ent.Schema
}

func (PayloadProperty) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id"),
		field.JSON("value_json", json.RawMessage{}).Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

func (PayloadProperty) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("payload", Payload.Type).Ref("properties").Unique().Required(),
		edge.From("definition", PayloadPropertyDefinition.Type).Ref("properties").Unique().Required(),
	}
}
