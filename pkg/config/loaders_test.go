package config

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

func testValues() (Config, map[string]string) {
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
		envDatabasePort:     "5432",
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
	expected, env := testValues()
	setupEnv(t, env)

	var cfg Config
	if err := LoadConfigFromEnv(&cfg); err != nil {
		t.Fatalf("LoadConfigFromEnv returned error: %v", err)
	}

	if cfg.IsDebug != expected.IsDebug {
		t.Errorf("IsDebug mismatch: got %v want %v", cfg.IsDebug, expected.IsDebug)
	}
	if cfg.ProjectName != expected.ProjectName {
		t.Errorf("ProjectName mismatch: got %q want %q", cfg.ProjectName, expected.ProjectName)
	}

	if cfg.Runtime.APILocale != expected.Runtime.APILocale {
		t.Errorf(
			"APILocale mismatch: got %q want %q",
			cfg.Runtime.APILocale, expected.Runtime.APILocale,
		)
	}

	if !reflect.DeepEqual(cfg.Database, expected.Database) {
		t.Errorf("Database mismatch:\n got  %+v\n want %+v", cfg.Database, expected.Database)
	}
}
