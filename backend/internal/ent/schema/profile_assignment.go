package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ProfileAssignment holds the schema definition for the ProfileAssignment entity.
type ProfileAssignment struct {
	ent.Schema
}

// Annotations of the ProfileAssignment.
func (ProfileAssignment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "profile_assignments"},
	}
}

// Fields of the ProfileAssignment.
func (ProfileAssignment) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id"),
		field.Uint("profile_id"),
		field.Enum("target_type").Values("device", "group", "user"),
		field.String("target_id").NotEmpty(), // Can be device ID (string) or group/user ID (uint as string)
		field.Enum("schedule_type").Values("immediate", "scheduled").Default("immediate"),
		field.Time("scheduled_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the ProfileAssignment.
func (ProfileAssignment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("profile", Profile.Type).
			Ref("assignments").
			Unique().
			Required().
			Field("profile_id"),
	}
}
