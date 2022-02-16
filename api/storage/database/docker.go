package database

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/andyklimenko/testify-usage-example/api/storage/migrations"
	"github.com/fortytw2/dockertest"
	"github.com/jmoiron/sqlx"
)

type dockerRepo struct {
	db    *sqlx.DB
	once  sync.Once
	close func()
}

var repo dockerRepo

func DB() *sqlx.DB {
	return repo.db
}

func InitDockerDB() (func(), error) {
	repo.once.Do(func() {
		db, closer, err := withDocker()
		if err != nil {
			fmt.Println(err)
			return
		}

		repo.db = db
		repo.close = closer
	})

	return repo.close, nil
}

func withDocker() (*sqlx.DB, func(), error) {
	var db *sqlx.DB
	container, runErr := dockertest.RunContainer("circleci/postgres:9.6-alpine", "5432", func(addr string) error {
		hostPort := strings.Split(addr, ":")
		if len(hostPort) != 2 {
			return errors.New("wrong addr format")
		}

		port, err := strconv.Atoi(hostPort[1])
		if err != nil {
			return err
		}

		dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			"postgres", "postgres", hostPort[0], port, "postgres")

		db, err = setupDB("postgres", dsn)
		return err
	}, "--rm")
	if runErr != nil {
		return nil, func() { _ = db.Close() }, runErr
	}

	return db, func() {
		_ = db.Close()
		container.Shutdown()
	}, nil
}

func setupDB(driver string, dsn string) (*sqlx.DB, error) {
	db, err := DbConnect(driver, dsn)
	if err != nil {
		return nil, err
	}

	if err := migrations.Up(db, driver); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			return nil, closeErr
		}
		return nil, err
	}

	return db, err
}
