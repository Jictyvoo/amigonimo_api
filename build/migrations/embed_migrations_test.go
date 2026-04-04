package migrations_test

import (
	"path/filepath"
	"slices"
	"sort"
	"testing"

	"github.com/jictyvoo/amigonimo_api/build/migrations"
)

func TestVersionedMigrationsFSMatchesDisk(t *testing.T) {
	// Collect .sql files from the real directory (test runs in the package dir).
	diskFiles, err := filepath.Glob("*.sql")
	if err != nil {
		t.Fatalf("glob disk files: %v", err)
	}
	sort.Strings(diskFiles)

	// Collect .sql files from the embedded FS.
	fsys := migrations.VersionedMigrationsFS()
	entries, err := fsys.ReadDir(".")
	if err != nil {
		t.Fatalf("ReadDir embedded FS: %v", err)
	}

	embeddedFiles := make([]string, 0, len(entries))
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".sql" {
			t.Fatal("Only .sql files should be inside the embed filesystem")
		}
		embeddedFiles = append(embeddedFiles, e.Name())
	}
	slices.Sort(embeddedFiles)

	if len(diskFiles) != len(embeddedFiles) {
		t.Fatalf(
			"file count mismatch: disk=%d embedded=%d\n  disk:     %v\n  embedded: %v",
			len(diskFiles), len(embeddedFiles), diskFiles, embeddedFiles,
		)
	}

	for i, name := range diskFiles {
		if name != embeddedFiles[i] {
			t.Errorf("file[%d] mismatch: disk=%q embedded=%q", i, name, embeddedFiles[i])
		}
	}
}
