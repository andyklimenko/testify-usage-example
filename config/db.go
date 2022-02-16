package config

import "errors"

var (
	ErrNoDbDSN = errors.New("no db dsn")
)

type DB struct {
	Driver string
	DSN    string
}

func (db *DB) Load(envPrefix string) error {
	v := setupViper(envPrefix)

	v.SetDefault("driver", "postgres")
	db.Driver = v.GetString("driver")

	db.DSN = v.GetString("dsn")
	if db.DSN == "" {
		return ErrNoDbDSN
	}

	return nil
}
