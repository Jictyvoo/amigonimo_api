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

// UserProfile holds the schema definition for user profile data.
type UserProfile struct{ ent.Schema }

// Fields of the UserProfile.
func (UserProfile) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("user_id", uuid.UUID{}),
		field.String("fullname").MaxLen(255).Optional(),
		field.String("nickname").MaxLen(120).Optional(),
		field.String("image_link").MaxLen(2048).Optional(),
		field.Time("birthday").Optional(),
		field.String("address").MaxLen(255).Optional(),
	}
}

// Indexes of the UserProfile.
func (UserProfile) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").Unique(),
		index.Fields("user_id").Unique(),
		index.Fields("nickname"),
		index.Fields("birthday"),
	}
}

// Edges of the UserProfile.
func (UserProfile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("profile").
			Field("user_id").
			Required().
			Unique().
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

// Mixin of the UserProfile.
func (UserProfile) Mixin() []ent.Mixin {
	return []ent.Mixin{
		customixins.UUIDMixin{},
		customixins.TimestampsMixin{},
	}
}
