package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBLoad(t *testing.T) {
	t.Parallel()

	var db DB
	assert.EqualError(t, db.Load("test.db"), ErrNoDbDSN.Error())

	require.NoError(t, os.Setenv("TEST_DB_DSN", "localhost"))
	assert.NoError(t, db.Load("test.db"))

	assert.Equal(t, "postgres", db.Driver)
	assert.Equal(t, "localhost", db.DSN)
}
