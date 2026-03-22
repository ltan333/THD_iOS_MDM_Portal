package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// DepProfile holds the schema definition for the DepProfile entity.
type DepProfile struct {
	ent.Schema
}

// Fields of the DepProfile.
func (DepProfile) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id"),
		field.String("profile_name").Unique().NotEmpty(),
		field.String("profile_uuid").Unique().Optional(),
		field.Bool("allow_pairing").Default(true),
		field.JSON("anchor_certs", []string{}).Optional(),
		field.Bool("auto_advance_setup").Default(false),
		field.Bool("await_device_configured").Default(false),
		field.String("configuration_web_url").Optional(),
		field.String("department").Optional(),
		field.JSON("devices", []string{}).Optional(),
		field.Bool("do_not_use_profile_from_backup").Default(false),
		field.Bool("is_return_to_service").Default(false),
		field.Bool("is_mandatory").Default(false),
		field.Bool("is_mdm_removable").Default(true),
		field.Bool("is_multi_user").Default(false),
		field.Bool("is_supervised").Default(false),
		field.String("language").Optional(),
		field.String("org_magic").Optional(),
		field.String("region").Optional(),
		field.JSON("skip_setup_items", []string{}).Optional(),
		field.JSON("supervising_host_certs", []string{}).Optional(),
		field.String("support_email_address").Optional(),
		field.String("support_phone_number").Optional(),
		field.String("url").Optional(),
		field.JSON("profile_data", map[string]interface{}{}).Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the DepProfile.
func (DepProfile) Edges() []ent.Edge {
	return nil
}
