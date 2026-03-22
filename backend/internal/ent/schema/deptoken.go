package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// DEPToken holds the schema definition for the DEPToken entity.
type DEPToken struct {
	ent.Schema
}

// Annotations of the DEPToken.
func (DEPToken) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "dep_names"},
	}
}

// Fields of the DEPToken.
func (DEPToken) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("name").Unique().NotEmpty(),
		field.Text("consumer_key").Optional(),
		field.Text("consumer_secret").Optional(),
		field.Text("access_token").Optional(),
		field.Text("access_secret").Optional(),
		field.Time("access_token_expiry").Optional(),
		field.String("config_base_url").Optional(),
		field.Text("tokenpki_cert_pem").Optional(),
		field.Text("tokenpki_key_pem").Optional(),
		field.String("syncer_cursor").Optional(),
		field.Text("assigner_profile_uuid").Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the DEPToken.
func (DEPToken) Edges() []ent.Edge {
	return nil
}
