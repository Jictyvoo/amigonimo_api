package entschema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/index"

	"github.com/jictyvoo/amigonimo_api/build/entschema/customixins"
	"github.com/jictyvoo/amigonimo_api/build/entschema/uuidfield"
)

// DrawResult holds the schema definition for the DrawResult entity.
type DrawResult struct{ ent.Schema }

// Fields of the DrawResult.
func (DrawResult) Fields() []ent.Field {
	return []ent.Field{
		uuidfield.Field("secret_friend_id"),
		uuidfield.Field("giver_participant_id"),
		uuidfield.Field("receiver_participant_id"),
	}
}

// Indexes of the DrawResult.
func (DrawResult) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").Unique(),
		index.Fields("secret_friend_id"),
		index.Fields("giver_participant_id"),
		index.Fields("receiver_participant_id"),
		index.Fields("secret_friend_id", "giver_participant_id").Unique(),
		index.Fields("secret_friend_id", "receiver_participant_id").Unique(),
	}
}

// Edges of the DrawResult.
func (DrawResult) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("secret_friend", SecretFriend.Type).
			Ref("draw_results").
			Field("secret_friend_id").
			Required().
			Unique().
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.From("giver", Participant.Type).
			Ref("given_results").
			Field("giver_participant_id").
			Required().
			Unique().
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.From("receiver", Participant.Type).
			Ref("received_results").
			Field("receiver_participant_id").
			Required().
			Unique().
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

// Mixin of the DrawResult.
func (DrawResult) Mixin() []ent.Mixin {
	return []ent.Mixin{
		customixins.UUIDMixin{},
		customixins.TimestampsMixin{},
	}
}
