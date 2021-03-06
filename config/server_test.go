package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerLoad(t *testing.T) {
	t.Parallel()

	var cfg Server
	assert.EqualError(t, cfg.load("test"), ErrNoServerAddr.Error())

	require.NoError(t, os.Setenv("TEST_ADDRESS", "http://localhost/hello/there"))
	require.NoError(t, cfg.load("test"))
	assert.Equal(t, "http://localhost/hello/there", cfg.Addr)
}
