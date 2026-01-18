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

// AuthToken holds the schema definition for the AuthToken entity.
type AuthToken struct{ ent.Schema }

// Fields of the AuthToken.
func (AuthToken) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("user_id", uuid.UUID{}).Unique(),
		field.String("token").MaxLen(52),
		field.UUID("refresh_token", uuid.UUID{}).
			// This is needed because SQLc doesn't recognize the UUID type on MariaDB dialect
			SchemaType(customixins.SchemaTypeUUID()).
			Optional(),
		field.Time("expires_at"),
	}
}

// Indexes of the AuthToken.
func (AuthToken) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").Unique(),
		index.Fields("user_id"),
		index.Fields("token"),
		index.Fields("refresh_token"),
	}
}

// Edges of the AuthToken.
func (AuthToken) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("auth_token").
			Field("user_id").
			Required().
			Unique().
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

// Mixin of the AuthToken.
func (AuthToken) Mixin() []ent.Mixin {
	return []ent.Mixin{
		customixins.UUIDMixin{},
		customixins.TimestampsMixin{},
	}
}
