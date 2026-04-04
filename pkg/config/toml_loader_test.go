package config

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
)

func writeTempConfig(t *testing.T, content []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "conf.toml")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("writeTempConfig: %v", err)
	}
	return path
}

// testConfigFixture marshals the canonical Config from testEnvValues to TOML,
// returning both the raw bytes and the expected struct simultaneously.
func testConfigFixture(t *testing.T) ([]byte, Config) {
	t.Helper()
	expected, _ := testEnvValues()
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(expected); err != nil {
		t.Fatalf("testConfigFixture: marshal: %v", err)
	}
	return buf.Bytes(), expected
}

func TestLoadTOML(t *testing.T) {
	content, want := testConfigFixture(t)

	tests := []struct {
		name      string
		setupPath func(t *testing.T) string
		wantErr   bool
		errCheck  func(t *testing.T, err error)
		checkFn   func(t *testing.T, cfg Config)
	}{
		{
			name: "loads all values from file",
			setupPath: func(t *testing.T) string {
				return writeTempConfig(t, content)
			},
			checkFn: func(t *testing.T, cfg Config) {
				if !reflect.DeepEqual(cfg, want) {
					t.Errorf("config mismatch:\n got  %+v\n want %+v", cfg, want)
				}
			},
		},
		{
			name: "file not found returns wrapped not-exist error",
			setupPath: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "missing.toml")
			},
			wantErr: true,
			errCheck: func(t *testing.T, err error) {
				if !os.IsNotExist(errors.Unwrap(err)) {
					t.Errorf("expected wrapped os.IsNotExist, got: %v", err)
				}
			},
		},
		{
			name: "starts from DefaultConfig before overlaying file",
			setupPath: func(t *testing.T) string {
				// Write a partial TOML that only overrides project_name.
				partial := []byte(`project_name = "partial_project"`)
				return writeTempConfig(t, partial)
			},
			checkFn: func(t *testing.T, cfg Config) {
				if cfg.ProjectName != "partial_project" {
					t.Errorf("ProjectName: got %q want %q", cfg.ProjectName, "partial_project")
				}
				def := DefaultConfig()
				if cfg.Database.Port != def.Database.Port {
					t.Errorf(
						"Database.Port: got %d want %d (default)",
						cfg.Database.Port,
						def.Database.Port,
					)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				path := tt.setupPath(t)
				cfg, err := LoadTOML(path)
				if (err != nil) != tt.wantErr {
					t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
				}
				if tt.errCheck != nil {
					tt.errCheck(t, err)
				}
				if tt.checkFn != nil {
					tt.checkFn(t, cfg)
				}
			},
		)
	}
}
