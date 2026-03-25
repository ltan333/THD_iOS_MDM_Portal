package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Profile holds the schema definition for the MDM Profile entity.
type Profile struct {
	ent.Schema
}

// Annotations of the Profile.
func (Profile) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "profiles"},
	}
}

// Fields of the Profile.
func (Profile) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id"),
		field.String("name").Unique().NotEmpty().MaxLen(255),
		field.Enum("platform").Values("ios", "android", "windows", "macos", "all").Default("all"),
		field.Enum("scope").Values("device", "user", "group").Default("device"),
		field.Enum("status").Values("active", "draft", "archived").Default("draft"),
		field.JSON("security_settings", map[string]interface{}{}).Optional(),
		field.JSON("network_config", map[string]interface{}{}).Optional(),
		field.JSON("restrictions", map[string]interface{}{}).Optional(),
		field.JSON("content_filter", map[string]interface{}{}).Optional(),
		field.JSON("compliance_rules", map[string]interface{}{}).Optional(),
		field.JSON("payloads", map[string]interface{}{}).Optional(),
		field.Int("version").Default(1),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Profile.
func (Profile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("device_groups", DeviceGroup.Type).
			Ref("profiles"),
		edge.To("assignments", ProfileAssignment.Type),
		edge.To("versions", ProfileVersion.Type),
		edge.To("deployment_statuses", ProfileDeploymentStatus.Type),
	}
}
