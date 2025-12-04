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
		field.String("fullname"),
		field.String("email"),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").Unique(),
		index.Fields("email").Unique(),
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
	}
}

// Mixin of the User.
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		customixins.UUIDMixin{},
		customixins.TimestampsMixin{},
	}
}
