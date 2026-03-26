package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ProfileDeploymentStatus holds the schema definition for the ProfileDeploymentStatus entity.
type ProfileDeploymentStatus struct {
	ent.Schema
}

// Annotations of the ProfileDeploymentStatus.
func (ProfileDeploymentStatus) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "profile_deployment_statuses"},
	}
}

// Fields of the ProfileDeploymentStatus.
func (ProfileDeploymentStatus) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id"),
		field.Uint("profile_id"),
		field.String("device_id").NotEmpty(),
		field.Enum("status").Values("pending", "success", "failed").Default("pending"),
		field.String("error_message").Optional(),
		field.Time("applied_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the ProfileDeploymentStatus.
func (ProfileDeploymentStatus) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("profile", Profile.Type).
			Ref("deployment_statuses").
			Unique().
			Required().
			Field("profile_id"),
		edge.To("device", Device.Type).
			Unique().
			Required().
			Field("device_id"),
	}
}
