package entschema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/jictyvoo/amigonimo_api/build/entschema/customixins"
	"github.com/jictyvoo/amigonimo_api/build/entschema/uuidfield"
)

// Participant holds the schema definition for the Participant entity.
type Participant struct{ ent.Schema }

// Fields of the Participant.
func (Participant) Fields() []ent.Field {
	return []ent.Field{
		uuidfield.Field("user_id"),
		uuidfield.Field("secret_friend_id"),
		field.Time("joined_at").Default(time.Now),
		field.Bool("is_ready").Default(false),
	}
}

// Indexes of the Participant.
func (Participant) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").Unique(),
		index.Fields("user_id"),
		index.Fields("secret_friend_id"),
		index.Fields("user_id", "secret_friend_id").Unique(),
	}
}

// Edges of the Participant.
func (Participant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("wishlist_items", WishlistItem.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("denylist", Denylist.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("given_results", DrawResult.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("received_results", DrawResult.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.From("user", User.Type).
			Ref("participants").
			Field("user_id").
			Required().
			Unique(),
		edge.From("secret_friend", SecretFriend.Type).
			Ref("participants").
			Field("secret_friend_id").
			Required().
			Unique(),
	}
}

// Mixin of the Participant.
func (Participant) Mixin() []ent.Mixin {
	return []ent.Mixin{
		customixins.UUIDMixin{},
		customixins.TimestampsMixin{},
	}
}
