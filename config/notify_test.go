package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotifyLoad(t *testing.T) {
	t.Parallel()

	var cfg Notify
	assert.EqualError(t, cfg.load("notify"), ErrNoNotificationAddr.Error())

	require.NoError(t, os.Setenv("NOTIFY_ADDRESS", "test"))
	require.NoError(t, cfg.load("notify"))
	assert.Equal(t, "test", cfg.Addr)
}
