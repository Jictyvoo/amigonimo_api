package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/google/uuid"

	"github.com/jictyvoo/amigonimo_api/build/entschema/customixins"
)

// SecretFriend holds the schema definition for the SecretFriend entity.
type SecretFriend struct{ ent.Schema }

// Fields of the SecretFriend.
func (SecretFriend) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Time("datetime"),
		field.String("location").Optional(),
		field.UUID("owner_id", uuid.UUID{}),
		field.Uint8("max_deny_list_size").Default(0),
		field.String("invite_code"),
		field.String("invite_link").Optional(),
		field.String("status"),
	}
}

// Indexes of the SecretFriend.
func (SecretFriend) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").Unique(),
		index.Fields("invite_code").Unique(),
		index.Fields("owner_id"),
	}
}

// Edges of the SecretFriend.
func (SecretFriend) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("participants", Participant.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("draw_results", DrawResult.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.From("owner", User.Type).
			Ref("secret_friends").
			Field("owner_id").
			Unique().
			Required(),
	}
}

// Mixin of the SecretFriend.
func (SecretFriend) Mixin() []ent.Mixin {
	return []ent.Mixin{
		customixins.UUIDMixin{},
		customixins.TimestampsMixin{},
	}
}
