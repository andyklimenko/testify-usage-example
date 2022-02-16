package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/andyklimenko/testify-usage-example/api/entity"
	"github.com/jmoiron/sqlx"
)

const (
	qInsertUser          = "INSERT INTO users(first_name, last_name) VALUES($1, $2) RETURNING *"
	qGetUserByID         = "SELECT * FROM users WHERE id=$1"
	qGetUserByIdWithLock = "SELECT * FROM users WHERE id=$1 FOR UPDATE"
	qUpdateUser          = "UPDATE users SET first_name=$1, last_name=$2 WHERE id=$3 RETURNING *"
)

type dbUser struct {
	ID        string    `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	CreatedAt time.Time `db:"created_at"`
}

func (u dbUser) entity() entity.User {
	return entity.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt,
	}
}

func (s *Storage) InsertUser(ctx context.Context, u entity.User) (entity.User, error) {
	var res dbUser
	if err := s.db.GetContext(ctx, &res, qInsertUser, u.FirstName, u.LastName); err != nil {
		return entity.User{}, err
	}

	u.ID = res.ID
	u.CreatedAt = res.CreatedAt
	return u, nil
}

func (s *Storage) UserByID(ctx context.Context, id string) (entity.User, error) {
	var res dbUser
	if err := s.db.GetContext(ctx, &res, qGetUserByID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = entity.ErrNotFound
		}
		return entity.User{}, err
	}

	return res.entity(), nil
}

func (s *Storage) userByIDTx(ctx context.Context, tx *sqlx.Tx, id string) (dbUser, error) {
	var res dbUser
	if err := tx.GetContext(ctx, &res, qGetUserByIdWithLock, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = entity.ErrNotFound
		}
		return dbUser{}, fmt.Errorf("looking for existing user %s: %w", id, err)
	}

	return res, nil
}

func (s *Storage) UpdateUser(ctx context.Context, id string, u entity.User) (entity.User, error) {
	var updated entity.User
	txErr := runInTx(s.db, func(tx *sqlx.Tx) error {
		if _, err := s.userByIDTx(ctx, tx, id); err != nil {
			return err
		}

		var res dbUser
		if err := tx.GetContext(ctx, &res, qUpdateUser, u.FirstName, u.LastName, id); err != nil {
			return fmt.Errorf("execute update: %w", err)
		}

		updated = res.entity()
		return nil
	})

	return updated, txErr
}
