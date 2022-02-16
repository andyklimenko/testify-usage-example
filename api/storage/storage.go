package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

type dbExecutor func(tx *sqlx.Tx) error

func runInTx(db *sqlx.DB, executor dbExecutor) error {
	tx, err := db.BeginTxx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	if err := executor(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("%s %w", rollbackErr.Error(), err)
		}
		return err
	}
	return tx.Commit()
}
