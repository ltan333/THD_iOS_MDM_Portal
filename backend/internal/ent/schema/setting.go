package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type Setting struct {
	ent.Schema
}

func (Setting) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "settings"},
	}
}

func (Setting) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id"),
		field.String("key").NotEmpty().Unique(),
		field.Text("value").NotEmpty(), // JSON or simple string
		field.String("description").Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Setting) Edges() []ent.Edge {
	return nil
}
