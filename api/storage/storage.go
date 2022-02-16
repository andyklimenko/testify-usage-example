package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

type dbExecutor func(ctx context.Context, tx *sqlx.Tx) error

func runInTx(db *sqlx.DB, executor dbExecutor, isoLevel sql.IsolationLevel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.BeginTxx(ctx, &sql.TxOptions{Isolation: isoLevel})
	if err != nil {
		return err
	}

	if err := executor(ctx, tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("%s %w", rollbackErr.Error(), err)
		}
		return err
	}
	return tx.Commit()
}
