package config

import (
	"maps"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func testEnvValues() (Config, map[string]string) {
	expected := Config{
		IsDebug:     true,
		ProjectName: "my_project",
		Runtime: Runtime{
			APILocale:     "pt-br",
			Host:          "localH0$T.0.0.1",
			Port:          5648,
			AuthSecretKey: "my_secret_key",
		},
		Database: Database{
			Host:     "localhost",
			Port:     5432,
			User:     "admin",
			Password: "secret",
			Database: "appdb",
			Timeout:  time.Second << 2,
		},
	}

	env := map[string]string{
		envUseDebug:         "true",
		envProjectName:      expected.ProjectName,
		envAPILocale:        expected.Runtime.APILocale,
		envAPIHost:          expected.Runtime.Host,
		envAPIPort:          strconv.Itoa(int(expected.Runtime.Port)),
		envAuthSecretKey:    expected.Runtime.AuthSecretKey,
		envDatabaseHost:     expected.Database.Host,
		envDatabasePort:     strconv.Itoa(int(expected.Database.Port)),
		envDatabaseUser:     expected.Database.User,
		envDatabasePassword: expected.Database.Password,
		envDatabaseName:     expected.Database.Database,
		envDatabaseTimeout:  expected.Database.Timeout.String(),
	}

	return expected, env
}

func setupEnv(t *testing.T, env map[string]string) {
	t.Helper()
	for k, v := range env {
		t.Setenv(k, v)
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	allExpected, allEnv := testEnvValues()

	tests := []struct {
		name string
		env  map[string]string
		want Config
	}{
		{
			name: "all env vars populated",
			env:  allEnv,
			want: allExpected,
		},
		{
			name: "only database env vars populated",
			env: func() map[string]string {
				dbOnly := maps.Clone(allEnv)
				for key := range dbOnly {
					if !strings.HasPrefix(key, "DATABASE") {
						delete(dbOnly, key)
					}
				}
				return dbOnly
			}(),
			want: Config{Database: allExpected.Database},
		},
		{
			name: "no env vars leaves config unchanged",
			env:  map[string]string{},
			want: Config{},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				setupEnv(t, tt.env)
				var cfg Config
				if err := LoadConfigFromEnv(&cfg); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(cfg, tt.want) {
					t.Errorf("config mismatch:\n got  %+v\n want %+v", cfg, tt.want)
				}
			},
		)
	}
}

func TestLoad(t *testing.T) {
	content, want := testConfigFixture(t)

	tests := []struct {
		name    string
		setup   func(t *testing.T) string
		wantErr bool
		checkFn func(t *testing.T, cfg Config)
	}{
		{
			name: "loads from explicit filename",
			setup: func(t *testing.T) string {
				return writeTempConfig(t, content)
			},
			checkFn: func(t *testing.T, cfg Config) {
				if !reflect.DeepEqual(cfg, want) {
					t.Errorf("config mismatch:\n got  %+v\n want %+v", cfg, want)
				}
			},
		},
		{
			name: "CONF_FILE env overrides filename argument",
			setup: func(t *testing.T) string {
				realPath := writeTempConfig(t, content)
				t.Setenv(envConfFile, realPath)
				return filepath.Join(t.TempDir(), "should_not_be_used.toml")
			},
			checkFn: func(t *testing.T, cfg Config) {
				if cfg.ProjectName != want.ProjectName {
					t.Errorf("ProjectName: got %q want %q", cfg.ProjectName, want.ProjectName)
				}
			},
		},
		{
			name: "missing file is not an error",
			setup: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "absent.toml")
			},
		},
		{
			name: "env var overlays file value",
			setup: func(t *testing.T) string {
				t.Setenv(envProjectName, "overridden_project")
				return writeTempConfig(t, content)
			},
			checkFn: func(t *testing.T, cfg Config) {
				if cfg.ProjectName != "overridden_project" {
					t.Errorf("env overlay: got %q want %q", cfg.ProjectName, "overridden_project")
				}
				if cfg.Runtime.Host != want.Runtime.Host {
					t.Errorf("file value lost: got %q want %q", cfg.Runtime.Host, want.Runtime.Host)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				filename := tt.setup(t)
				cfg, err := Load(filename)
				if (err != nil) != tt.wantErr {
					t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
				}
				if tt.checkFn != nil {
					tt.checkFn(t, cfg)
				}
			},
		)
	}
}
