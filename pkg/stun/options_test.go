package stun

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStunOptions(t *testing.T) {
	assert := require.New(t)

	t.Run("test_new_stun_options", func(t *testing.T) {
		logging := true

		options := NewStunOptions(logging)

		assert.Equal(logging, options.logging)
	})

	t.Run("test_new_stun_default_options", func(t *testing.T) {
		options := DefaultStunOptions()

		assert.Equal(DEFAULT_LOGGING, options.logging)
	})
}

func TestClientStunOptions(t *testing.T) {
	assert := require.New(t)

	t.Run("test_new_client_stun_options", func(t *testing.T) {
		logging := true
		maxMsgInQueue := 5

		options := NewClientStunOptions(logging, maxMsgInQueue)

		assert.Equal(logging, options.logging)
		assert.Equal(maxMsgInQueue, options.maxMsgInQueue)
	})

	t.Run("test_new_client_default_options", func(t *testing.T) {
		options := DefaultClientStunOptions()

		assert.Equal(DEFAULT_LOGGING, options.logging)
		assert.Equal(DEFAULT_MAX_MSG_IN_QUEUE, options.maxMsgInQueue)
	})
}
