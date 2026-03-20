package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// APNSConfig holds the schema definition for the APNSConfig entity.
type APNSConfig struct {
	ent.Schema
}

// Fields of the APNSConfig.
func (APNSConfig) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id"),
		field.String("topic").Unique().NotEmpty(),
		field.String("cert_file_path").NotEmpty(),
		field.String("key_file_path").NotEmpty(),
		field.Time("expiry").Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the APNSConfig.
func (APNSConfig) Edges() []ent.Edge {
	return nil
}
