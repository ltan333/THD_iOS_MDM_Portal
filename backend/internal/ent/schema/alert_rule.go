package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// AlertRule holds the schema definition for the AlertRule entity.
type AlertRule struct {
	ent.Schema
}

// Annotations of the AlertRule.
func (AlertRule) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "alert_rules"},
	}
}

// Fields of the AlertRule.
func (AlertRule) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id"),
		field.String("name").NotEmpty().MaxLen(255),
		field.String("description").Optional().MaxLen(500),
		field.JSON("condition", map[string]interface{}{}).Optional(), // e.g., {"type": "device_offline", "threshold": "24h"}
		field.JSON("actions", map[string]interface{}{}).Optional(),   // e.g., {"trigger_alert": true, "send_email": true}
		field.Bool("enabled").Default(true),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the AlertRule.
func (AlertRule) Edges() []ent.Edge {
	return []ent.Edge{}
}
