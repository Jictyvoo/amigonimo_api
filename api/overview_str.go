package api

import _ "embed"

//go:embed OVERVIEW_README.md
var overviewReadme string

func Description() string {
	return overviewReadme
}
