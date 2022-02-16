package migrations

import (
	migrate "github.com/rubenv/sql-migrate"
)

var migrations = &migrate.MemoryMigrationSource{
	Migrations: []*migrate.Migration{
		{
			Id: "01-initial",
			Up: []string{
				`CREATE EXTENSION "uuid-ossp";
				CREATE TABLE users(
					id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
					first_name VARCHAR(128) NOT NULL,
					last_name VARCHAR(128) NOT NULL,
					created_at timestamp NOT NULL DEFAULT now()
				);`,
			},
			Down: []string{
				`DELETE TABLE users;
				DROP EXTENSION "uuid-ossp";`,
			},
		},
	},
}
