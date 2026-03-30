package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type AppDeployment struct {
	ent.Schema
}

func (AppDeployment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "app_deployments"},
	}
}

func (AppDeployment) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id"),
		field.Uint("app_version_id"),
		field.Enum("target_type").Values("device", "group", "user"),
		field.String("target_id").NotEmpty(),
		field.Enum("status").Values("pending", "installing", "success", "failed").Default("pending"),
		field.String("error_message").Optional(),
		field.Time("installed_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (AppDeployment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("version", AppVersion.Type).
			Ref("deployments").
			Field("app_version_id").
			Unique().
			Required(),
	}
}
