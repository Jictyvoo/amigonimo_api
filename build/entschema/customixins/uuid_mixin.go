package customixins

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// UUIDMixin provides a UUID primary key stored as BINARY(16).
type UUIDMixin struct {
	mixin.Schema
}

func (UUIDMixin) Fields() []ent.Field {
	// Default expression for MySQL: use UNHEX(REPLACE(UUID(), '-', '')) if needed
	defaultExpr := entsql.DefaultExprs(
		map[string]string{
			dialect.MySQL:    "UUID_TO_BIN(UUID())", // native function to generate BINARY(16)
			dialect.Postgres: "gen_random_uuid()",   // requires pgcrypto extension
			dialect.SQLite:   "",                    // handled in Go with uuid.New
		},
	)

	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Unique().
			Immutable().
			Default(uuid.New). // generate in Go for SQLite and fallback
			Annotations(defaultExpr).
			SchemaType(
				map[string]string{
					dialect.MySQL:    "BINARY(16)",
					dialect.Postgres: "uuid",
					dialect.SQLite:   "blob",
				},
			),
	}
}
