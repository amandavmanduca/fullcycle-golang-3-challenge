package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func MigrateUp(db *sql.DB, migrationFolder string) error {
	if err := goose.SetDialect("mysql"); err != nil {
		panic(err)
	}
	return goose.Up(db, migrationFolder)
}
