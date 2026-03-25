package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ProfileVersion holds the schema definition for the ProfileVersion entity.
type ProfileVersion struct {
	ent.Schema
}

// Annotations of the ProfileVersion.
func (ProfileVersion) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "profile_versions"},
	}
}

// Fields of the ProfileVersion.
func (ProfileVersion) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id"),
		field.Uint("profile_id"),
		field.Int("version"),
		field.JSON("data", map[string]interface{}{}).Optional(), // Snapshot of profile data
		field.String("change_notes").Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the ProfileVersion.
func (ProfileVersion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("profile", Profile.Type).
			Ref("versions").
			Unique().
			Required().
			Field("profile_id"),
	}
}
