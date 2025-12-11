package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/jictyvoo/amigonimo_api/build/entschema/customixins"
)

// User holds the schema definition for the User entity.
type User struct{ ent.Schema }

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("fullname").MaxLen(255),
		field.String("email").MaxLen(255),
		field.String("username").Unique().MaxLen(78),
		field.String("password").MaxLen(76),
		field.Time("verified_at").Optional(),
		field.String("remember_token").Optional().MaxLen(52),
		field.String("verification_code").Optional().MaxLen(116),
		field.String("recovery_code").Optional().MaxLen(27),
		field.Time("recovery_code_expires_at").Optional(),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").Unique(),
		index.Fields("email").Unique(),
		index.Fields("username").Unique(),
		index.Fields("verification_code"),
		index.Fields("recovery_code"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("secret_friends", SecretFriend.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("participants", Participant.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("denied_entries", Denylist.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("auth_token", AuthToken.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

// Mixin of the User.
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		customixins.UUIDMixin{},
		customixins.TimestampsMixin{},
	}
}
