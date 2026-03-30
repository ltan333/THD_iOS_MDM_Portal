package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Alert holds the schema definition for the Alert entity.
type Alert struct {
	ent.Schema
}

// Annotations of the Alert.
func (Alert) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "alerts"},
	}
}

// Fields of the Alert.
func (Alert) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id"),
		field.Enum("severity").Values("critical", "high", "medium", "low").Default("medium"),
		field.String("title").NotEmpty().MaxLen(255),
		field.Enum("type").Values("security", "compliance", "connectivity", "application", "device_health").Default("security"),
		field.Enum("status").Values("open", "acknowledged", "resolved").Default("open"),
		field.String("device_id").Optional(),
		field.Uint("user_id").Optional().Nillable(),
		field.JSON("details", map[string]any{}).Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("acknowledged_at").Optional().Nillable(),
		field.Time("resolved_at").Optional().Nillable(),
	}
}

// Edges of the Alert.
func (Alert) Edges() []ent.Edge {
	return []ent.Edge{}
}
