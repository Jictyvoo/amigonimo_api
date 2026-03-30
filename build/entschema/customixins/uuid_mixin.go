package customixins

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/mixin"

	"github.com/jictyvoo/amigonimo_api/build/entschema/uuidfield"
)

// UUIDMixin provides a UUID primary key stored as BINARY(16).
type UUIDMixin struct {
	mixin.Schema
}

func (UUIDMixin) Fields() []ent.Field {
	// MariaDB's UUID_v7() returns the textual UUID form, so convert it to 16 bytes for BINARY(16).
	defaultExpr := entsql.DefaultExprs(
		map[string]string{
			dialect.MySQL:    "UNHEX(REPLACE(UUID_v7(), '-', ''))",
			dialect.Postgres: "gen_random_uuid()", // requires pgcrypto extension
			dialect.SQLite:   "",                  // handled in Go with uuid.New
		},
	)

	return []ent.Field{
		uuidfield.ID(defaultExpr),
	}
}
