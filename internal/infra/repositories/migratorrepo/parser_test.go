package migratorrepo

import (
	"testing"
)

func TestParseFilename(t *testing.T) {
	tests := []struct {
		name        string
		wantVersion string
		wantDesc    string
	}{
		{"20251204025235_create_user.sql", "20251204025235", "create_user"},
		{"20251204025518_add_auth_token.sql", "20251204025518", "add_auth_token"},
		{
			"20260311050148_add_user_profile_and_participant_ready.sql",
			"20260311050148",
			"add_user_profile_and_participant_ready",
		},
		{"20260118173823_initial_app_schema.sql", "20260118173823", "initial_app_schema"},
		{"nodescription.sql", "nodescription", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVersion, gotDesc := parseFilename(tt.name)
			if gotVersion != tt.wantVersion {
				t.Errorf("version = %q, want %q", gotVersion, tt.wantVersion)
			}
			if gotDesc != tt.wantDesc {
				t.Errorf("description = %q, want %q", gotDesc, tt.wantDesc)
			}
		})
	}
}

func TestSplitStatements(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
	}{
		{"empty", "", 0},
		{"single statement", "CREATE TABLE foo (id INT)", 1},
		{"two statements", "CREATE TABLE foo (id INT);\nCREATE TABLE bar (id INT)", 2},
		{"blank lines between", "CREATE TABLE foo (id INT);\n\n\nCREATE TABLE bar (id INT)", 2},
		{"comment only lines", "-- comment\nCREATE TABLE foo (id INT)", 1},
		{"trailing semicolon", "CREATE TABLE foo (id INT);", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitStatements(tt.content)
			if len(got) != tt.want {
				t.Errorf("len = %d, want %d (stmts: %v)", len(got), tt.want, got)
			}
		})
	}
}
