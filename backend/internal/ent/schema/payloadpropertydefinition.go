package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type PayloadPropertyDefinition struct {
	ent.Schema
}

func (PayloadPropertyDefinition) Fields() []ent.Field {
	return []ent.Field{
		field.String("payload_type").NotEmpty(),
		field.String("key").NotEmpty(),
		field.String("value_type").NotEmpty(),
		field.JSON("default_value", map[string]interface{}{}).Optional(),
		field.JSON("enum_values", []interface{}{}).Optional(),
		field.Bool("deprecated").Default(false),
		field.String("description").Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

func (PayloadPropertyDefinition) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("properties", PayloadProperty.Type),
	}
}

func (PayloadPropertyDefinition) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("payload_type", "key").Unique(),
	}
}
