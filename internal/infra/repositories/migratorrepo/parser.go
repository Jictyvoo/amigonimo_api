package migratorrepo

import (
	"strings"
)

// parseFilename splits a migration filename into its version (timestamp) and
// description components. For example:
//
// "20251204025235_create_user.sql" → ("20251204025235", "create_user")
func parseFilename(name string) (version, description string) {
	base := strings.TrimSuffix(name, ".sql")
	before, after, ok := strings.Cut(base, "_")
	if !ok {
		return base, ""
	}
	return before, after
}

// splitStatements splits raw SQL content into individual executable statements,
// stripping standalone comment lines and skipping blank entries.
func splitStatements(content string) []string {
	var stmts []string
	for raw := range strings.SplitSeq(content, ";") {
		stmt := stripCommentLines(raw)
		if stmt == "" {
			continue
		}
		stmts = append(stmts, stmt)
	}
	return stmts
}

// stripCommentLines removes full-line SQL comments (lines starting with --)
// and blank lines from a statement fragment, then trims surrounding whitespace.
func stripCommentLines(raw string) string {
	var lines []string
	for line := range strings.SplitSeq(raw, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "--") {
			lines = append(lines, line)
		}
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}
