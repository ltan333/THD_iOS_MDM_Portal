package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type AppVersion struct {
	ent.Schema
}

func (AppVersion) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "app_versions"},
	}
}

func (AppVersion) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id"),
		field.Uint("application_id"),
		field.String("version").NotEmpty(),
		field.String("build_number").NotEmpty(),
		field.String("minimum_os_version").Optional(),
		field.String("file_url").Optional(), // For enterprise apps
		field.Int64("size").Optional(), // Binary size in bytes
		field.JSON("metadata", map[string]interface{}{}).Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (AppVersion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("application", Application.Type).
			Ref("versions").
			Field("application_id").
			Unique().
			Required(),
		edge.To("deployments", AppDeployment.Type),
	}
}
