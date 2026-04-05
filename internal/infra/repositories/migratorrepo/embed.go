package migratorrepo

import _ "embed"

//go:embed schema.sql
var createTableSQL string
