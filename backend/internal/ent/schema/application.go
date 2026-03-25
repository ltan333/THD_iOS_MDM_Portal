package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Application struct {
	ent.Schema
}

func (Application) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "applications"},
	}
}

func (Application) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id"),
		field.String("name").NotEmpty().MaxLen(255),
		field.String("bundle_id").NotEmpty().Unique().MaxLen(255),
		field.Enum("platform").Values("ios", "android", "windows", "macos"),
		field.Enum("type").Values("app_store", "enterprise", "web_clip").Default("app_store"),
		field.String("description").Optional(),
		field.String("icon_url").Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Application) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("versions", AppVersion.Type),
	}
}
