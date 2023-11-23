package config

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBLoad(t *testing.T) {
	t.Parallel()

	var db DB
	assert.ErrorIs(t, db.Load("test.db"), ErrNoDbDSN)

	require.NoError(t, os.Setenv("TEST_DB_DSN", "localhost"))
	assert.NoError(t, db.Load("test.db"))

	assert.Equal(t, "postgres", db.Driver)
	assert.Equal(t, "localhost", db.DSN)
}

func TestDBLoadOldWay(t *testing.T) {
	t.Parallel()

	var db DB
	if err := db.Load("test.db"); errors.Is(err, ErrNoDbDSN) {
		t.Fatal(err)
	}

	if err := os.Setenv("TEST_DB_DSN", "localhost"); err != nil {
		t.Fatal(err)
	}

	if err := db.Load("test.db"); err != nil {
		t.Fatal(err)
	}

	if db.Driver != "postgres" {
		t.Errorf("%s was expected byt %s got", "postgres", db.Driver)
	}

	if db.DSN != "localhost" {
		t.Errorf("%s was expected byt %s got", "localhost", db.DSN)
	}
}
