package migrations

import (
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
)

func Up(db *sqlx.DB, driver string) error {
	_, err := migrate.Exec(db.DB, driver, migrations, migrate.Up)
	return err
}
