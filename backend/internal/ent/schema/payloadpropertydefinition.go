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
		field.String("payload_variant").Default(""),
		field.String("key").NotEmpty(),
		field.String("value_type").NotEmpty(),
		field.String("items_type").Optional().Nillable(),
		field.JSON("default_value", map[string]interface{}{}).Optional(),
		field.JSON("enum_values", []interface{}{}).Optional(),
		field.String("title").Optional().Nillable(),
		field.Text("description").Optional(),
		field.String("presence").Default("optional"),
		field.Bool("deprecated").Default(false),
		field.Bool("is_nested").Default(false),
		field.String("nested_reference").Optional().Nillable(),
		field.String("items_reference").Optional().Nillable(),
		field.JSON("supported_os", map[string]interface{}{}).Optional(),
		field.JSON("conditions", map[string]interface{}{}).Optional(),
		field.Int("order_index").Default(0),
		field.String("yaml_source_file").Optional().Nillable(),
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
		index.Fields("payload_type", "payload_variant", "key").Unique(),
		index.Fields("payload_type", "payload_variant", "order_index"),
		index.Fields("payload_type", "payload_variant", "is_nested"),
		index.Fields("payload_type"),
	}
}
