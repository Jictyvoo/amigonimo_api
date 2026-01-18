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

// WishlistItem holds the schema definition for the WishlistItem entity.
type WishlistItem struct{ ent.Schema }

// Fields of the WishlistItem.
func (WishlistItem) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("participant_id", uuid.UUID{}),
		field.String("label"),
		field.Text("comments").
			Optional().
			Nillable(),
	}
}

// Indexes of the WishlistItem.
func (WishlistItem) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").Unique(),
		index.Fields("participant_id"),
	}
}

// Edges of the WishlistItem.
func (WishlistItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("participant", Participant.Type).
			Ref("wishlist_items").
			Field("participant_id").
			Required().
			Unique().
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

// Mixin of the WishlistItem.
func (WishlistItem) Mixin() []ent.Mixin {
	return []ent.Mixin{
		customixins.UUIDMixin{},
		customixins.TimestampsMixin{},
	}
}
