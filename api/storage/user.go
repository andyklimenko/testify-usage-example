package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/andyklimenko/testify-usage-example/api/entity"
)

const (
	qInsertUser  = "INSERT INTO users(first_name, last_name) VALUES($1, $2) RETURNING *"
	qGetUserByID = "SELECT * FROM users WHERE id=$1"
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
