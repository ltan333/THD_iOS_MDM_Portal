package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// APNSConfig holds the schema definition for the APNSConfig entity.
type APNSConfig struct {
	ent.Schema
}

// Annotations of the APNSConfig.
func (APNSConfig) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "portal_push_certs"},
	}
}

// Fields of the APNSConfig.
func (APNSConfig) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("topic").Unique().NotEmpty(),
		field.Text("cert_pem").StorageKey("cert_pem").NotEmpty(),
		field.Text("key_pem").StorageKey("key_pem").NotEmpty(),
		field.Int("stale_token").Default(0),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the APNSConfig.
func (APNSConfig) Edges() []ent.Edge {
	return nil
}
