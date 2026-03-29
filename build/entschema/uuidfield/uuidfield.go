package uuidfield

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

func Field(name string) ent.Field {
	return field.UUID(name, uuid.UUID{}).
		SchemaType(schemaType())
}

func Unique(name string) ent.Field {
	return field.UUID(name, uuid.UUID{}).
		Unique().
		SchemaType(schemaType())
}

func Optional(name string) ent.Field {
	return field.UUID(name, uuid.UUID{}).
		Optional().
		SchemaType(schemaType())
}

func ID(defaultExpr *entsql.Annotation) ent.Field {
	return field.UUID("id", uuid.UUID{}).
		Unique().
		Immutable().
		Default(uuid.New).
		Annotations(defaultExpr).
		SchemaType(schemaType())
}

func schemaType() map[string]string {
	return map[string]string{
		dialect.MySQL:    "BINARY(16)",
		dialect.Postgres: "uuid",
		dialect.SQLite:   "blob",
	}
}
