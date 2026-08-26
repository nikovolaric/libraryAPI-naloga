package devops

import "embed"

// Migrations holds the embedded SQL migration files so the app can run them on
// startup without the files on disk or the migrate CLI installed.
//
//go:embed migrations/*.sql
var Migrations embed.FS
