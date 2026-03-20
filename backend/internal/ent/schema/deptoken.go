package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// DEPToken holds the schema definition for the DEPToken entity.
type DEPToken struct {
	ent.Schema
}

// Fields of the DEPToken.
func (DEPToken) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id"),
		field.String("name").Unique().NotEmpty(),
		field.String("p7m_file_path").NotEmpty(),
		field.Time("expiry").Optional(),
		field.Time("last_used").Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the DEPToken.
func (DEPToken) Edges() []ent.Edge {
	return nil
}
