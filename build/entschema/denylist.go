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

// Denylist holds the schema definition for the Denylist entity.
type Denylist struct{ ent.Schema }

// Fields of the Denylist.
func (Denylist) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("participant_id", uuid.UUID{}),
		field.UUID("denied_user_id", uuid.UUID{}),
	}
}

// Indexes of the Denylist.
func (Denylist) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").Unique(),
		index.Fields("participant_id"),
		index.Fields("denied_user_id"),
		index.Fields("participant_id", "denied_user_id").Unique(),
	}
}

// Edges of the Denylist.
func (Denylist) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("participant", Participant.Type).
			Ref("denylist").
			Field("participant_id").
			Required().
			Unique().
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.From("denied_user", User.Type).
			Ref("denied_entries").
			Field("denied_user_id").
			Required().
			Unique().
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

// Mixin of the Denylist.
func (Denylist) Mixin() []ent.Mixin {
	return []ent.Mixin{
		customixins.UUIDMixin{},
		customixins.TimestampsMixin{},
	}
}
