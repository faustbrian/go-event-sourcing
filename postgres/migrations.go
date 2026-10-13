package postgres

import (
	"embed"
	"io/fs"
)

// migrationFiles contains hand-maintained forward-only schema migrations.
// The SQL files intentionally have no Down directives; destructive reversal is
// not supplied by this package and embedded bytes are the owned producer input.
//
//go:embed migrations/*.sql
var migrationFiles embed.FS

// Migrations returns versioned, engine-neutral SQL migration files rooted at
// their filenames. The files use the migrations package directive format but
// do not require or expose a migration engine.
func Migrations() fs.FS {
	files, _ := fs.Sub(migrationFiles, "migrations")

	return files
}
